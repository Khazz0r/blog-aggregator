package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("expected single word, a username, for login")
	}

	s.cfg.SetUser(cmd.args[0])
	fmt.Println("user has been set with the provided username")

	return nil
}
