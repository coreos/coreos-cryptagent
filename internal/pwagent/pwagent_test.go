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
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

//TestNewAgent ensures NewAgent properly handles corner-cases.
func TestNewAgent(t *testing.T) {
	agent, err := newAgent("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if agent != nil {
		t.Fatalf("expected nil agent, got %#v", agent)
	}

	agent, err = newAgent("/invalid-non-existing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	tmpDir, err := ioutil.TempDir("", "test_agent")
	if err != nil {
		t.Fatalf("failed to create temporary test directory: %s", err)
	}
	defer os.RemoveAll(tmpDir)
	agent, err = newAgent(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if agent == nil {
		t.Fatal("unexpected nil agent")
	}
}

// TestAgentServerNil ensure nil receivers are properly handled.
func TestAgentServerNil(t *testing.T) {
	var agent *AgentServer
	errCh := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent.ServeRequests(ctx, nil, errCh)
	err := <-errCh
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestParseRequest provides basic testing for ParseRequest
func TestParseRequest(t *testing.T) {
	content := `
# discarded comment
Key1=value1
# next one is an invalid line
Key2
`
	rd := bufio.NewReader(strings.NewReader(content))
	kv, err := ParseRequest(rd)
	if err != nil {
		t.Fatal("unexpected error")
	}
	if value := kv["Key1"]; value != "value1" {
		t.Fatalf("got unexpected value: %s", value)
	}
	if len(kv) != 1 {
		t.Fatalf("got unexpected length: %d", len(kv))
	}
}
