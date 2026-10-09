package main

import (
	"fmt"
	"os"

	"github.com/muhammad-romli/ArachNet/internal/config"
)

func readLoop(maxRetries int) *config.Config {
	var data *config.Config
	var err error
	for attempt := 0; attempt < maxRetries; attempt++ {
		data, err = config.Read()
		if err == nil {
			break
		}
		fmt.Printf("Retrying reading\n")
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return data
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Enter the command to start using the CLI applications")
		fmt.Println("usage: app <command> <param>")
		return
	}
	commandName := os.Args[1]
	commandParams := os.Args[2:]

	maxRetries := 2
	data := readLoop(maxRetries)
	if data == nil {
		fmt.Printf("There is no data inside targeted file")
		return
	}
	// The data is just placeholder in the first readLoop

	mainCommands := commands{CommandsMap: make(map[string]func(*state, command) error)}
	mainCommands.register("login", handlerLogin)
	mainState := state{}
	mainState.ConfigState = data

	cmd := command{commandName, commandParams}
	mainCommands.run(&mainState, cmd)
}
