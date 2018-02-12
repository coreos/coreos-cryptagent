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
	"testing"
	"time"
)

// TestContentNilProvider ensures nil receivers do not incur panics.
func TestContentNilProvider(t *testing.T) {
	var c *content
	ch := make(chan Result)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c.SetCiphertext("")

	canEncrypt := c.CanEncrypt()
	if canEncrypt {
		t.Errorf("unexpected canEncrypt: %t", canEncrypt)
	}

	json, err := c.ToProviderJSON()
	if err == nil {
		t.Errorf("unexpected ToProviderJSON: %#v", json)
	}

	go c.GetCleartext(ctx, nil, ch)
	clear := <-ch
	if clear.Err == nil {
		t.Errorf("unexpected GetCleartext: %#v", clear.Ok)
	}

	go c.Encrypt(ctx, nil, "", ch)
	cipher := <-ch
	if cipher.Err == nil {
		t.Errorf("unexpected Encrypt: %#v", cipher.Ok)
	}
}
