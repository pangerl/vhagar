// Package tenant @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package tenant

import (
	"github.com/olivere/elastic/v7"
	"vhagar/libs"
)

// 每日巡检版本
var version = "v4.6"

type Inspect struct {
	ProjectName string
	ProxyURL    string
	Rocketmq    libs.Rocketmq
	Notifier    map[string]Notifier
	Tenant      *Tenant
}

type Tenant struct {
	Corp     []*Corp
	ESClient *elastic.Client
	PGClient *libs.PGClient
}

type Corp struct {
	Corpid               string `json:"corpid"`
	Convenabled          bool   `json:"convenabled"`
	CorpName             string `json:"corpName"`
	MessageNum           int64  `json:"messageNum"`
	UserNum              int    `json:"userNum"`
	CustomerNum          int64  `json:"customerNum"`
	CustomerGroupNum     int    `json:"customerGroupNum"`
	CustomerGroupUserNum int    `json:"customerGroupUserNum"`
	DauNum               int64  `json:"dauNum"`
	WauNum               int64  `json:"wauNum"`
	MauNum               int64  `json:"mauNum"`
}

type Notifier struct {
	Robotkey []string `json:"robotkey"`
	Userlist []string `json:"userlist"`
	IsPush   bool     `json:"ispush"`
}
