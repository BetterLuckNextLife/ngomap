package cmd

import (
	"fmt"
	"ngomap/scanners"

	"github.com/spf13/cobra"
)

var threadsPerHost int

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	networkCmd.Flags().IntVar(&threadsPerHost, "threads-per-host", 100, "Number of goroutines to use per host")
	networkCmd.Flags().IntVar(&amount, "amount", 1, "The amount of hosts to scan simultaneously")
	rootCmd.AddCommand(networkCmd)
}
