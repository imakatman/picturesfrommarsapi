// 09-23-2018
// @TODO #7  Figure out why curiosity, opportunity, and spirit files exist... Same data exists in manifest.json
// might be because we need it for ReturnRoversStruct which is for
// @TODO #7  Figure out if the files are necessary
// we dont. instead if you want printed records of what sin the variables, print them in loggers

package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
)

// InitializeData converts the data that exists in the JSON files in /data
// into variables of the Manifest, Rover, and Pictures structs. Then it will
// launch WatchFile() functions for each file to keep track of whether the files
// are modified.
func InitializeData() {
	//launched := make(chan bool)

	var rovers []string

	didInitData := make(chan float64, 3)

	// Manifest go routine
	go func() {
		var bytes []byte

		reader, responseReceived, err := ReturnLatestManifest()
		<-responseReceived
		Check(err)

		bytes, readerErr := ioutil.ReadAll(reader)
		fmt.Println(string(bytes))
		Check(readerErr)

		// Unmarshal the manifest bytes into Rovers variable
		json.Unmarshal(bytes, &Rovers)

		fmt.Println("Rovers", Rovers)

		rovers = make([]string, 0, len(Rovers.AllRovers))
		// Range over each slice in the AllRovers field in the Rovers struct variable
		// Each slice of data is a Rover struct
		for _, r := range Rovers.AllRovers {
			rovers = append(rovers, r.Name)
			fmt.Println(rovers, len(rovers))
			// Set the data in the slice as the value of the empty rover variable
			roverStruct := ReturnRoverStruct(r.Name)
			*roverStruct = r
		}
		didInitData <- 0

	}()

	// Go routine for the rovers
	// They run sequentially when the channel, didInitData returns a value
	for i := range didInitData {
		// If the index is the last index of the dataDrawers slice, close didInitData and exit out of for loop
		if i == float64(len(rovers))-1 {
			close(didInitData)
			return
		}
		go func(rover string) {
			roverData := ReturnRoverStruct(rover)
			datesStruct := ReturnRoverDatesStruct(rover)
			//picturesStruct := ReturnRoverPicturesStruct(rover)
			var x float64
			for x = 0; x < 10; x++ {
				// Make API request to grab latest rover pictures data
				// @TODO: Figure out how to handle api error during initialization
				sol := roverData.MaxSol - x
				reader, responseReceived, err := ReturnLatestRoverPictures(rover, sol)
				<-responseReceived
				Check(err)

				bytes, picturesReaderErr := ioutil.ReadAll(reader)
				Check(picturesReaderErr)
				datesStruct.AddDate(bytes)
			}

			// Unmarshall the returned data into the rovers pictures struct
			didInitData <- i + 1
		}(rovers[int(i)])
	}
}