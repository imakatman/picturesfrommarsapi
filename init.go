// 09-23-2018
// @TODO #7  Figure out why curiosity, opportunity, and spirit files exist... Same data exists in manifest.json
// might be because we need it for ReturnRoversStruct which is for

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
		file string
		obj  interface{}
	}

	dataDrawers := []dataDrawer{
		{"Manifest", "data/manifest.json", &Rovers},
		{"Curiosity", "data/curiosity.json", &Curiosity},
		{"Opportunity", "data/opportunity.json", &Opportunity},
		{"Spirit", "data/spirit.json", &Spirit},
		{"Curiosity's Pictures", "data/curiosityPictures.json", &CuriosityPictures},
		{"Opportunity Pictures", "data/opportunityPictures.json", &OpportunityPictures},
		{"Spirit Pictures", "data/spiritPictures.json", &SpiritPictures},
	}

	didInitalizeDataDrawer := make(chan int, 7)

	go func(dd dataDrawer) {
		fmt.Println("manifest dd.name", dd.name)
		var bytes []byte
		if isFileEmpty(dd.file) {
			// If file is empty, grab latest manifest data from NASA api
			body, err := ReturnLatestManifest()
			Check(err)
			bodyBytes, readerErr := ioutil.ReadAll(body)
			Check(readerErr)
			bytes = bodyBytes
		} else {
			// Slurp the bytes out of the manifest file
			fileBytes, err := SlurpFile(dd.file)
			bytes = fileBytes
			Check(err)
		}

		// Unmarshal the manifest bytes into Rovers variable
		json.Unmarshal(bytes, &Rovers)

		fmt.Println("Rovers", Rovers)

		// Range over the slice of data in the AllRovers field in the Rovers struct variable
		// Each slice of data is a Rover struct
		//@TODO: Finish this. Unmarshalling rover data into the variable
		for _, r := range Rovers.AllRovers {
			fmt.Println("for _, v := range roverSlices", r)
			json.Unmarshal(bytes, &Rovers)
			//roverStruct := ReturnRoverStruct(r.Name)
			//bytes, err := json.Marshal(roverStruct)
			//Check(err)
			//json.Unmarshal(bytes, r)
		}

		InitAndWatch(dd.file, dd.obj)
		didInitalizeDataDrawer <- 0

	}(dataDrawers[0])

	for i := range didInitalizeDataDrawer {
		go func(dd dataDrawer) {
			fmt.Println("dd.name", dd.name)
			r := ReturnRoverStruct(dd.name)
			bytes, err := json.Marshal(r)
			Check(err)
			WriteFile(dd.file, bytes)
			//// @TODO Might need to write a function that fills the rover files with data in case the files are empty!
			fmt.Println("InitializeData", dd.name, dd.obj)
			InitAndWatch(dd.file, dd.obj)
			didInitalizeDataDrawer <- i+1
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
