// 09-23-2018
// @TODO #7  Figure out why curiosity, opportunity, and spirit files exist... Same data exists in manifest.json
// might be because we need it for ReturnRoversStruct which is for

package main

import (
	"fmt"
	"encoding/json"
	"os"
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

	for _, v := range dataDrawers {
		go func(dd dataDrawer) {
			fmt.Println("dd.name", dd.name)

			if dd.name == "Manifest" {
				//if isFileEmpty(dd.file){
				//
				//}
				bytes, err := SlurpFile(dd.file)
				Check(err)
				json.Unmarshal(bytes, &dd.obj)
				roverSlices := Rovers.AllRovers
				fmt.Println("Rovers", Rovers)
				fmt.Println("roverSlices", roverSlices)
				for _, r := range roverSlices {
					r := ReturnRoverStruct(r.Name)
					bytes, err := json.Marshal(r)
					Check(err)
					json.Unmarshal(bytes, r)
					fmt.Println("for _, v := range roverSlices", v)
				}
				fmt.Println("InitializeData", dd.name, dd.obj)
				// @TODO Might need to write a function that fills the rover files with data in case the files are empty!
				InitAndWatch(dd.file, dd.obj)
			} else {
				fmt.Println("InitializeData", dd.name, dd.obj)
				InitAndWatch(dd.file, dd.obj)
			}
		}(v)
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

func isFileEmpty(path string) bool{
	fi, e := os.Stat(path)
	if e != nil {
		panic(e)
	}
	// get the size
	if fi.Size() != 0{
		return false
	} else {
		return true
	}
}