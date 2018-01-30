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

// TODO(lucab): double-check, stabilize and make this public.

// AzureVaultV1 is the v1 configuration for an Azure Vault provider.
type azureVaultV1 struct {
	BaseURL             string                    `json:"baseURL"`
	EncryptionAlgorithm string                    `json:"encryptionAlgorithm"`
	KeyName             string                    `json:"keyName"`
	KeyVersion          string                    `json:"keyVersion"`
	Ciphertext          string                    `json:"ciphertext"`
	PasswordAuth        *azureVaultV1PasswordAuth `json:"passwordAuth"`
}

// AzureVaultV1PasswordAuth is the password authentication stanza for AzureVaultV1.
type azureVaultV1PasswordAuth struct {
	AppID    string `json:"appID"`
	Password string `json:"password"`
}
