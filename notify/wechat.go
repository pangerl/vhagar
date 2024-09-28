// Package notify @Author lanpang
// @Date 2024/8/8 下午5:14:00
// @Desc
package notify

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type WeChatMarkdown struct {
	MsgType  string    `json:"msgtype"`
	Markdown *Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

func sendWecom(markdown *WeChatMarkdown, robotKey, proxyURL string) error {

	jsonStr, _ := json.Marshal(markdown)
	//fmt.Println("jsonStr长度：", len(jsonStr))
	robotURL := wechatRobotURL + robotKey

	req, err := http.NewRequest("POST", robotURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return err
		}
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed info: %s \n", err)
		}
	}(resp.Body)
	log.Print("推送企微机器人 response Status:", resp.Status)
	//log.Print("response Headers:", resp.Header)
	return nil
}