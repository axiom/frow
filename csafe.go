package main

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	ExtendedStartFlag = 0xf0
	StandardStartFlag = 0xf1
	StopFlag          = 0xf2
	StuffFlag         = 0xf3
)

type FrameContents interface {
	Bytes() []byte
	Len() int
}

type rawContents []byte

func (fc rawContents) Bytes() []byte {
	return []byte(fc)
}

func (fc rawContents) Len() int {
	return len([]byte(fc))
}

func (fc rawContents) String() string {
	fancyBytes := make([]string, len(fc))
	for i, b := range []byte(fc) {
		fancyBytes[i] = fmt.Sprintf("0x%x", b)
	}
	return fmt.Sprintf("[%v]", strings.Join(fancyBytes, " "))
}

type Frame struct {
	contents FrameContents
}

func NewFrame(contents FrameContents) Frame {
	return Frame{contents: contents}
}

func (f Frame) String() string {
	return rawContents(f.contents.Bytes()).String()
}

func Checksum(bytes []byte) byte {
	var checksum byte
	for _, b := range bytes {
		checksum ^= b
	}
	return checksum
}

func (f Frame) Bytes() ([]byte, error) {
	var buffer bytes.Buffer

	if err := buffer.WriteByte(StandardStartFlag); err != nil {
		return nil, err
	}

	for _, b := range f.contents.Bytes() {
		switch b {
		case ExtendedStartFlag:
			fallthrough
		case StandardStartFlag:
			fallthrough
		case StopFlag:
			fallthrough
		case StuffFlag:
			if err := buffer.WriteByte(StuffFlag); err != nil {
				return nil, err
			}
			if err := buffer.WriteByte(0x03 & b); err != nil {
				return nil, err
			}
		default:
			if err := buffer.WriteByte(b); err != nil {
				return nil, err
			}
		}
	}

	switch checksum := Checksum(f.contents.Bytes()); checksum {
	case ExtendedStartFlag:
		fallthrough
	case StandardStartFlag:
		fallthrough
	case StopFlag:
		fallthrough
	case StuffFlag:
		if err := buffer.WriteByte(StuffFlag); err != nil {
			return nil, err
		}
		if err := buffer.WriteByte(0x03 & checksum); err != nil {
			return nil, err
		}
	default:
		if err := buffer.WriteByte(checksum); err != nil {
			return nil, err
		}
	}

	if err := buffer.WriteByte(StopFlag); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
