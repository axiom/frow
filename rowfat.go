package main

import (
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
	e.Frame(CmdGetID)
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
	log.Printf(
		"weight %v %v, age %v, gender %v",
		uint16(response.DataStructures[0].Data[0]|response.DataStructures[0].Data[1]<<1),
		uint8(response.DataStructures[0].Data[2]),
		uint8(response.DataStructures[0].Data[3]),
		uint8(response.DataStructures[0].Data[4]),
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

	time.Sleep(25 * time.Millisecond)
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
	defer dev.Close()

	manufacturer, err := dev.ManufacturerString()
	if err != nil {
		log.Fatalln(err)
	}

	product, err := dev.ProductString()
	if err != nil {
		log.Fatalln(err)
	}

	serial, err := dev.SerialNumberString()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(manufacturer, product, serial)

	erg := Erg{dev}
	erg.Status()
	erg.Version()
	erg.GetID()
	erg.Work()
	erg.UserInfo()
	erg.Frame(CmdGetUtilization)
	erg.Frame(CmdReset)
}
