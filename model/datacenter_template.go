/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package model

// DatacenterTemplate ...
type DatacenterTemplate struct {
	URL      string `yaml:"vcloud-url"`
	Network  string `yaml:"public-network"`
	Org      string `yaml:"org"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Token    string `yaml:"secret_access_key"`
	Secret   string `yaml:"access_key_id "`
	Region   string `yaml:"region"`
	Fake     bool   `yaml:"fake"`
	// Azure
	SubscriptionID string `yaml:"subscription_id"`
	ClientID       string `yaml:"client_id"`
	ClientSecret   string `yaml:"client_secret"`
	TenantID       string `yaml:"tenant_id"`
	Environment    string `yaml:"environment"`
}
