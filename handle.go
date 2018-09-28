package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

/*
 Declaration of pointers to variables that
 are to be returned in JSON response
*/
var picturesData *Pictures
var manifestData *Manifest

func initVariables(rover string) {
	// Return data that are in variables
	manifestData = &Rovers
	// ReturnRoverPictures returns a pointer to Pictures
	picturesData = ReturnRoverPicturesStruct(rover)
}

func HandleRoverGet(c *gin.Context) {
	fmt.Println("handleRoverGet")
	roverParam := c.Param("rover")

	if c.Query("update") == "" {
		initVariables(roverParam)
		c.JSON(http.StatusOK, gin.H{
			"manifest": *manifestData,
			"rover":    *picturesData,
		})
	} else if c.Query("update") == "true" {
		// Update Manifest with latest data to use when returning latest rover pictures
		manifestBody, mResponseReceived, manErr := ReturnLatestManifest()
		<-mResponseReceived
		if manErr != nil {
			c.Status(http.StatusServiceUnavailable)
		}
		bodyBytes, readerErr := ioutil.ReadAll(manifestBody)
		Check(readerErr)

		fmt.Println("ReturnLatestManifest was successful")
		WriteFile("data/manifest.json", bodyBytes)
		// data/manifest.json has finished being updated
		<-FileChange

		fmt.Println("first channel came back")
		fmt.Println(" //////////////////////////// ")
		fmt.Println("roverParam", roverParam)
		pictureFile := fmt.Sprintf("data/%sPictures.json", roverParam)
		/*
		 now that we've received latest data from nasa api
		 initialize manifestData again
		*/
		manifestData = &Rovers

		picturesStruct := ReturnRoverPicturesStruct(roverParam)
		picturesBody, pResponseReceived, picsApiErr := ReturnLatestRoverPictures(roverParam)
		<-pResponseReceived
		if picsApiErr != nil {
			c.Status(http.StatusServiceUnavailable)
		}

		picturesBytes, picturesReaderErr := ioutil.ReadAll(picturesBody)
		Check(picturesReaderErr)


		// unmarshal latest pictures data into pictures struct
		fmt.Println("before unmarshall picturesStruct is", picturesStruct)
		fmt.Println(" //////////////////////////// ")
		json.Unmarshal(picturesBytes, picturesStruct)
		fmt.Println("after unmarshall picturesStruct is", picturesStruct)
		fmt.Println(" //////////////////////////// ")
		// update picture file with latest picture data
		WriteFile(pictureFile, picturesBytes)
		<-FileChange

		fmt.Println("second channel came back")
		/*
		 now that pictures file has been updated with
		 latest data, reinitialize picturesData variable
		 with latest data
		*/
		fmt.Println("2 if status == true")
		fmt.Println(" //////////////////////////// ")
		picturesData = picturesStruct

		c.JSON(http.StatusOK, gin.H{
			"data": *picturesData,
		})
	}
}


