package main

import (
	"github.com/gin-gonic/gin"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"fmt"
)

var roverData *Rover
var datesData *Dates

func initVariables(rover string) {
	// ReturnRoverPictures returns a pointer to Pictures
	//picturesData = ReturnRoverPicturesStruct(rover)
	roverData = ReturnRoverStruct(rover)
	datesData = ReturnRoverDatesStruct(rover)

	return
}

func HandleRoverGet(c *gin.Context) {
	roverParam := c.Param("rover")
	initVariables(roverParam)

	if c.Query("update") == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   *datesData,
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

		/*
		 now that we've received latest data from nasa api
		 re-initialize Rovers again
		*/
		json.Unmarshal(manifestBytes, &Rovers)

		fmt.Println(c.Query("date"))

		//if roverData.MaxDate == c.Query("date") {

		picturesBody, pResponseReceived, picsApiErr := roverData.ReturnRoverPicturesFromApi(roverData.MaxSol)
		<-pResponseReceived
		if picsApiErr != nil {
			c.Status(http.StatusServiceUnavailable)
		}

		picturesBytes, picturesReaderErr := ioutil.ReadAll(picturesBody)
		Check(picturesReaderErr)

		photosAvailable := datesData.FirstOutLastIn(picturesBytes)

		//fmt.Println("============HandleRoverGet============")
		fmt.Println(photosAvailable)

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   *datesData,
		})
		//} else {
		//	message := fmt.Sprintf("Photos for %s is not avaialable", c.Query("date"))
		//	c.JSON(http.StatusOK, gin.H{
		//		"status": http.StatusNotAcceptable,
		//		"message": message,
		//	})
		//}
	}
}
