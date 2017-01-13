package azure

import "fmt"

func newID(resource string, method string, name string) string {
	return fmt.Sprintf("%s.%s:%s", resource, method, name)
}
