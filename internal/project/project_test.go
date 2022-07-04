package project

import (
	"testing"
)

func TestReadProjectsFromDir(t *testing.T) {
	projects, err := GetProjects(map[string]string{
		"php":    "./test_data/php",
		"js":     "./test_data/js",
		"random": "./test_data/other",
	})

	if err != nil {
		t.Fatalf("Error returned %v", err)
	}

	if len(projects) != 3 {
		t.Fatalf("expecting 3 group of projects but got %d", len(projects))
	}

	if len(projects["php"]) != 3 {
		t.Fatalf("expecting 3 projects for php got %d", len(projects["php"]))
	}

	if len(projects["js"]) != 3 {
		t.Fatalf("expecting 3 projects for js got %d", len(projects["js"]))
	}

	if len(projects["random"]) != 1 {
		t.Fatalf("expecting 3 projects for random got %d", len(projects["random"]))
	}
}
