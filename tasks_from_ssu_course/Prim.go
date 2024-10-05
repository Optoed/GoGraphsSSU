package main

import (
	"container/heap"
	"fmt"
	"math"
)

type Edge struct {
	to, weight int
}

type Item struct {
	vertex, weight int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].weight < pq[j].weight
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func Prim(graph [][]Edge, N int) int {
	totalWeight := 0

	visited := make([]bool, N)
	minEdge := make([]int, N)
	for i := range minEdge {
		minEdge[i] = math.MaxInt64
	}

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, Item{vertex: 0, weight: 0})
	minEdge[0] = 0

	for pq.Len() > 0 {
		item := heap.Pop(pq).(Item)
		v := item.vertex

		if visited[v] {
			continue
		}

		visited[v] = true
		totalWeight += item.weight

		for _, edge := range graph[v] {
			if !visited[edge.to] && edge.weight < minEdge[edge.to] {
				minEdge[edge.to] = edge.weight
				heap.Push(pq, Item{vertex: edge.to, weight: edge.weight})
			}
		}
	}

	return totalWeight
}

func main() {
	N := 5 // количество вершин
	graph := make([][]Edge, N)
	// Добавляем рёбра (пример графа)
	graph[0] = append(graph[0], Edge{to: 1, weight: 2}, Edge{to: 3, weight: 6})
	graph[1] = append(graph[1], Edge{to: 0, weight: 2}, Edge{to: 2, weight: 3}, Edge{to: 3, weight: 8})
	graph[2] = append(graph[2], Edge{to: 1, weight: 3}, Edge{to: 3, weight: 7})
	graph[3] = append(graph[3], Edge{to: 0, weight: 6}, Edge{to: 1, weight: 8}, Edge{to: 2, weight: 7})

	result := Prim(graph, N)
	fmt.Printf("Минимальный вес остовного дерева: %d\n", result)
}
