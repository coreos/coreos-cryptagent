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
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/coreos/coreos-cryptagent/pkg/config"
	"github.com/coreos/go-systemd/unit"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const devBlockPath = "/dev/block/"

// LookupBlockdev translates a block device path into its `/dev/block` entry.
//
// `path` must be an existing absolute path. The resulting string is an absolute
// path rooted at `/dev/block/`.
func LookupBlockdev(pathIn string) (string, error) {
	logrus.Debugf("looking up block device for %s", pathIn)
	if pathIn == "" {
		return "", errors.New("empty path to lookup")
	}
	realPath, err := filepath.EvalSymlinks(pathIn)
	if err != nil {
		return "", err
	}

	fis, err := ioutil.ReadDir(devBlockPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to list %s", devBlockPath)
	}
	for _, fi := range fis {
		entry := filepath.Join(devBlockPath, fi.Name())
		resolved, err := filepath.EvalSymlinks(entry)
		if err == nil && realPath == resolved {
			return entry, nil
		}
	}

	return "", errors.Errorf("unable to lookup %s", pathIn)
}

// LookupVolName translates a block device path into its LUKS volume name.
//
// `path` must be an existing absolute path. The resulting string is a volume name.
func LookupVolName(pathIn string) (string, error) {
	// To ease testing, inject the base directory to a private function.
	return lookupVolName(config.DevConfigDir, pathIn)
}

func lookupVolName(devConfigDir string, pathIn string) (string, error) {
	logrus.Debugf("looking up volume name for device %s", pathIn)
	if pathIn == "" {
		return "", errors.New("empty path to lookup")
	}
	confDir, err := lookupConfigDir(devConfigDir, pathIn)
	if err != nil {
		return "", err
	}

	fp, err := os.Open(filepath.Join(confDir, "volume.json"))
	if err != nil {
		return "", err
	}
	defer fp.Close()
	var vj config.VolumeJSON
	if err := json.NewDecoder(bufio.NewReader(fp)).Decode(&vj); err != nil {
		return "", err
	}

	if luks1, ok := vj.Value.(config.CryptsetupLUKS1V1); ok {
		if luks1.Name == "" {
			return "", errors.New("empty volume name in configuration")
		}
		return luks1.Name, nil
	}

	return "", errors.New("unable to decode volume name from configuration")
}

// LookupConfigDir translates a block device path into its base config directory entry.
//
// `path` must be an existing absolute path to a device. `devConfigDir` is the default
// base config directory for coreos-cryptagent. The resulting string is the absolute
// path to the device configuration directory.
func LookupConfigDir(pathIn string) (string, error) {
	// To ease testing, inject the base directory to a private function.
	return lookupConfigDir(config.DevConfigDir, pathIn)
}

func lookupConfigDir(devConfigDir string, pathIn string) (string, error) {
	if pathIn == "" {
		return "", errors.New("empty device id")
	}

	dev := pathIn
	if !strings.HasPrefix(dev, devBlockPath) {
		res, err := LookupBlockdev(pathIn)
		if err != nil {
			return "", err
		}
		dev = res
	}

	fis, err := ioutil.ReadDir(devConfigDir)
	if err != nil {
		return "", errors.Wrapf(err, "failed to list %s", devConfigDir)
	}
	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		plain := unit.UnitNamePathUnescape(fi.Name())
		plainDev, err := LookupBlockdev(plain)
		if err != nil {
			return "", err
		}
		if plainDev == dev {
			path := filepath.Join(devConfigDir, fi.Name())
			logrus.Debugf("found config directory %q for %q", path, dev)
			return path, nil
		}
	}

	return "", errors.Errorf("no config directory found for %q", dev)
}
