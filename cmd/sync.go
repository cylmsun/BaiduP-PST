package cmd

import (
	"BaiduP-PST/pkg/model"
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
		checkFolder()
	},
}

func initSync() {
	rootCmd.AddCommand(syncCmd)
}

func checkFolder() {
	clouds := checkCloudFolder()
	locals := checkLocalFolder()

	c := make(chan map[string][]model.DicInfo, 3)
	wg.Add(3)
	compareDetail(clouds, locals, c)
	wg.Wait()

	close(c)
	for v := range c {
		if _, ok := v["LocalNew"]; ok {
			fmt.Printf("LocalNew:%d \n", len(v["LocalNew"]))
		}
		if _, ok := v["CloudNew"]; ok {
			fmt.Printf("CloudNew:%d \n", len(v["CloudNew"]))
		}
		if _, ok := v["Conflict"]; ok {
			fmt.Printf("Conflict:%d \n", len(v["Conflict"]))
		}
	}
}
