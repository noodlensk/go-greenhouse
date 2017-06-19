package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"log"
	"hydroponics/pkg/arduino"
	"fmt"
	"strings"
	"os"
)
var ArduinoBoard = arduino.Arduino{}
var ArduinoData arduino.ArduinoData
func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Static("/bower_components", "./bower_components")
	router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		fmt.Println(ArduinoData)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Hydroponics",
			"releList": ArduinoData.ReleList,
			"isManualHandling": ArduinoData.IsManualHandling,
			"releFirstIsOn": ArduinoData.ReleList[0].IsOn,
			"releSecondIsOn": ArduinoData.ReleList[1].IsOn,
		})
	})
	router.POST("/switch", func(c *gin.Context) {
		switchId, _ := c.GetPostForm("switchId")
		state, _ := c.GetPostForm("state")
		log.Printf("SwitchId: %s, state: %s", switchId, state)
		fmt.Println(ArduinoBoard.GetData())
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	router.POST("/setManual", func(c *gin.Context) {
		stateString, _ := c.GetPostForm("state")
		state := false
		if stateString == "true" {
			state = true
		}
		err := ArduinoBoard.SetManual(state)
		if err != nil {
			log.Fatal(err)
		}
		ArduinoData = ArduinoBoard.GetData()
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
//go func () {
//	ArduinoData = ArduinoBoard.GetData()
//	time.Sleep(10 * time.Second)
//}()


	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
func init() {
	portList := ArduinoBoard.PortList()
	fmt.Printf("Avaible ports: %s\n", strings.Join(portList, ","))
	ArduinoBoard.Connect(os.Getenv("HYDROPONICS_PORT"))
	ArduinoData = ArduinoBoard.GetData()
}