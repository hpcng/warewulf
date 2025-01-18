=================================
Development Environment (Vagrant)
=================================

Create Rocky Linux 9 virtual machine for Warewulf testbed using Vagrant


Host system requirements
========================

#. CPU supports H/W virtualization.
#. KVM kernel module available and loaded.


CPU H/W Virtualization support
------------------------------

Check CPU virtualization capability using following command.
If your system has Intel CPU, you will see :code:`Intel VT` here, and if your system has AMD CPU, you will see :code:`AMD-V` here.

.. code-block:: bash

   lscpu | grep Virtualization
   Virtualization:                  AMD-V
   Virtualization type:             full


KVM kernel module
-----------------

.. code-block:: bash

   lsmod | grep kvm
   ccp                   118784  1 kvm_amd
   kvm                  1105920  1 kvm_amd
   irqbypass              16384  1 kvm


Setup development environment on Rocky Linux 9
==============================================

Install QEMU, libvirt
---------------------

.. code-block:: bash

    # Install packages
    sudo dnf install -y libvirt qemu-kvm \
        libguestfs virtio-win guestfs-tools libguestfs-inspect-icons virt-win-reg \
        virt-install virt-top

    # Enable and start libvirtd
    sudo systemctl enable --now libvirtd

    # Add user to libvirt group
    sudo usermod -aG libvirt rocky


Install Cockpit (Optional)
--------------------------

.. code-block:: bash

    # Install packages
    sudo dnf install -y cockpit cockpit-machines

    # Enable and start cockpit (http://localhost:9090)
    sudo systemctl enable --now cockpit.socket


Install Vagrant, vagrant-libvirt plug-in and vagrant-reload plug-in
-------------------------------------------------------------------

.. code-block:: bash

    sudo dnf config-manager --add-repo https://rpm.releases.hashicorp.com/RHEL/hashicorp.repo
    sudo dnf install -y vagrant

    sudo dnf group install -y "Development tools"
    sudo dnf config-manager --set-enabled crb
    sudo dnf install -y libvirt-devel

    vagrant plugin install vagrant-libvirt
    vagrant plugin install vagrant-reload


Vagrant box and Vagrantfile for Warewulf sandbox
================================================

Create Rocky Linux 9.2 vagrant box
----------------------------------

.. code-block:: bash

    cat << 'EOF' > box-metadata.json
    {
    "name" : "rockylinux/9",
    "description" : "Rocky Linux 9 2.0.0",
    "versions" : [
        {
        "version" : "2.0.0-20230513.0",
        "providers" : [
            {
            "name" : "libvirt",
            "url" : "https://dl.rockylinux.org/pub/rocky/9.2/images/x86_64/Rocky-9-Vagrant-Libvirt-9.2-20230513.0.x86_64.box"
            }
        ]
        }
    ]
    }
    EOF

    vagrant box add box-metadata.json


Vagrantfile
-----------

.. code-block:: bash

    mkdir -p ~/warewulf-sandbox
    cd ~/warewulf-sandbox

    cat << 'EOF' > Vagrantfile
    Vagrant.configure("2") do |config|
        number_of_node = ENV["NODES"] || 2
        branch = ENV["BRANCH"] || "v4.5.8"

        config.vm.define :head do |head|
            head.vm.box = "rockylinux/9"
            head.vm.box_version = "2.0.0-20230513.0"
            head.vm.hostname = "warewulf"

            head.vm.network "private_network",
                ip: "192.168.200.254",
                netmask: "255.255.255.0",
                libvirt__network_name: "pxe",
                libvirt__dhcp_enabled: false

            head.vm.synced_folder ".", "/vagrant", type: "nfs", nfs_version: 4, nfs_udp: false

            head.vm.provider :libvirt do |libvirt|
                libvirt.cpu_mode = "host-passthrough"
                libvirt.memory = '8192'
                libvirt.cpus = '2'
                libvirt.machine_virtual_size = 40
            end

            head.vm.provision "shell", inline: <<-SHELL
                dnf install -y cloud-utils-growpart
                growpart /dev/vda 5
                xfs_growfs /dev/vda5
            SHELL

            head.vm.provision "shell", inline: <<-SHELL
                dnf groupinstall -y "Development Tools"
                dnf install -y epel-release
                dnf config-manager --set-enabled crb
                dnf install -y golang tftp-server dhcp-server nfs-utils gpgme-devel libassuan-devel

                cd /tmp
                git clone https://github.com/warewulf/warewulf.git
                cd warewulf
                git checkout v4.5.8
                make clean defaults \
                    PREFIX=/usr \
                    BINDIR=/usr/bin \
                    SYSCONFDIR=/etc \
                    DATADIR=/usr/share \
                    LOCALSTATEDIR=/var/lib \
                    SHAREDSTATEDIR=/var/lib \
                    MANDIR=/usr/share/man \
                    INFODIR=/usr/share/info \
                    DOCDIR=/usr/share/doc \
                    SRVDIR=/var/lib \
                    TFTPDIR=/var/lib/tftpboot \
                    SYSTEMDDIR=/usr/lib/systemd/system \
                    BASHCOMPDIR=/etc/bash_completion.d/ \
                    FIREWALLDDIR=/usr/lib/firewalld/services \
                    WWCLIENTDIR=/warewulf
                make all
                make install

                systemctl disable --now firewalld

                sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config
            SHELL

            head.vm.provision "reload"

            head.vm.provision "shell", inline: <<-SHELL
                cat << 'CONF' | sudo tee /etc/warewulf/warewulf.conf
    ipaddr: 192.168.200.254
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
    CONF

                sed -i 's@ExecStart=/usr/bin/wwctl server@ExecStart=/usr/bin/wwctl server -d -v@' /usr/lib/systemd/system/warewulfd.service
                systemctl enable --now warewulfd

                wwctl configure --all

                wwctl container import docker://ghcr.io/warewulf/warewulf-rockylinux:9 rocky-9
                wwctl profile set --yes --container rocky-9 "default"
                wwctl profile set --yes --netdev eth1 --netmask 255.255.255.0 --gateway 192.168.200.254 "default"

                wwctl node add n0001.cluster -I 192.168.200.101 --discoverable true
                wwctl node add n0002.cluster -I 192.168.200.102 --discoverable true
            SHELL
        end

        (1..number_of_node).each do |i|
            config.vm.define :"n000#{i}", autostart: false do |node|
                node.vm.hostname = "n000#{i}"
                node.vm.network "private_network",
                libvirt__network_name: "pxe"

                node.vm.provider :libvirt do |compute|
                    compute.cpu_mode = 'host-passthrough'
                    compute.memory = '8192'
                    compute.cpus = '2'
                    boot_network = {'network' => 'pxe'}
                    compute.boot boot_network
                end
            end
        end
    end
    EOF


Spin up head node
=================

.. code-block:: bash

    vagrant up


Spin up compute nodes
=====================

.. code-block:: bash

    vagrant up n0001

    # Wait until n0001 becomes ready

    vagrant up n0002
