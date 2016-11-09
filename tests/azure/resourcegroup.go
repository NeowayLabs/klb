package azure

import (
	"io/ioutil"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/go-autorest/autorest/to"
)

type Resources struct {
	client  resources.Client
	session *Session
}

func aborterr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func (r *Resources) List(t *testing.T) string {
	res, err := r.client.List("some", "some", nil)
	aborterr(t, err)
	body, err := ioutil.ReadAll(res.Body)
	aborterr(t, err)
	return string(body)
}

func (r *Resources) TryShit(t *testing.T) {
	// Copied from examples
	ac := storage.NewAccountsClient(r.session.subscriptionID)
	ac.Authorizer = r.session.token
	name := "justSomeShit"
	cna, err := ac.CheckNameAvailability(
		storage.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(name),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts")})
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !to.Bool(cna.NameAvailable) {
		t.Fatalf("%s is unavailable -- try with another name\n", name)
	}
}

func NewResources(t *testing.T, s *Session) *Resources {
	rg := &Resources{
		client:  resources.NewClient(s.subscriptionID),
		session: s,
	}
	rg.client.Authorizer = s.token
	return rg
}
