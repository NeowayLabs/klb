package azure

import "fmt"

func getID(resource string, method string, name string) string {
	return fmt.Sprintf("%s.%s:%s", resource, method, name)
}
