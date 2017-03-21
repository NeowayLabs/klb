package apimanagement

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
// Code generated by Microsoft (R) AutoRest Code Generator 1.0.1.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// ReportsClient is the use these REST APIs for performing operations on
// entities like API, Product, and Subscription associated with your Azure API
// Management deployment.
type ReportsClient struct {
	ManagementClient
}

// NewReportsClient creates an instance of the ReportsClient client.
func NewReportsClient(subscriptionID string) ReportsClient {
	return NewReportsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewReportsClientWithBaseURI creates an instance of the ReportsClient client.
func NewReportsClientWithBaseURI(baseURI string, subscriptionID string) ReportsClient {
	return ReportsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// ListByService lists report records.
//
// resourceGroupName is the name of the resource group. serviceName is the name
// of the API Management service. aggregation is report aggregation. filter is
// the filter to apply on the operation. top is number of records to return.
// skip is number of records to skip. interval is by time interval. This value
// is only applicable to ByTime aggregation. Interval must be multiple of 15
// minutes and may not be zero. The value should be in ISO  8601 format
// (http://en.wikipedia.org/wiki/ISO_8601#Durations).This code can be used to
// convert TimSpan to a valid interval string: XmlConvert.ToString(new
// TimeSpan(hours, minutes, secconds))
func (client ReportsClient) ListByService(resourceGroupName string, serviceName string, aggregation ReportsAggregation, filter string, top *int32, skip *int32, interval string) (result ReportCollection, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: serviceName,
			Constraints: []validation.Constraint{{Target: "serviceName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "serviceName", Name: validation.MinLength, Rule: 1, Chain: nil},
				{Target: "serviceName", Name: validation.Pattern, Rule: `^[a-zA-Z](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?$`, Chain: nil}}},
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil}}}}},
		{TargetValue: skip,
			Constraints: []validation.Constraint{{Target: "skip", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "skip", Name: validation.InclusiveMinimum, Rule: 0, Chain: nil}}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "apimanagement.ReportsClient", "ListByService")
	}

	req, err := client.ListByServicePreparer(resourceGroupName, serviceName, aggregation, filter, top, skip, interval)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "apimanagement.ReportsClient", "ListByService", nil, "Failure preparing request")
	}

	resp, err := client.ListByServiceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "apimanagement.ReportsClient", "ListByService", resp, "Failure sending request")
	}

	result, err = client.ListByServiceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.ReportsClient", "ListByService", resp, "Failure responding to request")
	}

	return
}

// ListByServicePreparer prepares the ListByService request.
func (client ReportsClient) ListByServicePreparer(resourceGroupName string, serviceName string, aggregation ReportsAggregation, filter string, top *int32, skip *int32, interval string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"aggregation":       autorest.Encode("path", aggregation),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"serviceName":       autorest.Encode("path", serviceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": client.APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if skip != nil {
		queryParameters["$skip"] = autorest.Encode("query", *skip)
	}
	if len(interval) > 0 {
		queryParameters["interval"] = autorest.Encode("query", interval)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/reports/{aggregation}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListByServiceSender sends the ListByService request. The method will close the
// http.Response Body if it receives an error.
func (client ReportsClient) ListByServiceSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// ListByServiceResponder handles the response to the ListByService request. The method always
// closes the http.Response Body.
func (client ReportsClient) ListByServiceResponder(resp *http.Response) (result ReportCollection, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByServiceNextResults retrieves the next set of results, if any.
func (client ReportsClient) ListByServiceNextResults(lastResults ReportCollection) (result ReportCollection, err error) {
	req, err := lastResults.ReportCollectionPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "apimanagement.ReportsClient", "ListByService", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListByServiceSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "apimanagement.ReportsClient", "ListByService", resp, "Failure sending next results request")
	}

	result, err = client.ListByServiceResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apimanagement.ReportsClient", "ListByService", resp, "Failure responding to next results request")
	}

	return
}
