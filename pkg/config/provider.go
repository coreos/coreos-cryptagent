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

// ProviderJSON is the top-level configuration container for a provider.
type ProviderJSON struct {
	Kind  ProviderKind `json:"kind"`
	Value interface{}  `json:"value"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (pj *ProviderJSON) UnmarshalJSON(b []byte) error {
	type tmps struct {
		Kind  ProviderKind     `json:"kind"`
		Value *json.RawMessage `json:"value"`
	}
	var tmp tmps
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	switch tmp.Kind {
	case ProviderContentV1:
		var v ContentV1
		if err := json.Unmarshal(*tmp.Value, &v); err != nil {
			return err
		}
		pj.Kind = tmp.Kind
		pj.Value = v
	case ProviderAzureVaultV1:
		return errors.New("azure-vault unimplemented")
	case ProviderHcVaultV1:
		return errors.New("hc-vault unimplemented")
	default:
		return errors.New("unknown kind")
	}

	return nil
}
