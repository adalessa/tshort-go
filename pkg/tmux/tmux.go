package tmux

import (
	"os"
	"os/exec"
)

func SessionExists(name string) bool {
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

func ChangeToSession(path string, name string) {
	os.Chdir(path)
	if !SessionExists(name) {
		exec.Command("tmux", "new-session", "-d", "-s", name, "nvim").Run()
	}
	exec.Command("tmux", "switch-client", "-t", name).Run()
}
