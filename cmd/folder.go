package cmd

import (
	"BaiduP-PST/config"
	"BaiduP-PST/pkg/model"
	"BaiduP-PST/util"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"sync"
)

var wg sync.WaitGroup

// check folder
// check http connection
var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "compare folder and files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Prepare seeking different local and cloud folder and file")

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
	},
}

func initFolderCheck() {
	rootCmd.AddCommand(folderCmd)
}

// check local folder,exists and detail
func checkLocalFolder() (infos []model.DicInfo) {
	infos = util.GetLocalFolderDetails(config.Setting.DefaultFolder)
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
	infos, err := util.GetCloudFolderDetails("/")
	if err != nil {
		fmt.Printf("checkCloudFolder errpr:%s", err.Error())
	}

	return
}

// logic
// three types: local new,cloud new, conflict
func compareDetail(clouds []model.DicInfo, locals []model.DicInfo, c chan<- map[string][]model.DicInfo) {

	// todo maybe there is a better way to acc
	go getLocalNew(clouds, locals, c)
	go getCloudNew(clouds, locals, c)
	go getConflict(clouds, locals, c)

	return
}

func getLocalNew(clouds []model.DicInfo, locals []model.DicInfo, c chan<- map[string][]model.DicInfo) {
	var list []model.DicInfo
	for _, local := range locals {
		i := 0
		for _, cloud := range clouds {
			if local.Path == cloud.Path {
				i++
				break
			}
		}
		if i == 0 {
			list = append(list, model.DicInfo{
				Path:    local.Path,
				Name:    local.Name,
				ModTime: local.ModTime,
				IsDir:   local.IsDir,
				RWMode:  local.RWMode,
			})
		}
	}
	m := make(map[string][]model.DicInfo)
	m["LocalNew"] = list
	c <- m
	defer wg.Done()
}

func getCloudNew(clouds []model.DicInfo, locals []model.DicInfo, c chan<- map[string][]model.DicInfo) {
	var list []model.DicInfo
	for _, cloud := range clouds {
		i := 0
		for _, local := range locals {
			if local.Path == cloud.Path {
				i++
				break
			}
		}
		if i == 0 {
			list = append(list, model.DicInfo{
				Path:    cloud.Path,
				Name:    cloud.Name,
				ModTime: cloud.ModTime,
				IsDir:   cloud.IsDir,
			})
		}
	}
	m := make(map[string][]model.DicInfo)
	m["CloudNew"] = list
	c <- m
	defer wg.Done()
}

func getConflict(clouds []model.DicInfo, locals []model.DicInfo, c chan<- map[string][]model.DicInfo) {
	var list []model.DicInfo
	for _, cloud := range clouds {
		for _, local := range locals {
			if local.Path == cloud.Path {
				list = append(list, model.DicInfo{
					Path:    cloud.Path,
					Name:    cloud.Name,
					ModTime: cloud.ModTime,
					IsDir:   cloud.IsDir,
					RWMode:  local.RWMode,
				})
				break
			}
		}
	}
	m := make(map[string][]model.DicInfo)
	m["Conflict"] = list
	c <- m
	defer wg.Done()
}
