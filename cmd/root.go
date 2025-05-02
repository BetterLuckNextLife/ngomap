package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var threads int
var timeout int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ngomap",
	Short: "My own attempt at writing an efficient network scanner in Go",
	//TODO: Add long description
	Long: `ADD DESCRIPTION HERE For example:

	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ngomap.yaml)")
	// TODO: make threads work

	rootCmd.PersistentFlags().IntVar(&threads, "threads", 1000, "The amout of threads (gorutines) to use")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 1000, "The amout of threads (gorutines) to use")
	//  Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
