package main

import (
	"errors"

	"github.com/Khazz0r/blog-aggregator/internal/config"
	"github.com/Khazz0r/blog-aggregator/internal/database"
)

type state struct {
	cfg *config.Config
	db *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registeredCommands[cmd.name]
	if !exists {
		return errors.New("command does not exist")
	}

	return handler(s, cmd)
}
