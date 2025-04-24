package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

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
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
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

// handler for the addfeed command that creates a feed attached to the current user
func handlerCreateFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("expected a name and url of the feed")
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		log.Fatalf("error getting user from database")
	}

	feed, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      cmd.args[0],
			Url:       cmd.args[1],
			UserID:    user.ID,
		})
	if err != nil {
		log.Fatalf("error adding feed: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		log.Fatalf("error creating feed follow: %v", err)
	}

	fmt.Println("feed was created successfully")
	printFeed(feed, user)
	fmt.Println()
	fmt.Printf("now following feed you just created\n")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("======================================================")

	return nil
}

// handler for the feeds command that prints out all the feeds and the users who created them
func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalf("unable to retrieve all feeds from database")
	}

	if len(feeds) == 0 {
		log.Fatalf("No feeds found")
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			log.Fatalf("failed to retrieve user by ID")
		}
		printFeed(feed, user)
		fmt.Println("======================================================")
	}

	return nil
}

// handler for the follow command that creates a feed follow for the current user
func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("no url found in follow command, please try again")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		log.Fatalf("error retrieving user for follow command: %v", err)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		log.Fatalf("error retrieving feed for follow command: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		log.Fatalf("error creating feed follow: %v", err)
	}

	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)

	return nil
}

// handler for the following command that prints out all the feeds that the current user is following
func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		log.Fatalf("error retrieving user for following command: %v", err)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		log.Fatalf("error getting feeds followed for current user: %v", err)
	}
	fmt.Printf("current feeds you are following:\n")
	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.FeedName)
	}

	return nil
}

// helper function to print out all the feed info to help reduce repeated code
func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
