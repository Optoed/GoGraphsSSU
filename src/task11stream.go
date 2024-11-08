package main

import (
	"math"
)

//В худшем случае сложность алгоритма Форда-Фалкерсона составляет O(E * F)
// E — количество рёбер в графе,
// F — максимальный поток от истока до стока

// MaxFlowFordFulkerson структура для хранения графа
type MaxFlowFordFulkerson struct {
	graph    [][]int        // список смежности
	capacity map[[2]int]int // ёмкости рёбер
}

// NewMaxFlowFordFulkerson создает граф с n вершинами
func NewMaxFlowFordFulkerson(n int) *MaxFlowFordFulkerson {
	return &MaxFlowFordFulkerson{
		graph:    make([][]int, n),
		capacity: make(map[[2]int]int),
	}
}

// AddEdge добавляет ребро с ёмкостью
func (m *MaxFlowFordFulkerson) AddEdge(u, v, cap int) {
	m.graph[u] = append(m.graph[u], v)
	m.graph[v] = append(m.graph[v], u) // обратное ребро для остаточной сети
	m.capacity[[2]int{u, v}] = cap
	m.capacity[[2]int{v, u}] = 0 // начальная ёмкость обратного ребра
}

// dfs выполняет поиск увеличивающего пути и возвращает поток
func (m *MaxFlowFordFulkerson) dfs(s, t, flow int, visited []bool) int {
	if s == t {
		return flow
	}
	visited[s] = true

	for _, neighbor := range m.graph[s] {
		if !visited[neighbor] && m.capacity[[2]int{s, neighbor}] > 0 {
			minCap := int(math.Min(float64(flow), float64(m.capacity[[2]int{s, neighbor}])))
			result := m.dfs(neighbor, t, minCap, visited)
			if result > 0 {
				// обновляем ёмкости рёбер
				m.capacity[[2]int{s, neighbor}] -= result
				m.capacity[[2]int{neighbor, s}] += result
				return result
			}
		}
	}
	return 0
}

// MaxFlow вычисляет максимальный поток от истока s до стока t
func (m *MaxFlowFordFulkerson) MaxFlow(s, t int) int {
	maxFlow := 0
	for {
		visited := make([]bool, len(m.graph))
		flow := m.dfs(s, t, math.MaxInt32, visited)
		if flow == 0 {
			break
		}
		maxFlow += flow
	}
	return maxFlow
}
