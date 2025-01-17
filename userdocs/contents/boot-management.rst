===============
Boot Management
===============

Warewulf uses iPXE to for network boot by default. As a tech preview, support
for GRUB is also available, which adds support for secure boot.

Also as a tech preview, Warewulf may also use iPXE to boot a dracut
initramfs as an initial stage before loading the container image.

Booting with iPXE
=================

.. graphviz::

  digraph G{
      node [shape=box];
      compound=true;
      edge [label2node=true]
      bios [shape=record label="{BIOS | boots from DHCP/next-server via TFTP}"]

      subgraph cluster0 {
       label="iPXE boot"
       iPXE;
       ipxe_cfg [shape=record label="{ipxe.cfg|generated for each node}"];
       iPXE -> ipxe_cfg [label="http"];
      }

      bios->iPXE [lhead=cluster0,label="iPXE.efi"];

      kernel [shape=record label="{kernel|ramdisk (root fs)|wwinit overlay}|extracted from node container"];
      ipxe_cfg->kernel[ltail=cluster0,label="http"];
  }

Starting in v4.5.0, Warewulf no longer includes an iPXE binary. In stead, by
default Warewulf uses the iPXE that comes with the host OS.

Unfortunately, we’ve encountered a few instances where bugs in the OS-provided
iPXE that sometimes make booting a full OS image as an "initrd" unreliable.

:ref:`Building iPXE locally`, using a more recent "version" of the iPXE source
code, can alleviate some of these issues.

Another alternative is :ref:`booting with dracut`, which uses the Linux kernel
to load the full OS image, avoiding the issue entirely.

.. _Building iPXE locally:

Building iPXE locally
---------------------

By default (as of v4.5.0) Warewulf packages use iPXE from the host operating system rather than bundling iPXE binaries with Warewulf.
However, sometimes the specific build included in the host OS has bugs or missing features, and a local build of iPXE is necessary.

The Warewulf project provides a `build-ipxe.sh`_ script to simplify the process of building iPXE locally.

.. _build-ipxe.sh: https://github.com/warewulf/warewulf/blob/main/scripts/build-ipxe.sh

.. code-block:: console

   # curl -LO https://raw.githubusercontent.com/warewulf/warewulf/main/scripts/build-ipxe.sh
   # bash build-ipxe.sh -h
   Usage: build-ipxe.sh
            [-h] (help)
   TARGETS: bin-x86_64-pcbios/undionly.kpxe bin-x86_64-efi/snponly.efi bin-arm64-efi/snponly.efi
   IPXE_BRANCH: master
   DESTDIR: /usr/local/share/ipxe

Running build-ipxe.sh
^^^^^^^^^^^^^^^^^^^^^

The script, by default, builds iPXE for x86_64 BIOS, x86_64 EFI, and arm64 EFI from the master branch on the iPXE project GitHub and stores the resultant builds in ``/usr/local/share/ipxe/``.
(These parameters can be adjusted by setting ``TARGETS``, ``IPXE_BRANCH``, and ``DESTDIR`` environment variables, with the current values shown in the ``-h`` output for reference.)

.. code-block:: console

   # mkdir -p /usr/local/share/ipxe
   # bash build-ipxe.sh
   [...]
   # ls -1 /usr/local/share/ipxe/
   bin-arm64-efi-snponly.efi
   bin-x86_64-efi-snponly.efi
   bin-x86_64-pcbios-undionly.kpxe

.. note::

   Building for aarch64 requires the package ``gcc-aarch64-linux-gnu``.

Build options
^^^^^^^^^^^^^

By default, ``build-ipxe.sh`` enables support for `ZLIB`_ and `GZIP`_ images, as well as commands for managing `VLANs`_ and the `framebuffer console`_.
The x86_64 build also enables support for the `serial console`_.

.. _ZLIB: https://ipxe.org/buildcfg/image_zlib

.. _GZIP: https://ipxe.org/buildcfg/image_gzip

.. _VLANs: https://ipxe.org/buildcfg/vlan_cmd

.. _framebuffer console: https://ipxe.org/buildcfg/console_framebuffer

.. _serial console: https://ipxe.org/buildcfg/console_serial

Additional `build options`_ can be configured by editing the ``build-ipxe.sh`` script.
For example, the x86_64 build is configured in the ``configure_x86_64`` function.

.. _build options: https://ipxe.org/buildcfg

.. code-block:: bash

   function configure_x86_64 {
     sed -i.bak \
         -e 's,//\(#define.*CONSOLE_SERIAL.*\),\1,' \
         -e 's,//\(#define.*CONSOLE_FRAMEBUFFER.*\),\1,' \
         config/console.h
     sed -i.bak \
         -e 's,//\(#define.*IMAGE_ZLIB.*\),\1,' \
         -e 's,//\(#define.*IMAGE_GZIP.*\),\1,' \
         -e 's,//\(#define.*VLAN_CMD.*\),\1,' \
         config/general.h
   }

For example, the ``imgextract`` command can be `explicitly enabled`_.

.. _explicitly enabled: https://ipxe.org/buildcfg/image_archive_cmd

.. code-block:: bash

   function configure_x86_64 {
     sed -i.bak \
         -e 's,//\(#define.*CONSOLE_SERIAL.*\),\1,' \
         -e 's,//\(#define.*CONSOLE_FRAMEBUFFER.*\),\1,' \
         config/console.h
     sed -i.bak \
         -e 's,//\(#define.*IMAGE_ZLIB.*\),\1,' \
         -e 's,//\(#define.*IMAGE_GZIP.*\),\1,' \
         -e 's,//\(#define.*VLAN_CMD.*\),\1,' \
         -e 's,//\(#define.*IMAGE_ARCHIVE_CMD.*\),\1,' \
         config/general.h
   }

.. note::

   ``IMG_ARCHIVE_CMD`` is already enabled by default in the iPXE master branch, but only takes effect when at least one archive image format is configured.
   This is the case in the default state of ``build-ipxe.sh``, which enables support for ZLIB and GZIP archive image formats.

Configuring Warewulf (≥ v4.5.0)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

In Warewulf v4.5.0, Warewulf can be configured to use these files using the ``tftp.ipxe`` and ``paths.ipxesource`` configuration parameters in ``warewulf.conf``.

.. code-block:: yaml

   # warewulf.conf
   tftp:
     ipxe:
       "00:00": bin-x86_64-pcbios-undionly.kpxe
       "00:07": bin-x86_64-efi-snponly.efi
       "00:09": bin-x86_64-efi-snponly.efi
       "00:0B": bin-arm64-efi-snponly.efi
   paths:
     ipxesource: /usr/local/share/ipxe

Restart ``warewulfd`` following the change to ``warewulf.conf``.
Then remove any previously-provisioned files from ``/var/lib/tftpboot/warewulf/`` and use ``wwctl configure tftp`` and ``wwctl configure dhcp`` to re-provision the TFTP files and update the DHCP configuration.

.. code-block:: console

   # sudo systemctl restart warewulfd
   # rm /var/lib/tftpboot/warewulf/*
   # wwctl configure tftp
   Writing PXE files to: /var/lib/tftpboot/warewulf
   Enabling and restarting the TFTP services
   # wwctl configure dhcp
   Building overlay for wwctl1: host
   Enabling and restarting the DHCP services

Configuring Warewulf (< v4.5.0)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

Prior to v4.5.0, Warewulf packages included bundled builds of iPXE and did not provide a mechanism for configuring which iPXE to use.
To use a custom iPXE before v4.5.0, replace the bundled builds included with Warewulf.
After that, remove any previously-provisioned files from ``/var/lib/tftpboot/warewulf/`` and use ``wwctl configure tftp`` to re-provision the TFTP files.

.. code-block:: console

   # cp /usr/local/share/ipxe/bin-arm64-efi-snponly.efi /usr/share/warewulf/ipxe/arm64.efi
   # cp /usr/local/share/ipxe/bin-x86_64-efi-snponly.efi /usr/share/warewulf/ipxe/x86_64.efi
   # cp /usr/local/share/ipxe/bin-x86_64-pcbios-undionly.kpxe /usr/share/warewulf/ipxe/x86_64.kpxe
   # rm /var/lib/tftpboot/warewulf/*
   # wwctl configure tftp
   Writing PXE files to: /var/lib/tftpboot/warewulf
   Enabling and restarting the TFTP services

Booting with GRUB
=================

Support for GRUB as a network bootloader (replacing iPXE) is available in
Warewulf as a technology preview.

.. graphviz::

  digraph G{
      node [shape=box];
      compound=true;
      edge [label2node=true]
      bios [shape=record label="{BIOS | boots from DHCP/next-server via TFTP}"]

      bios->shim [lhead=cluster1,label="shim.efi"];
      subgraph cluster1{
        label="Grub boot"
        shim[shape=record label="{shim.efi|from ww4 host}"];
        grub[shape=record label="{grubx64.efi | name hardcoded in shim.efi|from ww4 host}"]
        shim->grub[label="TFTP"];
        grubcfg[shape=record label="{grub.cfg|static under TFTP root}"];
        grub->grubcfg[label="TFTP"];
      }
      kernel [shape=record label="{kernel|ramdisk (root fs)|wwinit overlay}|extracted from node container"];
      grubcfg->kernel[ltail=cluster1,label="http"];
  }

Instead of the iPXE starter a combination of `shim and GRUB
<https://www.suse.com/c/uefi-secure-boot-details/>`_ can be used with the
advantage that secure boot can be used. That means that only the signed kernel
of a distribution can be booted. This can be a huge security benefit for some
scenarios.

In order to enable the grub boot method it has to be enabled in `warewulf.conf`.

.. code-block:: yaml

   warewulf:
     grubboot: true

Nodes which are not known to Warewulf are booted with the shim/grub from the
Warewulf server host.

Secure boot
-----------

.. graphviz::

   digraph foo {
      node [shape=box];
      subgraph boot {
        "EFI" [label="EFI",row=boot];
        "Shim" [label="Shim",row=boot];
        "Grub" [label="Grub",row=boot];
        "Kernel" [label="kernel",row=boot];
        EFI -> Shim[label="Check for Microsoft signature"];
        Shim -> Grub[label="Check for Distribution signature"];
        Grub->Kernel[label="Check for Distribution or MOK signature"];
      }
    }

If secure boot is enabled at every step a signature is checked and the boot
process fails if this check fails. The shim typically only includes the key for
a single operating system, which means that each distribution needs separate
`shim` and `grub` executables. Warewulf extracts these binaries from the
containers. If the node is unknown to Warewulf or can't be identified during
the TFTP boot phase, the shim/grub binaries of the host in which Warewulf is
running are used.

Install shim and efi
--------------------

`shim.efi` and `grub.efi` must be installed in the container for it to be
booted by GRUB.

.. code-block:: console

  # wwctl container shell leap15.5
  [leap15.5] Warewulf> zypper install grub2 shim

  # wwctl container shell rocky9
  [rocky9] Warewulf> dnf install shim-x64.x86_64 grub2-efi-x64.x86_64

These packages must also be installed on the Warewulf server host to enable
node discovery using GRUB.

http boot
---------

Modern EFI systems have the possibility to directly boot per http. The flow diagram
is the following:

.. graphviz::

  digraph G{
      node [shape=box];
      efi [shape=record label="{EFI|boots from URI defined in filename}"];
      shim [shape=record label="{shim.efi|replaces shim.efi with grubx64.efi in URI|extracted from node container}"];
      grub [shape=record label="{grub.efi|checks for grub.cfg|extracted from node container}"]
      kernel [shape=record label="{kernel|ramdisk (root fs)|wwinit overlay}|extracted from node container"];
      efi->shim [label="http"];
      shim->grub [label="http"];
      grub->kernel [label="http"];
    }

Warewulf delivers the initial `shim.efi` and `grub.efi` via http as taken
directly from the node's assigned container.

.. _booting with dracut:

Booting with dracut
===================

Some systems, typically due to limitations in their BIOS or EFI
firmware, are unable to load container image of a certain size
directly with a traditional bootloader, either iPXE or GRUB. As a
workaround for such systems, Warewulf can be configured to load a
dracut initramfs from the container and to use that initramfs to load
the full container image.

Warewulf provides a dracut module to configure the dracut initramfs to
load the container image. This module is available in the
``warewulf-dracut`` subpackage, which must be installed in the
container image.

With the ``warewulf-dracut`` package installed, you can build an
initramfs inside the container.

EL installation
---------------
.. code-block:: shell

   dnf -y install warewulf-dracut
   dracut --force --no-hostonly --add wwinit --regenerate-all


SUSE installation
-----------------
.. code-block:: shell

   zypper -y install warewulf-dracut
   dracut --force --no-hostonly --add wwinit --regenerate-all


.. note::

   In some systems, such as ``rockylinux:8``, it may be 
   necessary to remove ``/etc/machine-id`` for dracut to properly generate 
   the initramfs in the location that Warewulf is expecting.

Set the node's iPXE template to ``dracut`` to direct iPXE to fetch the
node's initramfs image and boot with dracut semantics, rather than
booting the node image directly.

.. note::

   Warewulf iPXE templates are located at ``/etc/warewulf/ipxe/`` when
   Warewulf is installed via official packages. You can learn more
   about how dracut booting works by inspecting its iPXE template at
   ``/etc/warewulf/ipxe/dracut.ipxe``.

.. code-block:: shell

   wwctl node set wwnode1 --ipxe dracut

.. note::

   The iPXE template may be set at the node or profile level.

Alternatively, to direct GRUB to fetch the node's initramfs image and boot with
dracut semantics, set a ``GrubMenuEntry`` tag for the node.

.. note::

   Warewulf configures GRUB with a template located at
   ``/etc/warewulf/grub/grub.cfg.ww``. Inspect the template to learn more about
   the dracut booting process.

.. code-block:: shell

   wwctl node set wwnode1 --tagadd GrubMenuEntry=dracut

.. note::

   The ``GrubMenuEntry`` variable may be set at the node or profile level.

During boot, ``warewulfd`` will detect and dynamically serve an
initramfs from a node's container image in much the same way that it
can serve a kernel from a container image. This image is loaded by
iPXE (or GRUB) which directs dracut to fetch the node's container image
during boot.

The wwinit module provisions to tmpfs. By default, tmpfs is permitted
to use up to 50% of physical memory. This size limit may be adjustd
using the kernel argument `wwinit.tmpfs.size`. (This parameter is
passed to the `size` option during tmpfs mount. See ``tmpfs(5)`` for
more details.)

.. warning::

   Kernel overrides are not currently fully supported during dracut initramfs boot.


Persisitent installation
========================
With the `dracut` installation enabled warewulf can also install 
the node image to a harddrive. On the first boot of the node the
compressed node image is simply dumped onto the configured partition.
Subsequent reboots will update only image with `rsync` except the
node image has changed. Then the configured partition will be erased and 
the configured node image will be dumped on the disk.

.. warning::
   warewulf doesn't install the bootloader to the disk and add UEFI 
   entries. In order to boot the node network booting is required and 
   at every boot the kernel and the initrd is transfered over the network.

Configuration
-------------

As the first step you will have to enable the `rsyncd` (which must be installed) 
Management in `warewulf.conf`. For this simply add

.. code-block:: yaml
   rsync:
     systemd name: rsyncd
     enabled: true

and configure the service with `wwctl configure rsync` which will
also generate a appropritate `rsyncd.conf`.

The node image will be installed to the partition called `rootfs`. You 
can add add a rootfs with e.g. following command

.. code-block:: shell

   wwctl node set n01 \
     --diskname /dev/vda --diskwipe \
     --partname rootfs --partcreate --partnumber 1 \
     --fsname rootfs --fsformat ext4 --fspath /rootfs

Now the node boot method has to be set to persitent with following command

.. code-block:: shell

   wwctl node set n01 --ipxe persistent

Finally whe have create the dracut image for the persten install with 

.. code-block:: shell
   wwctl container shell leap15.6 --bind /usr/lib/dracut/modules.d/90wwinit/
   zypper install-y jq dmidecode
   dracut --force --no-hostonly --add wwinit --regenerate-all

After a reboot the image should be on the disk.

.. note::

   If the boot mode is persitent the configured partion labeled `rootfs` will
   be mounted as `/`. With any other boot method mount point confiured (here `/roofs`)
   will be used.


