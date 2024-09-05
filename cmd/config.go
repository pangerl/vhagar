// Package cmd @Author lanpang
// @Date 2024/8/1 下午2:47:00
// @Desc
package cmd

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"vhagar/inspect"
	"vhagar/libs"
	"vhagar/metric"
	"vhagar/nacos"
)

var (
	CONFIG = &Config{
		ProjectName: "测试项目",
	}
	cfgFile string
)

type Config struct {
	ProjectName string                      `toml:"projectname"`
	ProxyURL    string                      `toml:"proxyurl"`
	Cron        map[string]crontab          `toml:"cron"`
	Notifier    map[string]inspect.Notifier `toml:"notifier"`
	Nacos       nacos.Config                `toml:"nacos"`
	Tenant      inspect.Tenant              `toml:"tenant"`
	PG          libs.DB                     `toml:"pg"`
	ES          libs.DB                     `toml:"es"`
	Doris       inspect.DorisCfg            `toml:"doris"`
	Rocketmq    libs.Rocketmq               `toml:"rocketmq"`
	Metric      metric.Config               `toml:"metric"`
	Customer    libs.DB                     `toml:"customer"`
}
type dorisCfg struct {
	libs.DB
	HttpPort int `toml:"httpport"`
}
type crontab struct {
	Crontab    bool   `toml:"crontab"`
	Scheducron string `toml:"scheducron"`
}

func createTempConfig() {
	log.Printf("创建配置文件模板")
	config := &Config{
		Cron: map[string]crontab{
			"tenant": {Crontab: false, Scheducron: "30 09 * * *"},
		},
		Notifier: map[string]inspect.Notifier{
			"tenant": {Robotkey: []string{"xxx"}, Userlist: []string{}},
		},
		Tenant: inspect.Tenant{
			Corp: []*inspect.Corp{
				{Corpid: "xxx", Convenabled: false},
			},
		},
	}
	// 创建并打开文件
	file, err := os.Create(cfgFile)
	if err != nil {
		log.Fatalf("Error creating config file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	// 序列化结构体到 TOML 格式并写入文件
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		log.Fatalf("Error encoding config to TOML: %v", err)
	}
	log.Printf("配置文件创建成功： %s", cfgFile)

}
