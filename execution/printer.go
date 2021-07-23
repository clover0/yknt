package execution

import (
	"bytes"
	"fmt"
	"sync"

	"yknk/task"
)

type Printer struct {
	printerEntries sync.Map
}

type printEntry struct {
	buff bytes.Buffer
	info string
}

func (re *printEntry) Write(b []byte) (n int, err error) {
	line := re.info
	line += string(b)
	fmt.Print(line)
	return len(b), nil
}

func (p *Printer) MustGet(key string) *printEntry {
	e := &printEntry{info: fmt.Sprintf("[%s]", key)}
	if v, ok := p.printerEntries.LoadOrStore(key, e); ok {
		if r, ok := v.(*printEntry); ok {
			return r
		}
		panic(v.(*printEntry))
	}
	return e
}

func (pe *printEntry) OnFail(t task.Task) {
}

func (pe *printEntry) OnSuccess(t task.Task) {
}

func (pe *printEntry) OnTimeout(t task.Task) {
}
