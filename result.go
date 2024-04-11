package ipInfo

import (
	"errors"
	"net"
	"strings"
)

type Result struct {
	Success bool
	Message string
	Service string
	Data    Data
}

type Data struct {
	IP       string
	Type     string
	CNIP     bool
	Country  string
	Province string
	City     string
	ISP      string
}

func getIPType(ip string) (*Result, error) {
	var r = &Result{
		Data: Data{
			IP: ip,
		},
	}

	netIP := net.ParseIP(ip)
	if netIP == nil {
		r.Message = "IP地址格式错误"
		r.Data.Type = "未知"
		return r, errors.New(r.Message)
	}
	if netIP.To4() != nil {
		r.Data.Type = "IPv4"
	} else {
		r.Data.Type = "IPv6"
	}

	return r, nil
}

func (r *Result) oneChina() {
	if strings.HasPrefix(r.Data.Country, "中国") {
		r.Data.CNIP = true
	} else {
		for _, p := range ChinaProvinces {
			if strings.HasPrefix(r.Data.Country, p) {
				r.Data.CNIP = true
				r.Data.City = r.Data.Province
				r.Data.Province = p
				r.Data.Country = "中国"
				break
			}
		}
	}
}
