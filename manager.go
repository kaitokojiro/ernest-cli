/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/fatih/color"
)

// Manager manages all api communications
type Manager struct {
	URL string `json:"url"`
}

// Token holds the JWT token that is received when authenticating
type Token struct {
	Token string `json:"token"`
}

// Service ...
type Service struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Datacenter int       `json:"datacenter_id"`
	Version    time.Time `json:"version"`
	Status     string    `json:"status"`
	Definition string    `json:"definition"`
	Result     string    `json:"result"`
	Endpoint   string    `json:"endpoint"`
}

// Datacenter ...
type Datacenter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Session ...
type Session struct {
	UserID  int    `json:"id"`
	GroupID int    `json:"group_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"admin"`
}

func (m *Manager) client() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	return client
}

func (m *Manager) doRequest(url, method string, payload []byte, token string, contentType string) (string, *http.Response, error) {
	url = m.URL + url
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := m.client().Do(req)

	if err != nil {
		return err.Error(), resp, err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}
	body := string(responseBody)

	if resp.StatusCode != 200 {
		return string(body), resp, errors.New(resp.Status)
	}
	return string(body), resp, nil
}

func (m *Manager) createClient(token string, name string) (string, error) {
	payload := []byte(`{"name":"` + name + `"}`)
	body, _, err := m.doRequest("/api/groups/", "POST", payload, token, "")
	if err != nil {
		return body, err
	}

	color.Green("SUCCESS: Group " + name + " created")

	var group struct {
		ID int `json:"id"`
	}
	err = json.Unmarshal([]byte(body), &group)
	if err != nil {
		return "", errors.New("ERROR: Couldn't read response from server")
	}
	return strconv.Itoa(group.ID), nil
}

func (m *Manager) getSession(token string) (session Session, err error) {
	res, _, err := m.doRequest("/api/session/", "GET", nil, token, "application/yaml")
	if err != nil {
		return session, err
	}
	err = json.Unmarshal([]byte(res), &session)
	if err != nil {
		return session, err
	}
	return session, nil
}

// ********************* Create *******************

// CreateDatacenter ...
func (m *Manager) CreateDatacenter(token string, name string, user string, password string, url string, network string, vseURL string) (string, error) {
	payload := []byte(`{"name": "` + name + `", "type": "vcloud", "region": "LON-001", "username":"` + user + `", "password":"` + password + `", "external_network":"` + network + `", "vcloud_url":"` + url + `", "vse_url":"` + vseURL + `"}`)
	body, _, err := m.doRequest("/api/datacenters/", "POST", payload, token, "")
	if err != nil {
		return body, err
	}
	color.Green("SUCCESS: Datacenter " + name + " created")
	return body, err
}

// ********************* Get *******************

// GetUUID ...
func (m *Manager) GetUUID(token string, payload []byte) string {
	id, err := buildServiceUUID(payload)
	if err != nil {
		log.Fatal(err)
	}
	body, _, _ := m.doRequest("/api/services/uuid/", "POST", []byte(`{"id":"`+id+`"}`), token, "")
	var dat map[string]interface{}
	err = json.Unmarshal([]byte(body), &dat)
	if err != nil {
		return ""
	}

	if str, ok := dat["uuid"].(string); ok {
		return str
	}
	return ""
}

// ********************* Apply *******************

// Apply ...
func (m *Manager) Apply(token string, path string, monit bool) (string, error) {
	payload, err := ioutil.ReadFile(path)
	if err != nil {
		color.Red(err.Error())
		return "", nil
	}

	color.Green("Environment creation requested")
	println("Ernest will show you all output from your requested service creation")
	println("You can cancel at any moment with Ctrl+C, even the service is still being created, you won't have any output")

	streamID := m.GetUUID(token, payload)
	if streamID == "" {
		color.Red("Please log in")
		return "", nil
	}

	if monit == true {
		go Monitorize(m.URL, token, streamID)
	} else {
		println("Additionally you can trace your service on ernest monitor tool with id: " + streamID)
	}

	if body, _, err := m.doRequest("/api/services/", "POST", payload, token, "application/yaml"); err != nil {
		return "", errors.New(body)
	}
	if monit == true {
		runtime.Goexit()
	}
	return streamID, nil
}

// ********************* Destroy *******************

// Destroy ...
func (m *Manager) Destroy(token string, name string, monit bool) error {
	body, _, err := m.doRequest("/api/services/"+name, "DELETE", nil, token, "application/yaml")
	if err != nil {
		return err
	}

	var res map[string]interface{}
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return err
	}

	if monit == true {
		if str, ok := res["stream_id"].(string); ok {
			Monitorize(m.URL, token, str)
			runtime.Goexit()
		}
	}

	return nil
}

// ********************* Reset *******************

// ResetService ...
func (m *Manager) ResetService(name string, token string) error {
	_, _, err := m.doRequest("/api/services/"+name+"/reset/", "POST", nil, token, "application/yaml")
	return err
}

// ********************* Status *******************

// ServiceStatus ...
func (m *Manager) ServiceStatus(token string, serviceName string) (service Service, err error) {
	body, _, err := m.doRequest("/api/services/"+serviceName, "GET", []byte(""), token, "")
	if err != nil {
		return service, err
	}
	err = json.Unmarshal([]byte(body), &service)
	if err != nil {
		return service, err
	}
	return service, nil
}

// ServiceBuildStatus ...
func (m *Manager) ServiceBuildStatus(token string, serviceName string, serviceID string) (service Service, err error) {
	body, _, err := m.doRequest("/api/services/"+serviceName+"/builds/"+serviceID, "GET", []byte(""), token, "")
	if err != nil {
		return service, err
	}
	err = json.Unmarshal([]byte(body), &service)
	if err != nil {
		return service, err
	}
	return service, nil
}

// ********************* List *********************

// ListDatacenters ...
func (m *Manager) ListDatacenters(token string) (datacenters []Datacenter, err error) {
	body, _, err := m.doRequest("/api/datacenters/", "GET", []byte(""), token, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &datacenters)
	if err != nil {
		return nil, err
	}
	return datacenters, err
}

// ListServices ...
func (m *Manager) ListServices(token string) (services []Service, err error) {
	body, _, err := m.doRequest("/api/services/", "GET", []byte(""), token, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &services)
	if err != nil {
		return nil, err
	}
	return services, err
}

// ListBuilds ...
func (m *Manager) ListBuilds(name string, token string) (builds []Service, err error) {
	body, _, err := m.doRequest("/api/services/"+name+"/builds/", "GET", []byte(""), token, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &builds)
	if err != nil {
		return nil, err
	}
	return builds, err
}
