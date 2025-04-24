package main

import (
	"log"
	"fmt"
	"context"
)

// NOT FOR PRODUCTION; handler for the reset command to allow easy testing by wiping users table
func handlerResetUsers(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		log.Fatalf("reset of users was not successful, please try again")
	}

	fmt.Println("reset of users was successful")
	return nil
}
