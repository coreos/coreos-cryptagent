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

package providers

import (
	"context"
	"errors"
	"net/http"

	"github.com/coreos/coreos-cryptagent/pkg/config"
	_ "github.com/coreos/ignition/config/types"
)

// Result represents the result of an async passphrase operation.
// It contains either an Ok string result or an operational Err error.
type Result struct {
	Ok  string
	Err error
}

// RemoteOptions contains optional parameters to tweak provider-specific
// options for remote operations.
type RemoteOptions struct {
	Headers http.Header
}

// PassProvider is an interface that can be used by library consumers
// without having direct access to each individual provider.
type PassProvider interface {
	// GetCleartext returns the cleartext passphrase for a volume.
	// It may entail multiple calls to a remote provider in order
	// to unwrap/decrypt a local ciphertext.
	GetCleartext(ctx context.Context, opts *RemoteOptions, doneCh chan<- Result)
	// SetCiphertext sets the ciphertext, so that it can be later serialized.
	SetCiphertext(string)
	// Encrypt encrypts an external cleartext.
	// It may entail multiple calls to a remote provider in order
	// to wrap/encrypt the cleartext.
	Encrypt(ctx context.Context, opts *RemoteOptions, cleartext string, doneCh chan<- Result)
	// ToProviderJSON serializes current PassProvider configuration
	// as a ProviderJSON struct.
	ToProviderJSON() (*config.ProviderJSON, error)
	// CanEncrypt signals whether the provider can encrypt external cleartext.
	CanEncrypt() bool
}

/*
// FromIgnitionV230 constructs an opaque PassProvider from an ignition-2.3.0
// keyslot configuration entry.
func FromIgnitionV230(ks types.LuksKeyslot, ign types.Ignition) (PassProvider, error) {
	switch {
	case ks.AzureVault != nil:
		return nil, errors.New("provider not implemented")
	case ks.Content != nil:
		return contentFromIgnitionV230(ks, ign)
	case ks.HcVault != nil:
		return nil, errors.New("provider not implemented")
	}

	return nil, errors.New("invalid keyslot")
}
*/

// FromProviderJSON constructs an opaque PassProvider from any ProviderJSON
// configuration.
func FromProviderJSON(cfg *config.ProviderJSON) (PassProvider, error) {
	if cfg == nil {
		return nil, errors.New("nil JSON provider configuration")
	}

	switch {
	case cfg.Kind == config.ProviderAzureVaultV1:
		return nil, errors.New("provider not implemented")
	case cfg.Kind == config.ProviderContentV1:
		return contentFromConfigV1(cfg)
	case cfg.Kind == config.ProviderHcVaultV1:
		return nil, errors.New("provider not implemented")
	}

	return nil, errors.New("invalid provider")
}
