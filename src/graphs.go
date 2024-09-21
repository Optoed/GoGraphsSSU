package main

import (
	"fmt"
	"os"
)

type Graph struct {
	adjList    map[string]map[string]int
	isDirected bool
}

func NewEmptyGraph(directed bool) *Graph {
	return &Graph{
		adjList:    make(map[string]map[string]int),
		isDirected: directed,
	}
}

func NewGraphCopy(g Graph) *Graph {
	newGraph := NewEmptyGraph(g.isDirected)

	// Создаем глубокую копию adjList
	newGraph.adjList = make(map[string]map[string]int)
	for u, neighbors := range g.adjList {
		newGraph.adjList[u] = make(map[string]int)
		for v, weight := range neighbors {
			newGraph.adjList[u][v] = weight
		}
	}

	return newGraph
}

const resourcesDir = "C:\\Users\\User\\Desktop\\Вуз домашка\\GraphsGoSSU\\resources\\"

func NewGraphFromFile(filename string, isDirected bool) (*Graph, error) {
	graph := NewEmptyGraph(isDirected)

	filePath := resourcesDir + filename
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var v, u string
	var weight int
	for {
		_, err := fmt.Fscanf(file, "%s %s %d\n", &u, &v, &weight)
		if err != nil {
			if err.Error() == "EOF" {
				break // Корректное завершение при достижении конца файла
			}
			return nil, err
		}
		graph.AddEdge(u, v, weight)
	}

	return graph, nil
}

func (g *Graph) SaveToFile(filename string) error {
	filePath := resourcesDir + filename
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for u, neighbors := range g.adjList {
		for v, weight := range neighbors {
			fmt.Fprintf(file, "%s %s %d\n", u, v, weight)
		}
	}
	return nil
}

func (g *Graph) AddVertex(v string) {
	if _, exist := g.adjList[v]; !exist {
		g.adjList[v] = make(map[string]int)
	}
}

func (g *Graph) AddEdge(u, v string, weight int) {
	g.AddVertex(v)
	g.AddVertex(u)
	g.adjList[u][v] = weight
	if !g.isDirected {
		g.adjList[v][u] = weight
	}
}

func (g *Graph) RemoveVertex(v string) {
	delete(g.adjList, v)
	for u := range g.adjList {
		delete(g.adjList[u], v)
	}
}

func (g *Graph) RemoveEdge(v, u string) {
	delete(g.adjList[v], u)
	if !g.isDirected {
		delete(g.adjList[u], v)
	}
}

// Вывод списка смежности в консоль
func (g *Graph) PrintAdjList() {
	for u, neighbors := range g.adjList {
		fmt.Printf("%s -> ", u)
		for v, weight := range neighbors {
			fmt.Printf("%s(%d) ", v, weight)
		}
		fmt.Println()
	}
}
