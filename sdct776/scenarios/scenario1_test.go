//go:build sdct776

package scenarios

import (
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func TestSingleFailure(t *testing.T) {
	gotesting.GetTest(t).Fatal("database connection timed out after 5000ms")
}
