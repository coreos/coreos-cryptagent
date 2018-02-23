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
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/coreos/coreos-cryptagent/pkg/config"
	_ "github.com/coreos/ignition/config/types"
)

type content struct {
	source                 string
	responseHeadersTimeout int
	totalTimeout           int
	certAuths              []string
}

func contentFromConfigV1(cfg *config.ProviderJSON) (*content, error) {
	if cfg == nil {
		return nil, errors.New("nil config")
	}
	if cfg.Kind != config.ProviderContentV1 {
		return nil, fmt.Errorf("expected kind %q, got %q", config.ProviderContentV1, cfg.Kind)
	}

	value, ok := cfg.Value.(config.ContentV1)
	if !ok {
		return nil, errors.New("not a ContentV1 value")
	}

	c := content{
		source: value.Source,
	}
	if value.Timeouts != nil {
		c.totalTimeout = value.Timeouts.HTTPTotal
		c.responseHeadersTimeout = value.Timeouts.HTTPResponseHeaders
	}
	for _, ca := range value.CertificateAuthorities {
		if ca.Authority != "" {
			c.certAuths = append(c.certAuths, ca.Authority)
		}
	}

	return &c, nil
}

/*
func contentFromIgnitionV230(ks types.LuksKeyslot, ign types.Ignition) (*content, error) {
	if ks.Content == nil {
		return nil, errors.New("nil Content keyslot")
	}
	if ks.Content.Source == "" {
		return nil, errors.New("empty source in Content keyslot")
	}

	c := content{
		source: ks.Content.Source,
	}
	if ign.Timeouts.HTTPTotal != nil {
		c.totalTimeout = *ign.Timeouts.HTTPTotal
	}
	if ign.Timeouts.HTTPResponseHeaders != nil {
		c.responseHeadersTimeout = *ign.Timeouts.HTTPResponseHeaders
	}

	// TODO(lucab): integrate CAs

	return &c, nil
}
*/

// GetCleartext implements the PassProvider interface.
func (c *content) GetCleartext(ctx context.Context, opts *RemoteOptions, doneCh chan<- Result) {
	if c == nil {
		doneCh <- Result{"", errors.New("nil content receiver")}
		return
	}
	if c.source == "" {
		doneCh <- Result{"", errors.New("missing source URL")}
		return
	}

	client, err := c.newClient()
	if err != nil {
		doneCh <- Result{"", err}
		return
	}
	headers := http.Header{}
	if opts != nil {
		headers = opts.Headers
	}
	rd, statusCode, err := c.getReaderWithHeader(ctx, client, c.source, headers)
	if err != nil {
		doneCh <- Result{"", err}
		return
	}
	defer rd.Close()
	if statusCode != 200 {
		doneCh <- Result{"", fmt.Errorf("%d", statusCode)}
		return
	}
	body, err := ioutil.ReadAll(rd)
	if err != nil {
		doneCh <- Result{"", err}
		return
	}

	doneCh <- Result{string(body), nil}
}

// Encrypt implements the PassProvider interface.
func (c *content) Encrypt(ctx context.Context, opts *RemoteOptions, cleartext string, doneCh chan<- Result) {
	// This provider cannot be used to encrypt an external cleartext.
	doneCh <- Result{"", errors.New("content provider does not support encryption")}
}

// ToProviderJSON implements the PassProvider interface.
func (c *content) ToProviderJSON() (*config.ProviderJSON, error) {
	if c == nil {
		return nil, errors.New("nil content receiver")
	}
	if c.source == "" {
		return nil, errors.New("empty source")
	}

	to := config.ContentV1Timeouts{
		HTTPResponseHeaders: c.responseHeadersTimeout,
		HTTPTotal:           c.totalTimeout,
	}
	cas := []config.ContentV1CertAuth{}
	for _, path := range c.certAuths {
		ca := config.ContentV1CertAuth{
			Authority: path,
		}
		cas = append(cas, ca)
	}
	v := config.ContentV1{
		Source:                 c.source,
		Timeouts:               &to,
		CertificateAuthorities: cas,
	}
	pj := config.ProviderJSON{
		Kind:  config.ProviderContentV1,
		Value: v,
	}

	return &pj, nil
}

// CanEncrypt implements the PassProvider interface.
func (c *content) CanEncrypt() bool {
	// This provider cannot be used to encrypt an external cleartext.
	return false
}

// SetCiphertext implements the PassProvider interface.
func (c *content) SetCiphertext(ciphertext string) {
	// This provider does not need to locally store any ciphertext.
}

// XXX(lucab): this code below comes from ignition `internal/resource`
// package. It should be unified, but the original is tangled with other
// ignition internal packages.

func (c *content) newClient() (*http.Client, error) {
	if c == nil {
		return nil, errors.New("nil content receiver")
	}

	transport := &http.Transport{
		ResponseHeaderTimeout: time.Duration(c.responseHeadersTimeout) * time.Second,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			Resolver: &net.Resolver{
				PreferGo: true,
			},
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := http.Client{
		Transport: transport,
	}

	return &client, nil
}

func (c *content) getReaderWithHeader(ctx context.Context, client *http.Client, url string, header http.Header) (io.ReadCloser, int, error) {
	initialBackoff := 100 * time.Millisecond
	maxBackoff := 5 * time.Second
	if c == nil {
		return nil, 0, errors.New("nil content receiver")
	}
	if client == nil {
		return nil, 0, errors.New("nil client")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	for key, values := range header {
		req.Header.Del(key)
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if c.totalTimeout != 0 {
		ctxTo, cancel := context.WithTimeout(ctx, time.Duration(c.totalTimeout)*time.Second)
		ctx = ctxTo
		defer cancel()
	}
	req = req.WithContext(ctx)
	duration := initialBackoff
	for attempt := 1; ; attempt++ {
		resp, err := client.Do(req)
		if err == nil {
			if resp.StatusCode < 500 {
				return resp.Body, resp.StatusCode, nil
			}
			resp.Body.Close()
		}

		duration = duration * 2
		if duration > maxBackoff {
			duration = maxBackoff
		}

		// Wait before next attempt or exit if we timeout while waiting
		select {
		case <-time.After(duration):
		case <-ctx.Done():
			return nil, 0, errors.New("unable to fetch resource in time")
		}
	}
}
