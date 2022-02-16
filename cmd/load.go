package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gokins/util"
	"strings"
)

var configCmd = &cobra.Command{
	Use:   "load",
	Short: "加载配置文件",
	Long: `通过http(s)协议加载远程配置文件，使用样例：
gokins load https://cdn.jsdelivr.net/gh/hxulin/gokins/config.yaml
加载的配置文件会保存到本地：$Home/.gokins/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(`该命令仅接收一个参数，为配置文件的url地址，如：
gokins load https://cdn.jsdelivr.net/gh/hxulin/gokins/config.yaml`)
			return
		}
		url := args[0]
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			fmt.Println("配置文件地址格式错误")
			return
		}
		if util.LoadConfig(url) != nil {
			fmt.Println("配置文件加载失败")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
