package main

import (
	"yknk/execution"
	"yknk/mapper"
)

func main(){

	m := mapper.ReadYaml("./examples/example.yaml")
	tasks := mapper.MkTasks(m.M.TaskMaps)

	e := execution.NewExecutor()

	e.Run(tasks, 0, m.M.Concurrent)
	e.Report(m.Keys)
}
