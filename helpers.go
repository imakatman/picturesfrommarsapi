package main

import (
	"github.com/tidwall/gjson"
	"encoding/json"
	"fmt"
)

func Check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

func (dates *Dates) AddDate(bytes []byte) bool{
	photosAvailable := len(gjson.GetBytes(bytes, "photos").Array()) > 0

	if photosAvailable {
		results := gjson.GetManyBytes(
			bytes,
			"photos.0.sol",
			"photos.0.earth_date",
			"photos.0.rover.name",
			"photos.0.rover.id",
		)

		var pictures Pictures
		json.Unmarshal(bytes, &pictures)

		date := Date{
			results[0].Num,
			MiniRover{
				results[2].Str,
				results[3].Num,
			},
			results[1].Str,
			pictures,
		}

		fmt.Println(len(date.Pictures.Day))

		// In order to prepent date into datesStruct.Days, you have to append to a slice of the new Date
		// the original Days slice.
		dateSlice := []Date{date}
		dates.Days = append(dateSlice, dates.Days...)

		//fmt.Println("============AddDate============")
		//fmt.Println(dates)
	}

	return photosAvailable
}

func (dates *Dates) FirstOutLastIn(bytes []byte) bool{
	photosAvailable := dates.AddDate(bytes)

	dates.Days = dates.Days[:8]

	return photosAvailable
}
