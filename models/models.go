package models

import "sync"

type Req struct {
	Links []string `json:"links"`
}

type RespData struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}

type PDFReq struct {
	LinksNum []int `json:"links_num"`
}

var LinksStorage = make(map[int]map[string]string)

var Mu sync.Mutex
