// Package task @Author lanpang
// @Date 2024/9/13 下午3:44:00
// @Desc
package task

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"vhagar/libs"
)

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func CallUser(users []string) string {
	var result string
	if len(users) == 0 {
		return result
	}
	for _, user := range users {
		result += fmt.Sprintf("<@%s>", user)
	}
	return result
}

func DoRequest(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			libs.Logger.Errorw("关闭响应失败", "err", err)
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)

	if err != nil {
		libs.Logger.Errorw("读取请求数据失败", "err", err)
		return nil
	} else {
		return body
	}
}

func echoPrompt(prompt string) {
	date := time.Now().Format("2006-01-02 15:04:05")
	taskPrompt := fmt.Sprintf(`
================================================================
%s %s
================================================================`, date, prompt)
	fmt.Fprintf(GetOutputWriter(), "\033[34m\033[1m%s\033[0m\n", taskPrompt)
}
