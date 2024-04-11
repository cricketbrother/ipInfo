// 文档地址：https://api.vore.top/doc/IPdata.html
// 支持IPv4和IPv6

package ipInfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AResponse struct {
	Code   int     `json:"code"`
	Msg    string  `json:"msg"`
	IPInfo AIPInfo `json:"ipinfo"`
	IPData AIPData `json:"ipdata"`
	Adcode AAdcode `json:"adcode"`
	Tips   string  `json:"tips"`
	Time   int64   `json:"time"`
}

type AIPInfo struct {
	Type string `json:"type"`
	Text string `json:"text"`
	CNIP bool   `json:"cnip"`
}

type AIPData struct {
	Info1 string `json:"info1"`
	Info2 string `json:"info2"`
	Info3 string `json:"info3"`
	ISP   string `json:"isp"`
}

type AAdcode struct {
	O string      `json:"o"`
	P string      `json:"p"`
	C string      `json:"c"`
	N string      `json:"n"`
	R interface{} `json:"r"`
	A interface{} `json:"a"`
	I bool        `json:"i"`
}

func getFromApiVeroCom(ip string) (*Result, error) {
	r, err := getIPType(ip)
	if err != nil {
		return r, err
	}

	r.Service = "api.vero.com"

	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.vore.top/api/IPdata?ip="+ip, nil)
	if err != nil {
		return r, fmt.Errorf("构建请求失败：%v", err)
	}

	res, err := c.Do(req)
	if err != nil {
		return r, fmt.Errorf("请求失败：%v", err)
	}
	defer res.Body.Close()

	var resp AResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return r, fmt.Errorf("解析失败：%v", err)
	}

	if resp.Code != 200 {
		return r, fmt.Errorf("从api.vero.com返回错误：%d %s", resp.Code, resp.Msg)
	}

	r.Success = true
	r.Message = "查询成功"

	r.Data.Country = resp.IPData.Info1
	r.Data.Province = resp.IPData.Info2
	r.Data.City = resp.IPData.Info3
	r.Data.ISP = resp.IPData.ISP

	r.oneChina()

	return r, nil
}

func GetFromApiVeroCom(ip string) *Result {
	r, err := getFromApiVeroCom(ip)
	if err != nil {
		r.Message = err.Error()
	}
	return r
}
