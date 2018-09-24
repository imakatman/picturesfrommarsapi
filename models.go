package main

import (
	"log"
	"fmt"
)

type Config struct {
	url   string
	token []string
}

type Manifest struct {
	AllRovers []Rover `json:"rovers"`
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
	Picture []Picture `json:"photos"`
}

type Camera struct {
	Id       int16  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	RoverId  int16  `json:"rover_id"`
}

type Picture struct {
	Rover struct {
		Name string `json:"name"`
		Id   int8   `json:"id"`
	} `json:"rover"`
	Day     string   `json:"earth_date"`
	Camera  []Camera `json:"camera"`
	ImgSrc string `json:"img_src"`
	Id     int16  `json:"id"`
}

func ReturnRoverStruct(rover string) *Rover{
	switch rover {
	case "curiosity", "Curiosity":
		fmt.Println("ReturnRoverPicturesStruct Curiosity", &CuriosityPictures)
		return &Curiosity
	case "opportunity", "Opportunity":
		fmt.Println("ReturnRoverPicturesStruct Opportunity", &OpportunityPictures)
		return &Opportunity
	case "spirit", "Spirit":
		fmt.Println("ReturnRoverPicturesStruct Spirit", &SpiritPictures)
		return &Spirit
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", rover)
		return &emptyRover
	}
}

func ReturnRoverPicturesStruct(rover string) *Pictures{
	switch rover {
	case "curiosity", "Curiosity":
		fmt.Println("ReturnRoverPicturesStruct Curiosity", &CuriosityPictures)
		return &CuriosityPictures
	case "opportunity", "Opportunity":
		fmt.Println("ReturnRoverPicturesStruct Opportunity", &OpportunityPictures)
		return &OpportunityPictures
	case "spirit", "Spirit":
		fmt.Println("ReturnRoverPicturesStruct Spirit", &SpiritPictures)
		return &SpiritPictures
	default:
		log.Println("Rover parameter provided was not of an expected kind: ", rover)
		return &emptyRoverPictures
	}
}

