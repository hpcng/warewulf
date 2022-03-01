# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed 
- Provision interface is not tied to 'eth0' any more. The provision interface must be named
  'default' now. The file `nodes.yaml' must be changed accordingly.
- Creating of '/etc/exports' can now be disabled, so that `wwctl configure -a` wont overwrite
  a existing '/etc/exports'.
- All configurations files for the host (/etc/exports, /etc/dhcpd.conf, /etc/hosts) are now
  popublated from the (OVERLAYDIR/host/etc/{exports|dhcpd|hosts}.ww . Also other configuration
  files like prometheus.yml.ww or slurm.conf.ww which depend on the cluster nodes can be
  placed. Also the new templated functions {{ abort }}, {{ IncludeBlock }} abd {{ no_backup }}
  are allowed now.


## [4.1.0] - 2021-07-29
### Added
- Support for ARM nodes
- firewalld service file for warewulfd
- `-y` option to skip "Are you sure" queries
- `wwctl kernel delete` command
- `wwctl vnfs` alias for `wwctl container`
- Support for authenticated OCI registries
- warewulfd can reload config on SIGHUP and when the config file changes
- Node database index to improve lookup speeds
- Kernels and containers can be imported from a chroot subdirectory
- Systemd service file

### Changed
- `wwctl node list` output beautification
- Log timestamps are more precise
- PID file and log files are now in `/var/run` and `/var/log`, respectively
- `make install` no longer overwrites preexisting configuration files
- Kernel modules and overlays are now compressed
- `rootfstype` now uses `rootfs` in default kernel arguments
- iPXE binaries updated
- Installed container directory is deleted when import fails
- Default iPXE script now reboots erroring nodes every 15 seconds
- Only open `/etc/hosts` when writing

### Removed
- `wwctl configure` `--persist` flags have been removed. `configure` commands persist changes by default unless `--show` is used
- In-repository documentation: has been moved to it's own repository

### Fixed
- Importing containers from directory
- Debug log verbosity option takes precedence over verbose option
- `wwctl node list -n` output is formatted corectly
- Container names can contain an underscore
- `wwctl overlay build --all` does not require an argument
- specfile date format works with older versions of rpmbuild
- Use SystemOverlay when building system overlay
- dhcpd template now references correct wwctl subcommand
- `wwctl node set kernelargs` and `wwctl profile set kernelargs` change kernel arguments
