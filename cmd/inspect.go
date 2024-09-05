// Package cmd @Author lanpang
// @Date 2024/8/6 下午4:49:00
// @Desc
package cmd

import (
	"context"
	"log"
	"vhagar/inspect"
	"vhagar/libs"

	"github.com/spf13/cobra"
)

var (
	rocketmq bool
	doris    bool
)

// inspectCmd represents the version command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "项目巡检",
	Long:  `获取项目的企业数据，活跃数，会话数`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化 inspect
		_inspect := NewInspect(CONFIG)
		switch {
		case rocketmq:
			_inspect.Rocketmq = CONFIG.Rocketmq
			inspect.RocketmqTask(_inspect)
		case doris:
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
				DorisCfg:    CONFIG.Doris,
			}
			inspect.DorisTask(_inspect, 0)
		default:
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
			inspect.TenantTask(_inspect, 0)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVarP(&rocketmq, "rocketmq", "m", false, "获取 rocketmq broker 信息")
	inspectCmd.Flags().BoolVarP(&doris, "doris", "d", false, "检查 doris 服务")

}

func NewInspect(cfg *Config) *inspect.Inspect {
	log.Println("初始化 Inspect 对象")

	_inspect := &inspect.Inspect{
		ProjectName: cfg.ProjectName,
		ProxyURL:    cfg.ProxyURL,
		Notifier:    cfg.Notifier,
	}
	return _inspect

}
