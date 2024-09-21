package main

import (
	"encoding/json"
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

func NewGraphFromFileJSON(filename string) (*Graph, error) {
	filePath := resourcesDir + filename
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Структура для временного хранения данных
	var data struct {
		IsDirected bool                      `json:"isDirected"`
		AdjList    map[string]map[string]int `json:"adjList"`
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	graph := NewEmptyGraph(data.IsDirected)
	graph.adjList = data.AdjList

	return graph, nil
}

func (g *Graph) SaveToFileJSON(filename string) error {
	filePath := resourcesDir + filename
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Структура для сериализации
	data := struct {
		IsDirected bool                      `json:"isDirected"`
		AdjList    map[string]map[string]int `json:"adjList"`
	}{
		IsDirected: g.isDirected,
		AdjList:    g.adjList,
	}

	// Сериализуем граф в JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Пустая строка для начала и два пробела для отступов
	err = encoder.Encode(data)
	if err != nil {
		return err
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

// PrintAdjList Вывод списка смежности в консоль
func (g *Graph) PrintAdjList() {
	fmt.Println("isDirected :", g.isDirected)
	for u, neighbors := range g.adjList {
		fmt.Printf("%s -> ", u)
		for v, weight := range neighbors {
			fmt.Printf("%s(%d) ", v, weight)
		}
		fmt.Println()
	}
}
