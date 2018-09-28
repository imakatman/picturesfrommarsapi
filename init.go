// 09-23-2018
// @TODO #7  Figure out why curiosity, opportunity, and spirit files exist... Same data exists in manifest.json
// might be because we need it for ReturnRoversStruct which is for
// @TODO #7  Figure out if the files are necessary
// we dont. instead if you want printed records of what sin the variables, print them in loggers

package main

import (
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
)

// InitializeData converts the data that exists in the JSON files in /data
// into variables of the Manifest, Rover, and Pictures structs. Then it will
// launch WatchFile() functions for each file to keep track of whether the files
// are modified.
func InitializeData() {
	//launched := make(chan bool)

	type dataDrawer struct {
		name string
		obj  interface{}
	}

	dataDrawers := []dataDrawer{
		{"manifest", &Rovers},
		{"curiosity", &CuriosityPictures},
		{"opportunity", &OpportunityPictures},
		{"spirit", &SpiritPictures},
	}

	didInitDD := make(chan int, 7)

	// Manifest go routine
	go func(dd dataDrawer) {
		fmt.Println("manifest dd.name", dd.name)
		var bytes []byte

		reader, responseReceived, err := ReturnLatestManifest()
		<-responseReceived
		Check(err)

		bytes, readerErr := ioutil.ReadAll(reader)
		Check(readerErr)

		// Unmarshal the manifest bytes into Rovers variable
		json.Unmarshal(bytes, &Rovers)

		fmt.Println("Rovers", Rovers)

		// Range over each slice in the AllRovers field in the Rovers struct variable
		// Each slice of data is a Rover struct
		for _, r := range Rovers.AllRovers {
			// Set the data in the slice as the value of the empty rover variable
			roverStruct := ReturnRoverStruct(r.Name)
			*roverStruct = r
			fmt.Println(fmt.Sprintf("after unmarshall %s is %v", r.Name, roverStruct))
		}
		didInitDD <- 0

	}(dataDrawers[0])

	// Go routine for the rovers
	// They run sequentially when the channel, didInitDD returns a value
	for i := range didInitDD {
		// If the index is the last index of the dataDrawers slice, close didInitDD and exit out of for loop
		if i == len(dataDrawers)-1{
			close(didInitDD)
			return
		}
		go func(dd dataDrawer) {
			fmt.Println("dd.name", dd.name)

			// Make API request to grab latet rover pictures data
			reader, responseReceived, err := ReturnLatestRoverPictures(dd.name)
			<-responseReceived
			// @TODO: Figure out how to handle api error during initialization
			Check(err)

			bytes, picturesReaderErr := ioutil.ReadAll(reader)
			Check(picturesReaderErr)

			picturesStruct := ReturnRoverPicturesStruct(dd.name)

			json.Unmarshal(bytes, picturesStruct)

			fmt.Println("picturesStruct", picturesStruct)
			didInitDD <- i + 1
		}(dataDrawers[i+1])
	}
}

func InitAndWatch(path string, obj interface{}) {
	bytes, err := SlurpFile(path)

	if err != nil {
		switch err.(type) {
		case *emptyFileErr:
			fmt.Println("api call should be made")
			// make api call and create file
		default:
			// default behavior should be to try and run this function again
			fmt.Println("default behavior for err switch in InitAndWatch")
			InitAndWatch(path, obj)
		}
	}

	json.Unmarshal(bytes, obj)
	// obj already is memory address so just pass it regularly
	WatchFile(path, obj)
}

func isFileEmpty(path string) bool {
	fi, e := os.Stat(path)
	if e != nil {
		panic(e)
	}
	// get the size
	if fi.Size() != 0 {
		return false
	} else {
		return true
	}
}
