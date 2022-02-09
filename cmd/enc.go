package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gokins/rsa"
	"gokins/util"
)

var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "字符串加密",
	Long:  "使用密钥 ($HOME/.gokins/public.bin)，加密字符串信息。",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			publicKey, err := util.ReadPublicKey()
			if err != nil {
				fmt.Println("密钥文件读取失败，请尝试使用：gokins init key 初始化密钥")
				return
			}
			for _, arg := range args {
				cipherText, err := rsa.Encrypt(publicKey, arg)
				if err != nil {
					fmt.Println("内容加密失败：", err)
				} else {
					fmt.Println("ENC(" + cipherText + ")")
				}
			}
		} else {
			fmt.Println(`enc 命令后面需要指定待加密的字符串，例如：
  gokins enc foo`)
		}
	},
}

func init() {
	rootCmd.AddCommand(encCmd)
}
