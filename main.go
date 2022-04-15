package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/adalessa/tshort/cmd"
	"github.com/adalessa/tshort/pkg/project"
)

const BINDED_SESSION_PATH = "/tmp/tmux-projects.json"

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	pm := project.NewProjectManager(
		BINDED_SESSION_PATH,
		fmt.Sprintf("/home/%s/code", os.Getenv("USER")),
	)

	cmds := []cmd.Runner{
		cmd.NewListCommand(pm),
		cmd.NewBindCommand(pm),
		cmd.NewSelectorCommand(pm),
	}

	subCommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subCommand {
			err := cmd.Init(os.Args[2:])
			if err != nil {
				return fmt.Errorf("Error init command %s %w", cmd.Name(), err)
			}

			return cmd.Run()
		}
	}
	return fmt.Errorf("Unknown sub-command %s", subCommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
