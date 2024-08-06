// Package inspect @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package inspect

import (
	"github.com/olivere/elastic/v7"
)

type Tenant struct {
	Corp     []Corp
	ESclient *elastic.Client
}

type Corp struct {
	Corpid      string
	Convenabled bool
	CorpName    string
	MessageNum  int
	UserNum     int
	CustomerNum int
	DauNum      int
	WauNum      int
	MauNum      int
}

type Db struct {
	Ip       string
	Port     int
	Username string
	Password string
	Ssl      bool
}