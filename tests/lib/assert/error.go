package assert

import "testing"

func NoError(t *testing.T, err error, details ...interface{}) {
	if err != nil {
		t.Fatalf("unexpected error[%s] %s", errordetails(details...))
	}
}
