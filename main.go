package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"log"
	"hydroponics/pkg/arduino"
	"fmt"
)
var ArduinoBoard = arduino.Arduino{}
func main() {
	ArduinoBoard.Connect()
	router := gin.Default()
	router.Static("/bower_components", "./bower_components")
	router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/", "index.html")
	router.POST("/switch", func(c *gin.Context) {
		switchId, _ := c.GetPostForm("switchId")
		state, _ := c.GetPostForm("state")
		log.Printf("SwitchId: %s, state: %s", switchId, state)
		fmt.Println(ArduinoBoard.GetData())
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})


	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
