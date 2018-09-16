// 09-10-2018
// @TODO #3: Create interface for manifest and rover?
// @TODO #5: In init(), write a function to just slurp and unmarshall Rover structs

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
var CuriosityPictures Pictures
var Opportunity Rover
var OpportunityPictures Pictures
var Spirit Rover
var SpiritPictures Pictures
var emptyRover Rover
var emptyRoverPictures Pictures

var FileChange chan bool

func init() {
	launched := make(chan bool)

	type dataPair struct {
		file string
		obj  interface{}

	}

	dataPairs := []dataPair{
		{"data/manifest.json", &manifest},
		{"data/curiosity.json", &CuriosityPictures},
		{"data/opportunity.json", &OpportunityPictures},
		{"data/spirit.json", &SpiritPictures},
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
		var picturesData *Pictures
		var manifestData *Manifest
		updateData := c.Query("update")
		roverParam := c.Param("rover")

		if updateData == "" {
			fmt.Println("updateData is empty string")
			manifestData = &manifest
			picturesStruct := ReturnRoverPicturesStruct(roverParam)
			// ReturnRoverPictures returns a pointer to Pictures
			picturesData = picturesStruct
		} else if updateData == "true" {
			fmt.Println("updateData == true")
			roverFile := fmt.Sprintf("data/%s.json", roverParam)
			// Update Manifest with latest data to use when returning latest rover pictures
			manifestBytes := ReturnLatestManifest(c)
			WriteFile("data/manifest.json", manifestBytes)
			status := <-FileChange
			if status == true {
				manifestData = &manifest
				picturesStruct := ReturnRoverPicturesStruct(roverParam)
				picturesBytes := ReturnLatestRoverPictures(c)
				json.Unmarshal(picturesBytes, picturesStruct)
				WriteFile(roverFile, picturesBytes)
				status := <-FileChange
				if status == true {
					fmt.Println("2 if status == true")
					picturesData = picturesStruct
				}
			}
		}

		fmt.Println("out of condition")
		c.JSON(http.StatusOK, gin.H{
			"manifest": *manifestData,
			"rover":    *picturesData,
		})
	})

	r.Run(":8080")
}

func ReturnRoverStruct(rover string) *Rover{
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

func ReturnRoverPicturesStruct(rover string) *Pictures{
	switch rover {
	case "curiosity":
		return &CuriosityPictures
	case "opportunity":
		return &OpportunityPictures
	case "spirit":
		return &SpiritPictures
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", rover)
		return &emptyRoverPictures
	}
}
