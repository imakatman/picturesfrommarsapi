// 8-24-2018
// Problem at the moment is that the file change event that is
// triggered when the json file itself is modified is RENAME.
//
// This in particular isn't the issue, the method of streaming the
// new data from the NASA API is going to dictate which file change
// events occur. And the method of choice will have to depend on
// which method is the most performant and cost-effective for the server.

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type Config struct {
	url   string
	token []string
}

var manifest *Manifest
var curiosity *Rover
var opportunity *Rover
var spirit *Rover

func main() {
	r := gin.Default()

	WatchFile("data/manfiest.json", &manifest)
	WatchFile("data/curiosityDates.json", &curiosity)
	WatchFile("data/opportunityDates.json", &opportunity)
	WatchFile("data/spiritDates.json", &spirit)

	r.GET("/manifest", func(c *gin.Context) {
		fmt.Println("get request for manifest")

		c.JSON(http.StatusOK, manifest)
	})

	r.Run(":8080")

	r.GET("/rover/:rover", func(c *gin.Context) {
		updateData := c.Query("update")

		if updateData == "" {
			returnRoverData(c)
		} else {
			roverParam := c.Param("rover")

			manifestReader := GetLatestManifest(c)
			//GetLatestRoverData(roverParam)
		}
	})
}

func returnRoverData(c *gin.Context) {
	roverParam := c.Param("rover")

	switch roverParam {
	case "curiosity":
		c.JSON(http.StatusOK, &curiosity)
	case "opportunity":
		c.JSON(http.StatusOK, &opportunity)
	case "spirit":
		c.JSON(http.StatusOK, &spirit)
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", roverParam)
	}
}
