//go:build sdct776

package scenarios

import (
	"os"
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/internal/civisibility"
)

func TestMain(m *testing.M) {
	os.Exit(civisibility.RunM(m))
}

func TestSDCT776OnboardPass(t *testing.T) {
	t.Log("sdct776-onboard-pass: intentional passing Test Visibility scenario")
}

func TestSDCT776OnboardFailure(t *testing.T) {
	t.Fatal("sdct776-onboard-failure: intentional deterministic CI and Test Visibility failure")
}
