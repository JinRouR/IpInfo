package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type GeoInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	TimeZone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func ValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func GetClientIP(r *http.Request) string {
	// check forwarded header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		ip := strings.TrimSpace(parts[0])
		if ValidIP(ip) {
			return ip
		}
	}

	// check x-real-ip header
	realIP := r.Header.Get("X-Real-Ip")

	if realIP != "" && ValidIP(realIP) {
		return realIP
	}

	// if header doesn't contain the IP, get it from the remote address
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	if ValidIP(remoteIP) {
		return remoteIP
	}
	return ""
}

func GetGeoInfo(ip string) (*GeoInfo, error) {
	// get geo info from ip
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,query", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geoInfo GeoInfo
	if err := json.Unmarshal(body, &geoInfo); err != nil {
		return nil, err
	}
	return &geoInfo, nil
}
