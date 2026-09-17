//go:build sdct776

package scenarios

import (
	"os"
	"testing"

	"github.com/DaniFdz/sdct-776-pr-signals-staging/main/civisibility/integrations/gotesting"
)

func TestMain(m *testing.M) {
	os.Exit(gotesting.RunM(m))
}

func TestObjectIsEquality(t *testing.T) {
	gotesting.GetTest(t).Fatal("expect(received).toBe(expected) // Object.is equality\n\nExpected: \"CVE-2023-44270\"\nReceived: undefined\n\nIgnored nodes: comments, script, style")
}
