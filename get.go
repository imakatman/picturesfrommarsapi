package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"io/ioutil"
)

var apiConfig = Config{
	url: "https://api.nasa.gov/mars-photos/api/v1/rovers",
	token: []string{
		"8m8bkcVYqxE5j0vQL2wk1bpiBGibgaqCrOvwZVyU",
		"a4q0jhngYKp9kn0cuwvKMHtKz7IrkKtFBRaiMv5t",
		"ef0eRn0aLh0Byb8q7wCniHbiqcjfDWITSIJVh9xy",
	},
}

func ReturnLatestManifest(c *gin.Context) []byte {
	// @TODO: Write a function that returns a different token if the one in use is invalid
	apiUrl := fmt.Sprintf("%s?api_key=%s", apiConfig.url, apiConfig.token[0])
	response, err := http.Get(apiUrl)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		fmt.Println(err)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println(string(bodyBytes))

	return bodyBytes
}

func ReturnLatestRoverData(c *gin.Context) []byte {
	// @TODO: Write a function that returns a different token if the one in use is invalid
	apiUrl := fmt.Sprintf("%s/%s/photos?sol=%d?api_key=%s", apiConfig.url, apiConfig.token[0])
	response, err := http.Get(apiUrl)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		fmt.Println(err)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println(string(bodyBytes))

	return bodyBytes
}