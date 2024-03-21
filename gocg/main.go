package main

import (
	"fmt"
	"gocg/commands"
	"gocg/parameters"
	"log"
	"strings"

	"github.com/borud/chatui"
)

func main() {
	outputCh := make(chan string, 10)
	commandCh := make(chan string)

	chatui := chatui.New(chatui.Config{
		OutputCh:     outputCh,
		CommandCh:    commandCh,
		DynamicColor: false,
		BlockCtrlC:   false,
		HistorySize:  10,
	})

	outputCh <- "CyberGear playground. The current settings are:"
	outputCh <- fmt.Sprintf("Host  CAN id is : 0x%02X", parameters.HostId)
	outputCh <- "Frame format: SLCAN"
	outputCh <- "When in doubt: Type 'help' for - wait for it - help."

	go func() {
		for command := range commandCh {
			if strings.ToLower(command) == "/quit" {
				chatui.Stop()
			}
			err := commands.Dispatch(command, outputCh)
			if err != nil {
				outputCh <- err.Error()
			}
			chatui.SetStatus("last command was: " + command)
		}
	}()

	go func() {
		// this is done in a goroutine because it will block if the UI is not running.
		chatui.SetStatus("type /quit to exit")
	}()

	err := chatui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
