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
var Curiosity *Rover
var Opportunity *Rover
var Spirit *Rover
var emptyRover *Rover

var FileChange chan bool

//func init() {
//	type dataPair struct {
//		file string
//		obj  interface{}
//	}
//
//	dataPairs := []dataPair{
//		{"data/manifest.json", &manifest},
//		{"data/curiosity.json", &Curiosity},
//		{"data/opportunity.json", &Opportunity},
//		{"data/spirit.json", &Spirit},
//	}
//
//	for n := range dataPairs {
//		fmt.Printf("%T", n)
//		//go func() {
//		//	InitAndWatch(n.file, n.obj)
//		//}()
//	}
//}

func main() {
	r := gin.Default()

	FileChange = make(chan bool)

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

	for n := range dataPairs {
		fmt.Printf("%T", n)
		//go func() {
		//	InitAndWatch(n.file, n.obj)
		//}()
	}

	r.GET("/manifest", func(c *gin.Context) {
		fmt.Println("get request for manifest")

		c.JSON(http.StatusOK, manifest)
	})

	r.GET("/rover/:rover", func(c *gin.Context) {
		var roverData string
		var manifestData string
		updateData := c.Query("update")
		roverParam := c.Param("rover")

		if updateData == "" {
			fmt.Println("updateData is empty string")
			unchangedManifestData, err := SlurpFile("data/manifest.json")
			manifestData = string(unchangedManifestData)
			roverStruct := ReturnRoverData(roverParam)
			bytes, err := json.Marshal(&roverStruct)
			Check(err)
			roverData = string(bytes)
			fmt.Println(roverData)
		} else if updateData == "true" {
			fmt.Println("updateData == true")
			roverFile := fmt.Sprintf("data/%s.json", roverParam)
			// Update Manifest
			manifestBytes := ReturnLatestManifest(c)
			WriteFile("data/manifest.json", manifestBytes)
			status := <-FileChange
			if status == true {
				roverBytes := ReturnLatestRoverData(c)
				WriteFile(roverFile, roverBytes)
				FileChange <- true
				if status == true {
					roverStruct := ReturnRoverData(roverParam)
					fmt.Println(roverStruct.MaxDate)
					bytes, err := json.Marshal(roverStruct)
					Check(err)
					roverData = string(bytes)
				}
			}
		}

		close(FileChange)

		fmt.Println("out of condition")
		c.JSON(http.StatusOK, gin.H{
			"manifest": &manifestData,
			"rover":    &roverData,
		})
	})

	r.Run(":8080")
}

func ReturnRoverData(rover string) *Rover {
	switch rover {
	case "curiosity":
		return Curiosity
	case "opportunity":
		return Opportunity
	case "spirit":
		return Spirit
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", rover)
		return emptyRover
	}
}
