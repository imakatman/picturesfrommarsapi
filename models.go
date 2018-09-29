package main

import (
	"log"
)

type Config struct {
	url   string
	token []string
}

type Manifest struct {
	AllRovers []Rover `json:"rovers"`
}

type Rover struct {
	Id          int        `json:"id"` // there can't be space between colon and "
	Name        string       `json:"name"`
	LandingDate string       `json:"landing_date"`
	LaunchDate  string       `json:"launch_date"`
	Status      string       `json:"status"`
	MaxSol      int        `json:"max_sol"`
	MaxDate     string       `json:"max_date"`
	TotalPhotos int        `json:"total_photos"`
	AllCameras  []AllCameras `json:"cameras"`
}

type AllCameras struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type Camera struct {
	Id       int  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	RoverId  int  `json:"rover_id"`
}

type Dates struct {
	Date []Date
}

type Date struct{
	Sol int
	MiniRover
	Day     string
	Pictures
}

type MiniRover struct {
	Name string
	Id   int
}

type Pictures struct {
	Day []Picture `json:"photos"`
}

type Picture struct {
	Id     int  `json:"id"`
	Camera  Camera `json:"camera"`
	ImgSrc string `json:"img_src"`
}

func ReturnRoverStruct(rover string) *Rover{
	switch rover {
	case "curiosity", "Curiosity":
		return &Curiosity
	case "opportunity", "Opportunity":
		return &Opportunity
	case "spirit", "Spirit":
		return &Spirit
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", rover)
		return &emptyRover
	}
}

func ReturnRoverPicturesStruct(rover string) *Pictures{
	switch rover {
	case "curiosity", "Curiosity":
		return &CuriosityPictures
	case "opportunity", "Opportunity":
		return &OpportunityPictures
	case "spirit", "Spirit":
		return &SpiritPictures
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", rover)
		return &emptyRoverPictures
	}
}


func ReturnRoverDatesStruct(rover string) *Dates{
	switch rover {
	case "curiosity", "Curiosity":
		return &CuriosityDates
	case "opportunity", "Opportunity":
		return &OpportunityDates
	case "spirit", "Spirit":
		return &SpiritDates
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", rover)
		return &emptyRoverDates
	}
}

