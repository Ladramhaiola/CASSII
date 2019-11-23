package main

import (
	"math/rand"
	"time"
)

type TaskSet struct {
	data map[int]*Task
}

func (set *TaskSet) Insert(t *Task) {
	if _, ok := set.data[t.id]; !ok {
		set.data[t.id] = t
	}
}

func (set *TaskSet) Pop() (t *Task) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(set.data))

	index := 0
	for id, task := range set.data {
		if index == r {
			t = task
			delete(set.data, id)
			return
		}
		index++
	}
	return
}

func (set *TaskSet) Remove(id int) {
	delete(set.data, id)
}

func (set *TaskSet) Len() int {
	return len(set.data)
}

type PriorityQueue struct {
	data map[int][2]int
}

func (queue *PriorityQueue) Insert(u, l, v int) {
	// todo: debug insertions
	if _, ok := queue.data[u]; !ok {
		queue.data[u] = [2]int{l, v}
	}
}

func (queue *PriorityQueue) DeleteMaxLValue() (u, l, v int) {
	if len(queue.data) == 0 {
		return -1, -1, -1
	}

	maxL := 0
	for parent, it := range queue.data {
		if it[0] > maxL {
			maxL = it[0]
			u = parent
			v = it[1]
		}
	}
	delete(queue.data, u)
	return
}

func (queue *PriorityQueue) Len() int {
	return len(queue.data)
}

func Queue() *PriorityQueue {
	return &PriorityQueue{data: make(map[int][2]int)}
}

func Set() *TaskSet {
	return &TaskSet{data: make(map[int]*Task)}
}
