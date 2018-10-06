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

func (dates *Dates) ParseDate(bytes []byte) (*Date, bool){
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

		//fmt.Println("============AddDate============")
		//fmt.Println(dates)

		return &date, photosAvailable
	}

	return nil, photosAvailable
}

func (dates *Dates) FirstOutLastIn(bytes []byte) bool{
	fmt.Println("dates", dates)

	date, photosAvailable := dates.ParseDate(bytes)

	dates.Days = append(dates.Days, *date)

	dates.Days = dates.Days[1:]

	return photosAvailable
}
