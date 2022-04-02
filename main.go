package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const BINDED_SESSION_PATH = "/tmp/tmux-projects.json"

type Project struct {
	Name  string
	Title string
	Path  string
}

type BindedProjects map[string]Project

func main() {
	projects, err := getBindedSessions()
	if err != nil {
		fmt.Println(err)
		return
	}
	projects = removedDeletedSessions(projects)

	switch {
	case len(os.Args) == 1:
		project, err := selectProject()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Switched to %s\n", project.Title)
	case len(os.Args) >= 2:
		action := os.Args[1]
		switch action {
		case "list":
			// json.NewEncoder(os.Stdout).Encode(projects)
			var bindingsStr []string
			for key, project := range projects {
				bindingsStr = append(bindingsStr, fmt.Sprintf("[#%s %s]", key, project.Name))
			}
			sort.Strings(bindingsStr)
			fmt.Println(strings.Join(bindingsStr, " | "))
		case "switch":
			if len(os.Args) > 3 {
				fmt.Println("Usage: tshort switch <key>")
				return
			}
			key := os.Args[2]
			if project, ok := projects[key]; ok {
				changeToSession(project)
			} else {
				project, err := selectProject()
				if err != nil {
					fmt.Println(err)
					return
				}
				projects[key] = project
			}
			// TODO: add bind command
			// the question is what project to bind ?
			// option create a project using the session name, search a project using the session name
			// if no project found create a simple with the name as it is
			// the path try to get cwd
		default:
			fmt.Println("Usage: tshort [list|switch]")
		}
	}
	saveBindedSessions(projects)
}

func saveBindedSessions(projects BindedProjects) error {
	content, err := json.Marshal(projects)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(BINDED_SESSION_PATH, content, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getBindedSessions() (projects BindedProjects, err error) {
	content, err := ioutil.ReadFile(BINDED_SESSION_PATH)
	if err != nil {
		if err.Error() == fmt.Sprintf("open %s: no such file or directory", BINDED_SESSION_PATH) {
			return make(BindedProjects), nil
		} else {
			fmt.Println(err)
			return projects, err
		}
	}

	err = json.Unmarshal(content, &projects)
	if err != nil {
		return projects, err
	}

	return projects, nil
}

func selectProject() (Project, error) {
	directory := fmt.Sprintf("/home/%s/code", os.Getenv("USER"))
	projects := getProjects(directory)

	projectName := withFilter("fzf -m", func(in io.WriteCloser) {
		for _, project := range projects {
			fmt.Fprintf(in, "%s\n", project.Title)
		}
	})

	if projectName == nil {
		return Project{}, errors.New("No project selected")
	}
	var project Project
	for _, p := range projects {
		if p.Title == projectName[0] {
			project = p
			break
		}
	}
	if project.Name == "" {
		return Project{}, errors.New("Project not found")
	}
	changeToSession(project)

	return project, nil
}

func changeToSession(project Project) {
	os.Chdir(project.Path)
	if !sessionExists(project.Name) {
		exec.Command("tmux", "new-session", "-d", "-s", project.Name, "nvim .").Run()
	}
	exec.Command("tmux", "switch-client", "-t", project.Name).Run()
}

func getProjects(directory string) (result []Project) {
	languages := getDirectories(directory)
	for _, language := range languages {
		if language != "go" {
			subdirs := getDirectories(fmt.Sprintf("%s/%s", directory, language))
			for _, project := range subdirs {
				result = append(result, Project{Name: project, Title: fmt.Sprintf("[%s] %s", strings.ToUpper(language), project), Path: fmt.Sprintf("%s/%s/%s", directory, language, project)})
			}
			continue
		}
		platforms := getDirectories(fmt.Sprintf("%s/%s/src", directory, language))
		for _, platform := range platforms {
			for _, account := range getDirectories(fmt.Sprintf("%s/%s/src/%s", directory, language, platform)) {
				for _, project := range getDirectories(fmt.Sprintf("%s/%s/src/%s/%s", directory, language, platform, account)) {
					result = append(result, Project{Name: project, Title: fmt.Sprintf("[%s %s/%s] %s", strings.ToUpper(language), platform, account, project), Path: fmt.Sprintf("%s/%s/src/%s/%s/%s", directory, language, platform, account, project)})
				}
			}
		}
	}
	return result
}

func getDirectories(directory string) []string {
	dirs, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var result []string
	for _, dir := range dirs {
		if dir.IsDir() {
			if dir.Name() != ".git" {
				result = append(result, dir.Name())
			}
		}
	}
	return result
}

func withFilter(command string, input func(in io.WriteCloser)) []string {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-c", command)
	cmd.Stderr = os.Stderr
	in, _ := cmd.StdinPipe()
	go func() {
		input(in)
		in.Close()
	}()
	result, _ := cmd.Output()
	return strings.Split(string(result), "\n")
}

func sessionExists(name string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", name)
	// var stdout, stderr bytes.Buffer
	// cmd.Stdout = &stdout
	// cmd.Stderr = &stderr
	err := cmd.Run()
	// outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	// fmt.Printf("Tmux out: %s\n", outStr)
	// fmt.Printf("Tmux err: %s\n", errStr)
	if err != nil {
		// fmt.Println(err)
		return false
	}
	return true
}

func removedDeletedSessions(bindedProjects BindedProjects) BindedProjects {
	for key, project := range bindedProjects {
		if !sessionExists(project.Name) {
			delete(bindedProjects, key)
		}
	}
	return bindedProjects
}
