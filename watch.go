package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

func WatchFile(path string, data interface{}) {
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
	Check(err)
	return data
}
