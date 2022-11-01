package model

type NetDicInfoResp struct {
	Errno  int       `json:"errno"`
	List   []DicInfo `json:"list"`
	Cursor int       `json:"cursor"`
}

//type NetDicInfo struct {
//	ServerMtime    string `json:"server_mtime"`
//	ServerFilename string `json:"server_filename"`
//	Path           string `json:"path"`
//	IsDir          int8   `json:"isdir"`
//}
