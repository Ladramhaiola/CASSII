package main

import (
	"fmt"
	"log"
	"math"
	"os"
)

func CASS(graph *TaskGraph) []*Cluster {
	var clusters []*Cluster
	var unscheduled = make(map[int]*Task)

	for id, task := range graph.nodes {
		unscheduled[id] = task
	}

	graph.SetInitialSLevel()

	// create cluster for each exit node
	for _, task := range unscheduled {
		if len(task.sinks) == 0 {
			cluster := &Cluster{graph: graph}
			cluster.Insert(task)
			clusters = append(clusters, cluster)

			delete(unscheduled, task.id)
		}
	}

	queue := Queue()
	for len(unscheduled) > 0 {
		for _, task := range unscheduled {
			if task.Ready() {
				ds, f := graph.DominantSuccessor(task.id)
				task.SetF(f)
				queue.Insert(task.id, task.l, ds.id)
			}
		}

		x, _, y := queue.DeleteMaxLValue()
		src := graph.nodes[x]
		dst := graph.nodes[y]

		if src.AllChildrenSinks() {
			// merge sink clusters if possible
			mergedF := 0
			for _, child := range src.sinks {
				if child != dst {
					mergedF += child.w
					if mergedF <= src.f {
						child.cluster.scheduled = nil
						dst.cluster.Insert(child)
					}
				}
			}
		}

		if dst.cluster.Acceptable(src) {
			dst.cluster.Insert(src)
		} else {
			cluster := &Cluster{graph: graph}
			cluster.Insert(src)
			clusters = append(clusters, cluster)
		}

		delete(unscheduled, x)
	}

	// compute final starting time for each task
	graph.SetInitialSLevel()

	res := []*Cluster{}
	for _, c := range clusters {
		if len(c.scheduled) > 0 {
			res = append(res, c)
		}
	}

	return res
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
	for _, cluster := range processors {
		fmt.Println(cluster)
	}

	// dynamic mapping
	pool := taskGraph.TopologicalList()
	for i := 0; i < len(pool); i++ {
		for j := 0; j < len(pool); j++ {
			if pool[i].s+pool[i].w < pool[j].s+pool[j].w {
				pool[i], pool[j] = pool[j], pool[i]
			}
		}
	}

	bridge := make(map[int][]*Task)
	transfers := make(map[[2]*Task]int)

	for _, task := range pool {
		if len(task.sinks) != 0 {
			clusters := make(map[*Cluster]bool)

			for i, child := range task.sinks {
				if child.cluster != task.cluster {
					if _, ok := clusters[child.cluster]; !ok {
						start := EarliestFree(bridge, task.s+task.w+i)
						bridge[start] = []*Task{task, child}
						transfers[[2]*Task{task, child}] = start
						clusters[child.cluster] = true
					}
				}
			}
		}
	}

	// final task start times setup
	for _, task := range taskGraph.TopologicalList() {
		task.SetS(taskGraph.BridgeS(task.id, transfers))
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

func EarliestFree(b map[int][]*Task, start int) int {
	for i := start; i < math.MaxInt32; i++ {
		if _, ok := b[i]; !ok {
			return i
		}
	}
	return -1
}
