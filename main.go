package main

import (
	"log"
	"os"

	"github.com/Khazz0r/blog-aggregator/internal/config"
	"github.com/Khazz0r/blog-aggregator/internal/database"
	"database/sql"
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

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerUserRegister)
	cmds.register("reset", handlerResetUsers)
	cmds.register("users", handlerGetUsers)

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
