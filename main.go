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

var manifest Manifest

func main() {
	r := gin.Default()

	go func() {
		watchFile("data/manfiest.json")
	}()

	r.GET("/manifest", func(c *gin.Context) {
		fmt.Println("get request for manifest")

		c.JSON(http.StatusOK, manifest)
	})

	r.Run(":8080")

	////r.GET("/:rover", func(c *gin.Context) {
	////
	////})
	////
	////r.GET("/:date", func(c *gin.Context) {
	////	date := c.Param("date")
	////
	////})
	//
}

func watchFile(path string) {
	watcher, err := fsnotify.NewWatcher()
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
				json.Unmarshal(slurpFile(path), &manifest)
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

func slurpFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	check(err)
	return data
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
