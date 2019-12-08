package main

import (
	"fmt"
	"strconv"
	"strings"
)

func Markdown(graph *TaskGraph, procs []*Cluster, transfers map[int][]*Task) string {
	solutionTime := 0
	for _, c := range procs {
		for _, task := range c.scheduled {
			if task.l > solutionTime {
				solutionTime = task.l
			}
		}
	}

	result := strings.Builder{}
	result.WriteString(Header(len(procs)))

	for i := 0; i < solutionTime; i++ {
		runningTasks := []string{}
		transfer := "|"

		for _, cluster := range procs {
			current := "-"
			for _, task := range cluster.scheduled {
				if task.s <= i && task.s+task.w > i {
					current = fmt.Sprintf("T%d", task.id)
				}
			}
			runningTasks = append(runningTasks, current)
		}
		if send, ok := transfers[i]; ok {
			var procid int
			for i, cluster := range procs {
				if cluster.Contains(send[1].id) {
					procid = i + 1
				}
			}

			transfer = fmt.Sprintf("T%d -> `P%d (%d)`", send[0].id, procid, graph.CommunicationCost(send[0].id, send[1].id))
		}
		result.WriteString(LineGen(i, runningTasks, transfer))
	}

	return result.String()
}

func Header(proccnt int) string {
	procnames := []string{}

	for i := 1; i <= proccnt; i++ {
		procnames = append(procnames, fmt.Sprintf("P%d", i))
	}

	header := "ticks |" + strings.Join(procnames, "|") + "|data transfer" + "\n"
	header += strings.Repeat(":---:|", proccnt+1) + ":---:\n"

	return header
}

func LineGen(tick int, tasks []string, transfer string) string {
	t := strconv.Itoa(tick)
	return t + "|" + strings.Join(tasks, "|") + "|" + transfer + "\n"
}
