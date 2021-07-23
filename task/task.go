package task

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"time"
)

type Task struct {
	Name      string
	Concurrent  bool
	Command   string
	Args      []string
	Tasks     []Task
	isFailed  bool
	isTimeout bool
}

type Callback interface {
	OnFail(e Task)
	OnTimeout(e Task)
	OnSuccess(e Task)
}

func (e *Task) RunCommand(stdoutWriter io.Writer, stderrWriter io.Writer, callbacks ...Callback) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.Command, e.Args...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		e.onFail(callbacks...)
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text() + "\n"
			if _, err := stdoutWriter.Write([]byte(line)); err != nil {
				e.onFail(callbacks...)
				return
			}
		}

	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text() + "\n"
			if _, err := stderrWriter.Write([]byte(line)); err != nil {
				e.onFail(callbacks...)
				return
			}
		}
	}()

	err = cmd.Wait()

	if ctx.Err() == context.DeadlineExceeded {
		e.onTimeout(callbacks...)
		return
	}

	if err != nil {
		e.onFail(callbacks...)
		return
	}
	e.onSuccess(callbacks...)
}

func (e *Task) onTimeout(callbacks ...Callback) {
	e.isFailed = true
	e.isTimeout = true
	for _, c := range callbacks {
		c.OnTimeout(*e)
	}
}

func (e *Task) onFail(callbacks ...Callback) {
	e.isFailed = true
	for _, c := range callbacks {
		c.OnFail(*e)
	}
}

func (e *Task) onSuccess(callbacks ...Callback) {
	e.isFailed = false
	for _, c := range callbacks {
		c.OnSuccess(*e)
	}
}
