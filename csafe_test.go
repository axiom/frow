package main

import "testing"

func TestFrame(t *testing.T) {
	testCases := []struct {
		contents FrameContents
		bytes    []byte
	}{
		{CmdGetStatus, []byte{0xf1, 0x80, 0x80, 0xf2}},
		{rawContents([]byte{0xf3}), []byte{0xf1, 0xf3, 0x03, 0xf3, 0x03, 0xf2}},
	}

	for _, testCase := range testCases {

		frame := NewFrame(testCase.contents)
		bytes, err := frame.Bytes()

		if err != nil {
			t.Fatal("Error getting bytes for frame", err)
		}

		if len(bytes) != len(testCase.bytes) {
			t.Errorf("Wrong length, %v want %v", len(bytes), len(testCase.bytes))
		} else {
			for i, b := range []byte(bytes) {
				if b != testCase.bytes[i] {
					t.Errorf("Bad frame contents at %v: 0x%x, want 0x%x", i, b, testCase.bytes[i])
				}
			}
		}
	}
}
