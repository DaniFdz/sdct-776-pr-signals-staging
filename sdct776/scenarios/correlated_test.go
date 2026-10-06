//go:build sdct776

package scenarios

import (
	"os"
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func isAttempt1() bool {
	attempt := os.Getenv("GITHUB_RUN_ATTEMPT")
	return attempt == "" || attempt == "1"
}

func TestRetryPolicy(t *testing.T) {
	if !isAttempt1() {
		t.Log("re-run attempt passed successfully")
		return
	}
	gotesting.GetTest(t).Fatal("expect(received).toBe(expected) // Object.is equality\n\nExpected: \"CVE-2023-44270\"\nReceived: undefined\n\nIgnored nodes: comments, script, style")
}

func TestRetryPolicyBackoff(t *testing.T) {
	if !isAttempt1() {
		t.Log("re-run attempt passed successfully")
		return
	}
	gotesting.GetTest(t).Fatal("retry backoff did not advance")
}

func TestRetryPolicyBudget(t *testing.T) {
	if !isAttempt1() {
		t.Log("re-run attempt passed successfully")
		return
	}
	gotesting.GetTest(t).Fatal("retry budget was exhausted too early")
}
