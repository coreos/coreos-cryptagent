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

// Package pwagent implements systemd PasswordAgent protocol as described in
// https://www.freedesktop.org/wiki/Software/systemd/PasswordAgents/.
package pwagent

// TODO(lucab): move to go-systemd once we are confident about the public API.

import (
	"bufio"
	"context"
	"errors"
	"net"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// Basedir is the base directory for the pwagent protocol file exchanges.
const basedir = "/run/systemd/ask-password/"

// ResponseFunc is a function type responsible for responding to password requests.
type ResponseFunc func(ctx context.Context, path string, errCh chan<- error)

// AgentServer is a systemd PasswordAgent server instance.
type AgentServer struct {
	w   *fsnotify.Watcher
	dir string
}

// NewAgent initialize a PasswordAgent.
func NewAgent() (*AgentServer, error) {
	// To ease testing, inject the base directory to a private function.
	return newAgent(basedir)
}

func newAgent(dir string) (*AgentServer, error) {
	if dir == "" {
		return nil, errors.New("empty path")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := w.Add(dir); err != nil {
		w.Close()
		return nil, err
	}

	as := AgentServer{
		w:   w,
		dir: dir,
	}
	return &as, nil
}

// ServeRequests processes password request files.
func (as *AgentServer) ServeRequests(ctx context.Context, fn ResponseFunc, errCh chan<- error) {
	if as == nil {
		errCh <- errors.New("got nil AgentServer")
		return
	}
	if fn == nil {
		errCh <- errors.New("nil ResponseFunc function")
		return
	}
	if as.w == nil {
		errCh <- errors.New("got nil inotify watcher")
		return
	}
	defer as.w.Close()
	defer as.w.Remove(as.dir)

	// Process pre-existing stale requests.
	files, err := filepath.Glob(filepath.Join(as.dir, "ask.*"))
	if err != nil {
		errCh <- err
		return
	}
	for _, path := range files {
		go fn(ctx, path, errCh)
	}

	// Process new incoming requests.
	for {
		select {
		case event := <-as.w.Events:
			if !strings.HasPrefix(event.Name, filepath.Join(as.dir, "ask.")) {
				continue
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				go fn(ctx, event.Name, errCh)
			}
		case err := <-as.w.Errors:
			errCh <- err
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		}
	}
}

// ParseRequest parses a password-request file and returns all indexed fields
// in it.
func ParseRequest(rd *bufio.Reader) (map[string]string, error) {
	ret := map[string]string{}

	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}
		ret[kv[0]] = kv[1]
	}
	if err := scan.Err(); err != nil {
		return ret, err
	}

	return ret, nil
}

// SendPassphrase replies to a password request.
func SendPassphrase(ctx context.Context, sock string, password string) error {
	if sock == "" {
		return errors.New("missing socket address")
	}
	if password == "" {
		return errors.New("empty passphrase")
	}
	socketAddr := &net.UnixAddr{
		Name: sock,
		Net:  "unixgram",
	}
	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	send := func(ch chan error, pass string) {
		_, err := conn.Write([]byte(pass))
		ch <- err
	}
	wirePass := "+" + password
	errCh := make(chan error, 1)
	go send(errCh, wirePass)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
