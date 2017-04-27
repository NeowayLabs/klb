package azure_test

import (
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

// TestExamples aims to just keep our examples working,
// but it does not validates if infrastructure has really been built.
// Only one parameter is supported on the examples, the resource group.
func TestExamples(t *testing.T) {
	t.Parallel()
	t.Skip()
	// FIXME: Try again must destroy the entire resgroup
	// Not sure if this is a good idea x_x
	examples := []struct {
		name   string
		script string
	}{
		{
			name:   "managedDisks",
			script: "../../examples/azure/managed-disks/build.sh",
		},
	}

	timeout := 30 * time.Minute

	for _, example := range examples {
		fixture.Run(t, example.name, timeout, location, func(t *testing.T, f fixture.F) {
			f.Shell.Run(example.script, f.ResGroupName)
		})
	}
}
