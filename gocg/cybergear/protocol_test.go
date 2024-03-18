package cybergear

import (
	"encoding/hex"
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

func TestSetSpeedMode(t *testing.T) {
	expected := cgSLCanFrame{}
	copy(expected.header[:], []byte{0x54, 0x31, 0x32, 0x30, 0x30, 0x30, 0x30, 0x37, 0x46, 0x38})
	copy(expected.data[:], []byte{0x30, 0x35, 0x37, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30})

	var hostId byte = 0x00
	var motorId byte = 0x7F

	actual, err := WriteRunMode(hostId, motorId, SPEED_MODE)

	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(expected.header[:], actual.header[:]) {
		t.Fatalf("Wrong header bytes. Expected: %+v, Actual: %+v", expected.header, actual.header)
	}

	if !slices.Equal(expected.data[:], actual.data[:]) {
		t.Fatalf("Wrong data bytes. Expected: %s, Actual: %s", hex.Dump(expected.data[:]), hex.Dump(actual.data[:]))
	}
}

func TestSetOperationControlMode(t *testing.T) {
	expected := cgSLCanFrame{}
	copy(expected.header[:], []byte{0x54, 0x31, 0x32, 0x30, 0x30, 0x30, 0x30, 0x37, 0x46, 0x38})
	copy(expected.data[:], []byte{0x30, 0x35, 0x37, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30})

	var hostId byte = 0x00
	var motorId byte = 0x7F

	actual, err := WriteRunMode(hostId, motorId, OPEARATION_CONTROL_MODE)

	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(expected.header[:], actual.header[:]) {
		t.Fatalf("Wrong header bytes. Expected: %+v, Actual: %+v", expected.header, actual.header)
	}

	if !slices.Equal(expected.data[:], actual.data[:]) {
		t.Fatalf("Wrong data bytes. Expected: %s, Actual: %s", hex.Dump(expected.data[:]), hex.Dump(actual.data[:]))
	}
}

func TestSetLocationMode(t *testing.T) {
	expected := cgSLCanFrame{}
	copy(expected.header[:], []byte{0x54, 0x31, 0x32, 0x30, 0x30, 0x30, 0x30, 0x37, 0x46, 0x38})
	copy(expected.data[:], []byte{0x30, 0x35, 0x37, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30})

	var hostId byte = 0x00
	var motorId byte = 0x7F

	actual, err := WriteRunMode(hostId, motorId, LOCATION_MODE)

	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(expected.header[:], actual.header[:]) {
		t.Fatalf("Wrong header bytes. Expected: %+v, Actual: %+v", expected.header, actual.header)
	}

	if !slices.Equal(expected.data[:], actual.data[:]) {
		t.Fatalf("Wrong data bytes. Expected: %s, Actual: %s", hex.Dump(expected.data[:]), hex.Dump(actual.data[:]))
	}
}

func TestSetCurrentMode(t *testing.T) {
	expected := cgSLCanFrame{}
	copy(expected.header[:], []byte{0x54, 0x31, 0x32, 0x30, 0x30, 0x30, 0x30, 0x37, 0x46, 0x38})
	copy(expected.data[:], []byte{0x30, 0x35, 0x37, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x33, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30})

	var hostId byte = 0x00
	var motorId byte = 0x7F

	actual, err := WriteRunMode(hostId, motorId, CURRENT_MODE)

	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(expected.header[:], actual.header[:]) {
		t.Fatalf("Wrong header bytes. Expected: %+v, Actual: %+v", expected.header, actual.header)
	}

	if !slices.Equal(expected.data[:], actual.data[:]) {
		t.Fatalf("Wrong data bytes. Expected: %s, Actual: %s", hex.Dump(expected.data[:]), hex.Dump(actual.data[:]))
	}
}

func TestWriteParameterCmd(t *testing.T) {

	var hostId byte = 0x00
	var motorId byte = 0x7F

	var speed float32 = 1.12 // rad/s

	// Speed mode - expected data
	expected := cgSLCanFrame{}
	copy(expected.header[:], []byte{0x54, 0x31, 0x32, 0x30, 0x30, 0x30, 0x30, 0x37, 0x46, 0x38})
	copy(expected.data[:], []byte{0x30, 0x41, 0x37, 0x30, 0x30, 0x30, 0x30, 0x30, 0x32, 0x39, 0x35, 0x43, 0x38, 0x46, 0x33, 0x46})

	actual, err := WriteParameterCmd(hostId, motorId, PARAMETER_SPD_REF, speed)

	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(expected.header[:], actual.header[:]) {
		t.Fatalf("Wrong header bytes. Expected: %+v, Actual: %+v", expected.header, actual.header)
	}

	if !slices.Equal(expected.data[:], actual.data[:]) {
		t.Fatalf("Wrong data bytes. Expected: %s, Actual: %s", hex.Dump(expected.data[:]), hex.Dump(actual.data[:]))
	}
}
