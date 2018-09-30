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
var manifestData *Manifest
var picturesData *Pictures
var roverData *Rover
var datesData *Dates

func initVariables(rover string) {
	// Return data that are in variables
	manifestData = &Rovers
	// ReturnRoverPictures returns a pointer to Pictures
	picturesData = ReturnRoverPicturesStruct(rover)
	roverData = ReturnRoverStruct(rover)
	datesData = ReturnRoverDatesStruct(rover)
}

func HandleRoverGet(c *gin.Context) {
	fmt.Println("handleRoverGet")
	roverParam := c.Param("rover")

	if c.Query("update") == "" {
		initVariables(roverParam)
		c.JSON(http.StatusOK, gin.H{
			"data": *picturesData,
		})
	} else if c.Query("update") == "true" {
		// Update Manifest with latest data to use when returning latest rover pictures
		manifestBody, mResponseReceived, manErr := ReturnLatestManifest()
		<-mResponseReceived
		if manErr != nil {
			c.Status(http.StatusServiceUnavailable)
		}

		manifestBytes, readerErr := ioutil.ReadAll(manifestBody)
		Check(readerErr)

		fmt.Println("first channel came back")
		fmt.Println(" //////////////////////////// ")
		fmt.Println("roverParam", roverParam)
		/*
		 now that we've received latest data from nasa api
		 initialize manifestData again
		*/
		json.Unmarshal(manifestBytes, manifestData)

		picturesStruct := ReturnRoverPicturesStruct(roverParam)
		picturesBody, pResponseReceived, picsApiErr := ReturnLatestRoverPictures(roverParam)
		<-pResponseReceived
		if picsApiErr != nil {
			c.Status(http.StatusServiceUnavailable)
		}

		picturesBytes, picturesReaderErr := ioutil.ReadAll(picturesBody)
		Check(picturesReaderErr)

		datesData.AddDate(picturesBytes)

		c.JSON(http.StatusOK, gin.H{
			"data": *datesData,
		})
	}
}
