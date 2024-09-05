// Package inspect @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package inspect

import (
	"database/sql"
	"github.com/jackc/pgx/v5"
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
	Doris       *Doris
}

type Tenant struct {
	Corp           []*Corp
	ESClient       *elastic.Client
	PGClient       *libs.PGClient
	CustomerClient *pgx.Conn
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

type Doris struct {
	DorisCfg
	MysqlClient        *sql.DB
	FailedJobs         []string
	StaffCount         int
	UseAnalyseCount    int
	CustomerGroupCount int
	OnlineBackendNum   int
	TotalBackendNum    int
}

type DorisCfg struct {
	libs.DB
	HttpPort int `toml:"httpport"`
}

type Notifier struct {
	Robotkey []string `json:"robotkey"`
	Userlist []string `json:"userlist"`
	IsPush   bool     `json:"ispush"`
}

type dorisResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		OnlineBackendNum int `json:"online_backend_num"`
		TotalBackendNum  int `json:"total_backend_num"`
	} `json:"data"`
	Count int `json:"count"`
}
