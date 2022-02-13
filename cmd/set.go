package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "set",
	Short: "初始化密钥或配置信息",
	//Long: "初始化密钥或配置信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`您是否想输入：
  gokins set username [username]
  gokins set token [token]`)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
