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
	"encoding/json"
	"log"
)

var manifest *Manifest
var curiosity *Rover
var opportunity *Rover
var spirit *Rover

func main() {
	r := gin.Default()

	go func() {
		InitAndWatch("data/manfiest.json", &manifest)
		InitAndWatch("data/curiosityDates.json", &curiosity)
		InitAndWatch("data/opportunityDates.json", &opportunity)
		InitAndWatch("data/spiritDates.json", &spirit)
	}()

	r.GET("/manifest", func(c *gin.Context) {
		fmt.Println("get request for manifest")

		c.JSON(http.StatusOK, manifest)
	})

	r.GET("/rover/:rover", func(c *gin.Context) {
		updateData := c.Query("update")

		if updateData == "" {
			returnRoverData(c)
		} else {
			reader := ReturnLatestManifestReader(c)
			json.NewDecoder(reader).Decode(&manifest)
			fmt.Sprintln(manifest)
			//GetLatestRoverData(roverParam)
		}
	})

	r.Run(":8080")
}

func returnRoverData(c *gin.Context) {
	switch c.Param("rover") {
	case "curiosity":
		c.JSON(http.StatusOK, &curiosity)
	case "opportunity":
		c.JSON(http.StatusOK, &opportunity)
	case "spirit":
		c.JSON(http.StatusOK, &spirit)
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", c.Param("rover"))
	}
}
