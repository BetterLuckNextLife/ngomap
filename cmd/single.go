package cmd

import (
	"fmt"
	"ngomap/scanners"

	"github.com/spf13/cobra"
)

var singleThreads int

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
		if singleThreads != 100 || timeout != 1000 {
			foundPorts = scanners.ScanHost(host, protocol, timeout, singleThreads)
		} else {
			foundPorts = scanners.ScanHost(host, protocol, 1000, 100)
		}

		for _, port := range foundPorts {
			fmt.Printf("%s:%d\n", args[0], port)
		}
	},
}

func init() {
	singleCmd.Flags().IntVar(&singleThreads, "threads", 100, "Total goroutines to use")
	rootCmd.AddCommand(singleCmd)
}
