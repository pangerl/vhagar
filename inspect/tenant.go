// Package inspect @Author lanpang
// @Date 2024/8/23 上午11:15:00
// @Desc
package inspect

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
	"vhagar/notifier"

	"github.com/jackc/pgx/v5"
	"github.com/olivere/elastic/v7"
)

func tenantDetail(tenant *Tenant) {
	// 当前时间
	dateNow := time.Now()
	log.Print("启动企微租户信息巡检任务")
	for _, corp := range tenant.Corp {
		// fmt.Println(corp.Corpid)
		if tenant.PGClient != nil {
			// 获取租户名
			tenant.SetCorpName(corp.Corpid)
			// 获取用户数
			tenant.SetUserNum(corp.Corpid)
			// 获取客户群
			tenant.SetCustomerGroupNum(corp.Corpid)
			// 获取客户群人数
			tenant.SetCustomerGroupUserNum(corp.Corpid)
		}
		if tenant.ESClient != nil {
			// 获取客户数
			tenant.SetCustomerNum(corp.Corpid)
			// 获取活跃数
			tenant.SetActiveNum(corp.Corpid, dateNow)
			// 获取会话数
			if corp.Convenabled {
				tenant.SetMessageNum(corp.Corpid, dateNow)
			}
		}
		//fmt.Println(*corp)
	}
}
func tenantNotifier(t *Tenant, name string, userlist []string) []*notifier.WeChatMarkdown {

	var inspectList []*notifier.WeChatMarkdown
	isalert = false

	headString := headCorpString(t, name)

	length := len(t.Corp)
	// 每次返回8个租户的信息
	chunkSize := 8

	for n := 0; n < length; n += chunkSize {
		end := n + chunkSize
		if end > length {
			end = length
		}
		slice := t.Corp[n:end]
		markdown := tenantMarkdown(headString, slice, userlist)
		inspectList = append(inspectList, markdown)
	}
	return inspectList
}
func tenantMarkdown(headString string, Corp []*Corp, users []string) *notifier.WeChatMarkdown {
	var builder strings.Builder
	// 添加巡检头文件
	builder.WriteString(headString)
	for _, corp := range Corp {
		// 组装租户巡检信息
		builder.WriteString(generateCorpString(corp))
	}
	if isalert {
		builder.WriteString("\n<font color='red'>**注意！巡检结果异常！**</font>" + calluser(users))
	}
	markdown := &notifier.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &notifier.Markdown{
			Content: builder.String(),
		},
	}

	// fmt.Println("调试信息", builder.String())
	return markdown
}
func generateCorpString(corp *Corp) string {
	var builder strings.Builder

	builder.WriteString("> 企业名称：<font color='info'>" + corp.CorpName + "</font>\n")
	if corp.Convenabled {
		builder.WriteString("> 昨天拉取会话数：<font color='info'>" + strconv.FormatInt(corp.MessageNum, 10) + "</font>\n")
		if corp.MessageNum <= 0 {
			isalert = true
		}
	}
	builder.WriteString("> 员工人数：<font color='info'>" + strconv.Itoa(corp.UserNum) + "</font>\n")
	builder.WriteString("> 客户人数：<font color='info'>" + strconv.FormatInt(corp.CustomerNum, 10) + "</font>\n")
	builder.WriteString("> 客户群数：<font color='info'>" + strconv.Itoa(corp.CustomerGroupNum) + "</font>\n")
	builder.WriteString("> 客户群人数：<font color='info'>" + strconv.Itoa(corp.CustomerGroupUserNum) + "</font>\n")
	builder.WriteString("> 日活跃数：<font color='info'>" + strconv.FormatInt(corp.DauNum, 10) + "</font>\n")
	builder.WriteString("> 周活跃数：<font color='info'>" + strconv.FormatInt(corp.WauNum, 10) + "</font>\n")
	builder.WriteString("> 月活跃数：<font color='info'>" + strconv.FormatInt(corp.MauNum, 10) + "</font>\n")
	builder.WriteString("==================\n")

	return builder.String()
}
func headCorpString(t *Tenant, name string) string {
	var builder strings.Builder
	// 组装巡检内容
	builder.WriteString("# 每日巡检报告 " + version + "\n")
	builder.WriteString("**项目名称：**<font color='info'>" + name + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**巡检内容：**\n")

	return builder.String()
}

// SetCustomerGroupUserNum 设置客户群人数
func (t *Tenant) SetCustomerGroupUserNum(corpid string) {
	customergroupusernum, _ := queryCustomerGroupUserNum(t.CustomerClient, corpid)
	for _, corp := range t.Corp {
		if corp.Corpid == corpid {
			corp.CustomerGroupUserNum = customergroupusernum
			return
		}
	}
}

// SetCustomerGroupNum 设置客户群数
func (t *Tenant) SetCustomerGroupNum(corpid string) {
	customergroupnum, _ := queryCustomerGroupNum(t.CustomerClient, corpid)
	for _, corp := range t.Corp {
		if corp.Corpid == corpid {
			corp.CustomerGroupNum = customergroupnum
			return
		}
	}
}

// SetMessageNum 统计昨天的会话数
func (t *Tenant) SetMessageNum(corpid string, dateNow time.Time) {
	date := dateNow.AddDate(0, 0, -1)
	startTime := getZeroTime(date).UnixNano() / 1e6
	endTime := getZeroTime(dateNow).UnixNano() / 1e6
	// 兼容加密 ID
	suiteId := corpid
	if strings.HasPrefix(corpid, "wpIaoBE") {
		//fmt.Println("加密 ID 为", corpid)
		suiteId, _ = querySuiteId(t.PGClient.Conn["qv30"], corpid)
		//fmt.Println("解密 ID 为", suiteId)
	}
	messagenum, _ := countMessageNum(t.ESClient, suiteId, startTime, endTime)
	//fmt.Println("会话数为", messagenum)
	for _, corp := range t.Corp {
		if corp.Corpid == corpid {
			corp.MessageNum = messagenum
			return
		}
	}
}

// SetCorpName 设置租户名称
func (t *Tenant) SetCorpName(corpid string) {
	corpName, _ := queryCorpName(t.PGClient.Conn["qv30"], corpid)
	for _, corp := range t.Corp {
		if corp.Corpid == corpid {
			corp.CorpName = corpName
			return
		}
	}
}

// SetCustomerNum 设置客户数
func (t *Tenant) SetCustomerNum(corpid string) {
	customernum, _ := searchCustomerNum(t.ESClient, corpid)
	for _, corp := range t.Corp {
		if corp.Corpid == corpid {
			corp.CustomerNum = customernum
			return
		}
	}
}

// SetUserNum 设置员工数
func (t *Tenant) SetUserNum(corpid string) {
	userNum, _ := queryUserNum(t.PGClient.Conn["user"], corpid)
	for _, corp := range t.Corp {
		if corp.Corpid == corpid {
			corp.UserNum = userNum
			return
		}
	}
}

// SetActiveNum 设置活跃数
func (t *Tenant) SetActiveNum(corpid string, dateNow time.Time) {
	dateDau := dateNow.AddDate(0, 0, -1)
	dateWau := dateNow.AddDate(0, 0, -7)
	dateMau := dateNow.AddDate(0, -1, 0)
	for _, corp := range t.Corp {
		if corp.Corpid == corpid {
			corp.DauNum, _ = searchActiveNum(t.ESClient, corpid, dateDau, dateNow)
			corp.WauNum, _ = searchActiveNum(t.ESClient, corpid, dateWau, dateNow)
			corp.MauNum, _ = searchActiveNum(t.ESClient, corpid, dateMau, dateNow)
			return
		}
	}
}

// 客户数
func searchCustomerNum(client *elastic.Client, corpid string) (int64, error) {
	// 创建 bool 查询
	query := elastic.NewBoolQuery().
		Filter(
			elastic.NewTermQuery("tenantId", corpid),
			elastic.NewTermQuery("relatedHiddenAt", 0),
			elastic.NewTermQuery("relatedDelAt", 0),
		)
	searchResult, err := client.Search().
		Index("customer_related_1"). // 设置索引名
		Query(query).                // 设置查询条件
		TrackTotalHits(true).
		Do(context.Background()) // 执行
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	//fmt.Printf("总客户数: %d\n", searchResult.TotalHits())
	return searchResult.TotalHits(), nil
}

// 会话数
func countMessageNum(client *elastic.Client, corpid string, startTime, endTime int64) (int64, error) {

	// Define the query
	query := elastic.NewBoolQuery().
		Must(elastic.NewRangeQuery("msgtime").
			From(startTime). // from timestamp for yesterday 0:00:00
			To(endTime),     // to timestamp for today 0:00:00
		)
	// Make the count request
	countResult, err := client.Count().
		Index("conversation_" + corpid).
		Query(query).
		Do(context.Background())
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	//fmt.Printf("昨天消息数: %d\n", countResult)
	return countResult, nil
}

// 活跃数
func searchActiveNum(client *elastic.Client, corpid string, startDate, endDate time.Time) (int64, error) {
	startTime := getZeroTime(startDate).UnixNano() / 1e6
	endTime := getZeroTime(endDate).UnixNano() / 1e6
	// 创建 bool 查询
	query := elastic.NewBoolQuery().
		Must(
			elastic.NewTermsQuery("where.entrance", "001", "002", "006"),
			elastic.NewMatchQuery("who.role", "02"),
			elastic.NewTermQuery("where.corpId.keyword", corpid),
			elastic.NewRangeQuery("when.start").Gte(startTime).Lte(endTime),
		)
	searchResult, err := client.Search().
		Index("text_event_index*"). // 设置索引名
		Query(query).               // 设置查询条件
		Aggregation("dau", elastic.NewCardinalityAggregation().Field("who.id.keyword")).
		Size(0).
		Do(context.Background()) // 执行
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	dauAgg, _ := searchResult.Aggregations["dau"]
	cardinalityAgg := &struct {
		Value int64 `json:"value"`
	}{}
	err = json.Unmarshal(dauAgg, cardinalityAgg)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	//fmt.Println("活跃数：", cardinalityAgg.Value)
	return cardinalityAgg.Value, nil
}

// 租户名称
func queryCorpName(conn *pgx.Conn, corpid string) (string, error) {
	var corpName string
	query := "SELECT corp_name FROM qw_base_tenant_corp_info WHERE tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&corpName)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return "-1", err
	}
	return corpName, nil
}

// 解密 ID
func querySuiteId(conn *pgx.Conn, corpid string) (string, error) {
	var suiteId string
	query := "SELECT org_corp_id FROM qw_base_tenant_corp_info WHERE tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&suiteId)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return "-1", err
	}
	return suiteId, nil
}

// 员工数
func queryUserNum(conn *pgx.Conn, corpid string) (int, error) {
	var userNum int
	query := "SELECT count(*) FROM qw_user WHERE deleted=0 AND tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&userNum)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	return userNum, nil
}

// 客户群数
func queryCustomerGroupNum(conn *pgx.Conn, corpid string) (int, error) {
	var customerGroupNum int
	query := "SELECT count(1) FROM co_saas_customer_group WHERE dismiss=false AND tenant_id=$1 AND deleted_at IS NULL"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&customerGroupNum)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	return customerGroupNum, nil
}

// 客户群人数
func queryCustomerGroupUserNum(conn *pgx.Conn, corpid string) (int, error) {
	var customerGroupUserNum int
	query := "SELECT count(1) FROM co_saas_customer_group_user WHERE type = 2 AND loss = false AND deleted_at IS NULL AND tenant_id=$1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&customerGroupUserNum)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	return customerGroupUserNum, nil
}
