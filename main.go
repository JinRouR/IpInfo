package main

import (
	"fmt"
	utils "go-ip-info/pkgs"
	"net/http"

	"github.com/gin-gonic/gin"
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

func main() {
	router := gin.Default()

	router.GET("/ip-info", func(c *gin.Context) {
		// set ip
		c.Request.Header.Set("X-Forwarded-For", "182.130.13.20")
		var ip string
		ip = c.ClientIP()
		if ip != "" {
			fmt.Println("Client IP:", ip)
		} else {
			ip = utils.GetClientIP(c.Request)
			if ip == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get IP address"})
				return
			}
		}

		geoInfo, err := utils.GetGeoInfo(ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if geoInfo.Status != "success" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get IP information"})
			return
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
		c.JSON(http.StatusOK, response)
	})
	fmt.Println("Server is running on port 8080...")
	router.Run(":8080")
}
