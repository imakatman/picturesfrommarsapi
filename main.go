// 10-02-2018
// @TODO #8 Set a token, so that only requests with a valid token can be made
// @TODO #9 Set up firstInLastOut function
// 10-03-2018
// @TODO #10 Figure out how to optimize the slices, using capacity, for each struct

package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"os"
	"log"
	"time"
)

// Declare variables
var Rovers Manifest
var Curiosity Rover
var CuriosityPictures Pictures
var CuriosityDates Dates
var Opportunity Rover
var OpportunityPictures Pictures
var OpportunityDates Dates
var Spirit Rover
var SpiritPictures Pictures
var SpiritDates Dates
var emptyRover Rover
var emptyRoverPictures Pictures
var emptyRoverDates Dates

var now string

func init() {
	/*
	Set a value to now variable
	*/
	now = time.Now().Format(time.RFC850)

	/*
	Create a logger
	*/
	fname := fmt.Sprintf("%s-log.txt", now)
	f, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(f)

	/*
	Set gin mode to debug
	*/
	gin.SetMode(gin.DebugMode)

	/*
	Initialize:
	- Rovers
	- Curiosity
	- CuriosityPictures
	- CuriosityDates
	- Opportunity
	- OpportunityPictures
	- OpportunityDates
	- Spirit
	- SpiritPictures
	- SpiritDates
	*/
	InitializeData()
}

func main() {
	r := gin.Default()

	r.GET("/manifest", func(c *gin.Context) {
		c.JSON(http.StatusOK, Rovers)
	})

	r.GET("/rover/:rover", HandleRoverGet)

	r.Run()
}

/*
logEntrance runs any time:
1) A GET request for route /:type/:calendar is made
2) A POST request for route /submit/:calendar is made

It prints to log.txt:
1) The IP address that the request is coming from
2) The time that the request was made
3) Description of the request
4) IF there was an error, it will print the error
*/
func writeToLog(c *gin.Context, msg string) {
	log.Println("====================")
	log.Println(fmt.Sprintf("%v // %s", now, msg))
}
