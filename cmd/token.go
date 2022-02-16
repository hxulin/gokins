package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gokins/util"
)

var resetCmd = &cobra.Command{
	Use:   "token",
	Short: "设置Jenkins的用户token信息",
	Long: `设置Jenkins的用户token信息，使用样例：
gokins token 1267423623afd221271ce93319ba2e98b0`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(`该命令接收一个参数，是Jenkins的用户token，使用样例：
gokins token 1267423623afd221271ce93319ba2e98b0`)
			return
		}
		authConfig, err := util.ReadAuthInfo()
		if err != nil {
			fmt.Println("用户token设置失败")
			return
		}
		cipherText, err := util.Encrypt(args[0])
		if err != nil {
			fmt.Println("token加密失败")
			return
		}
		authConfig.Auth.Token = cipherText
		err = util.SaveAuthInfo(authConfig)
		if err != nil {
			fmt.Println("用户token设置失败")
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
