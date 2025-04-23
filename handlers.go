package main

import (
	"context"
	"errors"
	"fmt"
	"time"
	"log"

	"github.com/Khazz0r/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

// handler for the CLI that deals with logging as long as it is an existing user
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("expected single word, a name, for login")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		log.Fatalf("user does not exist in database")
	}

	s.cfg.SetUser(cmd.args[0])
	if err != nil {
		log.Fatalf("error setting user")
	}
	fmt.Println("user has been set with the provided name")

	return nil
}

// handler for the CLI that deals with registering new users to the database
func handlerUserRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("expected single word, a name, for registering")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		log.Fatalf("user with that name already exists")
	}

	dbUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(), 
		CreatedAt: time.Now(), 
		UpdatedAt: time.Now(), 
		Name: cmd.args[0],
	})
	if err != nil {
		log.Fatalf("user already exists")
	}

	s.cfg.CurrentUserName = dbUser.Name
	s.cfg.SetUser(cmd.args[0])
	if err != nil {
		log.Fatalf("error setting user")
	}

	fmt.Printf("user %s was created successfully\n", dbUser.Name)

	// Test prints of created user's data for debugging purposes
	// fmt.Printf("User ID: %d\nUser created_at: %v\nUser updated_at: %v\nUser name: %s\n", dbUser.ID, dbUser.CreatedAt, dbUser.UpdatedAt, dbUser.Name)
	return nil
}
