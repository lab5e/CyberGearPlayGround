package slcan

import (
	"fmt"
	"gocg/cybergear"
	"strconv"
)

const CYBERGEAR_FRAME_SIZE = 27

type CANFrameType byte

const (
	EXTENDED_FRAME     CANFrameType = 'T'
	STANDARD_FRAME     CANFrameType = 't'
	EXTENDED_RTR_FRAME CANFrameType = 'R'
	STANDARD_RTR_FRAME CANFrameType = 'r'
)

type frameOffset int

const (
	CAN_FRAME_TYPE_INDEX           frameOffset = 0
	CYBERGEAR_FRAME_TYPE_INDEX     frameOffset = 1
	HOST_ID_INDEX                  frameOffset = 3
	MOTOR_ID_INDEX                 frameOffset = 5
	DLC_INDEX                      frameOffset = 8
	CURRENT_ANGLE_INDEX            frameOffset = 10
	CURRENT_ANGULAR_VELOCITY_INDEX frameOffset = 14
	CURRENT_TORQUE_INDEX           frameOffset = 18
	CURRENT_TEMPERATURE_INDEX      frameOffset = 22
)

type Frame interface {
	CyberGearFrameType() cybergear.CommunicationType
	HostId() byte
	MotorId() byte
	String() string
	Unmarshal(frameBuffer []byte) error
}

func HandleIncomingFrame(frameBuffer []byte) (Frame, error) {
	if len(frameBuffer) != CYBERGEAR_FRAME_SIZE {
		return nil, fmt.Errorf("invalid frame received : %+v", frameBuffer)
	}

	var slCANFrameSize int
	if len(frameBuffer) > 0 {
		switch CANFrameType(frameBuffer[CAN_FRAME_TYPE_INDEX]) {
		case EXTENDED_FRAME:
			slCANFrameSize = CYBERGEAR_FRAME_SIZE
		case STANDARD_FRAME:
			return nil, fmt.Errorf("standard CAN frame not supported")
		case EXTENDED_RTR_FRAME:
			return nil, fmt.Errorf("extended RTR CAN frame not supported")
		case STANDARD_RTR_FRAME:
			return nil, fmt.Errorf("standard RTR CAN frame not supported")
		}
	}

	if slCANFrameSize == 0 {
		return nil, fmt.Errorf("no response frame received")
	}

	if slCANFrameSize != CYBERGEAR_FRAME_SIZE {
		return nil, fmt.Errorf("so far, we're only playing with cybergear extended SLCAN frames of %d characters", CYBERGEAR_FRAME_SIZE)
	}

	cgFrameType, err := strconv.ParseInt(string(frameBuffer[CYBERGEAR_FRAME_TYPE_INDEX:CYBERGEAR_FRAME_TYPE_INDEX+2]), 16, 16)
	if err != nil {
		return nil, err
	}

	switch cgFrameType {
	case int64(cybergear.COMMUNICATION_FETCH_DEVICE_ID): // Motor broadcast frame
		return nil, fmt.Errorf(">>>> small TODO here - don't forget to unmarshal broadcast frames <<<<")
	case int64(cybergear.COMMUNICATION_STATUS_REPORT):
		f := MotorFeedback{}
		err = f.Unmarshal(frameBuffer)
		return &f, err
	case int64(cybergear.COMMUNICATION_READ_SINGLE_PARAM):
		f := ParameterFrame{}
		err = f.Unmarshal(frameBuffer)
		return &f, err

	default:
		return nil, fmt.Errorf("unexpected cybergear frame type : %d", cgFrameType)
	}
}
