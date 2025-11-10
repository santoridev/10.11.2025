package models

type Req struct {
	Links []string `json:"links"`
}

type Res struct {
	Resp RespData
}

type RespData struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}
type PDFReq struct {
	LinksNum []int `json:"links_num"`
}

var LinksStorage = make(map[int]map[string]string)
