package main

import (
	"fmt"
	"encoding/json"
)

// InitializeData converts the data that exists in the JSON files in /data
// into variables of the Manifest, Rover, and Pictures structs. Then it will
// launch WatchFile() functions for each file to keep track of whether the files
// are modified.
func InitializeData() {
	launched := make(chan bool)

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

	for i, v := range dataDrawers {
		 if i == 0 {
			go func() {
				fmt.Println(i)
				InitAndWatch(v.file, &v.obj, launched)
			}()
		}

		// This condition will only run if the previous condition
		// has run. With that assumption, in the first invocation of
		// this condition, range over the Rovers variable and Unmarshal
		// the data into the appropriate Rover variables.
		if <-launched == true {
			go func() {
				fmt.Println(i)
				if i == 1 {
					roverSlices := Rovers.AllRovers
					for _, v := range roverSlices{
						r := ReturnRoverStruct(v.Name)
						bytes, err := json.Marshal(v)
						Check(err)
						json.Unmarshal(bytes, r)
						fmt.Println(r)
					}
				}
				InitAndWatch(v.file, &v.obj, launched)
			}()
		}
	}
}

func InitAndWatch(path string, obj *interface{}, c chan bool) {
	bytes, err := SlurpFile(path)

	if err != nil {
		switch err.(type) {
		case *emptyFileErr:
			fmt.Println("api call should be made")
			// make api call and create file
		default:
			// default behavior should be to try and run this function again
			fmt.Println("default behavior for err switch in InitAndWatch")
			InitAndWatch(path, obj, c)
		}
	}

	json.Unmarshal(bytes, *obj)

	fmt.Println("InitAndWatch", *obj)
	c <- true
	// obj already is memory address so just pass it regularly
	WatchFile(path, obj)
}