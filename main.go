// 09-10-2018
// @TODO #3: Create interface for manifest and rover?

package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"os"
	"log"
	"time"
)

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

var FileChange chan bool

func init() {
	// Create a logger
	f, err := os.Create("log.txt")
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(f)

	gin.SetMode(gin.DebugMode)
	InitializeData()
}

func main() {
	r := gin.Default()

	FileChange = make(chan bool)

	r.GET("/manifest", func(c *gin.Context) {
		fmt.Println("get request for manifest")

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
	now := time.Now().Format(time.RFC850)

	log.Println("====================")
	log.Println(fmt.Sprintf("%v // %s", now, msg))
}
