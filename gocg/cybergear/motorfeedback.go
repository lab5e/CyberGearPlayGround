package cybergear

import "fmt"

type MotorMode int32

const (
	ResetMode       MotorMode = 0
	CalibrationMode MotorMode = 1
	OperatingMode   MotorMode = 2
)

type MotorFeedback struct {
	hostId                byte    // Host CAN Id
	motorId               byte    // Motor CAN Id
	currentTorque         float32 // [-12, 12] N/m
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

func (f *MotorFeedback) Parse(frame []byte) error {
	if frame[0] != EXTENDED_FRAME_TYPE {
		return fmt.Errorf("unknown frame type")
	}

	return nil
}
