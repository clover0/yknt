package execution

import (
	"io"
	"sync"

	"yknk/task"
)

type Executor struct {
	printer  Printer
	reporter Reporter
}

func NewExecutor() Executor {
	return Executor{
		printer: Printer{},
	}
}

func (c *Executor) Run(ts []task.Task, nest int, concurrent bool) {
	var wg sync.WaitGroup
	for _, e := range ts {
		reporter := c.reporter.MustGet(e.Name)
		reporter.Init(e, nest)
		printer := c.printer.MustGet(e.Name)

		if e.Command != "" {
			wg.Add(1)
			go func(t task.Task, reporter *reportEntry, printer *printEntry) {
				defer wg.Done()

				wout := io.MultiWriter(printer)
				werr := io.MultiWriter(printer, reporter)

				t.RunCommand(wout, werr, printer, reporter)
			}(e, reporter, printer)

			if !concurrent {
				wg.Wait()
			}
		}
		c.Run(e.Tasks, nest+1, e.Concurrent)
	}
	wg.Wait()
}

func (e *Executor) Report(order []string) {
	e.reporter.Report(order)
}
