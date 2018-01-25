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
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	attachCmd = &cobra.Command{
		Use:          "attach",
		RunE:         runAttachCmd,
		Short:        "Attach a crypsetup volume by device path",
		SilenceUsage: true,
	}
)

func runAttachCmd(cmd *cobra.Command, args []string) error {
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

	//TODO(lucab): volume attachment logic

	return nil
}
