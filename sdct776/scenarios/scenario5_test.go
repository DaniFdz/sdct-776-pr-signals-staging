//go:build sdct776

package scenarios

import (
	"os"
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func TestNewFlakyPaymentAuth(t *testing.T) {
	if _, err := os.Stat("/tmp/flaky_payment_auth.flag"); os.IsNotExist(err) {
		_ = os.WriteFile("/tmp/flaky_payment_auth.flag", []byte("1"), 0644)
		gotesting.GetTest(t).Fatal("payment authorization service timed out on attempt 1")
	}
	t.Log("payment authorization succeeded on retry")
}
