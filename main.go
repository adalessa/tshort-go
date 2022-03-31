package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Project struct {
	Name  string
	Title string
	Path  string
}

type Projects map[string]Project

func main() {
	// move to a function to clean the code just get the bindedProjects
	file := "/tmp/tmux-projects.json"
	content, err := ioutil.ReadFile(file)
	bindProjects := Projects{}
	if err != nil {
		if err.Error() == "open /tmp/tmux-projects.json: no such file or directory" {
			// create file
			f, err := os.Create(file)
			if err != nil {
				fmt.Println(err)
				return
			}
			// write empty projects
			json.NewEncoder(f).Encode(bindProjects)
			f.Close()
			content, _ = ioutil.ReadFile(file)
		} else {
			fmt.Println(err)
			return
		}
	}
	err = json.Unmarshal(content, &bindProjects)
	if err != nil {
		fmt.Println(err)
		return
	}

	// improve with a switch to select action
	if len(os.Args) == 1 {
		_, err := runSelection()
		if err != nil {
			fmt.Println(err)
		}
	}

	// read parameter from command line
	if len(os.Args) > 3 {
		fmt.Println("Usage: tmux-projects <action> <key>")
	}
	action := os.Args[1]
	// actions = [list, switch, bind]
	if action == "list" {
		fmt.Println(string(content))
	} else if action == "switch" {
		if len(os.Args) > 3 {
			fmt.Println("Usage: tmux-projects switch <key>")
		}
		if len(os.Args) == 3 {
			key := os.Args[2]
			if project, ok := bindProjects[key]; ok {
				switchOrCreateSession(project)
			} else {
				fmt.Println("Project not found")
			}
		}
	} else if action == "bind" {
		if len(os.Args) > 3 {
			fmt.Println("Usage: tmux-projects bind <key>")
		}
		key := os.Args[2]
		if len(os.Args) == 3 {
			project, err := runSelection()
			if err != nil {
				fmt.Println(err)
				return
			}
			bindProjects[key] = project

			content, err := json.Marshal(bindProjects)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = ioutil.WriteFile(file, content, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func runSelection() (Project, error) {
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
	switchOrCreateSession(project)

	return project, nil
}

func switchOrCreateSession(project Project) {
	os.Chdir(project.Path)
	cmd := exec.Command("tmux", "has-session", "-t", project.Name)
	result, _ := cmd.Output()
	if string(result) != "0\n" {
		cmd = exec.Command("tmux", "new-session", "-d", "-s", project.Name, "nvim .")
		cmd.Run()
	}
	cmd = exec.Command("tmux", "switch-client", "-t", project.Name)
	cmd.Run()
	// TODO trying to create the window fails because there is no command has-window
	// cmd = exec.Command("tmux", "has-window", "-t", project.Name, "-n", "vim")
	// result, _ = cmd.Output()
	// fmt.Print(result)
	// if string(result) != "0\n" {
	// 	cmd = exec.Command("tmux", "new-window", "-t", project.Name, "-n", "vim", "nvim .")
	// 	cmd.Run()
	// }
	// cmd = exec.Command("tmux", "switch-client", "-t", project.Name, "-w", "vim")
	// cmd.Run()
}

func getProjects(directory string) []Project {
	var result []Project
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
