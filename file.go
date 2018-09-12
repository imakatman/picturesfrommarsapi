// 09-11-2018
// @TODO re:#2: Look into setting up error types. It will be useful to pass values other than messages into the error that is returned from SlurpFile

package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"errors"
)

func InitAndWatch(path string, obj interface{}) {
	bytes, err := SlurpFile(path)
	if err != nil {
		InitAndWatch(path, obj)
		return
	}

	json.Unmarshal(bytes, obj)
	watchFile(path, obj)
}

func watchFile(path string, obj interface{}) {
	fmt.Println("watchFile")
	bytesFromFile, _ := SlurpFile(path)

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

/*
First make sure the file is not empty.
*/
func SlurpFile(path string) ([]byte, error) {
	fi, err := os.Stat(path)
	if err != nil {
		msg := fmt.Sprintf("Could not obtain stat from %s", path)
		// Could not obtain stat, handle error
		panic(msg)
		return nil, errors.New(msg)
	}

	if fi.Size() == 0{

	}

	data, err := ioutil.ReadFile(path)
	Check(err)
	return data, nil
}

func WriteFile(file string, data []byte) {
	err := ioutil.WriteFile(file, data, 0644)
	Check(err)
}
