package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"io"
)

var apiConfig = Config{
	url: "https://api.nasa.gov/mars-photos/api/v1/rovers",
	token: []string{
		"8m8bkcVYqxE5j0vQL2wk1bpiBGibgaqCrOvwZVyU",
		"a4q0jhngYKp9kn0cuwvKMHtKz7IrkKtFBRaiMv5t",
		"ef0eRn0aLh0Byb8q7wCniHbiqcjfDWITSIJVh9xy",
	},
}

func ReturnLatestManifestReader(c *gin.Context) io.Reader {
	apiUrl := fmt.Sprintf("%s?api_key=%s", apiConfig.url, apiConfig.token)

	response, err := http.Get(apiUrl)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
	}

	return response.Body
}

func WriteToFile() {

}
