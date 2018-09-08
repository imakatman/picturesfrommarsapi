package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
)

func GetLatestManifest(c *gin.Context) interface{}{
	apiConfig := Config{
		url: "https://api.nasa.gov/mars-photos/api/v1/rovers",
		token: []string{
			"8m8bkcVYqxE5j0vQL2wk1bpiBGibgaqCrOvwZVyU",
			"a4q0jhngYKp9kn0cuwvKMHtKz7IrkKtFBRaiMv5t",
			"ef0eRn0aLh0Byb8q7wCniHbiqcjfDWITSIJVh9xy",
		},
	}

	apiUrl := fmt.Sprintf("%s?api_key=%s", apiConfig.url, apiConfig.token)
	response, err := http.Get(apiUrl)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="gopher.png"`,
	}

	return reader
}

func WriteToFile(){

}
