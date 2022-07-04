package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/adalessa/tshort/pkg/fzf"
	"github.com/adalessa/tshort/pkg/project"
	"github.com/adalessa/tshort/pkg/tmux"
)

func NewSelectorCommand(pm *project.ProjectManager) *SelectorCommand {
	gc := &SelectorCommand{
		fs: flag.NewFlagSet("select", flag.ContinueOnError),
		pm: pm,
	}

	return gc
}

type SelectorCommand struct {
	fs *flag.FlagSet
	pm *project.ProjectManager
}

func (c *SelectorCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *SelectorCommand) Name() string {
	return c.fs.Name()
}

func (c *SelectorCommand) Run() error {
	var project project.Project
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
	tmux.ChangeToSession(project.Path, project.Name)

	return nil
}
