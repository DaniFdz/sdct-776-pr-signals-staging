//go:build sdct776

package scenarios

import (
	"os"
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func TestFlakyPaymentGateway(t *testing.T) {
	if os.Getenv("SDCT776_FLAKY_FAIL") == "1" {
		gotesting.GetTest(t).Fatal("payment gateway connection reset by peer")
	}
	t.Log("payment gateway transaction succeeded")
}
