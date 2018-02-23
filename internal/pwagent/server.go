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

package pwagent

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/coreos-cryptagent/internal/common"
	"github.com/coreos/coreos-cryptagent/pkg/config"
	"github.com/coreos/coreos-cryptagent/pkg/providers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ProcessRequest handles a single password request.
func ProcessRequest(ctx context.Context, path string, errCh chan<- error) {
	fp, err := os.Open(path)
	if err != nil {
		errCh <- err
		return
	}
	defer fp.Close()

	props, err := ParseRequest(bufio.NewReader(fp))
	if err != nil {
		errCh <- err
		return
	}

	sock, ok := props["Socket"]
	if !ok {
		errCh <- errors.New("missing 'Socket' field")
		return
	}
	logrus.Debugln(sock)

	id := getCryptsetupID(props)
	if id == "" {
		errCh <- nil
		return
	}

	pass, err := GetPassphrase(ctx, id)
	if err != nil {
		errCh <- err
		return
	}

	// TODO(lucab): handle timeout conditions in requests.
	if err := SendPassphrase(ctx, sock, pass); err != nil {
		errCh <- err
		return
	}

	err = os.Remove(path)
	if err != nil {
		errCh <- err
		return
	}

	errCh <- nil
}

// GetPassphrase perform a remote call to retrieve/unlock a passphrase.
func GetPassphrase(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", errors.New("empty id")
	}
	dir, err := common.LookupConfigDir(id)
	if err != nil {
		return "", errors.Wrapf(err, "failed to find config directory for %q", id)
	}
	path := filepath.Join(dir, "0.json")
	fp, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open config file %q", path)
	}
	defer fp.Close()

	var cfg config.ProviderJSON
	if err := json.NewDecoder(bufio.NewReader(fp)).Decode(&cfg); err != nil {
		return "", errors.Wrapf(err, "failed to decode %q", path)
	}

	p, err := providers.FromProviderJSON(&cfg)
	if err != nil {
		return "", err
	}

	var res providers.Result
	tries := 5
	sleep := 4
	for tries > 0 {
		ch := make(chan providers.Result, 1)
		go p.GetCleartext(ctx, nil, ch)
		res = <-ch
		if res.Err == nil {
			break
		}
		tries--
		logrus.Warnf("unable to retrieve cleartext passphrase, retrying in %ds", sleep)
		time.Sleep(time.Duration(sleep) * time.Second)
	}
	return res.Ok, res.Err
}

func getCryptsetupID(props map[string]string) string {
	reqID, ok := props["Id"]
	if !ok {
		return ""
	}
	if !strings.HasPrefix(reqID, "cryptsetup:") {
		return ""
	}
	id := strings.TrimPrefix(reqID, "cryptsetup:")
	return id
}
