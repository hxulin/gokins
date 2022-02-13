package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gokins/model"
	"gopkg.in/yaml.v2"
)

var setCmd = &cobra.Command{
	Use:   "username",
	Short: "设置Jenkins的用户名",
	Long: `设置你的Jenkins用户名，如用户名为hxl，可使用：
gokins username hxl`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(`该命令接收Jenkins的用户名参数，如用户名为hxl，可使用：
gokins username hxl`)
			return
		}

		authConfig := model.AuthConfig{Auth: model.Auth{Username: "xx", Token: "xx"}}
		marshal, _ := yaml.Marshal(authConfig)
		fmt.Println(string(marshal))
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
