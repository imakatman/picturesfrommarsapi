package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"encoding/json"
	"net/http"
)

/*
 Declaration of pointers to variables that
 are to be returned in JSON response
*/
var picturesData *Pictures
var manifestData *Manifest

func handleRoverGet(c *gin.Context) {
	roverParam := c.Param("rover")

	if c.Query("update") == "" {
		initVariables(roverParam)
	} else if c.Query("update") == "true" {
		// Update Manifest with latest data to use when returning latest rover pictures
		manifestBytes := ReturnLatestManifest(c)
		WriteFile("data/manifest.json", manifestBytes)
		// data/manifest.json has finished being updated
		status := <-FileChange
		if status == true {
			pictureFile := fmt.Sprintf("data/%sPictures.json", roverParam)
			/*
			 now that we've received latest data from nasa api
			 initialize manifestData again
			*/
			manifestData = &Rovers
			picturesStruct := ReturnRoverPicturesStruct(roverParam)
			picturesBytes := ReturnLatestRoverPictures(c)

			// unmarshal latest pictures data into pictures struct
			json.Unmarshal(picturesBytes, picturesStruct)
			// update picture file with latest picture data
			WriteFile(pictureFile, picturesBytes)
			status := <-FileChange
			if status == true {
				/*
				 now that pictures file has been updated with
				 latest data, reinitialize picturesData variable
				 with latest data
				*/
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
}

func initVariables(rover string) {
	// Return data that are in variables
	manifestData = &Rovers
	// ReturnRoverPictures returns a pointer to Pictures
	picturesData = ReturnRoverPicturesStruct(rover)
}
