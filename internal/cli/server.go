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
	"context"

	"github.com/coreos/coreos-cryptagent/internal/pwagent"
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
	var err error
	ctx := context.Background()
	logrus.Infoln("starting coreos-cryptagent server")
	if len(args) != 0 {
		return errors.New("too many arguments")
	}

	agent, err := pwagent.NewAgent()
	if err != nil {
		return errors.Wrap(err, "failed to initialize password-agent server")

	}
	errCh := make(chan error, 2048)
	go agent.ServeRequests(ctx, pwagent.ProcessRequest, errCh)
	for {
		select {
		case <-ctx.Done():
			if ctxErr := ctx.Err(); ctxErr != nil {
				err = errors.Wrap(ctxErr, "request serving failed")
			}
			return err
		case cbErr := <-errCh:
			if cbErr != nil {
				logrus.Errorln(cbErr)
			}
		}
	}
}
