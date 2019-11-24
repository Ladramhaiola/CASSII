package main

import "fmt"

type Task struct {
	id int

	w, s, f, l int

	parents []*Task
	sinks   []*Task

	cluster *Cluster

	marked bool
}

func (t *Task) SetF(f int) *Task {
	t.f = f
	t.l = t.s + f
	return t
}

func (t *Task) SetS(s int) *Task {
	t.s = s
	t.l = t.s + t.f
	return t
}

func (t *Task) Ready() bool {
	if t.cluster != nil {
		return false
	}

	for _, sink := range t.sinks {
		if sink.cluster == nil {
			return false
		}
	}
	return true
}

func (t *Task) AllChildrenSinks() bool {
	if len(t.sinks) < 2 {
		return false
	}

	for _, sink := range t.sinks {
		if len(sink.sinks) != 0 {
			return false
		}
	}
	return true
}

func (t *Task) String() string {
	return fmt.Sprintf("T%d w:%d s:%d f:%d l:%d", t.id, t.w, t.s, t.f, t.l)
}

func newTask(id, w int) *Task {
	return &Task{
		id: id,
		w:  w,
	}
}
