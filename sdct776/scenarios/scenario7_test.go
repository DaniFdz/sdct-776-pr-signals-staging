//go:build sdct776

package scenarios

import (
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func TestUnitAuth(t *testing.T) {
	gotesting.GetTest(t).Fatal("unauthorized unit test")
}

func TestIntegrationDB(t *testing.T) {
	gotesting.GetTest(t).Fatal("database pool closed")
}

func TestE2ECheckout(t *testing.T) {
	gotesting.GetTest(t).Fatal("e2e checkout gateway 502")
}
