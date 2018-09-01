// 8-24-2018
// Problem at the moment is that the file change event that is
// triggered when the json file itself is modified is RENAME.
//
// This in particular isn't the issue, the method of streaming the
// new data from the NASA API is going to dictate which file change
// events occur. And the method of choice will have to depend on
// which method is the most performant and cost-effective for the server.

package main

import (
	"fmt"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
	"github.com/howeyc/fsnotify"
	"encoding/json"
)

type Config struct {
	url   string
	token []string
}

type Rover struct {
	Id          float64 `json: "id"`
	Name        string  `json: "name"`
	LandingDate string  `json: "landing_date"`
	LaunchDate  string  `json: "launch_date"`
	Status      string  `json: "status"`
	MaxSol      float64 `json: "max_sol"`
	MaxDate     string  `json: "max_date"`
	TotalPhotos float64 `json: "total_photos"`
	//Camera *Camera
}

type Manifest struct {
	Data struct {
		Rovers []Rover `json:"rovers"`
	} `json:"data"`
}

var manifest *Manifest
var curiosity *Rover
var opportunity *Rover
var spirit *Rover

func main() {
	r := gin.Default()

	WatchFile("data/manfiest.json", &manifest)
	WatchFile("data/curiosityDates.json", &curiosity)
	WatchFile("data/opportunityDates.json", &opportunity)
	WatchFile("data/spiritDates.json", &spirit)

	r.GET("/manifest", func(c *gin.Context) {
		fmt.Println("get request for manifest")

		c.JSON(http.StatusOK, manifest)
	})

	r.Run(":8080")

	r.GET("/rover/:rover", func(c *gin.Context) {
		updateData := c.Query("update")

		if updateData == "" {
			ReturnRoverData(c)
		} else {
			roverParam := c.Param("rover")

			manifestReader := GetLatestManifest(c)
			//GetLatestRoverData(roverParam)
		}
	})
}

func WatchFile(path string, data interface{}) {
	watcher, err := fsnotify.NewWatcher(
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				fmt.Println("event:", ev)
				json.Unmarshal(SlurpFile("data/manfiest.json"), data)
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(path)
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
}

func SlurpFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	check(err)
	return data
}

func ReturnRoverData(c *gin.Context) {
	roverParam := c.Param("rover")

	switch roverParam {
	case "curiosity":
		c.JSON(http.StatusOK, &curiosity)
	case "opportunity":
		c.JSON(http.StatusOK, &opportunity)
	case "spirit":
		c.JSON(http.StatusOK, &spirit)
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", roverParam)
	}
}

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
