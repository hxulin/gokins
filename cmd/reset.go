package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "重置密钥或配置信息",
	//Long: "重置密钥或配置信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reset called")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
