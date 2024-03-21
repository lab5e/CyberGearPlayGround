package commands

import (
	"fmt"
	"gocg/cybergear"
	"gocg/parameters"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

var serialPort *serial.Port

func SendFrame(frame *cybergear.SLCanFrame) error {
	bytesToSend := frame.Serialize()
	bytesToSend = append(bytesToSend, '\r')

	if nil == serialPort {
		return fmt.Errorf("It might be a good idea to open a serial port first...")
	}

	n, err := serialPort.Write(bytesToSend)
	if err != nil {
		return err
	}
	if n != len(bytesToSend) {
		return fmt.Errorf("Error sending frame. %d bytes sent of %d", n, len(bytesToSend))
	}
	serialPort.Flush()
	return nil
}

type dispatchFunc func(args []string, outputCh chan string) error

func executeHelpCmd(args []string, outputCh chan string) error {
	outputCh <- "Commands:"
	outputCh <- "\topen <serial port name> - opens serial port"
	outputCh <- "\tclose - close serial port"
	outputCh <- "\tenable <motor CAN id> - enable motor."
	outputCh <- "\tdisable <motor CAN id> - disable / stop motor."
	outputCh <- "\tspeed <motor CAN id> <rad/s> - set motor speed (-30~30rad/s)."
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

	err = SendFrame(frame)
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

	err = SendFrame(frame)
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

	serialConfig := &serial.Config{Name: args[1], Baud: 115200, Size: 8, Parity: serial.ParityNone, StopBits: 1}

	outputCh <- fmt.Sprintf("Opening %s", args[1])

	serialPort, err = serial.OpenPort(serialConfig)
	if err != nil {
		return fmt.Errorf("unable to open %s. Error %s", args[0], err)
	}

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
	err = SendFrame(frame)
	if err != nil {
		return err
	}

	var tmp float64
	tmp, err = strconv.ParseFloat(args[2], 64)
	var speed float32 = float32(tmp)
	if err != nil {
		return err
	}

	if speed < 30.0 || speed > 30.0 {
		return fmt.Errorf("invalid speed parameter: %2.2f. Valid values are in the interval [-30,30] rad/s", speed)
	}

	outputCh <- fmt.Sprintf("Setting current speed to %2.2f rad/s", speed)
	frame, err = cybergear.WriteParameterCmd(parameters.HostId, byte(motorId), cybergear.PARAMETER_SPD_REF, speed)
	if err != nil {
		return err
	}

	err = SendFrame(frame)
	if err != nil {
		return err
	}

	outputCh <- "OK"

	return nil
}

var dispatchMap = map[string]dispatchFunc{
	"help":    executeHelpCmd,
	"enable":  executeEnableCmd,
	"disable": executeDisableCmd,
	"open":    executeOpenCmd,
	"close":   executeCloseCmd,
	"speed":   executeSetSpeedCmd,
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
