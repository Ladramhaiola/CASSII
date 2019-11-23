package main

type TaskGraph struct {
	nodes map[int]*Task
	edges map[[2]int]int
}

func (graph *TaskGraph) AddTasks(tasks ...[2]int) *TaskGraph {
	for _, t := range tasks {
		graph.nodes[t[0]] = newTask(t[0], t[1])
	}

	return graph
}

func (graph *TaskGraph) AddConnections(conns ...[3]int) *TaskGraph {
	for _, conn := range conns {
		src := graph.nodes[conn[0]]
		dst := graph.nodes[conn[1]]

		if src == nil || dst == nil {
			continue
		}

		key := [2]int{conn[0], conn[1]}
		graph.edges[key] = conn[2]

		src.sinks = append(src.sinks, dst)
		dst.parents = append(dst.parents, src)
	}

	return graph
}

func (graph *TaskGraph) CommunicationCost(src, dst int) int {
	if cost, ok := graph.edges[[2]int{src, dst}]; ok {
		return cost
	}
	return -1
}

func (graph *TaskGraph) SetCommunicationCost(src, dst, w int) {
	if _, ok := graph.edges[[2]int{src, dst}]; ok {
		graph.edges[[2]int{src, dst}] = w
	}
}

// u - current node, v - immediate successor
func (graph *TaskGraph) F(u, v int) int {
	return graph.nodes[u].w + graph.CommunicationCost(u, v) + graph.nodes[v].f
}

// v - current node, u - immediate predecessor
// s = max{s(u) + w(u) + conn(u, v)|(u, v) in E}
func (graph *TaskGraph) S(v int) int {
	maxS := 0
	curr := graph.nodes[v]

	for _, parent := range curr.parents {
		s := parent.s + graph.CommunicationCost(parent.id, v) + parent.w
		if s > maxS {
			maxS = s
		}
	}

	return maxS
}

func (graph *TaskGraph) DominantSuccessor(u int) (*Task, int) {
	src, ok := graph.nodes[u]
	if !ok {
		return nil, -1
	}

	if len(src.sinks) == 0 {
		return nil, -1
	}

	maxF := 0
	successor := src.sinks[0]

	for _, sink := range src.sinks {
		f := graph.F(u, sink.id)

		if f >= maxF {
			maxF = f
			successor = sink
		}
	}

	return successor, maxF
}

func (graph *TaskGraph) TopologicalList() (L []*Task) {
	S := Set()

	for _, task := range graph.nodes {
		if len(task.parents) == 0 {
			S.Insert(task)
		}
	}

	for S.Len() > 0 {
		n := S.Pop()
		L = append(L, n)

		n.marked = true
		for _, sink := range n.sinks {
			if checkMarked(sink) {
				S.Insert(sink)
			}
		}
	}

	for _, task := range graph.nodes {
		task.marked = false
	}

	return L
}

func (graph *TaskGraph) SetInitialSLevel() {
	for _, task := range graph.TopologicalList() {
		task.s = graph.S(task.id)
	}
}

func checkMarked(sink *Task) bool {
	for _, parent := range sink.parents {
		if !parent.marked {
			return false
		}
	}
	return true
}

func Graph() *TaskGraph {
	return &TaskGraph{
		nodes: make(map[int]*Task),
		edges: make(map[[2]int]int),
	}
}
