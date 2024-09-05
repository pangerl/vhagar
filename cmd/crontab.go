// Package cmd @Author lanpang
// @Date 2024/8/9 下午4:25:00
// @Desc
package cmd

import (
	"context"
	"log"
	"vhagar/inspect"
	"vhagar/libs"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var crontabCmd = &cobra.Command{
	Use:   "crontab",
	Short: " 启动定时任务",
	Long: `可自定义周期性运行 job
相关配置见配置文件的 [crontab]
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("启动任务调度")
		crontabJob()
	},
}

func init() {
	rootCmd.AddCommand(crontabCmd)
}

func crontabJob() {
	c := cron.New() //创建一个cron实例
	// 获取等待时间
	duration := inspect.GetRandomDuration()
	// 初始化 inspect
	_inspect := NewInspect(CONFIG)
	// 租户巡检 job
	if CONFIG.Cron["tenant"].Crontab {
		// 创建ESClient，PGClient
		esClient, _ := libs.NewESClient(CONFIG.ES)
		pgClient, _ := libs.NewPGClient(CONFIG.PG)
		customerClient, _ := libs.NewCustomerClient(CONFIG.Customer)

		defer func() {
			if pgClient != nil {
				pgClient.Close()
			}
			if esClient != nil {
				esClient.Stop()
			}
			if customerClient != nil {
				if err := customerClient.Close(context.Background()); err != nil {
					log.Printf("Failed to close connection for database %s: %v", "customer", err)
				}
			}
		}()
		_inspect.Tenant = &inspect.Tenant{
			ESClient:       esClient,
			PGClient:       pgClient,
			Corp:           CONFIG.Tenant.Corp,
			CustomerClient: customerClient,
		}
		// 加入定时任务
		_, err := c.AddFunc(CONFIG.Cron["tenant"].Scheducron, func() {
			inspect.TenantTask(_inspect, duration)
		})
		if err != nil {
			log.Printf("Failed to add crontab job: %s \n", err)
		}
	}
	//  doris 巡检 job
	if CONFIG.Cron["doris"].Crontab {
		// 创建 mysqlClinet
		mysqlClinet, _ := libs.NewMysqlClient(CONFIG.Doris.DB, "wshoto")
		defer func() {
			if mysqlClinet != nil {
				err := mysqlClinet.Close()
				if err != nil {
					return
				}
			}
		}()
		_inspect.Doris = &inspect.Doris{
			MysqlClient: mysqlClinet,
		}
		_, err := c.AddFunc(CONFIG.Cron["doris"].Scheducron, func() {
			inspect.DorisTask(_inspect, duration)
		})
		if err != nil {
			log.Printf("Failed to add crontab job: %s \n", err)
		}
	}

	//启动/关闭
	c.Start()
	defer c.Stop()
	select {
	//查询语句，保持程序运行，在这里等同于for{}
	}
}

//func testjob() {
//	log.Printf("大王叫我来巡山，巡了南山巡北山。。。 \n")
//}
