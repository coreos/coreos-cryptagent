# Configuration

This projects defines a stable configuration interchange format in order for Ignition (or other provisioning systems) to persist LUKS volume parameters and for Cryptagent to retrieve them on reboots.

Configuration files are JSON-based and strictly type-versioned, in order to ensure compatibility across versions and to ease manual inspection.

# Paths

Configuration for Cryptagent is rooted at `/boot/etc/coreos-cryptagent/`. Such directory would not exist on systems without encrypted volumes.
After a cryptsetup volume is created, its configuration needs to be stored under `/boot/etc/cryptsetup-agent/dev/`.

In particular, each encrypted device will get a configuration directory at `/boot/etc/cryptsetup-agent/dev/$DEVNAME/` containing a `volume.json` and `$N.json`, where:
 * `$DEVNAME` is the `systemd-escape --path` encoded path of the user-specified encrypted device.
 * `$N` is currently hardcoded to 0.
 * `volume.json` contains volume configuration parameters.
 * `$N.json` contains parameters for keyslot number `$N` (currently only keyslot 0 is available).
 * configuration files are valid JSON documents, whose format is specified below.

# Schemas

TODO(lucab): add JSON schema for all public `pkg/config` structs.
