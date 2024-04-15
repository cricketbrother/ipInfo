// 文档地址：https://api.vvhan.com/article/ipinfo.html
// 支持IPv4和IPv6

package ipInfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BResponse struct {
	Success bool   `json:"success"`
	IP      string `json:"ip"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Info    BInfo  `json:"info"`
}

type BInfo struct {
	Country string `json:"country"`
	Prov    string `json:"prov"`
	City    string `json:"city"`
	ISP     string `json:"isp"`
}

func getFromApiVvhanCom(ip string) (*Result, error) {
	r, err := getIPType(ip)
	if err != nil {
		return r, err
	}

	r.Service = "api.vvhan.com"

	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.vvhan.com/api/ipInfo?ip="+ip, nil)
	if err != nil {
		return r, fmt.Errorf("构建请求失败：%v", err)
	}

	res, err := c.Do(req)
	if err != nil {
		return r, fmt.Errorf("请求失败：%v", err)
	}
	defer res.Body.Close()

	var resp BResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return r, fmt.Errorf("解析失败：%v", err)
	}

	if !resp.Success {
		return r, fmt.Errorf("从api.vvhan.com返回错误：%d %s", resp.Code, resp.Message)
	}

	r.Success = true
	r.Message = "查询成功"

	r.Data.Country = resp.Info.Country
	r.Data.Province = resp.Info.Prov
	r.Data.City = resp.Info.City
	r.Data.ISP = resp.Info.ISP

	r.oneChina()

	return r, nil
}

func GetFromApiVvhanCom(ip string) *Result {
	r, err := getFromApiVvhanCom(ip)
	if err != nil {
		r.Message = err.Error()
	}
	return r
}
