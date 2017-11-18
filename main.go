package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Static("/bower_components", "./assets/bower_components")
	router.LoadHTMLGlob("./assets/templates/*")
	router.GET("/", func(c *gin.Context) {
		ArduinoData = ArduinoBoard.GetData()
		fmt.Println(ArduinoData)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":            "Greenhouse",
			"releList":         ArduinoData.ReleList,
			"isManualHandling": ArduinoData.IsManualHandling,
			"releFirstIsOn":    ArduinoData.ReleList[0].IsOn,
			"releSecondIsOn":   ArduinoData.ReleList[1].IsOn,
			"currentTime":      ArduinoData.DateTime.Format("2006-01-02 15:04:05"),
			"temperature":      ArduinoData.Temperature,
			"humidity":         ArduinoData.Humidity,
		})
	})
	router.POST("/switch", func(c *gin.Context) {
		switchId, _ := c.GetPostForm("switchId")
		stateString, _ := c.GetPostForm("state")
		state := false
		if stateString == "true" {
			state = true
		}
		log.Printf("SwitchId: %s, state: %t", switchId, state)
		fmt.Println(ArduinoBoard.Switch(switchId, state))
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
	ArduinoBoard.Connect(os.Getenv("GO-GREENHOUSE_PORT"))
	ArduinoData = ArduinoBoard.GetData()
}
