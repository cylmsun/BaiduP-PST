package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 检查网络
// token可用性
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seems sync OK...")
		fmt.Printf("syncCmd获得参数:%s", Folder)
	},
}

func initSync() {
	rootCmd.AddCommand(syncCmd)
}
