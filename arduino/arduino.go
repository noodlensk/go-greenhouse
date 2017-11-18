// Package arduino allows you to communicate with Arduino Board via serial port
package arduino

import (
	"bufio"
	"io"
	"strings"
	"time"

	sEncode "github.com/noodlensk/go-greenhouse/encode/serial"
	"github.com/pkg/errors"
	"go.bug.st/serial.v1"
)

const (
	// used for determine end of command/data
	cmdEnd           = '\n'
	cmdArgsDelimiter = ';'
)

type command string

var (
	commandPing     command = "ping"
	commandGetData  command = "get_data"
	commandGetReles command = "get_reles"
)

// Clienter represents communication with arduino
type Clienter interface {
	State() (*Data, error)
	ReleTurnOn(id string) error
	ReleTurnOff(id string) error
}

// Client is Arduino client
type Client struct {
	port io.ReadWriteCloser
}

// Rele is rele config
type Rele struct {
	Name      string
	HourStart int
	HourEnd   int
	IsOn      bool
	Pin       int
}

// UnmarshalSerial setup obj from string
func (r *Rele) UnmarshalSerial(s string) error {
	data := strings.Split(s, "#")
	if len(data) != 5 {
		return errors.New("failed to parse string")
	}

	return nil
}

// Data - data from arduino
type Data struct {
	DateTime         time.Time
	Temperature      float64
	Humidity         float64
	IsManualHandling bool
	//ReleList         []Rele
}

// NewClient return new Arduino client
func NewClient(port string) (Clienter, error) {
	sPort, err := serial.Open(port, &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		return nil, err
	}
	c := &Client{port: sPort}
	return c, c.ping()
}

// State return current state of Arduino
func (c *Client) State() (*Data, error) {
	reply, err := c.call(commandGetData)
	if err != nil {
		return nil, err
	}
	data := &Data{}
	if err := sEncode.Unmarshal(reply, data, cmdArgsDelimiter); err != nil {
		return nil, errors.Wrap(err, "failed to umarshal result")
	}
	return data, nil
}

// Reles return list of rele
func (c *Client) Reles() ([]Rele, error) {
	reply, err := c.call(commandGetReles)
	if err != nil {
		return nil, err
	}
	data := &Data{}
	if err := sEncode.Unmarshal(reply, data, cmdArgsDelimiter); err != nil {
		return nil, errors.Wrap(err, "failed to umarshal result")
	}
	return data, nil
}
func (c *Client) ping() error {
	_, err := c.call(commandPing)
	return err
}

func (c *Client) call(cmd command, args ...string) ([]byte, error) {
	fullCmd := string(cmd)
	if len(args) > 0 {
		fullCmd += string(cmdArgsDelimiter) + strings.Join(args, string(cmdArgsDelimiter))
	}
	_, err := c.port.Write([]byte(fullCmd + string(cmdEnd)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to write to serial")
	}
	res, err := bufio.NewReader(c.port).ReadString(cmdEnd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from serial")
	}
	parsedData := strings.Split(strings.Trim(res, string(cmdEnd)), string(cmdArgsDelimiter))
	if len(parsedData) < 1 {
		return nil, errors.New("failed to parse response")
	}
	if parsedData[0] != "OK" {
		return nil, errors.New("unknown error")
	}
	// hack for []byte convertion
	return []byte(strings.Join(parsedData[1:], string(cmdArgsDelimiter))), nil
}

// ReleTurnOn return current state of Arduino
func (c *Client) ReleTurnOn(id string) error {
	return nil
}

// ReleTurnOff return current state of Arduino
func (c *Client) ReleTurnOff(id string) error {
	return nil
}

/**
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
		cleanReply := strings.TrimSpace(string(reply))
		parsedReply := strings.Split(cleanReply, ";")
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
	parsedStringData := strings.Split(stringData, ";")
	i, err := strconv.ParseInt(parsedStringData[0], 10, 64)
	if err != nil {
		panic(err)
	}
	resData.DateTime = time.Unix(i, 0)
	temp, err := strconv.ParseFloat(parsedStringData[1], 64)
	if err != nil {
		panic(err)
	}
	resData.Temperature = temp
	humd, err := strconv.ParseFloat(parsedStringData[2], 64)
	if err != nil {
		panic(err)
	}
	resData.Humidity = humd
	isManual, err := strconv.ParseBool(parsedStringData[3])
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
			Name:      ReleParsedData[0],
			HourStart: int(startHour),
			HourEnd:   int(endHour),
			IsOn:      isOn,
			Pin:       int(pin),
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
	fmt.Println(fmt.Sprintf("set_manual;%d\n", stateBit))
	c.WriteData(fmt.Sprintf("set_manual;%d\n", stateBit))
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
	fmt.Println(fmt.Sprintf(fmt.Sprintf("switch;%s,%d\n", switchId, stateBit)))
	c.WriteData(fmt.Sprintf("switch;%s;%d\n", switchId, stateBit))
	stringData := c.ReadData(true)
	fmt.Println(stringData)
	//if stringData != "OK" {
	//	return errors.New(stringData)
	//}
	return nil
}
*/
