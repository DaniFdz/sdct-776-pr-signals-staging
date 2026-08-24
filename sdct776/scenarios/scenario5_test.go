//go:build sdct776

package scenarios

import (
	"os"
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func TestNewFlakyPaymentGatewayV2(t *testing.T) {
	if os.Getenv("SDCT776_FLAKY_FAIL") == "1" {
		gotesting.GetTest(t).Fatal("payment gateway connection reset on initial handshake")
	}
	t.Log("payment gateway handshake succeeded on retry")
}

func TestPersistentFailure(t *testing.T) {
	gotesting.GetTest(t).Fatal("order processing service returned 500")
}
