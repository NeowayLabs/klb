package azure_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func execWithIPC(t *testing.T, f fixture.F, exec func(string)) string {
	outfile, err := ioutil.TempFile("", "klb_script_ipc")
	if err != nil {
		t.Fatalf("error creating output file: %s", err)
	}
	defer os.Remove(outfile.Name()) // clean up

	exec(outfile.Name())

	f.Logger.Println("executed script, reading output")
	out, err := ioutil.ReadAll(outfile)
	if err != nil {
		t.Fatalf("error reading output file: %s", err)
	}
	return strings.Trim(strings.TrimSpace(string(out)), "\n")
}
