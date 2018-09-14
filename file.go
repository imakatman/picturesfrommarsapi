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
)

type emptyFileErr struct {
	file string
	err  error
}

func (e *emptyFileErr) Error() string {
	return fmt.Sprintf("%s is empty", e.file)
}

func InitAndWatch(path string, obj interface{}) {
	fmt.Println("InitAndWatch")

	bytes, err := SlurpFile(path)

	fmt.Println(string(bytes))

	if err != nil{
		switch err.(type) {
		case *emptyFileErr:
			fmt.Println("api call should be made")
			// make api call and create file
		default:
			// default behavior should be to try and run this function again
			InitAndWatch(path, obj)
		}
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
	// Could not obtain stat, handle error
	if err != nil {
		return nil, fmt.Errorf("could not obtain stat from %s", path)
	}

	if fi.Size() == 0 {
		msg := fmt.Errorf("%s is empty", path)
		fErr := emptyFileErr{ path, msg}
		return nil, &fErr
	}

	data, err := ioutil.ReadFile(path)
	Check(err)
	return data, nil
}

func WriteFile(file string, data []byte) {
	err := ioutil.WriteFile(file, data, 0644)
	Check(err)
}
