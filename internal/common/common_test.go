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

package common

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"encoding/json"
	"github.com/coreos/coreos-cryptagent/pkg/config"
	"github.com/coreos/go-systemd/unit"
)

const (
	loop0Dev      = "/dev/loop0"
	loop0BlockDev = "/dev/block/7:0"
)

func loopModprobe() error {
	cmd := exec.Command("modprobe", "loop")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("'modprobe loop' error: %s - %s", err, out)
	}
	return nil
}

func TestLookupBlockdevLoop(t *testing.T) {
	if err := loopModprobe(); err != nil {
		t.Skipf("test setup failed: %s", err)
	}

	out, err := LookupBlockdev(loop0Dev)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}
	if out != loop0BlockDev {
		t.Fatalf("expected lookup result %q, got %s", loop0BlockDev, out)
	}
}

func TestLookupVolNameLoop(t *testing.T) {
	if err := loopModprobe(); err != nil {
		t.Skipf("test setup failed: %s", err)
	}
	tmpDir, err := ioutil.TempDir("", "common_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	escaped := unit.UnitNamePathEscape(loop0Dev)
	devDir := filepath.Join(tmpDir, escaped)
	if err := os.MkdirAll(devDir, 0755); err != nil {
		t.Fatal(err)
	}
	fp, err := os.Create(filepath.Join(devDir, "volume.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()

	volName := "luks_vol"
	luks := config.CryptsetupLUKS1V1{
		Name:   volName,
		Device: loop0Dev,
	}
	vj := config.VolumeJSON{
		Kind:  config.VolumeCryptsetupLUKS1V1,
		Value: &luks,
	}
	err = json.NewEncoder(fp).Encode(vj)
	if err != nil {
		t.Fatal(err)
	}
	fp.Sync()

	out, err := lookupVolName(tmpDir, loop0BlockDev)
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}
	if out != volName {
		t.Fatalf("expected lookup result %q, got %s", loop0BlockDev, out)
	}
}

func TestLookupBlockdev(t *testing.T) {
	tests := []struct {
		pathIn string
		expErr error
	}{
		{
			"",
			errors.New("empty path to lookup"),
		},
		{
			"/non-existing",
			errors.New("lstat /non-existing: no such file or directory"),
		},
	}

	for _, tt := range tests {
		_, err := LookupBlockdev(tt.pathIn)
		if err != nil && err.Error() != tt.expErr.Error() {
			t.Fatalf("expected error %q, got %q", tt.expErr, err)
		}
	}

}

func TestLookupVolName(t *testing.T) {
	tests := []struct {
		pathIn string
		expErr error
	}{
		{
			"",
			errors.New("empty path to lookup"),
		},
		{
			"/non-existing",
			errors.New("lstat /non-existing: no such file or directory"),
		},
	}

	for _, tt := range tests {
		_, err := LookupVolName(tt.pathIn)
		if err != nil && err.Error() != tt.expErr.Error() {
			t.Fatalf("expected error %q, got %q", tt.expErr, err)
		}
	}

}
