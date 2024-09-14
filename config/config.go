// Package config @Author lanpang
// @Date 2024/9/11 上午11:15:00
// @Desc
package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"time"
	"vhagar/libs"
	//"vhagar/task/nacos"
)

var (
	Config *CfgType
)

type CfgType struct {
	Global
	CfgFile         string
	VictoriaMetrics string             `toml:"victoriaMetrics"`
	Cron            map[string]crontab `toml:"cron"`
	Nacos           NacosCfg           `toml:"nacos"`
	Tenant          Tenant             `toml:"tenant"`
	PG              libs.DB            `toml:"pg"`
	ES              libs.DB            `toml:"es"`
	Doris           DorisCfg           `toml:"doris"`
	RocketMQ        RocketMQCfg        `toml:"rocketmq"`
	//Metric   metric.Config    `toml:"metric"`
}

type Global struct {
	ProjectName string              `toml:"projectname"`
	ProxyURL    string              `toml:"proxyurl"`
	Notifier    map[string]Notifier `toml:"notifier"`
	Watch       bool
	Report      bool
	Interval    time.Duration
}

type crontab struct {
	Crontab    bool   `toml:"crontab"`
	Scheducron string `toml:"scheducron"`
}

type Notifier struct {
	Robotkey []string `json:"robotkey"`
	Userlist []string `json:"userlist"`
	IsPush   bool     `json:"ispush"`
}

func InitConfig(cfgFile string) (*CfgType, error) {
	//configFile := path.Join(configDir, "config.toml")
	Config = &CfgType{
		CfgFile: cfgFile,
	}

	log.Printf("读取配置文件 %s \n", cfgFile)
	defer func() {
		if err := recover(); err != nil {
			//log.Fatalf("Failed Info: 配置文件格式错误 %s", err)
			log.Println("Recovered from panic:", err)
			return
		}
	}()
	if _, err := os.Stat(cfgFile); err != nil {
		if os.IsNotExist(err) {
			//log.Fatalf("读取配置文件 %s 失败，报错：%s", cfgFile, err)
			return nil, fmt.Errorf("configuration file(%s) not found", cfgFile)
		}
	} else {
		if _, err := toml.DecodeFile(cfgFile, Config); err != nil {
			//log.Fatalf("Failed Info: 配置文件格式错误 %s", err)
			return nil, fmt.Errorf("failed to load configs of dir: %s err:%s", cfgFile, err)
		}
		//log.Println(Config.VictoriaMetrics)
	}
	return Config, nil
}

//type Instances interface {
//	TableRender()
//}
//
//type Creator interface {
//	factoryMethod() Instances
//}

//func NewInspect(cfg *CfgType) *inspect.Inspect {
//	log.Println("初始化 Inspect 对象")
//
//	Inspect := &inspect.Inspect{
//		ProjectName: cfg.ProjectName,
//		ProxyURL:    cfg.ProxyURL,
//		Notifier:    cfg.Notifier,
//	}
//}