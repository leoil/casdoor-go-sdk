// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: add enforcer, syncer, subscription, session, resource, record, product,
// pricing, plan, payment, model, message, account, group, chat, cert, adapter,

package casdoorsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// client is a shared http Client.
var client HttpClient = &http.Client{}

// SetHttpClient sets custom http Client.
func SetHttpClient(httpClient HttpClient) {
	client = httpClient
}

// HttpClient interface has the method required to use a type as custom http client.
// The net/*http.Client type satisfies this interface.
type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Response struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Data2  interface{} `json:"data2"`
}

// DoGetResponse is a general function to get response from param url through HTTP Get method.
func (c *Client) DoGetResponse(url string) (*Response, error) {
	respBytes, err := c.doGetBytesRawWithoutCheck(url)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return nil, err
	}

	if response.Status != "ok" {
		return nil, fmt.Errorf(response.Msg)
	}

	return &response, nil
}

// DoGetBytes is a general function to get response data in bytes from param url through HTTP Get method.
func (c *Client) DoGetBytes(url string) ([]byte, error) {
	response, err := c.DoGetResponse(url)
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DoGetBytesRaw is a general function to get response from param url through HTTP Get method.
func (c *Client) DoGetBytesRaw(url string) ([]byte, error) {
	respBytes, err := c.doGetBytesRawWithoutCheck(url)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(respBytes, &response)
	if err == nil && response.Status == "error" {
		return nil, errors.New(response.Msg)
	}

	return respBytes, nil
}

func (c *Client) DoPost(action string, queryMap map[string]string, postBytes []byte, isForm, isFile bool) (*Response, error) {
	url := c.GetUrl(action, queryMap)

	var err error
	var contentType string
	var body io.Reader
	if isForm {
		if isFile {
			contentType, body, err = createFormFile(map[string][]byte{"file": postBytes})
			if err != nil {
				return nil, err
			}
		} else {
			var params map[string]string
			err = json.Unmarshal(postBytes, &params)
			if err != nil {
				return nil, err
			}

			contentType, body, err = createForm(params)
			if err != nil {
				return nil, err
			}
		}
	} else {
		contentType = "text/plain;charset=UTF-8"
		body = bytes.NewReader(postBytes)
	}

	respBytes, err := c.DoPostBytesRaw(url, contentType, body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return nil, err
	}

	if response.Status != "ok" {
		return nil, fmt.Errorf(response.Msg)
	}

	return &response, nil
}

// DoPostBytesRaw is a general function to post a request from url, body through HTTP Post method.
func (c *Client) DoPostBytesRaw(url string, contentType string, body io.Reader) ([]byte, error) {
	if contentType == "" {
		contentType = "text/plain;charset=UTF-8"
	}

	var resp *http.Response

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.ClientId, c.ClientSecret)
	req.Header.Set("Content-Type", contentType)

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}

// doGetBytesRawWithoutCheck is a general function to get response from param url through HTTP Get method without checking response status
func (c *Client) doGetBytesRawWithoutCheck(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.ClientId, c.ClientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

// modifyOrganization is an encapsulation of permission CUD(Create, Update, Delete) operations.
// possible actions are `add-organization`, `update-organization`, `delete-organization`,
func (c *Client) modifyOrganization(action string, organization *Organization, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", organization.Owner, organization.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	// organization.Owner = c.OrganizationName
	postBytes, err := json.Marshal(organization)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyApplication is an encapsulation of permission CUD(Create, Update, Delete) operations.
// possible actions are `add-application`, `update-application`, `delete-application`,
func (c *Client) modifyApplication(action string, application *Application, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", application.Owner, application.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	if application.Owner == "" {
		application.Owner = "admin"
	}
	postBytes, err := json.Marshal(application)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyProvider is an encapsulation of permission CUD(Create, Update, Delete) operations.
// possible actions are `add-provider`, `update-provider`, `delete-provider`,
func (c *Client) modifyProvider(action string, provider *Provider, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", provider.Owner, provider.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	provider.Owner = c.OrganizationName
	postBytes, err := json.Marshal(provider)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifySession is an encapsulation of permission CUD(Create, Update, Delete) operations.
// possible actions are `add-session`, `update-session`, `delete-session`,
func (c *Client) modifySession(action string, session *Session, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", session.Owner, session.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	session.Owner = c.OrganizationName
	postBytes, err := json.Marshal(session)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyUser is an encapsulation of user CUD(Create, Update, Delete) operations.
// possible actions are `add-user`, `update-user`, `delete-user`,
func (c *Client) modifyUser(action string, user *User, columns []string) (*Response, bool, error) {
	return c.modifyUserById(action, user.GetId(), user, columns)
}

func (c *Client) modifyUserById(action string, id string, user *User, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": id,
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	user.Owner = c.OrganizationName
	postBytes, err := json.Marshal(user)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyPermission is an encapsulation of permission CUD(Create, Update, Delete) operations.
// possible actions are `add-permission`, `update-permission`, `delete-permission`,
func (c *Client) modifyPermission(action string, permission *Permission, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", permission.Owner, permission.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	permission.Owner = c.OrganizationName
	postBytes, err := json.Marshal(permission)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyRole is an encapsulation of role CUD(Create, Update, Delete) operations.
// possible actions are `add-role`, `update-role`, `delete-role`,
func (c *Client) modifyRole(action string, role *Role, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", role.Owner, role.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	role.Owner = c.OrganizationName
	postBytes, err := json.Marshal(role)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyCert is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-cert`, `update-cert`, `delete-cert`,
func (c *Client) modifyCert(action string, cert *Cert, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", cert.Owner, cert.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	cert.Owner = c.OrganizationName
	postBytes, err := json.Marshal(cert)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyEnforcer is an encapsulation of cert CUD(Create, Update, Delete) operations.
func (c *Client) modifyEnforcer(action string, enforcer *Enforcer, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", enforcer.Owner, enforcer.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	enforcer.Owner = c.OrganizationName
	postBytes, err := json.Marshal(enforcer)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyEnforcer is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-group`, `update-group`, `delete-group`,
func (c *Client) modifyGroup(action string, group *Group, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", group.Owner, group.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	group.Owner = c.OrganizationName
	postBytes, err := json.Marshal(group)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyAdapter is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-adapter`, `update-adapter`, `delete-adapter`,
func (c *Client) modifyAdapter(action string, adapter *Adapter, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", adapter.Owner, adapter.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	adapter.Owner = c.OrganizationName
	postBytes, err := json.Marshal(adapter)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyModel is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-model`, `update-model`, `delete-model`,
func (c *Client) modifyModel(action string, model *Model, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", model.Owner, model.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	model.Owner = c.OrganizationName
	postBytes, err := json.Marshal(model)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyProduct is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-product`, `update-product`, `delete-product`,
func (c *Client) modifyProduct(action string, product *Product, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", product.Owner, product.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	product.Owner = c.OrganizationName
	postBytes, err := json.Marshal(product)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyPayment is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-payment`, `update-payment`, `delete-payment`,
func (c *Client) modifyPayment(action string, payment *Payment, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", payment.Owner, payment.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	payment.Owner = c.OrganizationName
	postBytes, err := json.Marshal(payment)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyPlan is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-plan`, `update-plan`, `delete-plan`,
func (c *Client) modifyPlan(action string, plan *Plan, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", plan.Owner, plan.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	plan.Owner = c.OrganizationName
	postBytes, err := json.Marshal(plan)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyPricing is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-pricing`, `update-pricing`, `delete-pricing`,
func (c *Client) modifyPricing(action string, pricing *Pricing, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", pricing.Owner, pricing.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	pricing.Owner = c.OrganizationName
	postBytes, err := json.Marshal(pricing)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifySubscription is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-subscription`, `update-subscription`, `delete-subscription`,
func (c *Client) modifySubscription(action string, subscription *Subscription, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", subscription.Owner, subscription.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	subscription.Owner = c.OrganizationName
	postBytes, err := json.Marshal(subscription)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifySyner is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-syncer`, `update-syncer`, `delete-syncer`,
func (c *Client) modifySyncer(action string, syncer *Syncer, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", syncer.Owner, syncer.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	syncer.Owner = c.OrganizationName
	postBytes, err := json.Marshal(syncer)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}

// modifyWebhook is an encapsulation of cert CUD(Create, Update, Delete) operations.
// possible actions are `add-webhook`, `update-webhook`, `delete-webhook`,
func (c *Client) modifyWebhook(action string, webhook *Webhook, columns []string) (*Response, bool, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", webhook.Owner, webhook.Name),
	}

	if len(columns) != 0 {
		queryMap["columns"] = strings.Join(columns, ",")
	}

	webhook.Owner = c.OrganizationName
	postBytes, err := json.Marshal(webhook)
	if err != nil {
		return nil, false, err
	}

	resp, err := c.DoPost(action, queryMap, postBytes, false, false)
	if err != nil {
		return nil, false, err
	}

	return resp, resp.Data == "Affected", nil
}
