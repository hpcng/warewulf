=======
Dnsmasq
=======

Usage
=====

As experimental feature its possible to use ``dnsmasq`` instead of the ISC ``dhcpd`` server in combination with a ``TFTP`` server.
The ``dnsmasq`` service is then acting as ``dhcp`` and ``TFTP`` server.
In order to keep the file ``/etc/dnsmasq.d/ww4-hosts.conf`` is created and must be included in the main ``dnsmasq.conf`` via the ``conf-dir=/etc/dnsmasq.d`` option.


Installation
------------

Before the installation, make sure that ``dhcpd`` and ``tftp`` are disabled.
You can do that with the commands:

.. code-block:: shell

   systemctl disable dhcpd
   systemctl stop dhcpd
   systemctl disable tftp
   systemctl stop tftp

Now you can install ``dnsmasq``.

.. code-block:: shell

   zypper install dnsmasq

After the installation you have to instruct ``warewulf`` to use ``dnsmasq`` as its ``dhcpd`` and ``tftp`` service.
``dnsmasq`` has to be specified in the configuration file ``/etc/warewulf/warewulf.conf``.

.. code-block:: shell

   tftp:
     systemd name: dnsmasq
   dhcp:
     systemd name: dnsmasq

The configuration of ``dnsmasq`` doesn't need to be changed, as the default configuration includes all files with following pattern ``/etc/dnsmasq.d/*conf`` into its configuration.
This configuration is created by the overlay template ``host:/etc/dnsmasq.d/ww4-hosts.conf.ww``.

.. note::

   In certain OS versions, such as Rocky Linux 9, ``dnsmasq`` is configured to listen locally via the ``interface=lo`` option by default in ``/etc/dnsmasq.conf``. 
   Replace this with the interface associated with your Warewulf network, or remove/comment out the interface option entirely to enable listening on all interfaces.

In order to build this template run

.. code-block:: shell

   wwctl overlay build -H

After that the ``dnsmasq`` service has to be enabled.
Either

.. code-block:: shell

   systemctl enable --now dnsmasq

or by (re)configuring warewulf with

.. code-block:: shell

   wwctl configure dhcp
   wwctl configure tftp
