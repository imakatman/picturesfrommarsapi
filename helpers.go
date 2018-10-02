package main

import (
	"github.com/tidwall/gjson"
	"encoding/json"
)

func Check(e error) {
	if e != nil {
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

		dates.Days = append(dates.Days, date)
	}

	return photosAvailable
}
