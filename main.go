package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"fmt"
	"io/ioutil"
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

func watchFile(path string, channel chan []byte) {
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
				channel <- slurpFile(path)
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
	mfChan := make(chan []byte)

	go func(){
		watchFile("data/manfiest.json", mfChan)
	}()

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

	type Manifest struct{
		Rover1 *Rover
		Rover2 *Rover
		Rover3 *Rover
	}

	type File struct{
		Manifest []*Rover
	}

	var manifest Manifest

	for bs := range mfChan{
		fmt.Println("something came off mfChan")
		json.Unmarshal(bs, &manifest)
		fmt.Println(manifest)
	}

	//r := gin.Default()
	//
	//r.GET("/manifest", func(c *gin.Context) {
	//	type Rover struct {
	//		Id          float64 `json: "id"`
	//		Name        string  `json: "name"`
	//		LandingDate string  `json: "landing_date"`
	//		LaunchDate  string  `json: "launch_date"`
	//		Status      string  `json: "status"`
	//		MaxSol      float64 `json: "max_sol"`
	//		MaxDate     string  `json: "max_date"`
	//		TotalPhotos float64 `json: "total_photos"`
	//		//Camera *Camera
	//	}
	//
	//	var rover Rover
	//
	//	json.Unmarshal(<-mfChan, &rover)
	//})

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
