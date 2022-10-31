package cmd

import (
	"BaiduP-PST/pkg/model"
	"BaiduP-PST/util"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

// check folder
// check http connection
var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "check folder",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//clouds := checkCloudFolder()
		//locals := checkLocalFolder()
	},
}

func initFolderCheck() {
	rootCmd.AddCommand(folderCmd)
}

// check local folder,exists and detail
func checkLocalFolder() (infos []model.DicInfo) {
	infos = util.ReadDirssss("D:\\temp")
	for _, a := range infos {
		marshal, err := json.Marshal(a)
		if err != nil {
			fmt.Printf("checkCloudFolder errpr:%s", err.Error())
			continue
		}
		fmt.Println(string(marshal))
	}
	return
}

// check cloud folder,exists and detail
func checkCloudFolder() (infos []model.DicInfo) {
	m, err := util.GetCloudFolderDetails("/")
	if err != nil {
		fmt.Printf("checkCloudFolder errpr:%s", err.Error())
	}

	return
}

// logic
// three types: local new,cloud new, conflict
func compare(clouds []model.NetDicInfo, locals []model.DicInfo) (ans map[string]any) {

	return
}
