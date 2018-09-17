package main

import (
	"fmt"
	"encoding/json"
)

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
					}
				}
				InitAndWatch(v.file, &v.obj, launched)
			}()
		}
	}
}

func InitAndWatch(path string, obj *interface{}, c chan bool) {
	fmt.Println("InitAndWatch", *obj)

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
	//fmt.Println("*obj", *obj)
	c <- true
	// obj already is memory address so just pass it regularly
	WatchFile(path, obj)
}