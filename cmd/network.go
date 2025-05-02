package cmd

import (
	"fmt"
	"ngomap/scanners"

	"github.com/spf13/cobra"
)

var threadsPerHost int

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network ip mask protocol",
	Short: "Scan all hosts in a network for open ports",
	Long:  `Scan a single host for any open ports using a specified protocol. For example: ngomap network 192.168.1.0 24 tcp`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		mask := (args[1])
		protocol := args[2]
		var ScanResults []scanners.ScanResult
		if threads != 100 || timeout != 1000 {
			fmt.Printf("Using custom parameters! Threads:%d Timeout:%d\n", threadsPerHost, timeout)
		}
		ScanResults = scanners.ScanNetwork(address, mask, protocol, timeout, amount, threadsPerHost)
		for _, result := range ScanResults {
			for _, port := range result.Ports {
				ipStr := scanners.Uint32ToIP(result.IP).String()
				fmt.Printf("%s:%d\n", ipStr, port)
			}
		}
	},
}

func init() {
	networkCmd.Flags().IntVar(&threadsPerHost, "threads-per-host", defaultThreads, "Number of goroutines to use per host")
	networkCmd.Flags().IntVar(&amount, "amount", defaultAmount, "The amount of hosts to scan simultaneously")
	rootCmd.AddCommand(networkCmd)
}
