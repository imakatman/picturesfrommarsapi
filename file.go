// 09-11-2018
// @TODO #6: When Manifest file is modified, the new data isn't being properly unmarshaled into the variable. Fix this.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"encoding/json"
	"github.com/howeyc/fsnotify"
)

type emptyFileErr struct {
	file string
	err  error
}

func (e *emptyFileErr) Error() string {
	return fmt.Sprintf("%s is empty", e.file)
}

func WatchFile(path string, obj interface{}) {
	fmt.Println("watchFile", obj)
	bytesFromFile, _ := SlurpFile(path)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		fmt.Println("WatchFile inside gofunc before for loop", path, obj)
		for {
			select {
			case ev := <-watcher.Event:
				fmt.Println(obj)
				fmt.Println("watcher event", ev)
				fmt.Println("event:", ev)
				json.Unmarshal(bytesFromFile, &obj)
				fmt.Println("FileChange <- true")
				fmt.Println(obj)
				FileChange <- true
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
	fmt.Println("WriteFile", file)
	err := ioutil.WriteFile(file, data, 0644)
	Check(err)
}
