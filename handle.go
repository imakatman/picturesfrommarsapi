package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"io"
	"errors"
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
		manifestBody, err := ReturnLatestManifest()
		if err != nil {
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
		picturesBytes, picturesReceived := ReturnLatestRoverPictures(c)

		<-picturesReceived
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

var apiConfig = Config{
	url: "https://api.nasa.gov/mars-photos/api/v1/rovers",
	token: []string{
		"8m8bkcVYqxE5j0vQL2wk1bpiBGibgaqCrOvwZVyU",
		"a4q0jhngYKp9kn0cuwvKMHtKz7IrkKtFBRaiMv5t",
		"ef0eRn0aLh0Byb8q7wCniHbiqcjfDWITSIJVh9xy",
	},
}

//type apiError struct {
//	msg       error
//	code      int
//}

func ReturnLatestManifest() (io.Reader, error) {
	// @TODO: Write a function that returns a different token if the one in use is invalid
	apiUrl := fmt.Sprintf("%s?api_key=%s", apiConfig.url, apiConfig.token[0])
	response, getErr := http.Get(apiUrl)

	if getErr != nil || response.StatusCode != http.StatusOK {
		return nil, errors.New("there was an error with the GET request")
	}

	return response.Body, nil
}

// @TODO: Refactor so it doesnt need conext as a paramter
func ReturnLatestRoverPictures(c *gin.Context) ([]byte, chan bool) {
	picturesReceived := make(chan bool, 1)
	roverParam := c.Param("rover")

	fmt.Println("ReturnLatestRoverPictures roverParam", roverParam)
	fmt.Println(" //////////////////////////// ")

	roverStruct := *ReturnRoverStruct(roverParam)

	fmt.Println("ReturnLatestRoverPictures roverStruct", roverStruct)
	fmt.Println(" //////////////////////////// ")
	apiUrl := fmt.Sprintf(
		"%s/%s/photos?api_key=%s&earth_date=%s",
		apiConfig.url,
		roverParam,
		apiConfig.token[0],
		roverStruct.MaxDate,
	)

	fmt.Println("apiUrl", apiUrl)
	fmt.Println(" //////////////////////////// ")
	response, err := http.Get(apiUrl)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		fmt.Println(err)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bodyBytes))
	fmt.Println(" //////////////////////////// ")
	fmt.Println(&picturesReceived)
	fmt.Println(" //////////////////////////// ")

	// VERY IMPORTANT TO CLOSE CHANNEL
	close(picturesReceived)

	return bodyBytes, picturesReceived
}
