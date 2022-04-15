package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/adalessa/tshort/pkg/tmux"
)

type ProjectManager struct {
	storagePath       string
	projectsDirectory string
}

type Project struct {
	Name  string
	Title string
	Path  string
}

func NewProjectManager(path string, projectsDir string) *ProjectManager {
	return &ProjectManager{
		storagePath:       path,
		projectsDirectory: projectsDir,
	}
}

func (pm *ProjectManager) Load() (map[string]Project, error) {
	projects, err := read(pm.storagePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading the saved projects %w", err)
	}

	shouldResave := false
	if len(projects) > 0 {
		for key, project := range projects {
			if !tmux.SessionExists(project.Name) {
				delete(projects, key)
				shouldResave = true
			}
		}
	}

	if shouldResave {
		err := pm.Save(projects)
		if err != nil {
			return nil, fmt.Errorf("Error saving the projects %w", err)
		}
	}

	return projects, nil
}

func (pm *ProjectManager) Save(projects map[string]Project) error {
	content, err := json.Marshal(projects)
	if err != nil {
		return fmt.Errorf("Error encoding projects to save %w", err)
	}

	err = ioutil.WriteFile(pm.storagePath, content, 0644)
	if err != nil {
		return fmt.Errorf("Error saving projects to disk %w", err)
	}

	return nil
}

func (pm *ProjectManager) NewProject(name string) (Project, error) {
	if !tmux.SessionExists(name) {
		return Project{}, errors.New("session does not exits on tmux")
	}

	return Project{
		Name: name,
		Path: name,
	}, nil
}

func read(path string) (map[string]Project, error) {
	projects := make(map[string]Project)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		if err.Error() == fmt.Sprintf("open %s: no such file or directory", path) {
			return projects, nil
		}
		return nil, fmt.Errorf("Error reading the file %w", err)
	}

	err = json.Unmarshal(content, &projects)
	if err != nil {
		return projects, err
	}

	return projects, nil
}

func (pm *ProjectManager) GetProjects() (result []Project, err error) {
	directory := pm.projectsDirectory
	languages, err := getDirectories(directory)
	if err != nil {
		return nil, fmt.Errorf("Error getting the directory %w", err)
	}
	for _, language := range languages {
		if language != "go" {
			projects, err := getDirectories(fmt.Sprintf("%s/%s", directory, language))
			if err != nil {
				return nil, fmt.Errorf("Error getting the sub-directory %w", err)
			}
			for _, project := range projects {
				result = append(result, Project{Name: project, Title: fmt.Sprintf("[%s] %s", strings.ToUpper(language), project), Path: fmt.Sprintf("%s/%s/%s", directory, language, project)})
			}
			continue
		}
		platforms, err := getDirectories(fmt.Sprintf("%s/%s/src", directory, language))
		if err != nil {
			return nil, fmt.Errorf("Error getting the sub-directory %w", err)
		}
		for _, platform := range platforms {
			accounts, err := getDirectories(fmt.Sprintf("%s/%s/src/%s", directory, language, platform))
			if err != nil {
				return nil, fmt.Errorf("Error getting the go account %w", err)
			}
			for _, account := range accounts {
				projects, err := getDirectories(fmt.Sprintf("%s/%s/src/%s/%s", directory, language, platform, account))
				if err != nil {
					return nil, fmt.Errorf("Error getting the projects %w", err)
				}
				for _, project := range projects {
					result = append(result, Project{Name: project, Title: fmt.Sprintf("[%s %s/%s] %s", strings.ToUpper(language), platform, account, project), Path: fmt.Sprintf("%s/%s/src/%s/%s/%s", directory, language, platform, account, project)})
				}
			}
		}
	}
	return result, nil
}

func getDirectories(directory string) ([]string, error) {
	dirs, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s (%w)", directory, err)
	}
	var result []string
	for _, dir := range dirs {
		if dir.IsDir() {
			if dir.Name() != ".git" {
				result = append(result, dir.Name())
			}
		}
	}
	return result, nil
}
