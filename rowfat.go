package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/GeertJohan/go.hid"
	"log"
	"strings"
	"time"
)

const (
	vendorID  = 0x17a4
	productID = 0x1
)

const (
	inEndpoint  = 0x83
	outEndpoint = 0x04
)

type hex []byte

func (h hex) String() string {
	fancyBytes := make([]string, len(h))
	for i, b := range []byte(h) {
		fancyBytes[i] = fmt.Sprintf("%2x", b)
	}
	return fmt.Sprintf("[%v]", strings.Join(fancyBytes, " "))
}

type Erg struct {
	*hid.Device
}

func (e *Erg) Frame(f Framer) (ResponseStructure, error) {
	log.Println("Sending frame:", f)
	response, err := e.interact(f.Frame())
	if err != nil {
		return ResponseStructure{}, err
	}

	rs, err := ParseResponse(response)
	log.Printf("%v\n", rs)
	return rs, err
}

func (e *Erg) Status() {
	e.Frame(CmdGetStatus)
}

func (e *Erg) GetID() {
	if rs, err := e.Frame(CmdGetID); err == nil {
		log.Printf(
			"ID: %v",
			string(rs.DataStructures[0].Data[0:5]),
		)
	}
}

func (e *Erg) GetSerial() (string, error) {
	if rs, err := e.Frame(CmdGetSerial); err == nil {
		return string(rs.DataStructures[0].Data[0:9]), nil
	} else {
		return "", err
	}
}

// Get odometer distance in meters
func (e *Erg) GetOdometer() (uint32, error) {
	if rs, err := e.Frame(CmdGetOdometer); err == nil {
		d := bytes.NewBuffer(rs.DataStructures[0].Data)
		var distance uint32
		err := binary.Read(d, binary.LittleEndian, &distance)
		return distance, err
	} else {
		return 0, err
	}
}

func (e *Erg) Work() {
	response, _ := e.Frame(CmdGetTWork)
	log.Printf(
		"%v hours %v minutes %v seconds",
		uint8(response.DataStructures[0].Data[0]),
		uint8(response.DataStructures[0].Data[1]),
		uint8(response.DataStructures[0].Data[2]),
	)
}

func (e *Erg) UserInfo() {
	response, _ := e.Frame(CmdGetUserInfo)

	var weight uint16
	var unit, age, gender uint8

	d := bytes.NewBuffer(response.DataStructures[0].Data)
	binary.Read(d, binary.LittleEndian, &weight)
	binary.Read(d, binary.LittleEndian, &unit)
	binary.Read(d, binary.LittleEndian, &age)
	binary.Read(d, binary.LittleEndian, &gender)

	fmt.Printf(
		"weight %v, unit %v, age %v, gender %v\n",
		weight, unit, age, gender,
	)
}

func (e *Erg) Version() {
	rs, _ := e.Frame(CmdGetVersion)
	response := rs.DataStructures[0].Data
	manufacturerId := uint8(response[0])
	cid := uint8(response[1])
	model := uint8(response[2])
	hwVersion := uint16(response[3] | response[4]<<1)
	swVersion := uint16(response[5] | response[6]<<1)
	fmt.Printf(
		`Manufacturer ID: %v
            CID: %v
          Model: %v
     HW Version: %v
     SW Version: %v
`,
		manufacturerId,
		cid,
		model,
		hwVersion,
		swVersion,
	)
}

func (e *Erg) interact(bs []byte) ([]byte, error) {
	if err := e.write(bs); err != nil {
		return nil, err
	}

	time.Sleep(2 * time.Millisecond)
	return e.read()
}

func (e *Erg) write(bs []byte) error {
	log.Println("Writing", hex(bs))
	n, err := e.Write(bs)
	if err != nil {
		return err
	} else if n != len(bs) {
		return errors.New("Sent different number of bytes, omg")
	}
	return nil
}

func (e *Erg) read() ([]byte, error) {
	response := make([]byte, 21)
	_, err := e.Read(response)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("   Read", hex(response))
	return response, nil
}

func main() {
	dev, err := hid.Open(vendorID, productID, "")
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		log.Println("Closing device")
		dev.Close()
	}()

	manufacturer, err := dev.ManufacturerString()
	if err != nil {
		log.Fatalln(err)
	}

	product, err := dev.ProductString()
	if err != nil {
		log.Fatalln(err)
	}

	usbSerial, err := dev.SerialNumberString()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(manufacturer, product, usbSerial)

	erg := Erg{dev}
	erg.Status()
	erg.Version()
	erg.GetID()
	if serial, err := erg.GetSerial(); err == nil {
		fmt.Printf("Serial #%v\n", serial)
	}
	erg.Work()
	erg.UserInfo()
	erg.Frame(CmdGetUtilization)
	erg.Frame(CmdReset)
	erg.Frame(CmdPM3GetDragFactor)
	if distance, err := erg.GetOdometer(); err == nil {
		fmt.Printf("Odometer distance: %v\n", fancyDistance(distance))
	}
}

func fancyDistance(distance uint32) string {
	output := ""
	if distance > 1000 {
		output += fmt.Sprintf("%vkm ", distance/1000)
		distance %= 1000
	}
	if distance > 0 {
		output += fmt.Sprintf("%vm", distance)
	}
	return output
}
