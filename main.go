package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"fmt"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"encoding/json"
)

// If release mode is to be set, it needs to be set when building the program
//var release string

//type Camera struct {
//	Name     string
//	FullName string
//}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func slurpFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	check(err)
	return data
}

type slurpFn func(string) []byte

func watchFile(method slurpFn, path string) {
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
				method(path)
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

func main() {
	go func() {
		data := make(chan []byte)
		watchFile(slurpFile, "/data/manfiest.json")
	}()

	r := gin.Default()

	r.GET("/manifest", func(c *gin.Context) {
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

		var rover Rover

		json.Unmarshal(bs, &rover)
	})

	////r.GET("/:rover", func(c *gin.Context) {
	////
	////})
	////
	////r.GET("/:date", func(c *gin.Context) {
	////	date := c.Param("date")
	////
	////})
	//
	//r.Run(":8080")
}
