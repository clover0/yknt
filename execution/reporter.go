package execution

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"yknk/task"
)

const indent = "    "

type Reporter struct {
	reportEntries sync.Map
}

type reportEntry struct {
	buff    bytes.Buffer
	info    string
	nest    int
	failed  bool
	timeout bool
}

func (re *reportEntry) Init(t task.Task, nest int) {
	re.nest = nest
	info := fmt.Sprintf("%sname: %s\n", strings.Repeat(indent, nest), t.Name)
	info += fmt.Sprintf("%sexec: %s %s \n", strings.Repeat(indent, nest), t.Command, t.Args)
	re.info = info
}

func (re *reportEntry) Write(b []byte) (n int, err error) {
	return re.buff.Write(b)
}

func (re *reportEntry) OnFail(e task.Task) {
	re.failed = true
}

func (re *reportEntry) OnSuccess(e task.Task) {
	re.failed = false
}

func (re *reportEntry) OnTimeout(e task.Task) {
	re.failed = true
	re.timeout = true
}

func (re *reportEntry) String() string {
	return re.buff.String()
}

func (r *Reporter) MustGet(key string) *reportEntry {
	e := &reportEntry{}
	if v, ok := r.reportEntries.LoadOrStore(key, e); ok {
		if r, ok := v.(*reportEntry); ok {
			return r
		}
		panic(v.(*reportEntry))
	}
	return e
}

func (r *Reporter) Report(order []string) {
	fmt.Println(strings.Repeat("=", 30))
	totalCount := len(order)
	errCount := 0
	for _, k := range order {
		v, ok := r.reportEntries.Load(k)
		if !ok {
			fmt.Printf("not found key '%s'\n", k)
			continue
		}
		re, ok := v.(*reportEntry) //FIXME: not use sync.Map
		if !ok {
			fmt.Println("can not cast v.(*reportEntry)")
			continue
		}
		fmt.Print(re.info)
		timeout := ""
		if re.timeout {
			timeout = "(timeout)\n"
		}
		if re.failed {
			errCount++
			fmt.Print(strings.Repeat(indent, re.nest) + fmt.Sprintf(" => Fail:%s", timeout))
			fmt.Print(re.String())
			continue
		}
		//fmt.Println(idt + "Success!")
	}
	fmt.Println(strings.Repeat("=", 30))
	fmt.Printf("Total: %d    Success: %d    Error: %d    ", totalCount, totalCount-errCount, errCount)
}
