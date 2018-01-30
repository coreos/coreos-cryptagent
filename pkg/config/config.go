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

// Package config contains all configuration types for data interchange with
// coreos-cryptagent.
package config

import (
	"encoding/json"
	"errors"
	"path/filepath"
)

const (
	// BaseConfigDir is the path of the base directory storing coreos-cryptagent config.
	BaseConfigDir = "/boot/etc/coreos-cryptagent/"
)

var (
	// DevConfigDir is the base directory for devices/volumes configuration files.
	DevConfigDir = filepath.Join(BaseConfigDir, "dev")
)

// VolumeKind is an enum of volume kinds.
type VolumeKind uint

const (
	// VolumeInvalid is the default invalid value for volume kind.
	VolumeInvalid VolumeKind = iota
	// VolumeCryptsetupLUKS1V1 represents a cryptsetup-LUKS1 (v1) volume config.
	VolumeCryptsetupLUKS1V1
)

// UnmarshalJSON is part of the json.Unmarshaler interface.
func (vk *VolumeKind) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "CryptsetupLUKS1V1":
		*vk = VolumeCryptsetupLUKS1V1
	default:
		return errors.New("unknown kind")
	}

	return nil
}

// MarshalJSON is part of the json.Marshaler interface.
func (vk VolumeKind) MarshalJSON() ([]byte, error) {
	var s string
	switch vk {
	case VolumeCryptsetupLUKS1V1:
		s = "CryptsetupLUKS1V1"
	default:
		return nil, errors.New("unknown kind")
	}

	return json.Marshal(s)
}

// ProviderKind is an enum of provider kinds.
type ProviderKind int

const (
	// ProviderInvalid is the nil value for ProviderKind
	ProviderInvalid ProviderKind = iota
	// ProviderContentV1 represents a plain Content (v1) config
	ProviderContentV1
	// ProviderAzureVaultV1 represents an Azure Vault (v1) config
	ProviderAzureVaultV1
	// ProviderHcVaultV1 represents an HashiCorp Vault (v1) config
	ProviderHcVaultV1
)

// UnmarshalJSON is part of the json.Unmarshaler interface.
func (vk *ProviderKind) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "ContentV1":
		*vk = ProviderContentV1
	case "AzureVaultV1":
		return errors.New("azure-vault unimplemented")
	case "HcVaultV1":
		return errors.New("hc-vault unimplemented")
	default:
		return errors.New("unknown kind")
	}

	return nil
}

// MarshalJSON is part of the json.Marshaler interface.
func (vk ProviderKind) MarshalJSON() ([]byte, error) {
	var s string
	switch vk {
	case ProviderContentV1:
		s = "ContentV1"
	case ProviderAzureVaultV1:
		return nil, errors.New("azure-vault unimplemented")
	case ProviderHcVaultV1:
		return nil, errors.New("hc-vault unimplemented")
	default:
		return nil, errors.New("unknown kind")
	}

	return json.Marshal(s)
}
