package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/hpcng/warewulf/internal/pkg/node"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/spf13/cobra"
)

func CobraRunE(vars *variables) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		if len(args) > 0 && strings.Contains(args[0], ",") {
			args = strings.FieldsFunc(args[0], func(r rune) bool { return r == ',' })
		}
		req := wwapiv1.GetProfileList{
			ShowAll:     vars.showAll,
			ShowFullAll: vars.showFullAll,
			Profiles:    args,
		}
		profileInfo, err := apiprofile.ProfileList(&req)
		if err != nil {
			return
		}

	nodeDB, err := node.New()
	if err != nil {
		wwlog.Error("Could not open node configuration: %s", err)
		os.Exit(1)
	}

	profiles, err := nodeDB.FindAllProfiles()
	if err != nil {
		wwlog.Error("Could not find all nodes: %s", err)
		os.Exit(1)
	}

	if ShowAll {
		for _, profile := range node.FilterByName(profiles, args) {
			fmt.Printf("################################################################################\n")
			fmt.Printf("%-20s %-18s %s\n", "PROFILE NAME", "FIELD", "VALUE")
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "Id", profile.Id.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "Comment", profile.Comment.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "Cluster", profile.ClusterName.Print())

			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "Discoverable", profile.Discoverable.PrintB())

			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "Container", profile.ContainerName.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "KernelOverride", profile.Kernel.Override.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "KernelArgs", profile.Kernel.Args.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "Init", profile.Init.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "Root", profile.Root.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "AssetKey", profile.AssetKey.Print())

			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "SystemOverlay", profile.SystemOverlay.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "RuntimeOverlay", profile.RuntimeOverlay.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "Ipxe", profile.Ipxe.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "IpmiNetmask", profile.Ipmi.Netmask.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "IpmiPort", profile.Ipmi.Port.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "IpmiGateway", profile.Ipmi.Gateway.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "IpmiUserName", profile.Ipmi.UserName.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "IpmiInterface", profile.Ipmi.Interface.Print())
			fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "IpmiWrite", profile.Ipmi.Write.PrintB())

			for keyname, key := range profile.Tags {
				fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), "Tag["+keyname+"]", key.Print())
			}

			for name, netdev := range profile.NetDevs {
				fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), name+":IPADDR", netdev.Ipaddr.Print())
				fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), name+":NETMASK", netdev.Netmask.Print())
				fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), name+":GATEWAY", netdev.Gateway.Print())
				fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), name+":HWADDR", netdev.Hwaddr.Print())
				fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), name+":TYPE", netdev.Hwaddr.Print())
				fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), name+":ONBOOT", netdev.OnBoot.PrintB())
				fmt.Printf("%-20s %-18s %s\n", profile.Id.Get(), name+":PRIMARY", netdev.Primary.PrintB())
				for keyname, key := range netdev.Tags {
					fmt.Printf("%-20s %-18s %-12s %s\n", profile.Id.Get(), name+":TAG["+keyname+"]", key.Source(), key.Print())
				}
			}
		}
	} else {
		fmt.Printf("%-20s %s\n", "PROFILE NAME", "COMMENT/DESCRIPTION")
		fmt.Printf(strings.Repeat("=", 80) + "\n")

		for _, profile := range node.FilterByName(profiles, args) {
			fmt.Printf("%-20s %s\n", profile.Id.Print(), profile.Comment.Print())
		}
	}

	return nil
}
