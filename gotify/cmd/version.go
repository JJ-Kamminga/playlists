package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of my app's",
	Long:  `All software has versions. This is my app's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("App version none")
	},
}
