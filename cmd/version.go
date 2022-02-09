package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看版本信息",
	Long:  "查看版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v2022.02.09")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
