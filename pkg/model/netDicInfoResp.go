package model

type NetDicInfoResp struct {
	Errno  int       `json:"errno"`
	List   []DicInfo `json:"list"`
	Cursor int       `json:"cursor"`
}

type PreCreateResp struct {
	ErrNo      int    `json:"errno"`
	Path       string `json:"path"`
	UploadId   string `json:"uploadid"`
	ReturnType int    `json:"return_type"`
	BlockList  string `json:"block_list"`
}
