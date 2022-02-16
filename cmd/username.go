package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gokins/util"
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
		authConfig, err := util.ReadAuthInfo()
		if err != nil {
			fmt.Println("用户名设置失败")
			return
		}
		cipherText, err := util.Encrypt(args[0])
		if err != nil {
			fmt.Println("用户名加密失败")
			return
		}
		authConfig.Auth.Username = cipherText
		err = util.SaveAuthInfo(authConfig)
		if err != nil {
			fmt.Println("用户名设置失败")
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
