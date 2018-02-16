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

package cli

import (
	"os/exec"
	"path/filepath"

	"github.com/coreos/coreos-cryptagent/internal/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const sdHelperBin = "/lib/systemd/systemd-cryptsetup"

var (
	attachCmd = &cobra.Command{
		Use:          "attach",
		RunE:         runAttachCmd,
		Short:        "Attach a crypsetup volume by device path",
		SilenceUsage: true,
	}
)

func runAttachCmd(cmd *cobra.Command, args []string) error {
	var err error
	if len(args) == 0 {
		return errors.New("device path missing")
	}
	if len(args) != 1 {
		return errors.New("too many arguments")
	}
	pathIn := args[0]

	if !filepath.IsAbs(pathIn) {
		return errors.Errorf("input path %s is not absolute", pathIn)
	}

	blockPath, err := common.LookupBlockdev(pathIn)
	if err != nil {
		return errors.Wrap(err, "failed reverse block lookup")
	}

	volName, err := common.LookupVolName(pathIn)
	if err != nil {
		return errors.Wrap(err, "failed volume name lookup")
	}

	logrus.Debugf("unlocking volume %s on device %s\n", volName, blockPath)
	opts := []string{"-"}
	err = sdHelper(volName, blockPath, opts)
	if err != nil {
		return errors.Wrap(err, "failed to run systemd-crypsetup")
	}

	return nil
}

func sdHelper(volume string, path string, opts []string) error {
	if volume == "" {
		return errors.New("empty input volume name")
	}
	if path == "" {
		return errors.New("empty input path")
	}

	args := []string{"attach", volume, path}
	args = append(args, opts...)
	cmd := exec.Command(sdHelperBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := errors.New(string(out))
		return errors.Wrap(msg, err.Error())
	}

	return nil
}
