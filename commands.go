package main

import (
	"fmt"
)

type ShortCommand byte

func (cmd ShortCommand) Frame() []byte {
	return []byte{
		0x01, // USB Shit.
		StandardStartFlag, byte(cmd), 0 ^ byte(cmd), StopFlag,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}

func (cmd ShortCommand) String() string {
	if name, ok := commandName[byte(cmd)]; ok {
		return name
	} else {
		return fmt.Sprintf("ShortCommand[% 2x]", byte(cmd))
	}
}

const (
	CmdGetStatus       ShortCommand = 0x80
	CmdReset           ShortCommand = 0x81
	CmdGoIdle          ShortCommand = 0x82
	CmdGoHaveID        ShortCommand = 0x83
	CmdGoInUse         ShortCommand = 0x85
	CmdGoFinished      ShortCommand = 0x86
	CmdGoReady         ShortCommand = 0x87
	CmdBadID           ShortCommand = 0x88
	CmdGetVersion      ShortCommand = 0x91
	CmdGetID           ShortCommand = 0x92
	CmdGetUnits        ShortCommand = 0x93
	CmdGetSerial       ShortCommand = 0x94
	CmdGetList         ShortCommand = 0x98
	CmdGetUtilization  ShortCommand = 0x99
	CmdGetMotorCurrent ShortCommand = 0x9A
	CmdGetOdometer     ShortCommand = 0x9B
	CmdGetErrorCode    ShortCommand = 0x9C
	CmdGetServiceCode  ShortCommand = 0x9D
	CmdGetUserCfg1     ShortCommand = 0x9E
	CmdGetUserCfg2     ShortCommand = 0x9F
	CmdGetTWork        ShortCommand = 0xA0
	CmdGetHorizontal   ShortCommand = 0xA1
	CmdGetVertical     ShortCommand = 0xA2
	CmdGetCalories     ShortCommand = 0xA3
	CmdGetProgram      ShortCommand = 0xA4
	CmdGetSpeed        ShortCommand = 0xA5
	CmdGetPace         ShortCommand = 0xA6
	CmdGetCadence      ShortCommand = 0xA7
	CmdGetGrade        ShortCommand = 0xA8
	CmdGetGear         ShortCommand = 0xA9
	CmdGetUpList       ShortCommand = 0xAA
	CmdGetUserInfo     ShortCommand = 0xAB
	CmdGetTorque       ShortCommand = 0xAC
	CmdGetHRCur        ShortCommand = 0xB0
	CmdGetHRTZone      ShortCommand = 0xB2
	CmdGetMETS         ShortCommand = 0xB3
	CmdGetPower        ShortCommand = 0xB4
	CmdGetHRAvg        ShortCommand = 0xB5
	CmdGetHRMax        ShortCommand = 0xB6
	CmdGetUserData1    ShortCommand = 0xBE
	CmdGetUserData2    ShortCommand = 0xBF
	CmdGetAudioChannel ShortCommand = 0xC0
	CmdGetAudioVolume  ShortCommand = 0xC1
	CmdGetAudioMute    ShortCommand = 0xC2
	CmdDisplayPopup7   ShortCommand = 0xE1
)

const (
	CmdPM3GetDragFactor ShortCommand = 0xC1
)

var (
	commandName = map[byte]string{
		0x80: "GetStatus",
		0x81: "Reset",
		0x82: "GoIdle",
		0x83: "GoHaveID",
		0x85: "GoInUse",
		0x86: "GoFinished",
		0x87: "GoReady",
		0x88: "BadID",
		0x91: "GetVersion",
		0x92: "GetID",
		0x93: "GetUnits",
		0x94: "GetSerial",
		0x98: "GetList",
		0x99: "GetUtilization",
		0x9A: "GetMotorCurrent",
		0x9B: "GetOdometer",
		0x9C: "GetErrorCode",
		0x9D: "GetServiceCode",
		0x9E: "GetUserCfg1",
		0x9F: "GetUserCfg2",
		0xA0: "GetTWork",
		0xA1: "GetHorizontal",
		0xA2: "GetVertical",
		0xA3: "GetCalories",
		0xA4: "GetProgram",
		0xA5: "GetSpeed",
		0xA6: "GetPace",
		0xA7: "GetCadence",
		0xA8: "GetGrade",
		0xA9: "GetGear",
		0xAA: "GetUpList",
		0xAB: "GetUserInfo",
		0xAC: "GetTorque",
		0xB0: "GetHRCur",
		0xB2: "GetHRTZone",
		0xB3: "GetMETS",
		0xB4: "GetPower",
		0xB5: "GetHRAvg",
		0xB6: "GetHRMax",
		0xBE: "GetUserData1",
		0xBF: "GetUserData2",
		0xC0: "GetAudioChannel",
		0xC1: "GetAudioVolume",
		0xC2: "GetAudioMute",
		0xE1: "DisplayPopup7",
	}
)
