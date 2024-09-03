// Package inspect @Author lanpang
// @Date 2024/8/7 下午3:43:00
// @Desc
package inspect

import (
	"log"
	"time"
	"vhagar/notifier"
)

var isalert = false

func TenantTask(inspect *Inspect, duration time.Duration) {
	tenant := inspect.Tenant
	// 填充租户信息
	tenantDetail(tenant)
	// 发送巡检报告
	markdownList := tenantNotifier(tenant, inspect.ProjectName, inspect.Notifier["tenant"].Userlist)
	log.Println("任务等待时间", duration)
	time.Sleep(duration)
	for _, markdown := range markdownList {
		for _, robotkey := range inspect.Notifier["tenant"].Robotkey {
			err := notifier.SendWecom(markdown, robotkey, inspect.ProxyURL)
			if err != nil {
				return
			}
		}
	}
	if inspect.Notifier["tenant"].IsPush {
		log.Println("推送微盛运营平台")
		// 将 []*Corp 转换为 []any
		var data []any = make([]any, len(tenant.Corp))
		for i, c := range tenant.Corp {
			data[i] = c
		}
		inspectBody := notifier.InspectBody{
			JobType: "tenant",
			Data:    data,
		}
		err := notifier.SendWshoto(&inspectBody, inspect.ProxyURL)
		if err != nil {
			return
		}
	}

}

func RocketmqTask(inspect *Inspect) {
	log.Print("启动 rocketmq 巡检任务")
	clusterdata, _ := GetMQDetail(inspect.Rocketmq.RocketmqDashboard)
	markdown := mqDetailMarkdown(clusterdata, inspect.ProjectName)
	for _, robotkey := range inspect.Notifier["rocketmq"].Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, inspect.ProxyURL)
	}
}

func DorisTask(inspect *Inspect, duration time.Duration) {
	log.Print("启动 doris 巡检任务")
	// 获取当前零点时间
	todayTime := getZeroTime(time.Now())
	yesterday := todayTime.AddDate(0, 0, -1)
	yesterdayTime := getZeroTime(yesterday)
	if inspect.Doris.MysqlClient != nil {
		// 失败任务
		failedJobs := selectFailedJob(todayTime.String(), inspect.Doris.MysqlClient)
		inspect.Doris.FailedJobs = failedJobs
		// 员工统计表
		staffCount := selectStaffCount(yesterdayTime.String(), inspect.Doris.MysqlClient)
		inspect.Doris.StaffCount = staffCount
		// 使用分析表
		useAnalyseCount := selectUseAnalyseCount(yesterdayTime.String(), inspect.Doris.MysqlClient)
		inspect.Doris.UseAnalyseCount = useAnalyseCount
		// 客户群统计表
		customerGroupCount := selectCustomerGroupCount(yesterdayTime.String(), inspect.Doris.MysqlClient)
		inspect.Doris.CustomerGroupCount = customerGroupCount
	}
	// 发送巡检报告
	markdown := dorisToMarkdown(inspect.Doris, inspect.ProjectName)
	log.Println("任务等待时间", duration)
	time.Sleep(duration)
	for _, robotkey := range inspect.Notifier["doris"].Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, inspect.ProxyURL)
	}
}
