package main

import (
	"log"
	"os"
	"sort"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type transfer struct {
	src   *Task
	dst   *Task
	start int
	end   int
}

func main() {
	taskGraph := Graph()
	taskGraph.AddTasks(
		[2]int{1, 3},
		[2]int{2, 5},
		[2]int{3, 4},
		[2]int{4, 6},
		[2]int{5, 2},
		[2]int{6, 2},
		[2]int{7, 4},
		[2]int{8, 3},
		[2]int{9, 6},
		[2]int{10, 2},
	)
	taskGraph.AddConnections(
		[3]int{1, 3, 8},
		[3]int{1, 4, 7},
		[3]int{1, 6, 6},

		[3]int{2, 3, 10},
		[3]int{2, 4, 14},
		[3]int{2, 5, 10},

		[3]int{3, 6, 5},
		[3]int{3, 8, 8},

		[3]int{4, 7, 12},
		[3]int{4, 9, 7},

		[3]int{5, 6, 4},
		[3]int{5, 7, 10},
		[3]int{5, 10, 14},

		[3]int{6, 8, 4},
		[3]int{6, 10, 8},

		[3]int{7, 9, 5},
	)

	processors := CASS(taskGraph)

	bridge := rbt.NewWithIntComparator()
	for _, task := range taskGraph.TopologicalList() {
		// sort by income time
		sort.SliceStable(task.sinks, func(i, j int) bool {
			return task.sinks[i].s < task.sinks[j].s
		})

		cache := make(map[*Cluster]bool)
		cache[task.cluster] = true

		for _, child := range task.sinks {
			if _, ok := cache[child.cluster]; !ok {
				start := EarliestFree(bridge, task.s+task.w, taskGraph.CommunicationCost(task.id, child.id))
				t := &transfer{src: task, dst: child, start: start, end: start + taskGraph.CommunicationCost(task.id, child.id)}
				bridge.Put(start, t)
				child.SetS(t.end)

				cache[child.cluster] = true
			}
		}
	}

	for _, task := range taskGraph.TopologicalList() {
		task.SetS(taskGraph.BridgeS(bridge, task.id))
	}

	f, err := os.Create("plan.md")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	table := Markdown(taskGraph, processors, bridge)

	if _, err := f.WriteString(table); err != nil {
		log.Fatalln(err)
	}
}

func EarliestFree(src *rbt.Tree, start, duration int) int {
	it := src.Iterator()
	for it.Next() {
		interval := it.Value().(*transfer)
		if interval.start-start >= duration {
			return start
		}

		if interval.end > start {
			start = interval.end
		}
	}
	return start
}
