//go:build sdct776

package scenarios

import (
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func TestUserAuth(t *testing.T) {
	gotesting.GetTest(t).Fatal("unauthorized token")
}

func TestUserPermissions(t *testing.T) {
	gotesting.GetTest(t).Fatal("permission denied for resource")
}

func TestPaymentCapture(t *testing.T) {
	gotesting.GetTest(t).Fatal("payment gateway unreachable")
}

func TestRefund(t *testing.T) {
	gotesting.GetTest(t).Fatal("invalid transaction state for refund")
}
