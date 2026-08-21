package civisibility

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	originTag          = "_dd.origin"
	originValue        = "ciapp-test"
	testSessionIDTag   = "test_session_id"
	testModuleIDTag    = "test_module_id"
	testSuiteIDTag     = "test_suite_id"
	testTypeTag        = "test.type"
	testTypeValue      = "test"
	testFrameworkTag   = "test.framework"
	testFrameworkValue = "golang.org/pkg/testing"
	testVersionTag     = "test.framework_version"
	testModuleTag      = "test.module"
	testSuiteTag       = "test.suite"
	testNameTag        = "test.name"
	testStatusTag      = "test.status"
	testSourceFileTag  = "test.source.file"
	testSourceLineTag  = "test.source.start"
	testCommandTag     = "test.command"
	testWorkingDirTag  = "test.working_directory"
	testExitCodeTag    = "test.exit_code"
)

type testRun struct {
	sessionID uint64
	session   tracer.Span
	tags      map[string]string
	workspace string
}

type testEvent struct {
	test   tracer.Span
	suite  tracer.Span
	module tracer.Span
}

func startRun() *testRun {
	// The reference wrapper enables CI Visibility before starting the tracer so
	// the agentless CI Visibility transport is selected.
	_ = os.Setenv("DD_CIVISIBILITY_ENABLED", "1")
	_ = os.Setenv("DD_TRACE_SAMPLE_RATE", "1")
	tracer.Start()

	tags := githubTags()
	workspace := tags["ci.workspace_path"]
	command := testCommand()
	tags[testCommandTag] = command
	tags[testWorkingDirTag] = relativeWorkingDirectory(workspace)

	options := commonOptions(tags,
		tracer.ResourceName("test_session."+command),
		tracer.SpanType("test_session_end"),
		tracer.Tag(testTypeTag, testTypeValue),
	)
	session, _ := tracer.StartSpanFromContext(context.Background(), "test_session", options...)
	sessionID := session.Context().SpanID()
	session.SetTag(testSessionIDTag, fmt.Sprint(sessionID))

	return &testRun{sessionID: sessionID, session: session, tags: tags, workspace: workspace}
}

func (run *testRun) startTest(info testInfo) *testEvent {
	moduleOptions := commonOptions(run.tags,
		tracer.ResourceName(info.module),
		tracer.SpanType("test_module_end"),
		tracer.Tag(testTypeTag, testTypeValue),
		tracer.Tag(testModuleTag, info.module),
		tracer.Tag(testFrameworkTag, testFrameworkValue),
		tracer.Tag(testVersionTag, runtime.Version()),
	)
	module, _ := tracer.StartSpanFromContext(context.Background(), "test_module", moduleOptions...)
	moduleID := module.Context().SpanID()
	module.SetTag(testSessionIDTag, fmt.Sprint(run.sessionID))
	module.SetTag(testModuleIDTag, fmt.Sprint(moduleID))

	suiteOptions := commonOptions(run.tags,
		tracer.ResourceName(info.suite),
		tracer.SpanType("test_suite_end"),
		tracer.Tag(testTypeTag, testTypeValue),
		tracer.Tag(testModuleTag, info.module),
		tracer.Tag(testFrameworkTag, testFrameworkValue),
		tracer.Tag(testVersionTag, runtime.Version()),
		tracer.Tag(testSuiteTag, info.suite),
	)
	suite, _ := tracer.StartSpanFromContext(context.Background(), "test_suite", suiteOptions...)
	suiteID := suite.Context().SpanID()
	suite.SetTag(testSessionIDTag, fmt.Sprint(run.sessionID))
	suite.SetTag(testModuleIDTag, fmt.Sprint(moduleID))
	suite.SetTag(testSuiteIDTag, fmt.Sprint(suiteID))

	testOptions := commonOptions(run.tags,
		tracer.ResourceName(info.suite+"."+info.name),
		tracer.SpanType("test"),
		tracer.Tag(testTypeTag, testTypeValue),
		tracer.Tag(testModuleTag, info.module),
		tracer.Tag(testFrameworkTag, testFrameworkValue),
		tracer.Tag(testVersionTag, runtime.Version()),
		tracer.Tag(testSuiteTag, info.suite),
		tracer.Tag(testNameTag, info.name),
		tracer.Tag(testSourceFileTag, info.sourceFile),
		tracer.Tag(testSourceLineTag, info.sourceLine),
	)
	test, _ := tracer.StartSpanFromContext(context.Background(), "test", testOptions...)
	test.SetTag(testSessionIDTag, fmt.Sprint(run.sessionID))
	test.SetTag(testModuleIDTag, fmt.Sprint(moduleID))
	test.SetTag(testSuiteIDTag, fmt.Sprint(suiteID))

	return &testEvent{test: test, suite: suite, module: module}
}

func (event *testEvent) pass() {
	event.test.SetTag(testStatusTag, "pass")
}

func (event *testEvent) skip() {
	event.test.SetTag(testStatusTag, "skip")
}

func (event *testEvent) fail(kind, message string) {
	event.test.SetTag(testStatusTag, "fail")
	event.test.SetTag(ext.Error, true)
	event.test.SetTag(ext.ErrorType, kind)
	event.test.SetTag(ext.ErrorMsg, message)
	event.suite.SetTag(ext.Error, true)
	event.module.SetTag(ext.Error, true)
}

func (event *testEvent) close() {
	now := time.Now()
	event.test.Finish(tracer.FinishTime(now))
	event.suite.Finish(tracer.FinishTime(now))
	event.module.Finish(tracer.FinishTime(now))
}

func (run *testRun) close(exitCode int) {
	run.session.SetTag(testExitCodeTag, exitCode)
	if exitCode == 0 {
		run.session.SetTag(testStatusTag, "pass")
	} else {
		run.session.SetTag(testStatusTag, "fail")
		run.session.SetTag(ext.Error, true)
	}
	run.session.Finish()
	tracer.Flush()
}

func (run *testRun) stop() {
	tracer.Stop()
}

func commonOptions(tags map[string]string, options ...tracer.StartSpanOption) []tracer.StartSpanOption {
	options = append(options,
		tracer.Tag(originTag, originValue),
		tracer.Tag(ext.ManualKeep, true),
	)
	for key, value := range tags {
		options = append(options, tracer.Tag(key, value))
	}
	return options
}

func relativeWorkingDirectory(workspace string) string {
	workingDirectory, err := os.Getwd()
	if err != nil || workspace == "" {
		return workingDirectory
	}
	relative, err := filepath.Rel(workspace, workingDirectory)
	if err != nil {
		return workingDirectory
	}
	return relative
}

func (run *testRun) relativeSourcePath(file string) string {
	if run.workspace == "" {
		return file
	}
	relative, err := filepath.Rel(run.workspace, file)
	if err != nil || strings.HasPrefix(relative, "..") {
		return file
	}
	return filepath.ToSlash(relative)
}
