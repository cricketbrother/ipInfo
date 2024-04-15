// 文档地址：https://api.mir6.com/doc/ip_json.html
// 支持IPv4和IPv6

package ipInfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data CData  `json:"data"`
}

type CData struct {
	IP          string `json:"ip"`
	Dec         string `json:"dec"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Districts   string `json:"districts"`
	IDC         string `json:"idc"`
	ISP         string `json:"isp"`
	NET         string `json:"net"`
	ZIPCcode    string `json:"zipcode"`
	AreaCode    string `json:"areaCode"`
	Protocol    string `json:"protocol"`
	Location    string `json:"location"`
	MyIP        string `json:"myip"`
	Time        string `json:"time"`
}

func getFromApiMir6Com(ip string) (*Result, error) {
	r, err := getIPType(ip)
	if err != nil {
		return r, err
	}

	r.Service = "api.mir6.com"

	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.mir6.com/api/ip_json?ip="+ip, nil)
	if err != nil {
		return r, fmt.Errorf("构建请求失败：%v", err)
	}

	res, err := c.Do(req)
	if err != nil {
		return r, fmt.Errorf("请求失败：%v", err)
	}
	defer res.Body.Close()

	var resp CResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return r, fmt.Errorf("解析失败：%v", err)
	}

	if resp.Code != 200 {
		return r, fmt.Errorf("从api.mir6.com返回错误：%d %s", resp.Code, resp.Msg)
	}

	r.Success = true
	r.Message = "查询成功"

	r.Data.Country = resp.Data.Country
	r.Data.Province = resp.Data.Province
	r.Data.City = resp.Data.City
	r.Data.ISP = resp.Data.ISP

	r.oneChina()

	return r, nil
}

func GetFromApiMir6Com(ip string) *Result {
	r, err := getFromApiMir6Com(ip)
	if err != nil {
		r.Message = err.Error()
	}
	return r
}
