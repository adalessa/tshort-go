package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/adalessa/tshort/internal/config"
	"github.com/adalessa/tshort/internal/project"
)

const BINDED_SESSION_PATH = "/tmp/tmux-projects.json"
const CONFIG_FILE = ".config/tshort/config.json"

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}
	// get the config
	config, err := getConfig()
	if err != nil {
		log.Fatalf("can not load the confi, error: %v", err)
	}

	projects, err := project.GetProjects(config.Directories)
	if err != nil {
		log.Fatalf("can not load the projects, error: %v", err)
	}

	// have the projects, can be sent to all
	// there is no need to modify this list
	// the "something" handling the bindings will not use it
	// before was the manager but I don't like that name
	// I want something bore related
	// also use the cache from home istead of the tmp directory
	fmt.Printf("%v", projects)

	return nil

	// all of them use the projects so make sence to load them, not much to lose not using them

	// TODO change the project manager
	// pm := project.NewProjectManager(
	// 	BINDED_SESSION_PATH,
	// 	fmt.Sprintf("/home/%s/code", os.Getenv("USER")),
	// )

	// cmds := []commands.Runner{
	// 	commands.NewListCommand(pm),
	// 	commands.NewBindCommand(pm),
	// 	commands.NewSelectorCommand(pm),
	// }

	// subCommand := os.Args[1]

	// for _, cmd := range cmds {
	// 	if cmd.Name() == subCommand {
	// 		err := cmd.Init(os.Args[2:])
	// 		if err != nil {
	// 			return fmt.Errorf("Error init command %s %w", cmd.Name(), err)
	// 		}

	// 		return cmd.Run()
	// 	}
	// }
	// return fmt.Errorf("Unknown sub-command %s", subCommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func getConfig() (config config.Config, err error) {
	home_dir := os.Getenv("HOME")
	jsonFile, err := os.Open(filepath.Join(home_dir, CONFIG_FILE))

	if err != nil {
		return config, err
	}

	defer jsonFile.Close()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return config, err
	}

	return config, json.Unmarshal(jsonBytes, &config)
}
