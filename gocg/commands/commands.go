package commands

import (
	"fmt"
	"gocg/cybergear"
	"gocg/parameters"
	"strconv"
	"strings"
)

type dispatchFunc func(args []string, outputCh chan string) error

func executeHelpCmd(args []string, outputCh chan string) error {
	outputCh <- "Commands:"
	outputCh <- "\tenable <motor CAN id> - enable motor."
	outputCh <- "\tdisable <motor CAN id> - disable / stop motor."
	outputCh <- "All arguments are base 16"

	return nil
}

func executeEnableCmd(args []string, outputCh chan string) error {

	if len(args) != 1 || args[0] == "" {
		err := fmt.Sprintf("Syntax error ('enable <motor ID>')' Args: '%+v'", args)
		outputCh <- err
		return fmt.Errorf(err)
	}

	motorId, err := strconv.ParseUint(args[0], 16, 8)
	if err != nil {
		err := fmt.Sprintf("Syntax error: <motor ID>: '%s'", args[0])
		outputCh <- err
		return fmt.Errorf(err)
	}

	outputCh <- fmt.Sprintf("[DEBUG]: MotorId is : %02X", motorId)

	frame, err := cybergear.EnableMotorCmd(parameters.HostId, byte(motorId))

	if err != nil {
		outputCh <- err.Error()
		return err
	}

	outputCh <- fmt.Sprintf("[DEBUG]: SLCAN Frame: %+v", frame)
	outputCh <- fmt.Sprintf("TODO: Send to serial")

	return nil
}

func executeDisableCmd(args []string, outputCh chan string) error {
	if len(args) != 1 || args[0] == "" {
		err := fmt.Sprintf("Syntax error ('enable <motor ID>')' Args: '%+v'", args)
		outputCh <- err
		return fmt.Errorf(err)
	}

	motorId, err := strconv.ParseUint(args[0], 16, 8)
	if err != nil {
		err := fmt.Sprintf("Syntax error: <motor ID>: '%s'", args[0])
		outputCh <- err
		return fmt.Errorf(err)
	}

	outputCh <- fmt.Sprintf("[DEBUG]: MotorId is : %02X", motorId)

	frame, err := cybergear.EnableMotorCmd(parameters.HostId, byte(motorId))

	if err != nil {
		outputCh <- err.Error()
		return err
	}

	outputCh <- fmt.Sprintf("[DEBUG]: SLCAN Frame: %+v", frame)
	outputCh <- "TODO: Send to serial"
	return nil
}

var dispatchMap = map[string]dispatchFunc{
	"help":    executeHelpCmd,
	"enable":  executeEnableCmd,
	"disable": executeDisableCmd,
}

func Dispatch(command string, outputCh chan string) error {

	for key, value := range dispatchMap {
		if len(command) >= len(key) && command[:len(key)] == key {
			return value(
				strings.Split(strings.TrimSpace(command[len(key):]), " "), outputCh)
		}
	}

	err := fmt.Sprintf("unknown command: '%s'", command)
	outputCh <- err
	return fmt.Errorf(err)
}
