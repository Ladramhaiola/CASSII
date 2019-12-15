package main

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
