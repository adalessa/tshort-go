package project

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Name  string
	Title string
	Path  string
}

func GetProjects(directories map[string]string) (projects map[string][]Project, err error) {
	projects = make(map[string][]Project)
	home_dir := os.Getenv("HOME")

	for name, directory := range directories {
		if strings.HasPrefix(directory, "~/") {
			directory = strings.Replace(directory, "~", home_dir, 1)
		}

		elements, err := os.ReadDir(directory)
		if err != nil {
			log.Default().Printf("Error reading directory %s %v", directory, err)
			return projects, err
		}
		projects[name] = getProjectsFromDir(directory, elements)
	}

	return projects, err
}

func getProjectsFromDir(directory string, elements []fs.DirEntry) (projects []Project) {

	projects = make([]Project, len(elements))
	for i, element := range elements {
		if !element.IsDir() {
			continue
		}

		projects[i] = Project{
			Name:  element.Name(),
			Title: element.Name(),
			Path:  filepath.Join(directory, element.Name()),
		}
	}

	return projects
}
