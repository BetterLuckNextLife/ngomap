package cmd

import (
	"fmt"
	"ngomap/scanners"

	"github.com/spf13/cobra"
)

// singleCmd represents the single command
var singleCmd = &cobra.Command{
	Use:   "single host protocol",
	Short: "Scan a single host",
	Long:  `Scan a single host for any open ports using a specified protocol. For example: ngomap single 127.0.0.1 tcp`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		protocol := args[1]
		var foundPorts []int
		if threads != 1000 || timeout != 1000 {
			fmt.Printf("Using custom parameters! Threads:%d Timeout:%d\n", threads, timeout)
			foundPorts = scanners.ScanHost(host, protocol, timeout, threads)
		} else {
			foundPorts = scanners.ScanHost(host, protocol, 1000, 100)
		}

		for _, port := range foundPorts {
			fmt.Printf("%s:%d\n", args[0], port)
		}
	},
}

func init() {
	rootCmd.AddCommand(singleCmd)
}
