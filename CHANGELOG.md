# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## v4.6.0, unreleased

### Added

- Add the ability to boot nodes with `wwid=[interface]`, which replaces
  `interface` with the interface MAC address
- Added https://github.com/Masterminds/sprig functions to templates #1030
- Add multiple output formats (yaml & json) support. #447
- More aliases for many wwctl commands
- Add support to render template using `host` or `$(uname -n)` as the value of `overlay show --render`. #623
- Added an alternative image URL /containers-by-name/{{.ContainerName}} that does not contain a MAC addresses

### Changed

- Locally defined `tr` has been dropped, templates updated to use Sprig replace.
- Bump github.com/opencontainers/image-spec to 1.1.0
- Bump github.com/containers/storage to 1.53.0
- Bump google.golang.org/grpc 1.62.1
- Bump google.golang.org/protobuf to 1.33.0
- Bump github.com/containers/image/v5 to 5.30.0
- Bump github.com/docker/docker to 25.0.5+incompatible
- Bump github.com/grpc-ecosystem/grpc-gateway/v2 from 2.18.0 to 2.19.1 #1165
- Bump github.com/spf13/cobra from 1.7.0 to 1.8.0 #1166
- Bump github.com/fatih/color from 1.15.0 to 1.17.0 #1224
- Bump github.com/coreos/ignition/v2 from 2.15.0 to 2.19.0 #1239

### Fixed

- Prevent Networkmanager from trying to optain IP address via DHCP
  on unused/unmanaged network interfaces.
- Systems with no SMBIOS (Raspberry Pi) will create a UUID from
  `/sys/firmware/devicetree/base/serial-number`
- Replace slice in templates with sprig substr. #1093
- Fix an invalid format issue for the GitHub nightly build action. #1258

## v4.5.6, unreleased

### Fixed

- Ensure autobuilt overlays include contextual overlay contents. #1296
- Fix the failure when updating overlay files existing on different partitions. #1312

## v4.5.5, 2024-07-05

### Fixed

- Support leading and trailing slashes in `/etc/warewulf/excludes`. #1266
- Fix a regression in overlay autobuild. #1216
- Fix wwclient not reading asset-tag. #1110
- Fix dhcp not passing asset tag or uuid to iPXE. #1110
- Restored previous static dhcp behavior. #1263
- Capture "broken" symlinks during container build. #1267
- Fix the issue that removing lines during wwctl overlay edit didn't work. #1235
- Fix the issue that new files created with wwctl overlay edit have 755 permissions. #1236
- Fix tab-completion for `wwctl overlay list`. #1260

### Changed

- Explicitly ignore compat-style NIS lines in passwd/group during syncuser. #1286
- Accept `+` within kernel version. #1268
- Mount `/sys` and `/run` during `wwctl container exec`. #1287

## v4.5.4, 2024-06-12

### Fixed

- Fix a regression that caused an error when passing flags to `wwctl container exec` and `wwctl container shell`. #1250

## v4.5.3, 2024-06-07

### Added

- Add examples for building overlays in parallel to documentation
- Add `stage=initramfs` to warewulfd provision to serve initramfs from container image. #1115
- Add `warewulf-dracut` package to support building Warewulf-compatible initramfs images with dracut. #1115
- Add iPXE template `dracut.ipxe` to boot a dracut initramfs. #1115
- Add dracut menuentry to `grub.cfg.ww` to boot a dracut initramfs. #1115
- Add `.NetDevs` variable to iPXE and GRUB templates, similar to overlay templates. #1115
- Add `.Tags` variable to iPXE and GRUB templates, similar to overlay templates. #1115

### Changed

- Replace reference to docusaurus with Sphinx
- `wwctl container import` now only runs syncuser if explicitly requested. #1212
- wwinit now configures NetworkManager to not retain configurations from dracut. #1115
- Improved detection of SELinux capable root fs #1093

### Fixed

- Block unprivileged requests for arbitrary overlays in secure mode. #1215
- Fix installation docs to use github.com/warewulf instead of github.com/hpcng. #1219
- Fix the issue that warewulf.conf parse does not support CIDR format. #1130
- Reduce the number of times syncuser walks the container file system. #1209
- Create ssh key also when calling `wwctl configure --all` #1250
- Remove the temporary overlayfs dir and create them besides rootfs #1180

### Security

- Bump golang.org/x/net from 0.22.0 to 0.23.0. #1223

## v4.5.2, 2024-05-13

### Added

- Allow specification of the ssh-keys to be to be created. #1185

### Changed

- The command `wwctl container exec` locks now this container during execution. #830

### Fixed

- Fix nightly release build failure issue. #1195
- Reorder dnsmasq config to put iPXE last. #1146
- Update a reference to `--addprofile` to be `--profile`. #1085
- Update a dependency to address CVE-2024-3727. #1221

## v4.5.1, 2024-04-30

### Added

- Document warewulf.conf:paths. #635
- New "Overlay" template variable contains the name of the overlay being built. #1052 
- Documented HTTP proxy environment variables for `wwctl container import`. #1214

### Changed

- Update the glossary. #819
- Upgrade the golang version to 1.20.
- Bump github.com/opencontainers/umoci to 0.4.7
- Bump github.com/containers/image/v5 to 5.30.0
- Bump github.com/docker/docker to 25.0.5+incompatible
- Bump github.com/go-jose/go-jose/v3 to 3.0.3
- Bump gopkg.in/go-jose/go-jose.v2 to 2.6.3
- Bump github.com/opencontainers/runc to 1.1.12
- Dynamically calculate version and release from Git. #1162
- Update quickstarts to configure firewalld for dhcp. #1133
- Omit building the API on EL7. #1171
- Syncuser only walks the file system if it is going to write. #1207

### Fixed

- Fix `wwctl profile list -a` format when kernerargs are set.
- Don't attempt to rebuild protocol buffers in offline mode. #1155
- Fix Suse package by moving yq command to `%install` section. #1169
- Fix a rendering bug in the documentation for GRUB boot support. #1132
- Fix a locking issue with concurrent read/writes for node status. #1174
- Fix shim and grub detection for aarch64. #1145
- wwctl [profile|node] list -a handles now slices correclty. #1113
- Fix parsing of /etc/group during syncuser. #1202

## 4.5.0, 2024-02-08

Official v4.5.0 release.

### Added

- Publish v4.5.x documentation separately from `main`. #919
- Update quickstart for Enterprise Linux. #394, #401, #977

### Fixed

- Fix `Requires: ipxe-bootimgs` for building an Enterprise Linux 7 RPM. #1126

## 4.5.0rc2, 2024-02-21

### Fixed

- Fix mounting local partitions into sub-directories with Ignition. #1073
- Fix a panic in `wwctl node set` when modifying a network device that is only defined in a profile. #1094

## 4.5.0rc1, 2024-02-08

### Added

- Start building packages for Rocky Linux 9. #951
- Start building nightly GitHub releases. #969
- Support building on Fedora. #988
- New command `wwctl container copy` to duplicate a container image. #130
- New command `wwctl container rename` to rename a container image. #583
- New command `wwctl genconf` generates shell completions, initial configuration files, and documentation. #721
- New flag `wwctl container syncuser --build` automatically (re)builds a container image after syncuser. #509
- New flag `wwctl overlay import --parents` automatically creates intermediate parent directories. #481, #608
- New flag `wwctl <node|profile> list --fullall` shows all attributes, including those which do not have a set value. #786
- New option `wwctl <node|profile> set --ipmiescapechar` changes the `ipmitool` escape character. #999
- New option `wwclient --warewulfconf` specifies the location of `warewulf.conf`.
- New configuration option `warewulf.conf:paths` to override compiled-in paths. #721, #960, #1037
- Support configuring full paths to iPXE binaries in `warewulf.conf:tftp:ipxe`. #784
- Use `DNS` network tags to configure DNS resolution in network configuration overlays. #922
- Support a command-separated list of nodes or profiles in `wwctl <node|profile> list`. #739
- Support specifying `warewulf.conf:ipaddr` in CIDR format, optionally inferring netmask and network. #1016
- Support bonded network interfaces with NetworkManager using new network device types `bond` and `bond-slave`. #1071
- Support tab completion for overlay files in `wwctl overlay`. #813
- New `localdisk.ipxe` script to boot from local disk. #1058
- New example template to create genders database for use with clush, pdsh, and others. #825
- Document the "hostlist" syntax used by `wwctl`. #611
- Document using Vagrant as a development environment. #855
- Document quickstart for Enterprise Linux 9. #849
- Document configuring nodes with multiple networks. #1043
- (preview) Added support for initializing file systems, partitions and disks with Ignition. #786
- (preview) Added support for dnsmasq as a tftp and dhcp provider. #727, #1041
- (preview) Support booting nodes with GNU GRUB as an alternative to iPXE. #859

### Fixed

- Fix configuration of network device MTU in network configuration overlays. #807
- Separate provisioned overlays from overlay sources. #972
- Display the correct header for `wwctl overlay list --all --long`. #485
- Prevent column overflow in `wwctl <subcommand> list` with dynamic tabular output. #690
- Support relative path to a container image archive in `wwctl container import`. #493
- Correctly configure `ONBOOT` in `wwinit:etc/sysconfig/network-scripts/ifcfg.ww`. #644
- Fix multiple bugs in `wwctl node edit`. #691, #902, #1024
- Fix formatting of kernel arguments in `wwctl <node|profile> list`. #828
- Fix the ability to build the Warewulf API with `make`. #854
- Fix the ability to set MTU with `wwctl`. #947
- Fix multiple bugs in the handling of node and profile tags. #884, #967
- Fix an error when using `wwctl container import --force` to replace an existing container. #474
- Fix configuration of network device names. #926
- Fix a bug that prevented `wwclient` from starting on systemd hosts. #1066
- Set default ownership of built-in overlay files to `root:root`. #1078
- Fix a SIGSEGV error when building Warewulf on a host without network access. #907
- Use `sysrq-trigger` to reboot during wwinit, as `reboot` may require systemd, which is not yet running. #871
- Only write IPMI configuration if `ipmiwrite` is explicitly set to true. #823
- Don't show an error during `wwctl container list` if images are missing. #933
- Fix a bug preventing containers with symlinked `/bin/sh` from being imported. #797
- Fix a bug that caused syncuser to panic if `/etc/passwd` was malformed. #527
- Fix inclusion of kernel modules for imported kernels. #836
- Fix a bug that caused `wwctl overlay show` to write to stderr. #1000
- Correct "shebang" lines in `wwinit:warewulf/init.d` scripts. #821

### Changed

- Rebase OpenSUSE Leap build to 15.5. #951
- Update default `Makefile` provision directory to match RPM specfile. #972
- Replace deprecated `io.utils` functions with new `os` functions. #883
- Sort node and profiles by name in `wwctl <node|profile> list`. #816
- Use the mountpoint file system for `wwctl exec` bind mounts. #897
- Reduce a warning message during `overlay build` from `Warn` to `Debug`. #1025
- Use the primary interface hostname and the Warewulf server FQDN as the canonical name in `/etc/hosts`. #693
- Refactor `wwctl profile add` to behave more like `wwctl node add`. #658 #659
- Write `wwctl` error messages to stderr rather than stdout. #758
- `wwctl` now validates IP and MAC address format.
- Refactor `Makefile` for clarity and reliability. #890
- Move primary network interface definition to `PrimaryNetDev` on the node. #682
- `wwctl <node|profile> list --all` now only shows attributes that have a set value. #786
- Name system and runtime overlay images as `__SYSTEM__` and `__RUNTIME__` by default. #852, #876, #896
- Always prefer a node's primary network interface during node discovery. #775
- Stop bundling iPXE with Warewulf and use iPXE from the host distribution. #784
- Update and simplify the iPXE build script, now at `scripts/build-ipxe.sh`. #1026
- Merge overlays from multiple profiles on a single node, and exclude overlays with a `~` prefix. #885
- Move `hostlist` package into `internal` directory alongside other Warewulf packages. #804
- Move built-in overlay files to a transparent `rootfs` directory. #1086
- Update quickstart guides. #847, #848

## 4.4.1, 2023-05-07-05

### Fixed

- Properly update container file GIDs during syncuser. #840
- Add a missing `.ww` extension to the `70-ww4-netname.rules` template in the wwinit overlay. #724
- Restrict access to `/warewulf/config` to root only. #728

## 4.4.0, 2023-01-18

### Added

- New Docker container node image for CentOS 7. #621

### Fixed

- Replaced an invalid variable name in a NetworkManager overlay
  template. #626
- The 'nodes' alias now correctly refers to 'node' rather than
  'profile'.
- Fixed a typo in a log message. #631
- Boolean attributes now correctly account for profile and default
  values. #630
- Kernel version is shown correctly for symlink'd kernels #640
- Changing a profile always adds an empty default interface. #661

## 4.4.0rc3, 2022-12-23

### Added

- New `defaults.conf` man page. #593
- A new debug overlay includes a template which demonstrates accessing
  all available variables. #599
- Distribute a README along with staticfiles. #189
- Add a `-y` flag to `wwctl profile add`. #610
- Distribute a source RPM with GitHub releases. #614

### Changed

- No longer ask for confirmation when deleting 0 nodes. #603
- Ask for confirmation during `wwctl container delete`. #606

### Fixed

- `wwctl profile set` now indicates "profiles" in output where it
  previously mistakenly indicated "nodes." #600
- Set correct overlay permissions for a NetworkManager configuration
  file. #591

### Fixed

- Directories within overlays no longer lose group/other write permissions #584

## 4.4.0rc2, 2022-12-09

### Added

- The environment variable `WW_CONTAINER_SHELL` is defined in a `wwctl
  container shell` environment to indicate the container in use. #579
- Network interface configuration (`ifcfg`) files now include the
  interface name and type. #457

### Fixed

- Work-around for older versions of gzip that lack a `--keep` flag
  during `wwctl container build`. #580
- The default ipxe template is once again specified as a built-in
  default and in `defaults.conf`. #581
- `wwctl container list` no longer segfaults when a container chroot
  is present without a built image. #585
- `wwctl configure hostfile` now correctly detects the presence of the
  hostfile overlay template. #571
- `wwctl overlay build` no longer panics when rendering an template
  for a node which has tags set. #568
- Minor typographical fixes. #569

## 4.4.0rc1, 2022-10-27

### Added

- iPXE binaries included with Warewulf now support VLAN tagging. #563
- `wwctl container list` now shows the container creation date,
  modification date, and size. #537
- `wwctl node edit` supports directly editing or defining node
  configuration YAML in an editor. #540
- `wwctl node export` and `wwctl node import` support importing and
  exporting node definitions as YAML or (for import) CSV. The CSV file
  must have a header in where the first field must always be the
  nodename, and the rest of the fields are the same as the long
  commandline options. Network device must have the form
  `net.$NETNAME.$NETOPTION`. (e.g., `net.default.ipaddr`) #540
- The `warewulfd.service` systemd unit file now supports `execreload`
  and `execstop`. #550
- Network interfaces now accept an `mtu` attribute. #549
- The `wwinit` overlay now supports network interface configuration
  via NetworkManager for Ethernet and InfiniBand interfaces. #539
- Default node attribute values (e.g., for kernel arguments) are now
  read in from a `defaults.conf` configuration file. If this file
  is not present, built-in default values are used. #539
- [Warewulf documentation](https://warewulf.org/docs/) is now managed
  alongside the Warewulf source code in a single code repository so
  that documentation may be updated alongside code changes.
- New man pages for `warewulf.conf` and `nodes.conf` #510
- An initial cut of the [Warewulf API](API.md) #471
- `wwctl show --render` shows overlay templates as they would be
  rendered on a given target node. #467
- `wwctl ssh` now supports Bash completion. #466

### Changed

- `wwctl overlay edit` no longer saves a new template to the overlay
  if the template is not modified from its initial state. #522
- The wwinit overlay now only sets a name for a network interface if
  that interface has a MAC address defined. #553
- `wwctl container delete` now also deletes the built images
  associated with that container. #214
- Unified internal code paths for `wwctl profile` and `wwctl node`
  commands, and between the on disk YAML format and the in memory
  format, enabling the command-line options to be autogenerated from
  the datastructures and ensuring that profile and node capabilities
  remain in sync. Multiple command line arguments have been updated or
  changed. #495, #637
- `wwctl power` commands no longer separates node output with
  additional whitespace. #514

### Fixed

- `/etc/warewulf/excludes` (read from the node image) once again
  excludes files from being included in the node image. #532
- `wwctl ssh` always uses a node's primary interface. #544
- `wwctl container show` now correctly shows the kernel version. #542
- System users are no longer prevented from logging into compute
  nodes. #538
- `wwctl overlay chown` now correctly handles uid and gid arguments. #530
- `wwctl overlay chown` no longer sets gid to `0` when unspecified. #531
- Corrected the path for `.wwbackup` files in some situations. #524
- Bypass `imgextract` for legacy BIOS machines to avoid 32-bit memory
  limitations. #497
- `warewulfd` no longer panics when network interface tags are
  defined. #468
- The wwinit overlay now configures the network device type. #465
- Minor typographical fixes. #528, #519

## 4.3.0, 2022-06-25

### Added

- All configurations files for the host (/etc/exports, /etc/dhcpd.conf, /etc/hosts) are now
  populated from the (OVERLAYDIR/host/etc/{exports|dhcpd|hosts}.ww . Also other configuration
  files like prometheus.yml.ww or slurm.conf.ww which depend on the cluster nodes can be
  placed. Also the new templated functions {{ abort }}, {{ IncludeBlock }} abd {{ no_backup }}
  are allowed now.
- nodes and profiles can now have multiple system and runtime overlays, as a comma separated list.
  The overlays of the profile and the nodes are combined.
- simple ipv6 support is now enabled. In `warewulf.conf` the option `ipaddr6`/`Ipv6net` must
  be set to enable ipv6. If enabled on of these options is set a node will get a derived
  ipv6 in the scheme `ipv6net:ipaddr4`. This address can also be overwritten for every
  node
- Multiple files can now created from a single `template.ww` file with the `{{ file FILENAME }}`
  command in the template. The command is expanded to the magic template command
  `{{ /* file FILENAME */}}` which is picked up by wwctl and everything which comes after this
  magic comment will be written to the file `FILENAME`. This mechanism is leveraged in the
  configuration files for the network, see `ifcfg.xml.ww` and `ifcgf.ww`.
- Networks can now have arbitrary keys value pairs in the profiles and on the node, so that
  things like bridges and mtu sizes can be set
- oci container tar balls can be imported with the 'file://$PATH' scheme
- uids and gids of a container can now get synced at import time, so that at least users with the
  same name have the same uid. This is not necessarily needed for warewulf, but services like
  munge.

### Changed

- Provision interface is not tied to 'eth0' any more. The provision interface must be have the
  'primary' flag now. The file `nodes.conf' must be changed accordingly.
- the provisioning network is now called primary and not default
- Creating of '/etc/exports' can now be disabled, so that `wwctl configure -a` wont overwrite
  a existing '/etc/exports'.
- The yaml format for nodes has now sub-keys for ipmi and kernel, old nodes.conf files have to
  to be changed accordingly
- host overlays can globaly disbaled, but are enabled per default
- `wwctl overlay build -H` will only build the overlays which are assigned to the nodes

## 4.1.0, 2021-07-29

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
