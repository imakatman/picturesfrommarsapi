// @TODO: refactor to reflect new changes made in init.go

package main
//
//import (
//	"github.com/gin-gonic/gin"
//	"fmt"
//	"encoding/json"
//	"net/http"
//	"io/ioutil"
//)
//
///*
// Declaration of pointers to variables that
// are to be returned in JSON response
//*/
//var picturesData *Pictures
//var manifestData *Manifest
//
//func initVariables(rover string) {
//	// Return data that are in variables
//	manifestData = &Rovers
//	// ReturnRoverPictures returns a pointer to Pictures
//	picturesData = ReturnRoverPicturesStruct(rover)
//}
//
//func HandleRoverGet(c *gin.Context) {
//	fmt.Println("handleRoverGet")
//	roverParam := c.Param("rover")
//
//	if c.Query("update") == "" {
//		initVariables(roverParam)
//		c.JSON(http.StatusOK, gin.H{
//			"data": *picturesData,
//		})
//	} else if c.Query("update") == "true" {
//		// Update Manifest with latest data to use when returning latest rover pictures
//		manifestBody, mResponseReceived, manErr := ReturnLatestManifest()
//		<-mResponseReceived
//		if manErr != nil {
//			c.Status(http.StatusServiceUnavailable)
//		}
//
//		manifestBytes, readerErr := ioutil.ReadAll(manifestBody)
//		Check(readerErr)
//
//		fmt.Println("first channel came back")
//		fmt.Println(" //////////////////////////// ")
//		fmt.Println("roverParam", roverParam)
//		/*
//		 now that we've received latest data from nasa api
//		 initialize manifestData again
//		*/
//		json.Unmarshal(manifestBytes, manifestData)
//
//		picturesStruct := ReturnRoverPicturesStruct(roverParam)
//		picturesBody, pResponseReceived, picsApiErr := ReturnLatestRoverPictures(roverParam)
//		<-pResponseReceived
//		if picsApiErr != nil {
//			c.Status(http.StatusServiceUnavailable)
//		}
//
//		parseReader(picturesBody)
//
//		// unmarshal latest pictures data into pictures struct
//		fmt.Println(" //////////////////////////// ")
//		//gjson.Parse(string(picturesBytes).getMany())
//		// picturesStruct = picturesStruct{}
//		json.Unmarshal(picturesBytes, picturesStruct)
//		fmt.Println(" //////////////////////////// ")
//		// update picture file with latest picture data
//
//		fmt.Println("second channel came back")
//		/*
//		 now that pictures file has been updated with
//		 latest data, reinitialize picturesData variable
//		 with latest data
//		*/
//		fmt.Println("2 if status == true")
//		fmt.Println(" //////////////////////////// ")
//		picturesData = picturesStruct
//
//		// days: {
//		//}
//		c.JSON(http.StatusOK, gin.H{
//			"photos": *picturesStruct,
//		})
//	}
//}
