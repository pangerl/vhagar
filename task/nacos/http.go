package nacos

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"vhagar/libs"
)

func (nacos *Nacos) get(apiurl string) []byte {
	u, err := url.Parse(apiurl)
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
	}
	if len(nacos.Config.Nacos.Username) != 0 && len(nacos.Config.Nacos.Password) != 0 {
		if len(u.RawQuery) == 0 {
			apiurl += "?accessToken=" + url.QueryEscape(nacos.Token)
		} else {
			apiurl += "&accessToken=" + url.QueryEscape(nacos.Token)
		}
	}
	req, _ := http.NewRequest("GET", apiurl, nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
	}
	if res.StatusCode != 200 {
		if res.StatusCode == 403 {
			libs.Logger.Errorf("%s请求状态码异常:nacos 请使用--username --password参数进行鉴权", apiurl, res.StatusCode)
		}
		libs.Logger.Errorf("%s请求状态码异常:nacos", apiurl, res.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			libs.Logger.Errorw("Failed info", "err", err)
		}
	}(res.Body)
	resp, _ := io.ReadAll(res.Body)
	return resp
}

func (nacos *Nacos) post(apiurl string, formData map[string]string) []byte {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for key, val := range formData {
		_ = bodyWriter.WriteField(key, val)
	}
	contentType := bodyWriter.FormDataContentType()
	err := bodyWriter.Close()
	if err != nil {
		return nil
	}
	var req *http.Request
	u, err := url.Parse(apiurl)
	if u.Path == "/nacos/v1/auth/login" {
		req, _ = http.NewRequest("POST", apiurl, bodyBuf)
		req.Header.Set("Content-Type", contentType)
	}
	res, err := nacos.Client.Do(req)
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
	}
	if res.StatusCode != 200 {
		if u.Path == "/nacos/v1/auth/login" && res.StatusCode == 403 {
			libs.Logger.Errorf("%s请求状态码异常,认证失败!:nacos", apiurl, res.StatusCode)
		}
		libs.Logger.Errorf("%s请求状态码异常:nacos", apiurl, res.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			libs.Logger.Errorw("Failed info", "err", err)
		}
	}(res.Body)
	resp, _ := io.ReadAll(res.Body)
	return resp
}
