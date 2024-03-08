package cybergear

import (
	"fmt"
	"slices"
	"testing"
)

func TestFrameEnable(t *testing.T) {
	expected := cgSLCanFrame{}

	expected.header[0] = 0x54
	expected.header[1] = 0x30
	expected.header[2] = 0x33
	expected.header[3] = 0x30
	expected.header[4] = 0x30
	expected.header[5] = 0x36
	expected.header[6] = 0x34
	expected.header[7] = 0x37
	expected.header[8] = 0x46
	expected.header[9] = 0x30

	frame, err := EnableMotorCmd(0x64, 0x7F)

	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(expected.header[:], frame.header[:]) {
		t.Errorf("Unextected frame header: expected: %+v - actual: %+v", expected.header, frame.header)
	}
}

func TestFrameDisable(t *testing.T) {
	expected := cgSLCanFrame{}

	expected.header[0] = 0x54
	expected.header[1] = 0x30
	expected.header[2] = 0x34
	expected.header[3] = 0x30
	expected.header[4] = 0x30
	expected.header[5] = 0x36
	expected.header[6] = 0x34
	expected.header[7] = 0x37
	expected.header[8] = 0x46
	expected.header[9] = 0x30

	frame, err := DisableMotorCmd(0x64, 0x7F)

	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(expected.header[:], frame.header[:]) {
		t.Errorf("Unextected frame header: expected: %+v - actual: %+v", expected.header, frame.header)
	}
}

func TestFrameSerializeNoPayload(t *testing.T) {
	f := cgSLCanFrame{}
	f.header[9] = 0
	b := f.Serialize()

	if len(b) != 10 {
		t.Fatal(fmt.Errorf("Unexpected length of serialized frame (%d). DLC=0 => only header should be serialized", len(b)))
	}
}

func TestFrameSerializeWithPayload(t *testing.T) {
	f := cgSLCanFrame{}
	var dlc byte = 3
	f.header[9] = dlc
	b := f.Serialize()

	if len(b) != int(10+2*dlc) {
		t.Fatal(fmt.Errorf("Unexpected length of serialized frame (%d). DLC=%d => only header + %d bytes should be serialized", len(b), dlc, 2*dlc))
	}
}

// func TestFrameSetSpeed(t *testing.T) {
// 	t.Fatal("WIP")
// }

// func TestFrameSpeedMode(t *testing.T) {
// 	t.Fatal("WIP")
// }

// func TestFramePositionMode(t *testing.T) {
// 	t.Fatal("WIP")
// }

// func TestFrameCurrentMode(t *testing.T) {
// 	t.Fatal("WIP")
// }
