package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "token",
	Short: "设置Jenkins用户的token信息",
	Long: `设置Jenkins用户的token信息，如token为，可使用：
gokins token hxl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("token called")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
