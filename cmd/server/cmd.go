package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "greeter",
	Short: "greeter is a demonstration of grpc service",
	Long:  "greeter is a full featured grpc service which has good code style",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.json)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show verbose output")
}

func main() {
	rootCmd.Execute()
}
