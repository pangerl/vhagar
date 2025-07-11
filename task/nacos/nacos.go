// Package nacos @Author lanpang
// @Date 2024/8/2 下午2:05:00
// @Desc
package nacos

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/task"

	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
)

var tablerow []string
var mutex sync.Mutex

func init() {
	task.Add(taskName, func() task.Tasker {
		return NewNacos(config.Config, libs.Logger)
	})
}

//func GetNacos() *Nacos {
//	cfg := config.Config
//	nacos := newNacos(cfg)
//	if !nacos.WithAuth() {
//		return nil
//	}
//	nacos.Gather()
//	return nacos
//}

func (nacos *Nacos) Init() error {
	err := nacos.WithAuth()
	if err != nil {
		return err
	}
	return nil
}

func (nacos *Nacos) Check() {
	//task.EchoPrompt("开始巡检微服务状态信息")
	if nacos.Config.Nacos.Writefile != "" {
		nacos.WriteFile()
		return
	}
	if config.Config.Report {
		// 发送机器人
		nacos.ReportRobot()
		return
	}
	if nacos.Config.Global.Watch {
		nacos.Logger.Info("Watch模式已开启")
		interval := nacos.Config.Global.Interval
		nacos.Logger.Infof("监控模式 刷新时间:%s/次", interval)
		for {
			libs.Logger.Info("")
			nacos.Gather()
			nacos.TableRender()
			time.Sleep(interval)
		}
	}
	nacos.TableRender()
}

func (nacos *Nacos) ReportRobot() {
	nacos.Logger.Info("暂不支持发送企微机器人")
}

func (nacos *Nacos) WithAuth() error {
	nacos.Logger.Info("更新 nacos 的 token")
	_url := fmt.Sprintf("%s/nacos/v1/auth/login", nacos.Config.Nacos.Server)
	formData := map[string]string{
		"username": nacos.Config.Nacos.Username,
		"password": nacos.Config.Nacos.Password,
	}
	res := nacos.post(_url, formData)
	if len(gjson.GetBytes(res, "accessToken").String()) != 0 {
		//log.Println("Authentication successful...")
		nacos.Token = gjson.GetBytes(res, "accessToken").String()
	} else {
		//log.Println("Authentication failed!")
		return fmt.Errorf("authentication failed")
	}
	return nil
}

func (nacos *Nacos) GetService(url string, namespaceId string, group string) []byte {
	_url := fmt.Sprintf("%s/nacos/v1/ns/service/list?pageNo=1&pageSize=500&namespaceId=%s&groupName=%s", url, namespaceId, group)
	res := nacos.get(_url)
	return res
}

func (nacos *Nacos) GetInstance(url string, servicename string, namespaceId string, group string) []byte {
	_url := fmt.Sprintf("%s/nacos/v1/ns/instance/list?serviceName=%s&namespaceId=%s&groupName=%s", url, servicename, namespaceId, group)
	//fmt.Println(_url)
	res := nacos.get(_url)
	return res
}

func InString(filed string, array []string) bool {
	for _, str := range array {
		if filed == str {
			return true
		}
	}
	return false
}

func (nacos *Nacos) tableAppend(table *tablewriter.Table, data []string) {
	datastr := strings.Join(data, "-")
	if !InString(datastr, tablerow) {
		tablerow = append(tablerow, datastr)
		table.Append(data)
	}
}

func (nacos *Nacos) TableRender() {
	tablerow = []string{}
	nacosServer := nacos.Clusterdata
	tabletitle := []string{"服务名称", "实例", "健康状态", "主机名", "权重", "组", "命名空间"}
	table := tablewriter.NewWriter(task.GetOutputWriter())
	table.SetHeader(tabletitle)
	table.SetRowLine(true)
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1})
	for _, v := range nacosServer.HealthInstance {
		tabledata := []string{v.ServiceName, v.IpAddr, v.Health, v.Hostname, v.Weight, v.GroupName, v.NamespaceName}
		nacos.tableAppend(table, tabledata)
	}
	caption := fmt.Sprintf("健康实例: %d .", table.NumLines())
	table.SetCaption(true, caption)
	table.Render()
}

func (nacos *Nacos) Gather() {
	var ser Service
	var cluster ClusterStatus
	_url := nacos.Config.Nacos.Server
	namespace := nacos.Config.Nacos.Namespace
	group := "DEFAULT_GROUP"
	res := nacos.GetService(_url, namespace, group)
	err := json.Unmarshal(res, &ser)
	if err != nil {
		fmt.Println(err)
	}
	for _, se := range ser.Doms {
		res := nacos.GetInstance(_url, se, namespace, group)
		var in Instance
		err := json.Unmarshal(res, &in)
		if err != nil {
			fmt.Printf("json序列化错误:%s\n", err)
		}
		for _, host := range in.Hosts {
			instance := ServerInstance{
				NamespaceName: namespace,
				ServiceName:   se,
				IpAddr:        fmt.Sprintf("%s:%d", host.Ip, host.Port),
				Health:        strconv.FormatBool(host.Healthy),
				Hostname:      host.Ip,
				Weight:        fmt.Sprintf("%.1f", host.Weight),
				Ip:            host.Ip,
				Port:          strconv.Itoa(host.Port),
				GroupName:     in.GroupName,
			}
			if host.Healthy {
				cluster.HealthInstance = append(cluster.HealthInstance, instance)
			} else {
				cluster.UnHealthInstance = append(cluster.UnHealthInstance, instance)
			}
			//fmt.Println(instance)
		}
		nacos.Clusterdata = cluster
	}
	//fmt.Println(nacos.Clusterdata.HealthInstance)
}

func (nacos *Nacos) GetJson(resultType string) (result interface{}, err error) {
	mutex.Lock()
	defer func() {
		mutex.Unlock()
		if funcErr := recover(); funcErr != nil {
			result = []byte("[]")
			err = errors.New("error")
		}
	}()
	var nacosfile Nacosfile
	nacosServer := nacos.Clusterdata
	if len(nacosServer.HealthInstance) != 0 {
		for _, na := range nacosServer.HealthInstance {
			var ta Nacostarget
			ta.Labels = make(map[string]string)
			ta.Targets = append(ta.Targets, na.IpAddr)
			ta.Labels["namespace"] = na.NamespaceName
			ta.Labels["service"] = na.ServiceName
			ta.Labels["hostname"] = na.Hostname
			ta.Labels["weight"] = na.Weight
			ta.Labels["ip"] = na.Ip
			ta.Labels["port"] = na.Port
			ta.Labels["group"] = na.GroupName
			nacosfile.Data = append(nacosfile.Data, ta)
		}

	}

	if resultType == "json" {
		result = nacosfile.Data
		return result, err
	}
	data, err := json.MarshalIndent(&nacosfile.Data, "", "  ")
	if err != nil {
		fmt.Println("json序列化失败!")
		result = []byte("[]")
		return result, err
	}
	result = data
	return result, err
}

func (nacos *Nacos) WriteFile() {
	var basedir string
	var filename string
	basedir = path.Dir(nacos.Config.Nacos.Writefile)
	filename = path.Base(nacos.Config.Nacos.Writefile)
	if err := os.MkdirAll(basedir, os.ModePerm); err != nil {
		os.Exit(1)
	}
	file, err := os.OpenFile(basedir+"/.nacos_tmp.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		nacos.Logger.Errorw("创建文件失败", "err", err)
		os.Exit(2)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	jsondata, err := nacos.GetJson("byte")
	data := make([]byte, 0)
	var check bool
	if data, check = jsondata.([]byte); !check {
		nacos.Logger.Error("转换失败")
	}
	if _, err := file.Write(data); err != nil {
		nacos.Logger.Errorw("写入失败", "err", err)
		os.Exit(1)
	}
	err = file.Close()
	if err != nil {
		return
	}
	if err := os.Rename(basedir+"/.nacos_tmp.json", basedir+"/"+filename); err != nil {
		nacos.Logger.Errorw("写入失败", "file", basedir+"/"+filename)
	} else {
		nacos.Logger.Infow("写入成功", "file", basedir+"/"+filename)
	}
}
