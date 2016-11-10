package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/azure"
)

const (
	ResourceGroupName = "klb-resource-group-test"
)

func TestHandleResourceGroupLifeCycle(t *testing.T) {
	// TODO: call nash stuff
	session := azure.NewSession(t)
	resources := azure.NewResources(t, session)
	//defer resources.Delete(t, ResourceGroupName)
	resources.Check(t, ResourceGroupName)
}
