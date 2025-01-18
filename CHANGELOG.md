# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## v4.6.0, unreleased

### Added

- Added Netplan NIC support for Debian/Ubuntu #1463
- Added documentation on ensuring `systemctl restart warewulfd` is ran when editing `nodes.conf` or `warewulf.conf`
- Add the ability to boot nodes with `wwid=[interface]`, which replaces
  `interface` with the interface MAC address
- Added https://github.com/Masterminds/sprig functions to templates #1030
- Add multiple output formats (yaml & json) support. #447
- More aliases for many wwctl commands
- Add support to render template using `host` or `$(uname -n)` as the value of `overlay show --render`. #623
- Added command line parameters for credentials of a container registry
- Add flag `--build` to `wwctl container copy`. #1378
- Add `wwctl clean` to remove OCI cache and overlays from deleted nodes
- Add `wwctl container import --platform`. #1381
- Read environment variables from `/etc/default/warewulfd` #725
- Add support for VLANs to NetworkManager, wicked, ifcfg, debian.network_interfaces overlays. #1257
- Add support for static routes to NetworkManager, wicked, ifcfg, debian.network_interfaces overlays. #1257
- Add `wwctl upgrade <config|nodes>`. #230, #517
- Better handling of InfiniBand udev net naming. #1227
- use templating mechanism for power commands. #1004
- Document "known issues."
- Add `wwctl <node|profile> <add|set> --kernelversion` to specify the desired kernel version or path. #1556
- Add `wwctl container kernels` to list discovered kernels from containers. #1556
- Add possibility to define a softlink target with an overlay template
- Support defining a symlink with an overlay template. #1303
- New "localtime" overlay to define the system time zone. #1303
- Add support for nested profiles. #1572, #1598
- Adds `wwctl container <exec|shell> --build=false` to prevent automatically (re)building the container. #1490, #1489
- Added resources as generic, arbitrary YAML data for nodes and profiles. #1568
- New `fstab` resource configures mounts in fstab overlay, including NFS mounts. #515

### Changed

- Locally defined `tr` has been dropped, templates updated to use Sprig replace.
- Bump github.com/opencontainers/image-spec to 1.1.0
- Bump google.golang.org/grpc 1.62.1
- Bump google.golang.org/protobuf to 1.33.0
- Bump github.com/containers/image/v5 to 5.30.0
- Bump github.com/docker/docker to 25.0.5+incompatible
- Bump github.com/grpc-ecosystem/grpc-gateway/v2 from 2.18.0 to 2.19.1 #1165
- Bump github.com/spf13/cobra from 1.7.0 to 1.8.0 #1166
- Bump github.com/fatih/color from 1.15.0 to 1.17.0 #1224
- Bump github.com/coreos/ignition/v2 from 2.15.0 to 2.19.0 #1239
- Bump github.com/spf13/cobra from 1.8.0 to 1.8.1 #1481
- Bump google.golang.org/protobuf from 1.34.1 to 1.35.1 #1480
- Bump golang.org/x/term from 0.20.0 to 0.25.0 #1476
- Bump github.com/grpc-ecosystem/grpc-gateway/v2 from 2.19.1 to 2.23.0 #1513
- Bump github.com/containers/image/v5 from 5.30.1 to 5.32.2 #1366
- Bump github.com/fatih/color from 1.17.0 to 1.18.0 #1523
- Disable building containers by default when calling `wwctl container copy`. #1378
- Split wwinit and generic overlays into discrete functionality. #987
- Updated IgnitionJson to sort filesystems. #1433
- `wwctl node set` requires mandatory pattern input. #502
- Remove NodeInfo (in-memory-only) data structure, consolidating onto NodeConf. #916
- Replace `defaults.conf` with settings on the default profile. #917
- Switched from yaml.v2 to yaml.v3 #1462
- Make OCIBlobCache a seperate path and point it to `/var/cache` #1459
- Updated various shell scripts for POSIX compatibility. #1464
- Update `wwctl server` to always run in the foreground #508
- Update `wwctl server` to log to stdout rather than a file #503
- Changed `wwctl server` to use "INFO" for send and receive logs #725
- Remove a 3-second sleep during iPXE boot. #1500
- Don't package the API in RPM packages by default. #1493
- Update default `warewulfd` port to match shipped configuration. #1448
- Replace `olekukonko/tablewriter` with `cheynewallace/tabby`. #1497, #1498
- replaced deprecated errors.Wrapf with fmr.Errorf. #1534
- Rename udev net naming file to 70-persistent-net.rules. #1227
- Manage warewulfd template data as a pointer. #1548
- Added test for sending grub.cfg.ww. #1548
- Use a sentinel file to determine container readonly state. #1447
- Bump github.com/Masterminds/sprig/v3 from 3.2.3 to 3.3.0 #1553
- Bump github.com/golang/glog from 1.2.0 to 1.2.3 #1527
- Bump github.com/opencontainers/runc from 1.1.12 to 1.1.14
- Repurpose Kernel.Override to specify the path to the desired kernel within the container. #1556
- Merge Kernel.Override into Kernel.Version to specify the desired kernel version or path. #1556
- Provide detected kernel version to overlay templates. #1556
- Bump github.com/containers/storage from 1.53.0 to 1.55.2 #1316, #892
- Process nodes.conf path dynamically from config. #1595, #1596, #1569
- Split overlays into distribution and site overlays. #831
- Added note to booting userdoc for removing machine-id. #1609
- Log cpio errors more prominently. #1615
- Improved syncuser conflict help text. #1614
- Parallelized overlay build. #1018
- Parallelized and optimized overlay build. #1018
- Added note about dnsmasq interface options in Rocky 9.
- Added retries to curl in wwinit dracut module. #1631

### Removed

- `wwctl node list --fullall` has been removed
- `wwctl profile list --fullall` has been removed
- Remove `wwctl server <start,stop,status,restart,reload>` #508
- Remove `wwctl overlay build --host` #1419
- Remove `wwctl overlay build --nodes` #1419
- Remove `wwctl kernel` #1556
- Remove `wwctl <node|profile> <add|set> --kerneloverride` #1556
- Remove `wwctl container <build|import> --setdefault` #1335
- Remove NFS mount options from warewulf.conf. #515

### Fixed

- Update links on contributing page to point to warewulf repo.
- Prevent Networkmanager from trying to optain IP address via DHCP
  on unused/unmanaged network interfaces.
- Systems with no SMBIOS (Raspberry Pi) will create a UUID from
  `/sys/firmware/devicetree/base/serial-number`
- Replace slice in templates with sprig substr. #1093
- Fix an invalid format issue for the GitHub nightly build action. #1258
- Return non-zero exit code on overlay build failure #1393
- Return non-zero exit code on container copy failure #1377
- Return non-zero exit code on container sub-commands #1414
- Fix excessive line spacing issue when listing nodes. #1241
- Return non-zero exit code on node sub-commands #1421
- Fix panic when getting a long container list before building the container. #1391
- Return non-zero exit code on power sub-commands #1439
- Fix issue that pattern matching broken on `node set` #964
- Fix issue that domain globs not supported during wwctl node delete. #1449
- Fix overlay permissions in /root/ and /root/.ssh/. #1452
- Return non-zero exit code on container sub-commands #1437
- Return non-zero exit code on profile sub-commands #1435
- Fix issue that NetworkManager marks managed interfaces "unmanaged" if they do
  not have a device specified. #1154
- Return non-zero exit code on overlay sub-commands #1423
- Simplify passing of arguments to commands through `wwctl container exec`. #253
- Don't update IPMI if password isn't set. #638
- Fix issue that `--nettagdel` does not work properly. #1503
- Fix test for dhcp static configuration #1536 #1537
- Fix issue that initrd fails at downloading runtime overlay with permission denied error,
  when warewulf secure option in warewulf.conf is enabled. #806
- Allow iPXE to continue booting without runtime overlay. #806
- Format errors in logs as strings. #1563
- Fix display of profiles during node list. #1496
- Fix internal DelProfile function to correctly operate on profiles rather than nodes. #1622
- Fix parsing of bool command line variables #1627

## v4.5.8, 2024-10-01

### Added

- Added `--syncuser` flag to `wwctl container shell`. #1358
- Added a troubleshooting guide. #1234
- Added documentation about `rootfstype=ramfs` for SELinux support. #1001
- Added workaround documentation for importing containers with sockets. #892
- Added documentation for building iPXE locally. #1114
- Documented that ignition is not available for Rocky Linux 8. #1373, #1272
- Additional help text when container `RunDir` already exists. #1389

### Changed

- Interleave tmpfs across all available NUMA nodes. #1347, #1348
- Syncuser watches for changes in mtime rather than ctime. #1358
- Change the default permissions for provisioned overlay images to `0750` (dirs) and `0660` (files). #1388

### Fixed

- Return an error during `wwctl container import` if the archive filename includes a colon. #1371
- Correctly extract smbios asset key during GRUB boot. #1291
- Refactor of `wwinit/init` to more properly address rootfs options. #1098
- Fix autodetected kernel sorting and filtering. #1332
- Avoid a panic during container import. #1244
- Make sure that tftp files have unmasked permissions at creation time. #674
- Fix "onboot" behavior for NetworkManager, Debian networking, and Suse wicked. #1278
- Clarified missing steps in Enterprise Linux quickstart. #1179
- Fix dhcpd.conf static template to include next-server and dhcp-range #1536
- Fix panic when adding tag with existing netdev #1546


## v4.5.7, 2024-09-11

### Added

- Added option for wwclient port number. #1349
- Additional helper directions during syncuser conflict. #1359
- Add `:copy` suffix to `wwctl container exec --bind` to temporarily copy files into the node image. #1365

### Changed

- Added a link to an example SELinux-enabled node image in documentation. #1305
- Refine error handling for `wwctl configure`. #1273
- Updated dracut guidance for building initramfs. #1369

### Fixed

- Fixed application of node overlays such that they override overlapping files from profile overlays. #1259
- Prevent overlays from being improperly used as format strings during `wwctl overlay show --render`. #1363
- Fix dracut booting with secure mode. #1261

## v4.5.6, 2024-08-05

### Added

- Show more information during `wwctl container <shell|exec>` about when and if the container image will be rebuilt. #1302
- Command-line completion for `wwctl overlay <edit|delete|chmod|chown>`. #1298
- Display an error during boot if no container is defined. #1295
- `wwctl conatiner list --kernel` shows the kernel detected for each container. #1283
- `wwctl container list --size` shows the uncompressed size of each container. `--compressed` shows the compressed size, and `--chroot` shows the size of the container source on the server. #954, #1117
- Add a logrotate config for `warewulfd.log`. #1311
### Changed

- Refactor URL handling in wwclient to consistently escape arguments.

### Fixed

- Ensure autobuilt overlays include contextual overlay contents. #1296
- Fix the failure when updating overlay files existing on different partitions. #1312
- Escape asset tag for `wwclient` query strings when pulling runtime overlays. #1310

### Changed

- `wwctl container list` only lists names by default. (`--long` shows all attributes.) #1117

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
