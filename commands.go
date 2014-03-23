package main

import (
	"fmt"
)

type ShortCommand byte

func (cmd ShortCommand) Bytes() []byte {
	return []byte{byte(cmd)}
}

func (cmd ShortCommand) Len() int {
	return 1
}

func (cmd ShortCommand) String() string {
	switch cmd {
	case CmdGetStatus:
		return "GetStatus"
	default:
		return fmt.Sprintf("ShortCommand 0x%x", byte(cmd))
	}
}

const (
	CmdGetStatus ShortCommand = 0x80
	CmdReset     ShortCommand = 0x81
	CmdGoIdle    ShortCommand = 0x82
	CmdGoHaveID  ShortCommand = 0x83
)
