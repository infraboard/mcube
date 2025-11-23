package cmd

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "应用配置工具",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var generateRandomEncryptKeyCmd = &cobra.Command{
	Use:   "generate-random-encrypt-key",
	Short: "生成随机密钥",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("应用加密密码: %s\n", application.Get().GenerateRandomEncryptKey())
	},
}

var keyInfoCmd = &cobra.Command{
	Use:   "key-info",
	Short: "当前应用配置的密钥信息",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := application.Get().GetKeyInfo()
		cobra.CheckErr(err)
		// 格式化为key: value进行打印
		fmt.Println("==== 当前应用配置的密钥信息 ====")
		for k, v := range info {
			fmt.Printf("%s: %v\n", k, v)
		}
	},
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "加密",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("请输入需要加密的字符串")
			return
		}

		planText := args[0]
		cipherText, err := application.Get().EncryptString(planText)
		cobra.CheckErr(err)
		fmt.Printf("加密完成: 明文[%s] -> 密文[%s]\n", planText, cipherText)

	},
}

func init() {
	appCmd.AddCommand(generateRandomEncryptKeyCmd)
	appCmd.AddCommand(keyInfoCmd)
	appCmd.AddCommand(encryptCmd)
	Root.AddCommand(appCmd)
}
