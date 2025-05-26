package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pveidmapper",
	Short: "A tool to manage PVE ID mappings for LXC containers",
	Long: `pveidmapper is a CLI tool that helps manage UID/GID mappings for Proxmox VE LXC containers.
It generates the necessary configuration for both the container and the host system.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.pveidmapper.yaml)")
}

func main() {
	Execute()
}
