package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var threads int
var amount int
var timeout int

var defaultThreads int = 100
var defaultAmount = 1
var defaultTimeout int = 1000

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ngomap",
	Short: "A network scanner written in go, focusing on the simplicity",
	Long:  `Scan individual hosts or whole networks for open ports`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 1000, "The maximum time (in ms) to wait for a port response")
}
