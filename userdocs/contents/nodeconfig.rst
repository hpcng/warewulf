==================
Node Configuration
==================

The Node Configuration DB
=========================

As mentioned in the :doc:`Warewulf Configuration <configuration>`
section, node configs are persisted to the ``nodes.conf`` YAML file,
but generally it is best not to edit this file directly (however that
is supported, it is just prone to errors).

This method of using a YAML configuration file as a backend datastore
is both scalable and very lightweight. We've tested this out to over
10,000 node entries which yielded update latencies under 1 second,
which we felt was both tolerable and advantageous.

Adding a New Node
=================

Creating a new node is as simple as running the following command:

.. code-block:: console

   # wwctl node add n001 -I 10.0.2.1
   Added node: n001

Adding several nodes
--------------------

Several nodes can be added with a single command if a node range is
given. An additional IP address will incremented. So the command

.. code-block:: console

  # wwctl node add n00[2-4] -I 10.0.2.2
  Added node: n002
  Added node: n003
  Added node: n004

  # wwctl node list -n n00[1-4]
  NODE NAME              NAME     HWADDR             IPADDR          GATEWAY         DEVICE
  n001                   default  --                 10.0.2.1        --              (eth0)
  n002                   default  --                 10.0.2.2        --              (eth0)
  n003                   default  --                 10.0.2.3        --              (eth0)
  n004                   default  --                 10.0.2.4        --              (eth0)

has added 4 nodes with the incremented IP addresses.

Node Names
----------

For small clusters, you can use simple names (e.g. ``n0000``); but for
larger, more complicated clusters that are comprised of multiple
clusters and roles it is highly recommended to use node names that
include a cluster descriptor. In Warewulf, this is generally done by
using a domain name (e.g. ``n001.cluster01``). Warewulf will
automatically assume that the domain is the equivalent of the cluster
name.

This also means that you can address groups of nodes by the cluster
descriptor with globs. For example, you are able to refer to all nodes
in "cluster01" with the following string: ``*.cluster01`` which is
valuable for other ``wwctl`` commands.

Listing Nodes
=============

Once you have configured one or more nodes, you can list them and
their attributes as follows:

.. code-block:: console

  # wwctl node list
  NODE NAME              PROFILES                   NETWORK
  n001                                              default

You can also see the node's full attribute list by specifying the
``-a`` option (all):

.. code-block:: console

  # wwctl node list -a n001
  NODE                 FIELD              PROFILE      VALUE
  n001                 Id                 --           n001
  n001                 comment            default      This profile is automatically included for each node
  n001                 cluster            --           --
  n001                 container          default      sle-micro-5.3
  n001                 ipxe               --           (default)
  n001                 runtime            --           (hosts,ssh.authorized_keys,syncuser)
  n001                 wwinit             --           (wwinit,wwclient,fstab,hostname,ssh.host_keys,issue,resolv,udev.netname,systemd.netname,ifcfg,NetworkManager,debian.interfaces,wicked,ignition)
  n001                 root               --           (initramfs)
  n001                 discoverable       --           --
  n001                 init               --           (/sbin/init)
  n001                 asset              --           --
  n001                 kerneloverride     --           --
  n001                 kernelargs         --           (quiet crashkernel=no net.ifnames=1)
  n001                 ipmiaddr           --           --
  n001                 ipminetmask        --           --
  n001                 ipmiport           --           --
  n001                 ipmigateway        --           --
  n001                 ipmiuser           --           --
  n001                 ipmipass           --           --
  n001                 ipmiinterface      --           --
  n001                 ipmiwrite          --           --
  n001                 profile            --           default
  n001                 default:type       --           (ethernet)
  n001                 default:onboot     --           --
  n001                 default:netdev     --           (eth0)
  n001                 default:hwaddr     --           --
  n001                 default:ipaddr     --           172.16.1.11
  n001                 default:ipaddr6    --           --
  n001                 default:netmask    --           (255.255.255.0)
  n001                 default:gateway    --           --
  n001                 default:mtu        --           --
  n001                 default:primary    --           true

.. note::

   The attribute values in parenthesis are default values and can be
   overridden in the next section, granted, the default values are
   generally usable.

Setting Node Attributes
=======================

In the above output we can see that there is no kernel or container
defined for this node. To provision a node, the minimum requirements
are a kernel and container, and for that node to be useful, we will
also need to configure the network so the nodes are reachable after
they boot.

Node configurations are set using the ``wwctl node set`` command. To
see a list of all configuration attributes, use the command ``wwctl
node set --help``.

Configuring the Node's Container Image
======================================

.. code-block:: console

   # wwctl node set --container rocky-8 n001
   Are you sure you want to modify 1 nodes(s): y

And you can check that the container name is set for ``n001``:

.. code-block:: console

   # wwctl node list -a  n001 | grep Container
   n0000                Container          --           rocky-8

Configuring the Node's Network
------------------------------

To configure the network, we have to pick a network device name and
provide the network information as follows:

.. code-block:: console

  # wwctl node set --netdev eno1 --hwaddr 11:22:33:44:55:66 --ipaddr 10.0.2.1 --netmask 255.255.252.0 n001
   Are you sure you want to modify 1 nodes(s): y

You can now see that the node contains configuration attributes for
container, kernel, and network:

.. code-block:: console

  # wwctl node list -a n001
  NODE                 FIELD              PROFILE      VALUE
  n001                 Id                 --           n001
  n001                 comment            default      This profile is automatically included for each node
  n001                 cluster            --           --
  n001                 container          default      sle-micro-5.3
  n001                 ipxe               --           (default)
  n001                 runtime            --           (hosts,ssh.authorized_keys,syncuser)
  n001                 wwinit             --           (wwinit,wwclient,fstab,hostname,ssh.host_keys,issue,resolv,udev.netname,systemd.netname,ifcfg,NetworkManager,debian.interfaces,wicked,ignition)
  n001                 root               --           (initramfs)
  n001                 discoverable       --           --
  n001                 init               --           (/sbin/init)
  n001                 asset              --           --
  n001                 kerneloverride     --           tw
  n001                 kernelargs         --           (quiet crashkernel=no net.ifnames=1)
  n001                 ipmiaddr           --           --
  n001                 ipminetmask        --           --
  n001                 ipmiport           --           --
  n001                 ipmigateway        --           --
  n001                 ipmiuser           --           --
  n001                 ipmipass           --           --
  n001                 ipmiinterface      --           --
  n001                 ipmiwrite          --           --
  n001                 profile            --           default
  n001                 default:type       --           (ethernet)
  n001                 default:onboot     --           --
  n001                 default:netdev     --           eno1
  n001                 default:hwaddr     --           11:22:33:44:55:66
  n001                 default:ipaddr     --           10.0.2.1
  n001                 default:ipaddr6    --           --
  n001                 default:netmask    --           255.255.252.0
  n001                 default:gateway    --           --
  n001                 default:mtu        --           --
  n001                 default:primary    --           true

  # wwctl node set --cluster cluster01 n001
  Are you sure you want to modify 1 nodes(s): y

  # wwctl node list -a n001 | grep cluster
  n001                 cluster            --           cluster01

.. note::
  Due to the way network interface names are assigned by the Linux kernel and overwritten by udev
  and systemd in the default warewulf configuration, the use of `eth0/1/...` as interface names can lead to issues.
  We recommend the use of the original predictable names assigned to the interfaces (`eno1, ...`),
  as otherwise an interface may remain unconfigured if its name conflicts with the name of an already existing interface during boot.

To configure a bonded (link aggregation) network interface the following commands can be used:

.. code-block:: console

  # wwctl node set --netname=bond0_member_1 --netdev=eth2 --type=bond-slave n001
  # wwctl node set --netname=bond0_member_2 --netdev=eth3 --type=bond-slave n001
  # wwctl node set --netname=bond0 --netdev=bond0 --onboot=true --type=bond --ipaddr 10.0.3.1 --netmask=255.255.255.0 --mtu=9000 n001

Note: the netnames of the member interterfaces need to match the "netname" of the bonded interface until the first "_" (in the example bond0)


Additional networks
-------------------

Additional networks for the node can also be configured.
You will have provide all the necessary network information.

.. code-block:: shell

   wwctl node set \
     --netdev ib0 \
     --hwaddr aa:bb:cc:dd:ee:ff \
     --ipaddr 10.0.20.1 \
     --netmask 255.255.252.0 \
     --netname iband \
     --type infiniband \
     n001

VLAN
----

You can set the type also to ``vlan``.

Some network configuration systems use the network device name
(e.g., of the form ``eno1.100``)
to configure VLANs.
Other network systems need additional network tags:

- ``vlan_id``: configures the VLAN ID of the interface
- ``parent_device``: configures which physical interface to use

.. code-block:: shell

   wwctl node set \
     --netdev vlan42 \
     --ipaddr 10.0.42.1 \
     --netmask 255.255.252.0 \
     --netname iband \
     --type vlan \
     --nettagadd "vlan_id=42,parent_device=eth0" \
     n001

Static Routes
-------------

The included Warewulf network overlays support the configuration of static routes
using a network tag of the form ``route<N>=<dest>,<gateway>``.

.. code-block:: shell

   wwctl node set n001 \
     --nettagadd "route1=192.168.2.0/24,192.168.1.254"

Node Discovery
--------------

The hwaddr of a node can be automatically discovered by setting
``--discoverable`` on a node. If a node attempts to provision against
Warewulf using an interface that is unknown to Warewulf, that address
is associated with the first discoverable node. (Multiple discoverable
nodes are sorted lexically, first by cluster, then by ID.)

Once a node has been discovered its "discoverable" flag is
automatically cleared.

Un-setting Node Attributes
==========================

If you wish to ``unset`` a particular value, set the value to
``UNDEF``. For example:

And to unset this configuration attribute:

.. code-block:: console

   # wwctl node set --cluster UNDEF n001
   Are you sure you want to modify 1 nodes(s): y

   # wwctl node list -a n001 | grep Cluster
   n001                Cluster            --           --

Resources
=========

Warewulf nodes (and profiles) support generic "resources" that may hold arbitrarily complex YAML
data. This data, along with tags, may be used by both distribution and site overlays.

.. code-block:: yaml

   nodeprofiles:
     default:
       resources:
         fstab:
           - spec: warewulf:/home
             file: /home
             vfstype: nfs
             mntops: defaults
             freq: 0
             passno: 0
           - spec: warewulf:/opt
             file: /opt
             vfstype: nfs
             mntops: defaults,ro
             freq: 0
             passno: 0

Due to the arbitrary nature of generic resource data, it can only be managed with `wwctl
<node|profile> edit`.
