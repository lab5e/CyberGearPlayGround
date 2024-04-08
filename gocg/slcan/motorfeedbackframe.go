package slcan

import (
	"fmt"
	"gocg/cybergear"
	"strconv"
)

type MotorMode int32

const (
	ResetMode       MotorMode = 0
	CalibrationMode MotorMode = 1
	OperatingMode   MotorMode = 2
)

type MotorFeedback struct {
	hostId                byte    // Host CAN Id
	motorId               byte    // Motor CAN Id
	currentTorque         float32 // [-12, 12] Nm
	currentAngle          float32 // [-4pi, 4pi]
	currentSpeed          float32 // [-30rad/s, 30rad/s]
	currentTemperature    float32 // Current temperature: Temp (degrees Celsius) * 10
	calibrationError      bool
	hallEncoderError      bool
	magneticEncodingError bool
	overtemperature       bool
	overcurrent           bool
	undervoltage          bool
	mode                  MotorMode
}

func (f *MotorFeedback) CyberGearFrameType() cybergear.CommunicationType {
	return cybergear.COMMUNICATION_STATUS_REPORT
}

func (f *MotorFeedback) HostId() byte {
	return f.hostId
}

func (f *MotorFeedback) MotorId() byte {
	return f.motorId
}

func (f *MotorFeedback) ParseByte(ascii []byte) (byte, error) {
	num, err := strconv.ParseInt(string(ascii), 16, 16)
	return byte(num), err
}

func (f *MotorFeedback) ParseInt(ascii []byte) (int32, error) {
	var err error
	num, err := strconv.ParseUint(string(ascii), 16, 16)
	return int32(num), err
}

func (f *MotorFeedback) Unmarshal(frameBuffer []byte) error {
	var err error

	f.hostId, err = f.ParseByte(frameBuffer[HOST_ID_INDEX : HOST_ID_INDEX+2])
	if err != nil {
		return err
	}
	f.motorId, err = f.ParseByte(frameBuffer[MOTOR_ID_INDEX : MOTOR_ID_INDEX+2])
	if err != nil {
		return err
	}
	var dlc byte
	dlc, err = f.ParseByte(frameBuffer[DLC_INDEX : DLC_INDEX+2])
	if err != nil {
		return err
	}
	if dlc != 8 {
		return fmt.Errorf("unexpected DLC (%d). Expected DLC of 8", dlc)
	}

	var num int32
	// Data 00-01: Current angle [0-65535] == [-4PI, 4PI]
	num, err = f.ParseInt(frameBuffer[CURRENT_ANGLE_INDEX : CURRENT_ANGLE_INDEX+4])
	if err != nil {
		return err
	}
	f.currentAngle = 8*float32(num)/65535 - 4

	// Data 02-03: Current angular velocity [0-65535] == [-30 rad/s, 30 rad/s]
	num, err = f.ParseInt(frameBuffer[CURRENT_ANGULAR_VELOCITY_INDEX : CURRENT_ANGULAR_VELOCITY_INDEX+4])
	if err != nil {
		return err
	}
	f.currentSpeed = 60*float32(num)/65535 - 30

	// Data 04-05: Current torque [0-65535] == [-12Nm, 12Nm]
	num, err = f.ParseInt(frameBuffer[CURRENT_TORQUE_INDEX : CURRENT_TORQUE_INDEX+4])
	if err != nil {
		return err
	}
	f.currentTorque = 24*float32(num)/65535 - 12

	// Data 06-07: Current temperature (C * 10)
	num, err = f.ParseInt(frameBuffer[CURRENT_TEMPERATURE_INDEX : CURRENT_TEMPERATURE_INDEX+4])
	if err != nil {
		return err
	}
	f.currentTemperature = float32(num) / 10

	/*
		// 	// bit 28 - 24 communication type
		// 	// bit 8 - 15 motor CAN ID

		TODO: Decode bit fields
		// 	// bit 21 - not calibrated
		// 	// bit 20 - hall encoding fault
		// 	// bit 19 - magnetic encoding error
		// 	// bit 18 - overtemperature
		// 	// bit 17 - overcurrent
		// 	// bit 16 - undervoltage
		// 	// bit 22-23 mode (0: reset, 1: calibration: 2: run mode )

	*/

	return nil

}

func (f *MotorFeedback) String() string {
	var s string

	s += fmt.Sprintf("torque : %02.2f Nm", f.currentTorque)

	// s += "Motor status:\n"
	// s += fmt.Sprintf("host id : 0x%02X\n", f.hostId)
	// s += fmt.Sprintf("motor id : 0x%02X\n", f.motorId)
	// s += fmt.Sprintf("angle : %02.2f rad\n", f.currentAngle)
	// s += fmt.Sprintf("angular velocity : %02.2f rad/s\n", f.currentSpeed)
	// s += fmt.Sprintf("torque : %02.2f Nm\n", f.currentTorque)
	// s += fmt.Sprintf("temperature : %02.2f C\n", f.currentTemperature)

	return s
}
