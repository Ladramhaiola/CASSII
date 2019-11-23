package main

import "fmt"

type Cluster struct {
	f int

	scheduled []*Task
}

func (c *Cluster) Insert(graph *TaskGraph, t *Task) {
	if len(c.scheduled) == 0 {
		t.SetF(t.w)
		c.f = t.f
	} else {
		// create new pseudo connection
		successor := c.scheduled[0]
		graph.AddConnections([3]int{t.id, successor.id, 0})

		t.SetF(t.w + c.ExecutionTime())

		// update all immediate successors s
		for _, sink := range t.sinks {
			sink.SetS(graph.S(sink.id))
		}
	}

	t.cluster = c
	c.scheduled = append([]*Task{t}, c.scheduled...)
}

func (c *Cluster) ExecutionTime() (f int) {
	for _, task := range c.scheduled {
		f += task.w
	}
	return
}

func (c *Cluster) Acceptable(t *Task) bool {
	insertionF := t.w + c.ExecutionTime()
	// hmm
	return insertionF <= t.f && !(t.s == 0 && c.scheduled[0].s == 0)
}

func (c *Cluster) String() string {
	res := "Cluster {\n"
	for _, task := range c.scheduled {
		res += "\t" + fmt.Sprintf("%+v\n", task)
	}
	return res
}