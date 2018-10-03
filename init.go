package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
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
		Check(readerErr)

		// Unmarshal the manifest bytes into Rovers variable
		json.Unmarshal(bytes, &Rovers)

		rovers = make([]string, 0, len(Rovers.AllRovers))
		// Range over each slice in the AllRovers field in the Rovers struct variable
		// Each slice of data is a Rover struct
		for _, r := range Rovers.AllRovers {
			rovers = append(rovers, r.Name)
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
		if i == float64(len(rovers)) {
			close(didInitData)
			return
		}
		go func(rover string) {
			roverData := ReturnRoverStruct(rover)
			datesStruct := ReturnRoverDatesStruct(rover)
			//picturesStruct := ReturnRoverPicturesStruct(rover)
			var x float64

			var tenAvailableDaysAdded int
			for x = 0; tenAvailableDaysAdded != 10; x++ {
				// Make API request to grab latest rover pictures data
				// @TODO: Figure out how to handle api error during initialization
				sol := roverData.MaxSol - x
				reader, responseReceived, err := ReturnLatestRoverPictures(rover, sol)
				<-responseReceived
				Check(err)

				bytes, picturesReaderErr := ioutil.ReadAll(reader)
				Check(picturesReaderErr)

				photosAvailable := datesStruct.AddDate(bytes)

				fmt.Println(photosAvailable)

				if photosAvailable {
					tenAvailableDaysAdded++
				} else {
					date := Date{
						sol,
						MiniRover{
							roverData.Name,
							roverData.Id,
						},
						"",
						Pictures{},
					}

					// In order to prepent date into datesStruct.Days, you have to append to a slice of the new Date
					// the original Days slice.
					dateSlice := []Date{date}
					datesStruct.Days = append(dateSlice, datesStruct.Days...)
				}
			}

			//fmt.Println("============InitializeData============")
			//fmt.Println(datesStruct)

			// Unmarshall the returned data into the rovers pictures struct
			didInitData <- i + 1
		}(rovers[int(i)])
	}
}
