package cmd

import (
	"BaiduP-PST/config"
	"BaiduP-PST/pkg/base"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var Folder string

var rootCmd = &cobra.Command{
	Use:   "bdptool",
	Short: "A Baidu NetDisk Tool",
	Long: `This is a Baidu NetDisk Tool for personal using.
Developing...
...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use bdptool -h or --help for help.")
		fmt.Printf("rootCmd获得参数:%s \n", Folder)
		fmt.Println("看下config")
		fmt.Printf("config:%s,%s,%s", config.Setting.SecretKey, config.Setting.AppKey, config.Setting.AppId)
	},
}

func initCommands() {
	initTestToken()
	rootCmd.PersistentFlags().StringVarP(&Folder, "dic", "d", "D:\\temp\\test", "sync dir")
	initCheck()
	initSync()
}

func Execute() {
	initCommands()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//tick := time.Tick(time.Second * 1)
	//i := 0
	//for range tick {
	//	i++
	//	fmt.Printf("定时执行:%d \n", i)
	//}
}

func initTestToken() {
	base.TokenBody.AccessToken = config.Setting.AccessToken
	base.TokenBody.RefreshToken = config.Setting.RefreshToken
	base.TokenBody.ExpiresIn = config.Setting.ExpiresIn
}
