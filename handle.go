package main

import (
	"github.com/gin-gonic/gin"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

func initVariables(rover string) (*Manifest, *Rover, *Dates) {
	// Return data that are in variables
	manifestData := &Rovers

	// ReturnRoverPictures returns a pointer to Pictures
	//picturesData := ReturnRoverPicturesStruct(rover)
	roverData := ReturnRoverStruct(rover)
	datesData := ReturnRoverDatesStruct(rover)

	return manifestData, roverData, datesData
}

func HandleRoverGet(c *gin.Context) {
	roverParam := c.Param("rover")
	manifestData, roverData, datesData := initVariables(roverParam)

	if c.Query("update") == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data": *datesData,
		})
	} else if c.Query("update") == "true" {
		if roverData.MaxDate == c.Query("date") {
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
			 initialize manifestData again
			*/
			json.Unmarshal(manifestBytes, manifestData)

			picturesBody, pResponseReceived, picsApiErr := ReturnLatestRoverPictures(roverParam, roverData.MaxSol)
			<-pResponseReceived
			if picsApiErr != nil {
				c.Status(http.StatusServiceUnavailable)
			}

			picturesBytes, picturesReaderErr := ioutil.ReadAll(picturesBody)
			Check(picturesReaderErr)

			datesData.AddDate(picturesBytes)

			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
				"data": *datesData,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusNotAcceptable,
			})
		}
	}
}
