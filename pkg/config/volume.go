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

import (
	"encoding/json"
	"errors"
)

// VolumeJSON is the top-level configuration container for an encrypted volume.
type VolumeJSON struct {
	Kind  VolumeKind  `json:"kind"`
	Value interface{} `json:"value"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (vj *VolumeJSON) UnmarshalJSON(b []byte) error {
	type tmps struct {
		Kind  VolumeKind       `json:"kind"`
		Value *json.RawMessage `json:"value"`
	}
	var tmp tmps
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	switch tmp.Kind {
	case VolumeCryptsetupLUKS1V1:
		var v CryptsetupLUKS1V1
		if err := json.Unmarshal(*tmp.Value, &v); err != nil {
			return err
		}
		vj.Kind = tmp.Kind
		vj.Value = v

	default:
		return errors.New("unknown kind")
	}

	return nil
}

// CryptsetupLUKS1V1 represents a cryptsetup-LUKS1 volume.
type CryptsetupLUKS1V1 struct {
	Name           string `json:"name"`
	Device         string `json:"device"`
	DisableDiscard *bool  `json:"disableDiscard,omitempty"`
}
