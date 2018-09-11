package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

func InitAndWatch(path string, obj interface{}){
	json.Unmarshal(SlurpFile(path), obj)
	watchFile(path, obj)
}

func watchFile(path string, obj interface{}) {
	fmt.Println("watchFile")
	bytesFromFile := SlurpFile(path)

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
				fmt.Println("watcher event", ev)
				fmt.Println("event:", ev)
				json.Unmarshal(bytesFromFile, obj)
				fmt.Println("FileChange <- true")
				FileChange <- true
				close(FileChange)
				//fmt.Sprintf("FileChange is %v", FileChange)
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

func WriteFile(file string, data []byte){
	err := ioutil.WriteFile(file, data, 0644)
	Check(err)
}
