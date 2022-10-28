package cmd

import (
	"BaiduP-PST/config"
	"BaiduP-PST/pkg/driver"
	"BaiduP-PST/util"
	"fmt"
	"github.com/spf13/cobra"
)

// check folder
// check http connection
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check status",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := util.CheckDir(&Folder)
		if err != nil {
			fmt.Printf("check error:%s", err.Error())
		} else {
			fmt.Println("dic check OK...")
		}

		if util.CheckTime(&config.Setting.LastRefresh, config.Setting.ExpiresIn) {
			fmt.Println("prepare refresh token...")
			driver.RefreshToken()
		}
		fmt.Println("prepare check token...")
		driver.CheckToken()
	},
}

func initCheck() {
	rootCmd.AddCommand(checkCmd)
}
