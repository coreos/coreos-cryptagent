// Copyright 2018 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

// ContentV1 is the v1 configuration for a generic remote content provider.
type ContentV1 struct {
	Source   string             `json:"source"`
	Timeouts *ContentV1Timeouts `json:"timeouts,omitempty"`
	//TODO(lucab): specify this better
	CertificateAuthorities []ContentV1CertAuth `json:"certificateAuthorities,omitempty"`
}

// ContentV1Timeouts records HTTPS client timeouts
type ContentV1Timeouts struct {
	HTTPResponseHeaders int `json:"httpResponseHeaders"`
	HTTPTotal           int `json:"httpTotal"`
}

// ContentV1CertAuth records HTTPS client custom CAs
type ContentV1CertAuth struct {
	//TODO(lucab): store PEM here? Or just a path?
	Authority string `json:"authority"`
}
