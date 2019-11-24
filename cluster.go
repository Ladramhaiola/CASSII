package main

import (
	"fmt"
	"strings"
)

type Cluster struct {
	f     int
	graph *TaskGraph

	scheduled []*Task
}

func (c *Cluster) Insert(t *Task) *Cluster {
	if len(c.scheduled) == 0 {
		t.SetF(t.w)
		c.f = t.f
	} else {
		// create new pseudo connections
		for _, successor := range c.scheduled {
			if c.graph.CommunicationCost(t.id, successor.id) != -1 {
				c.graph.SetCommunicationCost(t.id, successor.id, 0)
			} else {
				c.graph.AddConnections([3]int{t.id, successor.id, 0})
			}
		}

		t.SetF(t.w + c.ExecutionTime())
	}

	t.cluster = c
	c.scheduled = append([]*Task{t}, c.scheduled...)

	return c
}

func (c *Cluster) ExecutionTime() (f int) {
	for _, task := range c.scheduled {
		f += task.w
	}
	return
}

// check if task's insertion won't increase f(t)
// additional check to separate independent entry tasks into different clusters
func (c *Cluster) Acceptable(t *Task) bool {
	insertionF := t.w + c.ExecutionTime()
	return insertionF <= t.f && !(t.s == 0 && c.scheduled[0].s == 0)
}

func (c *Cluster) String() string {
	var res strings.Builder
	res.WriteString("Cluster {\n")

	for _, task := range c.scheduled {
		res.WriteString("\t")
		res.WriteString(task.String())

		// show parent task's from another clusters and communication cost
		if c.graph != nil {
			dependencies := make(map[int]int)
			for _, p := range task.parents {
				if !Contains(c.scheduled, p) {
					dependencies[p.id] = c.graph.CommunicationCost(p.id, task.id)
				}
			}

			if len(dependencies) > 0 {
				res.WriteString(fmt.Sprintf("\tdep: %v", dependencies))
			}
		}

		res.WriteString("\n")
	}

	res.WriteString("}")
	return res.String()
}
