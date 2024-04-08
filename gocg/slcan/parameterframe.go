package slcan

import (
	"fmt"
	"gocg/cybergear"
)

type ParameterFrame struct {
	hostId        byte // Host CAN Id
	motorId       byte // Motor CAN Id
	parameter     int16
	parameterData []byte
}

func (f *ParameterFrame) CyberGearFrameType() cybergear.CommunicationType {
	return cybergear.COMMUNICATION_STATUS_REPORT
}

func (f *ParameterFrame) HostId() byte {

	return f.hostId
}

func (f *ParameterFrame) MotorId() byte {
	return f.motorId
}

func (f *ParameterFrame) String() string {
	var s string

	s += "Parameter:\n"
	s += fmt.Sprintf("host id : 0x%02X\n", f.hostId)
	s += fmt.Sprintf("motor id : 0x%02X\n", f.motorId)
	s += fmt.Sprintf("parameter : %04X\n", f.parameter)
	s += fmt.Sprintf("value : %+v\n", f.parameterData)
	return s
}

func (f *ParameterFrame) Unmarshal(frameBuffer []byte) error {
	// TBD

	return fmt.Errorf("Parameter frame received: %s", f.String())

	// return nil
}
