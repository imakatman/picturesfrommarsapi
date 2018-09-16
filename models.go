// 9/14/2018
// @TODO #4: Create a method of adding latest rover data into _rover_.json file

package main

type Config struct {
	url   string
	token []string
}

type Manifest struct {
	Rovers []Rover `json:"rovers"`
}

type Rover struct {
	Id          int16        `json:"id"` // there can't be space between colon and "
	Name        string       `json:"name"`
	LandingDate string       `json:"landing_date"`
	LaunchDate  string       `json:"launch_date"`
	Status      string       `json:"status"`
	MaxSol      int16        `json:"max_sol"`
	MaxDate     string       `json:"max_date"`
	TotalPhotos int64        `json:"total_photos"`
	AllCameras  []AllCameras `json:"cameras"`
}

type AllCameras struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type Pictures struct {
	Rover struct {
		Name string `json:"name"`
		Id   int8   `json:"id"`
	} `json:"rover"`
	Day     string   `json:"earth_date"`
	Camera  []Camera `json:"camera"`
	Picture Picture
}

type Camera struct {
	Id       int16  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	RoverId  int16  `json:"rover_id"`
}

type Picture struct {
	ImgSrc string `json:"img_src"`
	Id     int16  `json:"id"`
}
