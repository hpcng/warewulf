package list

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hpcng/warewulf/internal/pkg/node"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/hpcng/warewulf/pkg/hostlist"
	"github.com/spf13/cobra"
)

func CobraRunE(cmd *cobra.Command, args []string) error {
	var err error

	nodeDB, err := node.New()
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Could not open node configuration: %s\n", err)
		os.Exit(1)
	}

	nodes, err := nodeDB.FindAllNodes()
	if err != nil {
		wwlog.Printf(wwlog.ERROR, "Could not get node list: %s\n", err)
		os.Exit(1)
	}

	args = hostlist.Expand(args)

	if ShowAll {
		for _, node := range node.FilterByName(nodes, args) {
			fmt.Printf("################################################################################\n")
			fmt.Printf("%-20s %-18s %-12s %s\n", "NODE", "FIELD", "PROFILE", "VALUE")
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Id", node.Id.Source(), node.Id.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Comment", node.Comment.Source(), node.Comment.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Cluster", node.ClusterName.Source(), node.ClusterName.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Profiles", "--", strings.Join(node.Profiles, ","))

			fmt.Printf("%-20s %-18s %-12s %t\n", node.Id.Get(), "Discoverable", node.Discoverable.Source(), node.Discoverable.PrintB())

			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Container", node.ContainerName.Source(), node.ContainerName.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "KernelOverride", node.Kernel.Override.Source(), node.Kernel.Override.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "KernelArgs", node.Kernel.Args.Source(), node.Kernel.Args.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "SystemOverlay", node.SystemOverlay.Source(), node.SystemOverlay.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "RuntimeOverlay", node.RuntimeOverlay.Source(), node.RuntimeOverlay.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Ipxe", node.Ipxe.Source(), node.Ipxe.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Init", node.Init.Source(), node.Init.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Root", node.Root.Source(), node.Root.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "AssetKey", node.AssetKey.Source(), node.AssetKey.Print())

			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "IpmiIpaddr", node.Ipmi.Ipaddr.Source(), node.Ipmi.Ipaddr.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "IpmiNetmask", node.Ipmi.Netmask.Source(), node.Ipmi.Netmask.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "IpmiPort", node.Ipmi.Port.Source(), node.Ipmi.Port.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "IpmiGateway", node.Ipmi.Gateway.Source(), node.Ipmi.Gateway.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "IpmiUserName", node.Ipmi.UserName.Source(), node.Ipmi.UserName.Print())
			fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "IpmiInterface", node.Ipmi.Interface.Source(), node.Ipmi.Interface.Print())

			for keyname, key := range node.Tags {
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), "Tag["+keyname+"]", key.Source(), key.Print())
			}

			for name, netdev := range node.NetDevs {
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":DEVICE", netdev.Device.Source(), netdev.Device.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":HWADDR", netdev.Hwaddr.Source(), netdev.Hwaddr.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":IPADDR", netdev.Ipaddr.Source(), netdev.Ipaddr.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":IPADDR6", netdev.Ipaddr.Source(), netdev.Ipaddr6.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":NETMASK", netdev.Netmask.Source(), netdev.Netmask.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":GATEWAY", netdev.Gateway.Source(), netdev.Gateway.Print())
				fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":TYPE", netdev.Type.Source(), netdev.Type.Print())
				fmt.Printf("%-20s %-18s %-12s %t\n", node.Id.Get(), name+":ONBOOT", netdev.OnBoot.Source(), netdev.OnBoot.PrintB())
				fmt.Printf("%-20s %-18s %-12s %t\n", node.Id.Get(), name+":DEFAULT", netdev.Default.Source(), netdev.Default.PrintB())
				for keyname, key := range netdev.Tags {
					fmt.Printf("%-20s %-18s %-12s %s\n", node.Id.Get(), name+":TAG["+keyname+"]", key.Source(), key.Print())
				}
			}

		}

	} else if ShowNet {
		fmt.Printf("%-22s %-6s %-18s %-15s %-15s %-15s\n", "NODE NAME", "DEVICE", "HWADDR", "IPADDR", "GATEWAY", "IPADDR6")
		fmt.Println(strings.Repeat("=", 80))

		for _, node := range node.FilterByName(nodes, args) {
			if len(node.NetDevs) > 0 {
				for name, dev := range node.NetDevs {
					fmt.Printf("%-22s %-6s %-18s %-15s %-15s %-15s\n",
						node.Id.Get(), name, dev.Hwaddr.Print(), dev.Ipaddr.Print(), dev.Gateway.Print(), dev.Ipaddr6.Print())
				}
			} else {
				fmt.Printf("%-22s %-6s %-18s %-15s %-15s %-15s\n", node.Id.Get(), "--", "--", "--", "--", "--")
			}
		}

	} else if ShowIpmi {
		fmt.Printf("%-22s %-16s %-10s %-20s %-20s %-14s\n", "NODE NAME", "IPMI IPADDR", "IPMI PORT", "IPMI USERNAME", "IPMI PASSWORD", "IPMI INTERFACE")
		fmt.Println(strings.Repeat("=", 108))

		for _, node := range node.FilterByName(nodes, args) {
			fmt.Printf("%-22s %-16s %-10s %-20s %-20s %-14s\n", node.Id.Get(), node.Ipmi.Ipaddr.Print(), node.Ipmi.Port.Print(), node.Ipmi.UserName.Print(), node.Ipmi.Password.Print(), node.Ipmi.Interface.Print())
		}

	} else if ShowLong {
		fmt.Printf("%-22s %-26s %-35s %s\n", "NODE NAME", "KERNEL OVERRIDE", "CONTAINER", "OVERLAYS (S/R)")
		fmt.Println(strings.Repeat("=", 120))

		for _, node := range node.FilterByName(nodes, args) {
			fmt.Printf("%-22s %-26s %-35s %s\n", node.Id.Get(), node.Kernel.Override.Print(), node.ContainerName.Print(), node.SystemOverlay.Print()+"/"+node.RuntimeOverlay.Print())
		}

	} else {
		fmt.Printf("%-22s %-26s %s\n", "NODE NAME", "PROFILES", "NETWORK")
		fmt.Println(strings.Repeat("=", 80))

		for _, node := range node.FilterByName(nodes, args) {
			var netdevs []string
			if len(node.NetDevs) > 0 {
				for name, dev := range node.NetDevs {
					netdevs = append(netdevs, fmt.Sprintf("%s:%s", name, dev.Ipaddr.Print()))
				}
			}
			sort.Strings(netdevs)

			fmt.Printf("%-22s %-26s %s\n", node.Id.Get(), strings.Join(node.Profiles, ","), strings.Join(netdevs, ", "))
		}

	}

	return nil
}
