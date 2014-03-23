package main

import (
	"errors"
	"fmt"
)

type ResponseStructure struct {
	StatusByte     uint8
	DataStructures []DataStructure
}

func (rs ResponseStructure) String() string {
	return fmt.Sprintf("ResponseStructure {StatusByte: %2x []DataStructures: %v}", rs.StatusByte, rs.DataStructures)
}

type DataStructure struct {
	Identifier uint8
	Data       []byte
}

func (ds DataStructure) String() string {
	return fmt.Sprintf("DataStructure {Identifier: %2x Data: %v}", ds.Identifier, hex(ds.Data))
}

// [ 1 f1 81 a0  3  0  0  0 22 f2 a6 f2 71 f2 f2  0  0  0  0  0  0]
func ParseResponse(frame []byte) (ResponseStructure, error) {
	c := 1 // Skip USB byte...
	rs := ResponseStructure{}

	if len(frame) < 4 {
		return rs, errors.New("Frame to small")
	}

	if frame[c] != StandardStartFlag {
		return rs, errors.New("Bad start of frame")
	}
	c++

	var runningChecksum byte
	runningChecksum ^= frame[c]

	rs.StatusByte = uint8(frame[c])
	c++

	var checksum byte

	rs.DataStructures = make([]DataStructure, 0)
	for frame[c] != StopFlag {
		// Unstuff
		if frame[c] == StuffFlag {
			c++
			frame[c] = (0xff & (0xff << 2)) | frame[c]
		}

		if frame[c+1] == StopFlag {
			checksum = frame[c]
			break
		}

		ds := DataStructure{}

		runningChecksum ^= frame[c]
		ds.Identifier = uint8(frame[c])
		c++

		runningChecksum ^= frame[c]
		dataByteCount := uint8(frame[c])
		c++

		if dataByteCount > 0 {
			ds.Data = frame[c : c+int(dataByteCount)]
			for _, b := range ds.Data {
				runningChecksum ^= b
			}
			c += int(dataByteCount)
		}

		rs.DataStructures = append(rs.DataStructures, ds)
	}

	if runningChecksum != checksum {
		return rs, errors.New("Checksum did not match")
	}

	return rs, nil
}
