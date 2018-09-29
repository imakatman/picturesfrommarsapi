package main

import (
	"io"
	"fmt"
	"net/http"
	"errors"
)

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

func ReturnLatestManifest() (io.Reader, chan bool, error) {
	responseReceived := make(chan bool, 1)
	// @TODO: Write a function that returns a different token if the one in use is invalid
	apiUrl := fmt.Sprintf("%s?api_key=%s", apiConfig.url, apiConfig.token[0])
	response, getErr := http.Get(apiUrl)
	responseReceived <- true
	close(responseReceived)
	if getErr != nil || response.StatusCode != http.StatusOK {
		return nil, responseReceived, errors.New("there was an error with the GET request")
	}

	fmt.Println(apiUrl)

	return response.Body, responseReceived, nil
}

func ReturnLatestRoverPictures(rover string, sol int) (io.Reader, chan bool, error) {
	responseReceived := make(chan bool, 1)

	apiUrl := fmt.Sprintf(
		"%s/%s/photos?api_key=%s&sol=%v",
		apiConfig.url,
		rover,
		apiConfig.token[0],
		sol,
	)

	fmt.Println(apiUrl)

	response, err := http.Get(apiUrl)

	// VERY IMPORTANT TO CLOSE CHANNEL
	responseReceived <- true
	close(responseReceived)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, responseReceived, errors.New("there was an error with the GET request")
	}

	return response.Body, responseReceived, nil
}
