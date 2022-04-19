package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/adalessa/tshort/pkg/fzf"
	"github.com/adalessa/tshort/pkg/project"
	"github.com/adalessa/tshort/pkg/tmux"
)

func NewBindCommand(pm *project.ProjectManager) *BindCommand {
	gc := &BindCommand{
		fs: flag.NewFlagSet("bind", flag.ContinueOnError),
		pm: pm,
	}

	gc.fs.StringVar(&gc.key, "key", "", "Key to bind")
	gc.fs.StringVar(&gc.session, "session", "", "session to bind")

	return gc
}

type BindCommand struct {
	fs *flag.FlagSet
	pm *project.ProjectManager

	key     string
	session string
}

func (c *BindCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *BindCommand) Name() string {
	return c.fs.Name()
}

func (c *BindCommand) Run() error {
	if c.key == "" {
		return errors.New("Should provide a key to bind")
	}
	bindedProjects, err := c.pm.Load()
	if err != nil {
		return fmt.Errorf("error loading projects %w", err)
	}

	if project, ok := bindedProjects[c.key]; ok {
		tmux.ChangeToSession(project.Path, project.Name)

		return nil
	}

	var project project.Project
	if c.session == "" {
		projects, err := c.pm.GetProjects()
		if err != nil {
			return fmt.Errorf("Error getting projects, %w", err)
		}

		list := fzf.FZF(func(in io.WriteCloser) {
			for _, project := range projects {
				fmt.Fprintf(in, "%s\n", project.Title)
			}
		})

		if list[0] == "" {
			return errors.New("No Project selected")
		}
		for _, p := range projects {
			if p.Title == list[0] {
				project = p
				break
			}
		}

		if project.Name == "" {
			return errors.New("Project not found")
		}

	} else {
		project, err = c.pm.NewProject(c.session)
		if err != nil {
			return fmt.Errorf("Error saving project %s, %w", c.session, err)
		}
	}

	bindedProjects[c.key] = project
	c.pm.Save(bindedProjects)

	tmux.ChangeToSession(project.Path, project.Name)

	return nil
}
