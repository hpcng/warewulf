====================================
openSUSE Leap and SLES 15 Quickstart
====================================

Install Warewulf and dependencies
=================================

.. code-block:: bash

   sudo zypper install -t pattern devel_basis
   sudo zypper install go
   sudo zypper install tftp dhcp-server nfs-kernel-server

   sudo systemctl stop firewalld
   sudo systemctl disable firewalld

   git clone https://github.com/warewulf/warewulf.git
   cd warewulf
   PREFIX=/usr SYSCONFDIR=/etc TFTPDIR=/srv/tftproot LOCALSTATEDIR=/var/lib make clean defaults
   make all
   sudo make install

The standard configuration template for the dhcpd service is installed
at the wrong location, you have to fix this with

.. code-block:: bash

   mv /var/lib/warewulf/overlays/host/etc/dhcp/dhcpd.conf.ww /var/lib/warewulf/overlays/host/etc/dhcpd.conf.ww

Install Warewulf from the open build service
============================================

You can also just install the 'warewulf4' package with ``zypper`` from
the openbuild service. Up to date versions are available on the devel
project

``https://build.opensuse.org/project/show/network:cluster``

Configure the controller
========================

Edit the file ``/etc/warewulf/warewulf.conf`` and ensure that you've
set the appropriate configuration paramaters. Here are some of the
defaults for reference assuming that ``192.168.200.1`` is the IP
address of your cluster's private network interface:

.. code-block:: yaml

    ipaddr: 192.168.200.1
    netmask: 255.255.255.0
    network: 192.168.200.0
    warewulf:
      port: 9873
      secure: false
      update interval: 60
      autobuild overlays: true
      host overlay: true
      syslog: false
    dhcp:
      enabled: true
      range start: 192.168.200.50
      range end: 192.168.200.99
      systemd name: dhcpd
    tftp:
      enabled: true
      systemd name: tftp
    nfs:
      enabled: true
      export paths:
      - path: /home
        export options: rw,sync
      - path: /opt
        export options: ro,sync,no_root_squash
      systemd name: nfs-server
    container mounts:
      - source: /etc/resolv.conf
        dest: /etc/resolv.conf
        readonly: true

.. note::

   The DHCP range ends at ``192.168.200.99`` and as you will see
   below, the first node static IP address (post boot) is configured
   to ``192.168.200.100``.

Start and enable the Warewulf service
=====================================

.. code-block:: bash

   # Start and enable the warewulfd service
   sudo systemctl enable --now warewulfd

Configure system services automatically
=======================================

There are a number of services and configurations that Warewulf relies
on to operate.  If you wish to configure all services, you can do so
individually (omitting the ``--all``) will print a help and usage
instructions.

.. note::

   If the ``dhcpd`` service was not used before you will have to add
   the interface on which the cluster network is running to the
   ``DHCP_INTERFACE`` in the file ``/etc/sysconfig/dhcpd``.

.. code-block:: bash

   sudo wwctl configure --all

Pull and build the VNFS container and kernel
============================================

This will pull a basic VNFS container from Docker Hub and import the
default running kernel from the controller node and set both in the
"default" node profile.

.. code-block:: bash

   $ sudo wwctl container import docker://registry.opensuse.org/science/warewulf/leap-15.4/containers/kernel:latest leap15.4 --setdefault

Set up the default node profile
===============================

The ``--setdefault`` arguments above will automatically set those
entries in the default profile, but if you wanted to set them by hand
to something different, you can do the following:

.. code-block:: bash

   sudo wwctl profile set -y -C leap15.4

Next we set some default networking configurations for the first
ethernet device. On modern Linux distributions, the name of the device
is not critical, as it will be setup according to the HW
address. Because all nodes will share the netmask and gateway
configuration, we can set them in the default profile as follows:

.. code-block:: bash

   sudo wwctl profile set -y default --netname default --netmask 255.255.255.0 --gateway 192.168.200.1
   sudo wwctl profile list -a

Add a node
==========

Adding nodes can be done while setting configurations in one
command. Here we are setting the IP address of ``eth0`` and setting
this node to be discoverable, which will then automatically have the
HW address added to the configuration as the node boots.

Node names must be unique. If you have node groups and/or multiple
clusters, designate them using dot notation.

Note that the full node configuration comes from both cascading
profiles and node configurations which always supersede profile
configurations.

.. code-block:: bash

   sudo wwctl node add n0000.cluster --netdev eth0 --ipaddr 192.168.200.100 --discoverable true
   sudo wwctl node list -a n0000.cluster

Warewulf Overlays
=================

There are two types of overlays: system and runtime overlays.

System overlays are provisioned to the node before ``/sbin/init`` is
called. This enables us to prepopulate node configurations with
content that is node specific like networking and service
configurations.

Runtime overlays are re-applied
periodically during the normal runtime of the node. Because these
overlays are provisioned at periodic intervals, they are very useful
for content that changes, like users and groups.

Overlays are generated from a template structure that is viewed using
the ``wwctl overlay`` commands. Files that end in the ``.ww`` suffix
are templates and abide by standard text/template rules. This supports
loops, arrays, variables, and functions making overlays extremely
flexible.

All overlays are compiled before being provisioned. This accelerates
the provisioning process because there is less to do when nodes are
being managed at scale.

Here are some of the common ``overlay`` commands:

.. code-block:: bash

   sudo wwctl overlay list -l
   sudo wwctl overlay list -ls
   sudo wwctl overlay edit default /etc/hello_world.ww
   sudo wwctl overlay build -a

Boot your compute node and watch it boot!
