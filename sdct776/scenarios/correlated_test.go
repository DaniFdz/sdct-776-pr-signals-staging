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

func TestRecoveryScenario(t *testing.T) {
	if !isAttempt1() {
		t.Log("re-run attempt passed successfully")
		return
	}
	gotesting.GetTest(t).Fatal("initial failure: verification of recovery reply")
}
