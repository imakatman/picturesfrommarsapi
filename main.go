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

var manifest *Manifest
var curiosity *Rover
var opportunity *Rover
var spirit *Rover
var emptyRover *Rover

var FileChange chan bool

func main() {
	r := gin.Default()

	go func() {
		InitAndWatch("data/manifest.json", &manifest)
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
			data := returnRoverData(c.Param("rover"))
			c.JSON(http.StatusOK, &data)
		} else if updateData == "true" {
			// Update Manifest
			manifestFile := fmt.Sprintf("data/manifest.json")
			manifestBytes := ReturnLatestManifest(c)
			WriteFile(manifestFile, manifestBytes)
			data := returnRoverData(c.Param("rover"))
			for n := range FileChange{
				fmt.Println(n)
				if n == true {
					fmt.Println("n is true")
				}
			}
			c.JSON(http.StatusOK, gin.H{
				"manifest": string(SlurpFile(manifestFile)),
				"rover": data,
			})
		}
	})

	r.Run(":8080")
}

func returnRoverData(rover string) *Rover {
	switch rover {
	case "curiosity":
		return curiosity
	case "opportunity":
		return opportunity
	case "spirit":
		return spirit
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", rover)
		return emptyRover
	}
}
