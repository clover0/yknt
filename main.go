package main

import (
	"flag"
	"yknk/execution"
	"yknk/mapper"
)

func main() {
	yamlpath := flag.String("in", "./examples/example.yaml", "Execution YAML file.")

	flag.Parse()

	m := mapper.ReadYaml(*yamlpath)
	tasks := mapper.MkTasks(m.M.TaskMaps)

	e := execution.NewExecutor()

	e.Run(tasks, 0, m.M.Concurrent)
	e.Report(m.Keys)
}
