// Package rocketmq @Author lanpang
// @Date 2024/9/13 下午5:53:00
// @Desc
package rocketmq

import (
	"vhagar/config"

	"go.uber.org/zap"
)

const taskName = "rocketmq"

func NewRocketMQ(cfg *config.CfgType, logger *zap.SugaredLogger) *RocketMQ {
	return &RocketMQ{
		Config:      cfg,
		Logger:      logger,
		RocketMQCfg: cfg.RocketMQ,
		BrokerMap:   make(map[string]*BrokerDetail),
	}
}

type RocketMQ struct {
	Config      *config.CfgType
	Logger      *zap.SugaredLogger
	RocketMQCfg config.RocketMQCfg
	BrokerMap   map[string]*BrokerDetail
}

type BrokerDetail struct {
	name              string
	role              string
	version           string
	addr              string
	runTime           string
	useDisk           string
	todayProduceCount int
	todayConsumeCount int
}

type BrokerData struct {
	RunTime                 string `json:"runtime"`
	CommitLogDirCapacity    string `json:"commitLogDirCapacity"`
	BrokerVersionDesc       string `json:"brokerVersionDesc"`
	MsgPutTotalTodayNow     string `json:"msgPutTotalTodayNow"`
	MsgPutTotalTodayMorning string `json:"msgPutTotalTodayMorning"`
	MsgGetTotalTodayNow     string `json:"msgGetTotalTodayNow"`
	MsgGetTotalTodayMorning string `json:"msgGetTotalTodayMorning"`
}

type Broker struct {
	BrokerName  string            `json:"brokerName"`
	BrokerAddrs map[string]string `json:"brokerAddrs"`
}

type ClusterInfo struct {
	BrokerAddrTable map[string]Broker `json:"brokerAddrTable"`
}

type ClusterData struct {
	BrokerServer map[string]map[string]BrokerData `json:"brokerServer"`
	ClusterInfo  ClusterInfo                      `json:"clusterInfo"`
}

type ResponseData struct {
	Status int         `json:"status"`
	Data   ClusterData `json:"data"`
}
