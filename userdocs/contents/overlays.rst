=================
Warewulf Overlays
=================

Warewulf is designed to scalably provision and manage thousands of cluster nodes by utilizing
identical stateless boot images. But if these boot images are completely identical, then how do we
configure things like hostnames? IP addresses? Or any other node-specific configurations?

Some of configuration can be managed by services like DHCP. You can also use traditional
configuration management on a provisioned Warewulf cluster node. But these are heavy-weight
solutions to a simple problem.

Warewulf addresses cluster node configuration with its **overlay** system. Overlays are collections
of files and templates that are rendered and built per-node and then applied over the container
image during the provisioning process.

System and runtime overlays
===========================

A node or profile can configure an overlay in two different ways:

* An overlay can be configured to apply only during boot as part of the ``wwinit`` process. These
  overlays are called **system overlays**.
* An overlay can be configured to also apply periodically while the system is running. These overlays
  are called **runtime overlays**.

Distribution and site overlays
==============================

Warewulf also distinguishes between **distribution** overlays, which are included with Warewulf, and
**site** overlays, which are created or added locally. A site overlay always takes precedence over a
distribution overlay with the same name.  Any modification of a distribution overlay with ``wwctl``
actually makes changes to an automatically-generated **site** overlay cloned from the distribution
overlay.

Site overlays are often stored at ``/var/lib/warewulf/overlays/``. Distribution overlays are often
stored at ``/usr/share/warewulf/overlays/``. But these paths are dependent on compilation,
distribution, packaging, and configuration settings.

Provided distribution overlays
------------------------------

These overlays are provided as part of Warewulf.

wwinit
------

The **wwinit** overlay performs initial configuration of the Warewulf node.
Its `wwinit` script runs before ``systemd`` or other init is called and
contains all configurations which are needed to boot.

In particular:

- Configure the loopback interface
- Configure the BMC based on the node's configuration
- Update PAM configuration to allow missing shadow entries
- Relabel the file system for SELinux

Other overlays may place additional scripts in ``/warewulf/init.d/`` to affect
node configuration in this pre-boot environment.

wwclient
--------

All configured overlays are provisioned initially along with the node image
itself; but **wwclient** periodically fetches and applies the runtime overlay
to allow configuration of some settings without a reboot.

Network interfaces
------------------

Warewulf ships with support for many different network interface configuration
systems. All of these are applied by default; but the list may be trimmed to
the desired system.

- ifcfg
- NetworkManager
- debian.interfaces
- wicked

Warewulf also configures both systemd and udev with the intended names of
configured network interfaces, typically based on a known MAC address.

- systemd.netname
- udev.netname

Basics
------

The **hostname** overlay sets the hostname based on the configured Warewulf
node name.

The **hosts** overlay configures ``/etc/hosts`` to include all Warewulf nodes.

The **issue** overlay configures a standard Warewulf status message for display
during login.

The **resolv** overlay configures ``/etc/resolv.conf`` based on the value of
"DNS" nettags. (In most situations this should be unnecessary, as the network
interface configuration should handle this dynamically.)

fstab
-----

The **fstab** overlay configures ``/etc/fstab`` based on the data provided in the "fstab"
resource. It also creates entries for file systems defined by Ignition.

.. code-block:: yaml

   nodeprofiles:
     default:
       resources:
         fstab:
           - spec: warewulf:/home
             file: /home
             vfstype: nfs
           - spec: warewulf:/opt
             file: /opt
             vfstype: nfs

ssh
---

Two SSH overlays configure host keys (one set for all node in the cluster) and
``authorized_keys`` for the root account.

- ssh.authorized_keys
- ssh.host_keys

syncuser
--------

The **syncuser** overlay updates ``/etc/passwd`` and ``/etc/group`` to include
all users on both the Warewulf server and from the container image. To function
properly, ``wwctl container syncuser`` must have also been run on the container
image to synchronize its user and group IDs with those of the server.

ignition
--------

The **ignition** overlay defines partitions and file systems on local disks.

debug
-----

The **debug** overlay is not intended to be used in configuration, but can be
used as an example and applied with ``wwctl overlay show --render <nodename>
debug /warewulf/template-variables.md.ww``.

host
----

Configuration files used for the configuration of the Warewulf host /
server are stored in the **host** overlay. Unlike other overlays, it
*must* have the name ``host``. Existing files on the host are copied
to backup files with a ``wwbackup`` suffix at the first
run. (Subsequent use of the host overlay won't overwrite existing
``wwbackup`` files.)

The following services get configuration files via the host overlay:

* ssh keys are created with the scrips ``ssh_setup.sh`` and
  ``ssh_setup.csh``
* hosts entries are created by manipulating ``/etc/hosts`` with the
  template ``hosts.ww``
* nfs kernel server receives its exports from the template
  ``exports.ww``
* the dhcpd service is configured with ``dhcpd.conf.ww``

Combining and overriding overlays
=================================

Multiple overlays can be applied to a single node, and overlays from multiple profiles are appended
together. The configuration fields for the system and runtime overlays are lists and can contain
several overlays. As an example for this, we will overwrite the ``/etc/issue`` file from the "issue"
overlay. For this we will create a new overlay called "welcome" and import the file ``/etc/issue``
from the host to it. This overlay is then combined with the existing overlays.

.. code-block:: console

  # wwctl overlay create welcome
  # wwctl overlay mkdir welcome /etc
  # wwctl overlay import welcome /etc/issue
  # wwctl profile set default --wwinit=wwinit,wwclient,welcome
  ? Are you sure you want to modify 1 profile(s)? [y/N] y
  # wwctl profile list default -a |grep welcome
  default              SystemOverlay      wwinit,wwclient,welcome

Templates
=========

Templates allow you to create dynamic content such that the files
downloaded for each node will be customized for that node. Templates
allow you to insert everything from variables, to including files from
the control node, as well as conditional content and loops.

Warewulf uses the ``text/template`` engine to facilitate implementing
dynamic content in a simple and standardized manner.

All template files will end with the suffix of ``.ww``. That tells
Warewulf that when building a file, that it should parse that file as
a template. When it does that, the resulting file is static and can
have node customizations that are obtained from the node configuration
attributes.

.. note::

   When the file is persisted within the built overlay, the ``.ww``
   will be dropped, so ``/etc/hosts.ww`` will end up being
   ``/etc/hosts``.

Managing overlays
=================

Warewulf includes a command group for manipulating overlays (``wwctl
overlay``). With this you can add, edit, remove, change ownership,
permissions, etc.

..
  note::
  It is not possible to delete files with an overlay.

Build
-----

.. code-block:: console

  wwctl overlay build [-H,--hosts|-N,--nodes|-o,--output directory|-O,--overlay-name] nodepattern

Without any arguments the command will interpret the templates for all
overlays for every compute node and also all the templates in the host
overlay. For every overlay of the compute nodes a gzip compressed cpio
archive is created. The range of the nodes can be restricted as the
last argument.  With the ``-H`` flag only the host overlay is
built. With the ``-N`` flag only compute node overlays are
built. Specific overlays can be selected with ``-O`` flag. For
debugging purposes the templates can be written to a directory given
via the ``-o`` flag.

Overlay images for multiple node are built in parallel. By default, each CPU in
the Warewulf server will build overlays independently. The number of workers
can be specified with the ``--workers`` option.

Warewulf will attempt to build/update overlays as needed
(configurable in the ``warewulf.conf``); but not all cases are detected,
and manual overlay builds are often necessary.

Chmod
-----

.. code-block:: console

  wwctl overlay chmod overlay-name filename mode

This subcommand changes the permissions of a single file within an
overlay. You can use any mode format supported by the chmod command.

Chown
-----

.. code-block:: console

  wwctl overlay chown overlay-name filename UID [GID]

With this command you can change the ownership of a file within a
given overlay to the user specified by UID. Optionally, it will also
change group ownership to GID.

Create
------

.. code-block:: console

  wwctl overlay create overlay-name

This command creates a new empty overlay with the given name.

Delete
------

.. code-block:: console

  wwctl overlay delete [-f,--force] overlay-name [File [File ...]]

Either the given overlay is deleted (must be empty or use the
``--force`` flag) or the specified file within the overlay is
deleted. With the ``--parents`` flag the directory of the deleted file
is also removed if no other file is in the directory.

Edit
----
.. code-block:: console

  wwctl overlay edit [--mode,-m MODE|--parents,-p] overlay-name file

Use this command to edit an existing or a new file in the given
overlay. If a the new file ends with a ``.ww`` suffix an example
template header is added to the file. With the ``--parents`` flag
necessary parent directories for a new file are created.

Import
------
.. code-block:: console

  wwctl overlay import [--mode,-m|--noupdate,-n] overlay-name file-name [new-file-name]

The given file is imported to the overlay. If no new-file-name is
given, the file will be placed in the overlay at the same path as on
the host. With the ``--noupdate`` flag you can block the rebuild of
the overlays.

List
----

.. code-block:: console

  wwctl overlay list [--all,-a|--long,-l] [overlay-name]

With this command all existing overlays and files in them can be
listed. Without any option only the overlay names and their number of
files are listed. With the ``--all`` switch also the every file is
shown. The ``--long`` option will also display the permissions, UID,
and GID of each file.

Show
----

.. code-block:: console

  wwctl overlay show [--quiet,-q|--render,-r nodename] overlay-name file

The content of the file for the given overlay is displayed with this
command. With the ``--render`` option a template is rendered as it
will be rendered for the given node. The node name is a mandatory
argument to the ``--render`` flag. Additional information for the file
can be suppressed with the ``--quiet`` option.
