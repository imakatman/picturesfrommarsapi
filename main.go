// 09-10-2018
// @TODO #3: Create interface for manifest and rover?
// @TODO #5: In init(), write a function to just slurp and unmarshall Rover structs

package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
)

var Rovers Manifest
var Curiosity Rover
var CuriosityPictures Pictures
var Opportunity Rover
var OpportunityPictures Pictures
var Spirit Rover
var SpiritPictures Pictures
var emptyRover Rover
var emptyRoverPictures Pictures

var FileChange chan bool

func init() {
	InitializeData()
}

func main() {
	r := gin.Default()

	FileChange = make(chan bool)

	r.GET("/manifest", func(c *gin.Context) {
		fmt.Println("get request for manifest")

		c.JSON(http.StatusOK, Rovers)
	})

	r.GET("/rover/:rover", handleRoverGet)

	fmt.Println("right before run")

	r.Run()
}
