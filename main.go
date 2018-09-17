// 09-10-2018
// @TODO #3: Create interface for manifest and rover?

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

	fmt.Println("right before rover route line")

	r.GET("/rover/:rover", HandleRoverGet)

	r.Run()
}
