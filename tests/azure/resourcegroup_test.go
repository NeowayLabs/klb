package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/azure"
)

func TestHandleResourceGroupLifeCycle(t *testing.T) {
	// TODO: call nash stuff
	session := azure.NewSession(t)
	t.Log(session)
}
