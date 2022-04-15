package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/adalessa/tshort/pkg/project"
)

func NewListCommand(pm *project.ProjectManager) *ListCommand {
	gc := &ListCommand{
		fs: flag.NewFlagSet("list", flag.ContinueOnError),
		pm: pm,
	}

	gc.fs.StringVar(&gc.format, "format", "json", "format to return")

	return gc
}

type ListCommand struct {
	fs     *flag.FlagSet
	pm     *project.ProjectManager
	format string
}

func (c *ListCommand) Init(args []string) error {
	return c.fs.Parse(args)
}

func (c *ListCommand) Name() string {
	return c.fs.Name()
}

func (c *ListCommand) Run() error {
	projects, err := c.pm.Load()
	if err != nil {
		return fmt.Errorf("error loading projects %w", err)
	}

	switch c.format {
	case "json":
		json.NewEncoder(os.Stdout).Encode(projects)
	case "text":
		var bindingsStr []string
		for key, project := range projects {
			bindingsStr = append(bindingsStr, fmt.Sprintf("[#%s %s]", key, project.Name))
		}
		sort.Strings(bindingsStr)
		fmt.Fprintf(os.Stdout, strings.Join(bindingsStr, " | "))
	}

	return nil
}
