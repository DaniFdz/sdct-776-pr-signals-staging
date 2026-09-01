//go:build sdct776

package scenarios

import (
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func TestRetryPolicy(t *testing.T) {
	gotesting.GetTest(t).Fatal("expected 3 retry attempts, got 1")
}

func TestRetryPolicyBackoff(t *testing.T) {
	gotesting.GetTest(t).Fatal("retry backoff did not advance")
}

func TestRetryPolicyBudget(t *testing.T) {
	gotesting.GetTest(t).Fatal("retry budget was exhausted too early")
}
