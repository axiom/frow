package main

import (
	"bytes"
	"github.com/GeertJohan/go.hid"
	"log"
)

const (
	vendorID  = 0x17a4
	productID = 0x1
)

const (
	inEndpoint  = 0x83
	outEndpoint = 0x04
)

type Erg struct {
	*hid.Device
}

func (e *Erg) Status() {
	e.sendFrame(NewFrame(CmdGetStatus))
}

func (e *Erg) sendFrame(frame Frame) ([]byte, error) {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x01)

	frameBytes, _ := frame.Bytes()
	buffer.Write(frameBytes)

	for i := 21 - buffer.Len(); i > 0; i-- {
		buffer.WriteByte(0x00)
	}

	response, err := e.write(buffer.Bytes())
	return response, err
}

func (e *Erg) write(command []byte) ([]byte, error) {
	log.Printf("Writing %v to device\n", command)
	e.Write(command)
	return e.read()
}

func (e *Erg) read() ([]byte, error) {
	log.Println("Reading from device")
	response := make([]byte, 21)
	n, err := e.Read(response)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Read", n, "bytes response:", response)
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
}
