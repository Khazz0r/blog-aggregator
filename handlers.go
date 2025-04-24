package main

import (
	"context"
	"errors"
	"fmt"
	"time"
	"log"

	"github.com/Khazz0r/blog-aggregator/internal/database"
	"github.com/Khazz0r/blog-aggregator/internal/rss"
	"github.com/google/uuid"
)

// handler for the login command that logs in the username provided as long as it is an existing user
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

// handler for the register command that registers new users to the database
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

// NOT FOR PRODUCTION; handler for the reset command to allow easy testing by wiping users table
func handlerResetUsers(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		log.Fatalf("reset of users was not successful, please try again")
	}

	fmt.Println("reset of users was successful")
	return nil
}

// handler for the users command that prints out all the users in the database and shows the currently logged in one
func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("unable to retrieve all users from database")
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}

	return nil
}

// handler for the agg command that gets data from url provided and prints it all out
func handlerFetchFeed(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatalf("error fetching feed")
	}

	fmt.Printf("Feed: %+v\n", feed)

	return nil
}
