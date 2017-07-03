package arduino

import (
	"go.bug.st/serial.v1"
	"bufio"
	"fmt"
	"time"
	"log"
	"strings"
	"strconv"
)
type Arduino struct {
	port serial.Port
}
type Rele struct {
	Name string
	HourStart int
	HourEnd int
	IsOn bool
	Pin int
}
type ArduinoData struct {
	DateTime time.Time
	Temperature float64
	Humidity float64
	IsManualHandling bool
	ReleList []Rele
}

func (c *Arduino) PortList() (portList []string) {
	portList, err  := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func (c *Arduino) Connect(port string) {
	connection, err := serial.Open(port, &serial.Mode{})
	if err != nil {
		log.Fatal(err)
	}
	c.port = connection
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	if err := c.port.SetMode(mode); err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)
}
func (c *Arduino) WriteData(data string) {
	n, err := c.port.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
}
func (c *Arduino) ReadData(shouldReturn bool) string {
	reader := bufio.NewReader(c.port)
	for {
		reply, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		cleanReply :=  strings.TrimSpace(string(reply))
		parsedReply := strings.Split(cleanReply,";")
		fmt.Printf("REPLY: %s\n", cleanReply)
		if parsedReply[0] == "OK" {
			return strings.Join(append(parsedReply[:0], parsedReply[1:]...), ";")
		}
	}
}
func (c *Arduino) GetData() (resData ArduinoData) {
	c.WriteData("get_data\n")
	stringData := c.ReadData(true)
	fmt.Println(stringData)
	parsedStringData := strings.Split(stringData,";")
	i, err := strconv.ParseInt(parsedStringData[0], 10, 64)
	if err != nil {
		panic(err)
	}
	resData.DateTime = time.Unix(i, 0)
	temp , err := strconv.ParseFloat(parsedStringData[1], 64)
	if err != nil {
		panic(err)
	}
	resData.Temperature = temp
	humd , err := strconv.ParseFloat(parsedStringData[2], 64)
	if err != nil {
		panic(err)
	}
	resData.Humidity = humd
	isManual , err := strconv.ParseBool(parsedStringData[3])
	if err != nil {
		panic(err)
	}
	resData.IsManualHandling = isManual

	for i := 4; i < len(parsedStringData); i++ {
		ReleParsedData := strings.Split(parsedStringData[i], "#")
		if len(ReleParsedData) < 4 {
			continue
		}
		startHour, err := strconv.ParseInt(ReleParsedData[1], 10, 64)
		if err != nil {
			panic(err)
		}
		endHour, err := strconv.ParseInt(ReleParsedData[2], 10, 64)
		if err != nil {
			panic(err)
		}
		isOn, err := strconv.ParseBool(ReleParsedData[3])
		if err != nil {
			panic(err)
		}
		pin, err := strconv.ParseInt(ReleParsedData[4], 10, 64)
		if err != nil {
			panic(err)
		}


		ReleItem := Rele{
			Name: ReleParsedData[0],
			HourStart: int(startHour),
			HourEnd:int(endHour),
			IsOn:isOn,
			Pin: int(pin),
		}
		resData.ReleList = append(resData.ReleList, ReleItem)
	}

	fmt.Println(resData)

	return
}
func (c *Arduino) SetManual(state bool) error {
	stateBit := 0
	if state {
		stateBit = 1
	}
	fmt.Println(fmt.Sprintf("set_manual;%d\n",stateBit ))
	c.WriteData(fmt.Sprintf("set_manual;%d\n",stateBit ))
	stringData := c.ReadData(true)
	fmt.Println(stringData)
	//if stringData != "OK" {
	//	return errors.New(stringData)
	//}
	return nil
}
func (c *Arduino) Switch(switchId string, state bool) error {
	stateBit := 0
	if state {
		stateBit = 1
	}
	fmt.Println(fmt.Sprintf(fmt.Sprintf("switch;%s,%d\n",switchId, stateBit )))
	c.WriteData(fmt.Sprintf("switch;%s;%d\n",switchId, stateBit ))
	stringData := c.ReadData(true)
	fmt.Println(stringData)
	//if stringData != "OK" {
	//	return errors.New(stringData)
	//}
	return nil
}
