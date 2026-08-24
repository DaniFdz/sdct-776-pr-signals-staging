//go:build sdct776

package scenarios

import (
	"os"
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func TestFreshFlakyCheckoutTransaction(t *testing.T) {
	test := gotesting.GetTest(t)
	if _, err := os.Stat("/tmp/flaky_checkout.flag"); os.IsNotExist(err) {
		_ = os.WriteFile("/tmp/flaky_checkout.flag", []byte("1"), 0644)
		test.Fatal("payment gateway connection reset on initial handshake")
	}
	test.SetTag("test.is_retry", "true")
	test.SetTag("test.retry_reason", "auto_test_retry")
	t.Log("payment gateway handshake succeeded on retry")
}
