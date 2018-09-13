// 09-10-2018
// @TODO #2: Make API calls to NASA server and save individual rover information in files
// @TODO #3: Create interface for manifest and rover?

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"encoding/json"
)

var manifest *Manifest
var curiosity *Rover
var opportunity *Rover
var spirit *Rover
var emptyRover *Rover

var FileChange chan bool

func main() {
	r := gin.Default()

	FileChange = make(chan bool, 1)

	go func() {
		InitAndWatch("data/manifest.json", &manifest)
		InitAndWatch("data/curiosity.json", &curiosity)
		InitAndWatch("data/opportunity.json", &opportunity)
		InitAndWatch("data/spirit.json", &spirit)
	}()

	r.GET("/manifest", func(c *gin.Context) {
		fmt.Println("get request for manifest")

		c.JSON(http.StatusOK, manifest)
	})

	r.GET("/rover/:rover", func(c *gin.Context) {
		var roverData string
		var manifestData string
		manifestFile := fmt.Sprintf("data/manifest.json")
		updateData := c.Query("update")
		roverParam := c.Param("rover")

		if updateData == "" {
			fmt.Println("updateData is empty string")
			manifestData, err = SlurpFile(manifestFile)
			roverStruct := returnRoverData(roverParam)
			bytes, err := json.Marshal(&roverStruct)
			Check(err)
			roverData = string(bytes)
			fmt.Println(roverData)
		} else if updateData == "true" {
			fmt.Println("updateData == true")
			// Update Manifest
			manifestBytes := ReturnLatestManifest(c)
			WriteFile(manifestFile, manifestBytes)
			for n := range FileChange{
				fmt.Println(n)
				if n == true {
					roverStruct := returnRoverData(roverParam)
					bytes, err := json.Marshal(roverStruct)
					Check(err)
					roverData = string(bytes)
				}
			}
		}

		fmt.Println("out of condition")
		c.JSON(http.StatusOK, gin.H{
			"manifest": &manifestData,
			"rover": &roverData,
		})
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
