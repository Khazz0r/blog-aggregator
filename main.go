package main

import (
	"log"
	"os"

	"database/sql"

	"github.com/Khazz0r/blog-aggregator/internal/config"
	"github.com/Khazz0r/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error opening postgres database: %v", err)
	}
	dbQueries := database.New(db)

	mainState := &state{
		&cfg,
		dbQueries,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// commands that don't require you to be logged in
	cmds.register("register", handlerUserRegister)
	cmds.register("agg", handlerFetchFeed)
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("login", handlerLogin)

	// commands that require you to be logged in
	cmds.register("users", middlewareLoggedIn(handlerGetUsers))
	cmds.register("addfeed", middlewareLoggedIn(handlerCreateFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowingListFeeds))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	// NOT FOR PROD USE
	cmds.register("reset", handlerResetUsers)

	if len(os.Args) < 2 {
		log.Fatal("error: not enough arguments were provided")
	}

	cmd := command{
		os.Args[1],
		os.Args[2:],
	}

	err = cmds.run(mainState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
