package main

import (
	"context"
	"log"

	"github.com/Khazz0r/blog-aggregator/internal/database"
)
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			log.Fatalf("Unable to retrieve logged in user from database")
		}

		return handler(s, cmd, user)
	}
}
