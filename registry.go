package main

import (
	"fmt"

	"github.com/muhammad-romli/ArachNet/internal/config"
)

type state struct {
	ConfigState *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	CommandsMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.CommandsMap[cmd.Name]
	if !ok {
		return fmt.Errorf("Value doesn't exist")
	}
	err := function(s, cmd)
	if err != nil {
		return fmt.Errorf("error happen when running command: %w", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.CommandsMap[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Login function need 1 parameter at least")
	}
	s.ConfigState.SetUser(cmd.Args[0])
	fmt.Printf("User has been set")
	return nil
}
