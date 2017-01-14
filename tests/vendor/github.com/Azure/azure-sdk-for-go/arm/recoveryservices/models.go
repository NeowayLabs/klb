package recoveryservices

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 0.17.0.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

// SkuName enumerates the values for sku name.
type SkuName string

const (
	// RS0 specifies the rs0 state for sku name.
	RS0 SkuName = "RS0"
	// Standard specifies the standard state for sku name.
	Standard SkuName = "Standard"
)

// Resource is
type Resource struct {
	ID       *string             `json:"id,omitempty"`
	Name     *string             `json:"name,omitempty"`
	Type     *string             `json:"type,omitempty"`
	Location *string             `json:"location,omitempty"`
	Sku      *Sku                `json:"sku,omitempty"`
	Tags     *map[string]*string `json:"tags,omitempty"`
}

// Sku is identifies the unique system identifier for each Azure resource.
type Sku struct {
	Name SkuName `json:"name,omitempty"`
}

// Vault is resource information, as returned by the resource provider.
type Vault struct {
	autorest.Response `json:"-"`
	ID                *string             `json:"id,omitempty"`
	Name              *string             `json:"name,omitempty"`
	Type              *string             `json:"type,omitempty"`
	Location          *string             `json:"location,omitempty"`
	Sku               *Sku                `json:"sku,omitempty"`
	Tags              *map[string]*string `json:"tags,omitempty"`
	Etag              *string             `json:"etag,omitempty"`
	Properties        *VaultProperties    `json:"properties,omitempty"`
}

// VaultList is the response model for Vault.
type VaultList struct {
	autorest.Response `json:"-"`
	Value             *[]Vault `json:"value,omitempty"`
	Null              *string  `json:",omitempty"`
}

// VaultListPreparer prepares a request to retrieve the next set of results. It returns
// nil if no more results exist.
func (client VaultList) VaultListPreparer() (*http.Request, error) {
	if client.Null == nil || len(to.String(client.Null)) <= 0 {
		return nil, nil
	}
	return autorest.Prepare(&http.Request{},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(client.Null)))
}

// VaultProperties is properties of the vault.
type VaultProperties struct {
	ProvisioningState *string `json:"provisioningState,omitempty"`
}
