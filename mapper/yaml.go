package mapper

import (
	"errors"
	"io/ioutil"

	"yknk/task"
)
import "gopkg.in/yaml.v2"

type Main struct {
	Version string `yaml:"version"`
	M       Flow   `yaml:"flow"`
	Keys    []string
}
type Flow struct {
	Name       string    `yaml:"name"`
	Concurrent bool      `yaml:"concurrent"`
	TaskMaps   []TaskMap `yaml:"tasks"`
}

type TaskMap struct {
	Name       string    `yaml:"name"`
	Concurrent bool      `yaml:"concurrent"`
	Command    string    `yaml:"command"`
	Args       []string  `yaml:"args"`
	TaskMaps   []TaskMap `yaml:"tasks"`
}

func MkTasks(tms []TaskMap) []task.Task {
	ts := make([]task.Task, len(tms))
	for i, t := range tms {
		ts[i] = task.Task{
			Name:       t.Name,
			Concurrent: t.Concurrent,
			Command:    t.Command,
			Args:       t.Args,
			Tasks:      MkTasks(t.TaskMaps),
		}
	}
	return ts
}

func ReadYaml(path string) Main {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	mapper := Main{}
	err = yaml.Unmarshal(f, &mapper)
	if err != nil {
		panic(err.Error())
	}
	keys := collectKeys(mapper.M.TaskMaps)
	mapper.Keys = keys
	return mapper
}

func collectKeys(e []TaskMap) []string {
	ret := make([]string, 0)
	for _, v := range e {
		if v.Name == "" {
			panic(errors.New("can't blank 'name'"))
		}
		ret = append(ret, v.Name)
		if len(v.TaskMaps) > 0 {
			ret = append(ret, collectKeys(v.TaskMaps)...)
		}
	}
	return ret
}
