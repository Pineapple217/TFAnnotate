package state

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os/exec"
)

type State struct {
	Version          int        `json:"version"`
	TerrafromVersion string     `json:"terraform_version"`
	Serial           int        `json:"serial"`
	Lineage          string     `json:"lineage"`
	Outputs          any        `json:"outputs"`
	Resources        []Resource `json:"resources"`
	CheckResults     any        `json:"check_results"`
}

type Resource struct {
	Module    string           `json:"module"`
	Mode      string           `json:"mode"`
	Type      string           `json:"type"`
	Name      string           `json:"name"`
	Provider  string           `json:"provider"`
	Instances []map[string]any `json:"instances"`
}

func Pull(path string) State {
	cmd := exec.Command("terraform", "state", "pull")
	cmd.Dir = path
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		println(out.String())
	}

	var state State
	err = json.Unmarshal(out.Bytes(), &state)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return state
}

type Query struct {
	Module string
	Type   string
	Name   string
}

func (state *State) GetResource(q Query) (Resource, error) {
	if q.Module != "" {
		q.Module = "module." + q.Module
		for _, r := range state.Resources {
			if r.Module == q.Module && r.Type == q.Type && r.Name == q.Name {
				return r, nil
			}
		}
	} else {
		for _, r := range state.Resources {
			if r.Type == q.Type && r.Name == q.Name {
				return r, nil
			}
		}
	}
	return Resource{}, errors.New("Resource not found")
}
