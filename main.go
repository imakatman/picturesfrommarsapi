// 09-10-2018
// @TODO #3: Create interface for manifest and rover?

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"encoding/json"
)

var manifest Manifest
var Curiosity Rover
var Opportunity Rover
var Spirit Rover
var emptyRover Rover

var FileChange chan bool

func init() {
	launched := make(chan bool)

	type dataPair struct {
		file string
		obj  interface{}
	}

	dataPairs := []dataPair{
		{"data/manifest.json", &manifest},
		{"data/curiosity.json", &Curiosity},
		{"data/opportunity.json", &Opportunity},
		{"data/spirit.json", &Spirit},
	}

	for i, v := range dataPairs {
		if i == 0 {
			go func() {
				fmt.Println("if i == 0")
				fmt.Println(v.file, i)
				InitAndWatch(v.file, &v.obj, launched)
			}()
		}

		if <-launched == true {
			go func() {
				fmt.Println("if <-launched == true")
				fmt.Println(v.file, i)
				InitAndWatch(v.file, &v.obj, launched)
			}()
		}
	}
}

func main() {
	r := gin.Default()

	FileChange = make(chan bool)

	r.GET("/manifest", func(c *gin.Context) {
		fmt.Println("get request for manifest")

		c.JSON(http.StatusOK, manifest)
	})

	r.GET("/rover/:rover", func(c *gin.Context) {
		var roverData *Rover
		var manifestData *Manifest
		updateData := c.Query("update")
		roverParam := c.Param("rover")

		if updateData == "" {
			fmt.Println("updateData is empty string")
			manifestData = &manifest
			roverStruct := ReturnRoverData(roverParam)
			// ReturnRoverData returns a pointer to a Rover
			roverData = roverStruct
		} else if updateData == "true" {
			fmt.Println("updateData == true")
			roverFile := fmt.Sprintf("data/%s.json", roverParam)
			// Update Manifest
			manifestBytes := ReturnLatestManifest(c)
			WriteFile("data/manifest.json", manifestBytes)
			status := <-FileChange
			if status == true {
				manifestData = &manifest
				roverBytes := ReturnLatestRoverData(c)
				roverStruct := ReturnRoverData(roverParam)
				json.Unmarshal(roverBytes, roverStruct)
				roverBytes, err := json.Marshal(roverStruct)
				Check(err)
				WriteFile(roverFile, roverBytes)
				status := <-FileChange
				if status == true {
					fmt.Println("2 if status == true")
					fmt.Println("roverStruct.MaxDate", roverStruct.MaxDate)
					roverData = roverStruct
				}
			}
		}

		fmt.Println("out of condition")
		c.JSON(http.StatusOK, gin.H{
			"manifest": *manifestData,
			"rover":    *roverData,
		})
	})

	r.Run(":8080")
}

func ReturnRoverData(rover string) *Rover {
	switch rover {
	case "curiosity":
		return &Curiosity
	case "opportunity":
		return &Opportunity
	case "spirit":
		return &Spirit
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", rover)
		return &emptyRover
	}
}
