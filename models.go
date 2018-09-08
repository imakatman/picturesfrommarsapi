package main

type Rover struct {
	Id          float64 `json: "id"`
	Name        string  `json: "name"`
	LandingDate string  `json: "landing_date"`
	LaunchDate  string  `json: "launch_date"`
	Status      string  `json: "status"`
	MaxSol      float64 `json: "max_sol"`
	MaxDate     string  `json: "max_date"`
	TotalPhotos float64 `json: "total_photos"`
	//Camera *Camera
}

type Manifest struct {
	Data struct {
		Rovers []Rover `json:"rovers"`
	} `json:"data"`
}
