package civisibility

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func githubTags() map[string]string {
	serverURL := strings.TrimSuffix(firstEnv("GITHUB_SERVER_URL"), "/")
	if serverURL == "" {
		serverURL = "https://github.com"
	}

	repository := firstEnv("GITHUB_REPOSITORY")
	rawRepositoryURL := serverURL + "/" + repository
	runID := firstEnv("GITHUB_RUN_ID")
	attempt := firstEnv("GITHUB_RUN_ATTEMPT")
	commitSHA := firstEnv("GITHUB_SHA")
	branch, tag := branchOrTag(firstEnv("GITHUB_HEAD_REF", "GITHUB_REF"))

	pipelineURL := fmt.Sprintf("%s/actions/runs/%s", rawRepositoryURL, runID)
	if attempt != "" {
		pipelineURL = fmt.Sprintf("%s/attempts/%s", pipelineURL, attempt)
	}

	tags := map[string]string{
		"ci.provider.name":   "github",
		"ci.workspace_path":  firstEnv("GITHUB_WORKSPACE"),
		"ci.pipeline.id":     runID,
		"ci.pipeline.name":   firstEnv("GITHUB_WORKFLOW"),
		"ci.pipeline.number": firstEnv("GITHUB_RUN_NUMBER"),
		"ci.pipeline.url":    pipelineURL,
		"ci.job.name":        firstEnv("GITHUB_JOB"),
		"ci.job.url":         fmt.Sprintf("%s/commit/%s/checks", rawRepositoryURL, commitSHA),
		"git.repository_url": rawRepositoryURL + ".git",
		"git.commit.sha":     commitSHA,
		"git.branch":         branch,
		"git.tag":            tag,
		"os.platform":        runtime.GOOS,
		"os.architecture":    runtime.GOARCH,
		"runtime.name":       runtime.Compiler,
		"runtime.version":    runtime.Version(),
	}

	if correlation, err := json.Marshal(map[string]string{
		"GITHUB_SERVER_URL":  firstEnv("GITHUB_SERVER_URL"),
		"GITHUB_REPOSITORY":  repository,
		"GITHUB_RUN_ID":      runID,
		"GITHUB_RUN_ATTEMPT": attempt,
	}); err == nil {
		tags["_dd.ci.env_vars"] = string(correlation)
	}

	for key, value := range tags {
		if value == "" {
			delete(tags, key)
		}
	}
	return tags
}

func testCommand() string {
	if len(os.Args) == 0 {
		return "go test"
	}
	return strings.TrimSpace(strings.Join(append([]string{filepath.Base(os.Args[0])}, os.Args[1:]...), " "))
}

func branchOrTag(ref string) (branch, tag string) {
	if strings.Contains(ref, "tags/") {
		return "", strings.TrimPrefix(ref, "refs/tags/")
	}
	return strings.TrimPrefix(ref, "refs/heads/"), ""
}

func moduleName(functionName string) string {
	lastSlash := strings.LastIndexByte(functionName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	firstDot := strings.IndexByte(functionName[lastSlash:], '.')
	if firstDot < 0 {
		return functionName
	}
	return functionName[:lastSlash+firstDot]
}

func baseName(path string) string {
	return filepath.Base(path)
}

func firstEnv(names ...string) string {
	for _, name := range names {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}
	return ""
}
