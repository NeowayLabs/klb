package azure

import (
	"io/ioutil"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

type Resources struct {
	client resources.Client
}

func aborterr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
func (r *Resources) List(t *testing.T) string {
	res, err := r.client.List("", "", nil)
	aborterr(t, err)
	body, err := ioutil.ReadAll(res.Body)
	aborterr(t, err)
	return string(body)
}

func NewResources(t *testing.T, s *Session) *Resources {
	rg := &Resources{client: resources.NewClient(s.subscriptionID)}
	rg.client.Authorizer = s.token
	return rg
}
