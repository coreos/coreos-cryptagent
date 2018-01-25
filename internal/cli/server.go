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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:          "server",
		RunE:         runServerCmd,
		Short:        "Runs the password agent server",
		SilenceUsage: true,
	}
)

func runServerCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return errors.New("too many arguments")
	}
	logrus.Infoln("starting coreos-cryptagent server")

	//TODO(lucab): agent server logic
	for {
	}
}
