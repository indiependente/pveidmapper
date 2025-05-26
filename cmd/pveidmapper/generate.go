package main

import (
	"fmt"
	"sort"

	"github.com/indiependente/pveidmapper/internal/mapper"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ID mappings for a PVE LXC container",
	Long: `Generate the necessary ID mappings for a Proxmox VE LXC container.
The command takes one or more ID mappings in the format:
  containeruid[:containergid][=hostuid[:hostgid]]

Examples:
  pveidmapper generate -i 1000=1000
  pveidmapper generate -i 1000:1000=1000:1000
  pveidmapper generate -i 1000=1000 -i 1001=1001`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ids, err := cmd.Flags().GetStringArray("id")
		if err != nil {
			return err
		}

		if len(ids) == 0 {
			return fmt.Errorf("no IDs provided. Use -i to specify mappings")
		}

		var uidList, gidList [][2]int
		for _, id := range ids {
			if err := mapper.ValidateMappingString(id); err != nil {
				return fmt.Errorf("invalid mapping format: %w", err)
			}

			mapping, err := mapper.ValidateInput(id)
			if err != nil {
				return err
			}

			if mapping.ContainerUID != -1 {
				uidList = append(uidList, [2]int{mapping.ContainerUID, mapping.HostUID})
			}
			if mapping.ContainerGID != -1 {
				gidList = append(gidList, [2]int{mapping.ContainerGID, mapping.HostGID})
			}
		}

		sort.Slice(uidList, func(i, j int) bool { return uidList[i][0] < uidList[j][0] })
		sort.Slice(gidList, func(i, j int) bool { return gidList[i][0] < gidList[j][0] })

		uidMap := mapper.CreateMap("u", uidList)
		gidMap := mapper.CreateMap("g", gidList)

		fmt.Println("# Add to /etc/pve/lxc/<container_id>.conf:")
		for _, line := range uidMap {
			fmt.Println(line)
		}
		for _, line := range gidMap {
			fmt.Println(line)
		}

		fmt.Println("\n# Add to /etc/subuid:")
		for _, uid := range uidList {
			fmt.Printf("root:%d:1\n", uid[1])
		}

		fmt.Println("\n# Add to /etc/subgid:")
		for _, gid := range gidList {
			fmt.Printf("root:%d:1\n", gid[1])
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringArrayP("id", "i", []string{}, "containeruid[:containergid][=hostuid[:hostgid]]")
	generateCmd.MarkFlagRequired("id")
}
