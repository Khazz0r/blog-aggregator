package main

import (
	"log"
	"os"

	"github.com/Khazz0r/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	mainState := &state{
		&cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

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
