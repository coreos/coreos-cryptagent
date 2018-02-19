# CoreOS Cryptagent

[![Apache](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Build Status (Travis)](https://travis-ci.org/coreos/coreos-cryptagent.svg?branch=master)](https://travis-ci.org/coreos/coreos-cryptagent)

`coreos-cryptagent` is the utility used by Container Linux to unlock encrypted disks at boot time, directly in initramfs.
It is meant to complement [Ignition][ignition] first-boot setup, taking care of activating volumes on subsequent boots.

## Usage

On typical a run, there is no direct user interaction. Unlocking is triggered via `udev` events, and volumes are automatically processed based on relevant [configuration entries](Documentation/devel/config.md).

To report bugs, please use the [common CoreOS bug tracker][issues].

## License

`coreos-cryptagent` is released under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.

[ignition]: https://github.com/coreos/ignition
[issues]: https://github.com/coreos/bugs/issues/new?labels=component/coreos-cryptagent
