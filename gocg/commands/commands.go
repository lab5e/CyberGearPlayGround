package commands

import (
	"fmt"
	"gocg/cybergear"
	"gocg/parameters"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

var serialPort *serial.Port

func ReadFrame(outputCh chan string) {
	var frameBuffer []byte
	readBuffer := make([]byte, 32)
	var n int
	for {
		n, _ = serialPort.Read(readBuffer)
		if n > 0 {
			frameBuffer = append(frameBuffer, readBuffer[:n]...)
		} else {
			break
		}
	}

	if len(frameBuffer) > 0 {
		// Parse frame (only two different response frame types)
		outputCh <- fmt.Sprintf("Response: %+v", frameBuffer)
		outputCh <- fmt.Sprintf("Response: %s", frameBuffer)

		// Move to motor feedback
		// 	// bit 28 - 24 communication type
		// 	// bit 8 - 15 motor CAN ID

		// 	// bit 21 - not calibrated
		// 	// bit 20 - hall encoding fault
		// 	// bit 19 - magnetic encoding error
		// 	// bit 18 - overtemperature
		// 	// bit 17 - overcurrent
		// 	// bit 16 - undervoltage
		// 	// bit 22-23 mode (0: reset, 1: calibration: 2: run mode )

		// 	// bit 0-7 host CAN id

		// 	// er B checksum ?

		// 	Example frame (ASCII): 	T02807f008F FF F7 FF 87 FF F0 12 B

		// 	B == checksum ?

		// 	T
		// 	02 - communication type
		// 	80 - 0101 0000
		// 	7f - motorID
		// 	00 - host canID
		// 	8F FF F7 FF 87 FF F0 12B

		// 30 02 80 7f 00 8f ff f7 fb b7 ff f0 13

	}

}

func SendFrame(frame *cybergear.SLCanFrame, outputCh chan string) error {
	bytesToSend := frame.Serialize()
	bytesToSend = append(bytesToSend, '\r')

	if nil == serialPort {
		return fmt.Errorf("it might be a good idea to open a serial port first")
	}

	n, err := serialPort.Write(bytesToSend)
	if err != nil {
		return err
	}
	if n != len(bytesToSend) {
		return fmt.Errorf("error sending frame. %d bytes sent of %d", n, len(bytesToSend))
	}
	serialPort.Flush()

	ReadFrame(outputCh)

	// readBuffer := make([]byte, 256)
	// for i := 0; i < 2; i++ {
	// 	_, err = serialPort.Read(readBuffer)
	// 	if err != nil {
	// 		outputCh <- err.Error()
	// 	}
	// 	// outputCh <- fmt.Sprintf("Reply: %d bytes", n)
	// 	// outputCh <- fmt.Sprintf("%+v", readBuffer[:n])
	// }
	return nil
}

type dispatchFunc func(args []string, outputCh chan string) error

func executeHelpCmd(args []string, outputCh chan string) error {
	outputCh <- "Commands:"
	outputCh <- "\topen <serial port name> - opens serial port"
	outputCh <- "\tclose - close serial port"
	outputCh <- "\tenable <motor CAN id> - enable motor."
	outputCh <- "\tdisable <motor CAN id> - disable / stop motor."
	outputCh <- "\tset_speed <motor CAN id> <rad/s> - set motor speed (-30~30rad/s)."
	outputCh <- "\tget_feedback <motor CAN id>"
	//	outputCh <- "\tmode <motor CAN id> <speed | position | current> - set operation mode"

	return nil
}

func executeEnableCmd(args []string, outputCh chan string) error {
	var frame *cybergear.SLCanFrame

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('enable <motor ID>')' Args: '%+v'", args)
	}

	motorId, err := strconv.ParseUint(args[1], 16, 8)
	if err != nil {
		return fmt.Errorf("syntax error: <motor ID>: '%s'", args[1])
	}

	if err != nil {
		return err
	}

	outputCh <- fmt.Sprintf("Enabling motor (CAN id: %02X)", motorId)
	frame, err = cybergear.EnableMotorCmd(parameters.HostId, byte(motorId))
	if err != nil {
		return err
	}

	err = SendFrame(frame, outputCh)
	if err != nil {
		return err
	}

	outputCh <- "OK"

	return nil
}

func executeDisableCmd(args []string, outputCh chan string) error {
	if len(args) != 2 {
		return fmt.Errorf("syntax error ('enable <motor ID>')' Args: '%+v'", args)
	}

	motorId, err := strconv.ParseUint(args[1], 16, 8)
	if err != nil {
		return fmt.Errorf("syntax error: disable <motor ID>: '%s'", args[1])
	}

	frame, err := cybergear.DisableMotorCmd(parameters.HostId, byte(motorId))

	if err != nil {
		return err
	}

	err = SendFrame(frame, outputCh)
	if err != nil {
		return err
	}

	return nil
}

func executeOpenCmd(args []string, outputCh chan string) error {
	var err error

	if len(args) != 2 {
		return fmt.Errorf("syntax error ('open <serial port name>')' Args: '%+v'", args)
	}

	serialConfig := &serial.Config{Name: args[1], Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1, ReadTimeout: time.Second * 1}

	outputCh <- fmt.Sprintf("Opening %s", args[1])

	serialPort, err = serial.OpenPort(serialConfig)

	if err != nil {
		return fmt.Errorf("unable to open %s. Error %s", args[0], err)
	}

	outputCh <- "Setting CAN bitrate to 1Mbit"
	setBitrateCmd := []byte{'S', '8', '\r'}
	serialPort.Write(setBitrateCmd)

	ReadFrame(outputCh)

	// readBuffer := make([]byte, 16)
	// var n int
	// for i := 0; i < 2; i++ {
	// 	n, _ = serialPort.Read(readBuffer)
	// 	if n > 0 {
	// 		outputCh <- fmt.Sprintf("Response: %+v", readBuffer[:n])
	// 	}
	// }

	time.Sleep(20 * time.Millisecond)

	outputCh <- "Opening CAN Channel in normal mode (send/recevie)"
	openCANcmd := []byte{'O', '\r'}
	serialPort.Write(openCANcmd)

	// for i := 0; i < 2; i++ {
	// 	n, _ = serialPort.Read(readBuffer)
	// 	if n > 0 {
	// 		outputCh <- fmt.Sprintf("Response: %+v", readBuffer[:n])
	// 	}
	// }
	ReadFrame(outputCh)

	outputCh <- "OK"

	return nil
}

func executeCloseCmd(args []string, outputCh chan string) error {
	if len(args) != 1 {
		return fmt.Errorf("syntax error ('close')' Args: '%s'", args)
	}

	outputCh <- "Closing serial port"

	if serialPort != nil {
		serialPort.Close()
	} else {
		outputCh <- "No worries, I'll close the serial port you never bothered to open in the first place..."
	}

	outputCh <- "OK"
	return nil
}

func executeSetSpeedCmd(args []string, outputCh chan string) error {
	var frame *cybergear.SLCanFrame
	var err error
	var motorId int64

	if len(args) != 3 {
		return fmt.Errorf("syntax error ('speed <motorId> <rad/s>')' Args: '%+v'", args)
	}

	motorId, err = strconv.ParseInt(args[1], 16, 8)
	if err != nil {
		return err
	}

	outputCh <- fmt.Sprintf("Setting run mode to SPEED MODE for motor %02X", motorId)
	frame, err = cybergear.SetRunMode(parameters.HostId, byte(motorId), cybergear.SPEED_MODE)
	if err != nil {
		return err
	}
	err = SendFrame(frame, outputCh)
	if err != nil {
		return err
	}

	var tmp float64
	tmp, err = strconv.ParseFloat(args[2], 64)
	var speed float32 = float32(tmp)
	if err != nil {
		return err
	}

	if speed < -30.0 || speed > 30.0 {
		return fmt.Errorf("invalid speed parameter: %2.2f. Valid values are in the interval [-30,30] rad/s", speed)
	}

	outputCh <- fmt.Sprintf("Setting current speed to %2.2f rad/s", speed)
	frame, err = cybergear.WriteParameterCmd(parameters.HostId, byte(motorId), cybergear.PARAMETER_SPD_REF, speed)
	if err != nil {
		return err
	}

	err = SendFrame(frame, outputCh)
	if err != nil {
		return err
	}

	outputCh <- "OK"

	return nil
}

var dispatchMap = map[string]dispatchFunc{
	"help":      executeHelpCmd,
	"enable":    executeEnableCmd,
	"disable":   executeDisableCmd,
	"open":      executeOpenCmd,
	"close":     executeCloseCmd,
	"set_speed": executeSetSpeedCmd,
	//	"get_feedback": executeGetFeedbackCmd,
}

func Dispatch(command string, outputCh chan string) error {

	command = strings.TrimSpace(command)
	for key, value := range dispatchMap {
		if len(command) >= len(key) && command[:len(key)] == key {
			return value(
				strings.Split(command, " "), outputCh)
		}
	}

	return fmt.Errorf("unknown command: '%s'", command)
}
