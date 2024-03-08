package cybergear

import (
	"fmt"
	"slices"
	"testing"
)

func TestFrameEnable(t *testing.T) {
	expected := cgSLCanFrame{}

	copy(expected.header[:], []byte{0x54, 0x30, 0x33, 0x30, 0x30, 0x36, 0x34, 0x37, 0x46, 0x30})

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

	copy(expected.header[:], []byte{0x54, 0x30, 0x34, 0x30, 0x30, 0x36, 0x34, 0x37, 0x46, 0x30})

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

func TestFramePositionMode(t *testing.T) {

	expected := cgSLCanFrame{}

	t.Fatal("WIP")

}

func TestWriteParameterCmd(t *testing.T) {
	t.Fatalf("WIP")

	/*
		Location mode:
		copy(expected.header[:], []byte{0x54, 0x31, 0x32, 0x30, 0x30, 0x36, 0x34, 0x37, 0x46, 0x38})
		copy(expected.data[:], []byte{0x30, 0x35, 0x37, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30})
	*/

	// const (
	// 	OPEARATION_CONTROL_MODE runModeType = 0
	// 	LOCATION_MODE           runModeType = 1
	// 	SPEED_MODE              runModeType = 2
	// 	CURRENT_MODE            runModeType = 3
	// )
}

// func TestFrameSetSpeed(t *testing.T) {
// 	t.Fatal("WIP")
// }

// func TestFrameSpeedMode(t *testing.T) {
// 	t.Fatal("WIP")
// }

// func TestFrameCurrentMode(t *testing.T) {
// 	t.Fatal("WIP")
// }
