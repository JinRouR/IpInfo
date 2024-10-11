package main

import (
	"encoding/json"
	"fmt"
	utils "go-ip-info/pkgs"
	"log"
	"net/http"
)

type Response struct {
	IP        string  `json:"ip"`
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Zip       string  `json:"zip"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	TimeZone  string  `json:"timezone"`
	ISP       string  `json:"isp"`
}

func ipInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("X-Forwarded-For", "182.130.13.11")
	ip := utils.GetClientIP(r)
	if ip == "" {
		http.Error(w, "Unable to get IP address", http.StatusBadRequest)
		return
	}

	geoInfo, err := utils.GetGeoInfo(ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if geoInfo.Status != "success" {
		http.Error(w, "Unable to get IP information", http.StatusBadRequest)
	}

	response := Response{
		IP:        geoInfo.Query,
		Country:   geoInfo.Country,
		Region:    geoInfo.RegionName,
		City:      geoInfo.City,
		Zip:       geoInfo.Zip,
		Latitude:  geoInfo.Lat,
		Longitude: geoInfo.Lon,
		TimeZone:  geoInfo.TimeZone,
		ISP:       geoInfo.ISP,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/ip-info", ipInfoHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
