package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var initConfCmd = &cobra.Command{
	Use:   "conf",
	Short: "初始化配置信息",
	Long: `该命令会为当前系统登录用户生成默认配置文件：
$HOME/.gokins/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("initConf called")
	},
}

func init() {
	initCmd.AddCommand(initConfCmd)
}
