package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gokins/util"
)

var initKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "初始化加密配置项的密钥",
	Long: `该命令会为当前系统登录用户生成一对密钥文件：
$HOME/.gokins/public.bin
$HOME/.gokins/private.bin
敏感的配置信息项可以使用此密钥加密。`,
	Run: func(cmd *cobra.Command, args []string) {
		keyFileIsExist, err := util.KeyFileIsExist()
		if err != nil {
			fmt.Println("密钥文件存储位置初始化失败")
			return
		}
		if keyFileIsExist {
			fmt.Println("密钥文件已存在")
		} else {
			util.CreateSecretKey()
		}
	},
}

func init() {
	initCmd.AddCommand(initKeyCmd)
}
