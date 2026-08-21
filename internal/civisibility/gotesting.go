// Package civisibility provides the narrow Go test wrapper used by the SDCT-776 scenarios.
package civisibility

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"unsafe"
)

const testFramework = "golang.org/pkg/testing"

type testInfo struct {
	name       string
	module     string
	suite      string
	function   *runtime.Func
	testFunc   func(*testing.T)
	sourceFile string
	sourceLine int
}

// RunM instruments root Go tests, runs them, and returns the test process exit code.
// It intentionally does not instrument subtests or benchmarks because this harness only
// needs one root pass and one root failure scenario.
func RunM(m *testing.M) int {
	run := startRun()
	defer run.stop()

	instrumentRootTests(m, run)
	exitCode := m.Run()
	run.close(exitCode)
	return exitCode
}

func instrumentRootTests(m *testing.M, run *testRun) {
	tests := internalTests(m)
	if tests == nil {
		return
	}

	wrapped := make([]testing.InternalTest, len(*tests))
	for i, test := range *tests {
		function := runtime.FuncForPC(reflect.ValueOf(test.F).Pointer())
		file, line := function.FileLine(function.Entry())
		info := testInfo{
			name:       test.Name,
			module:     moduleName(function.Name()),
			suite:      baseName(file),
			function:   function,
			testFunc:   test.F,
			sourceFile: run.relativeSourcePath(file),
			sourceLine: line,
		}
		wrapped[i] = testing.InternalTest{Name: test.Name, F: run.wrap(info)}
	}
	*tests = wrapped
}

func (run *testRun) wrap(info testInfo) func(*testing.T) {
	return func(t *testing.T) {
		event := run.startTest(info)
		defer func() {
			if recovered := recover(); recovered != nil {
				event.fail("panic", fmt.Sprint(recovered))
				event.close()
				panic(recovered)
			}

			switch {
			case t.Failed():
				event.fail("test failure", "Go test marked this scenario as failed")
			case t.Skipped():
				event.skip()
			default:
				event.pass()
			}
			event.close()
		}()
		info.testFunc(t)
	}
}

func internalTests(m *testing.M) *[]testing.InternalTest {
	value := reflect.Indirect(reflect.ValueOf(m)).FieldByName("tests")
	if !value.IsValid() {
		return nil
	}
	return (*[]testing.InternalTest)(unsafe.Pointer(value.UnsafeAddr()))
}
