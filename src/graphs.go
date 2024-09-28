package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Graph struct {
	adjList    map[string]map[string]int
	isDirected bool
	used       map[string]bool
}

func NewEmptyGraph(directed bool) *Graph {
	return &Graph{
		adjList:    make(map[string]map[string]int),
		isDirected: directed,
		used:       make(map[string]bool),
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
	fmt.Println("isDirected (Ориентированный?) :", g.isDirected)
	for u, neighbors := range g.adjList {
		fmt.Printf("%s -> ", u)
		for v, weight := range neighbors {
			fmt.Printf("%s(%d) ", v, weight)
		}
		fmt.Println()
	}
}

// VerticesFromUAndNotFromV task 2: adj list la 2
// Определить, существует ли вершина, в которую есть дуга из вершины u, но нет из v. Вывести такую вершину.
func (g *Graph) VerticesFromUAndNotFromV(u, v string) (bool, []string) {
	verticesList := make([]string, 0)
	for vertex, _ := range g.adjList[u] {
		if _, exist := g.adjList[v][vertex]; !exist {
			verticesList = append(verticesList, vertex)
		}
	}

	if len(verticesList) == 0 {
		return false, nil
	}
	return true, verticesList
}

// task 3 la: option 6
// Print all hanging vertices of the graph (of degree 1).
func (g *Graph) HangingVertices() (bool, []string) {
	hangingVertices := make([]string, 0)
	cntNeighbours := make(map[string]int)

	for u, neighbours := range g.adjList {
		cntNeighbours[u] += len(neighbours)
		if !g.isDirected {
			continue
		}
		for v := range neighbours {
			cntNeighbours[v]++
		}
	}

	for u, cnt := range cntNeighbours {
		if cnt == 1 {
			hangingVertices = append(hangingVertices, u)
		}
	}

	if len(hangingVertices) == 0 {
		return false, nil
	}
	return true, hangingVertices
}

// task4 : is g a subgraph of otherG
func (g *Graph) isSubgraphOf(otherG *Graph) bool {
	for v, neighbours := range g.adjList {
		if _, existVertex := otherG.adjList[v]; !existVertex {
			return false
		}
		for u, _ := range neighbours {
			if _, existEdge := otherG.adjList[v][u]; !existEdge {
				return false
			}
		}
	}
	return true
}

func (g *Graph) dfs(v string) {
	g.used[v] = true
	for u := range g.adjList[v] {
		if !g.used[u] {
			g.dfs(u)
		}
	}
}

func (g *Graph) countConnectedComponents() int {
	g.used = make(map[string]bool)
	cnt := 0
	for v := range g.adjList {
		if !g.used[v] {
			cnt++
			g.dfs(v)
		}
	}
	return cnt
}

func (g *Graph) countVerticesAndEdges() (int, int) {
	cntVertices, cntEdges := 0, 0
	for _, neighbours := range g.adjList {
		cntEdges++
		cntVertices += len(neighbours)
	}
	return cntVertices, cntEdges
}

// 19.	Проверить, можно ли из графа удалить какую-либо вершину так, чтобы получилось дерево.
func (g *Graph) isAlmostTree() (bool, error) {
	if g.isDirected == true {
		return false, errors.New("graph is directed")
	}
	if g.countConnectedComponents() > 2 {
		return false, errors.New("too much connected components")
	}
	for v := range g.adjList {
		gCopy := NewGraphCopy(*g)
		gCopy.RemoveVertex(v)
		vertices, edges := g.countVerticesAndEdges()
		if vertices == edges+1 && g.countConnectedComponents() == 1 {
			return true, nil
		}
	}
	return false, nil
}

func (g *Graph) isWithoutCyclesDFS(v string) bool {
	g.used[v] = true
	for u := range g.adjList[v] {
		if !g.used[u] {
			g.isWithoutCyclesDFS(u)
		} else {
			return false
		}
	}
	return true
}

func (g *Graph) isWithoutCyclesBFS(v string) bool {
	g.used[v] = true
	queue := make([]string, 0)
	queue = append(queue, v)
	for len(queue) > 0 {
		from := queue[0]
		queue = queue[1:]
		for u := range g.adjList[from] {
			if !g.used[u] {
				g.used[u] = true
				queue = append(queue, u)
			} else {
				return false
			}
		}
	}
	return true
}

// 22.	* Проверить, является ли орграф деревом, или лесом, или не является ни тем, ни другим.
func (g *Graph) isDirectedGraphTheTreeOrForest() (string, error) {
	if !g.isDirected {
		return "", errors.New("graph is not directed (граф не ориентированный)")
	}

	roots := make([]string, 0)
	for root := range g.adjList {
		isRoot := true
		for other := range g.adjList {
			if _, existEdge := g.adjList[other][root]; existEdge {
				isRoot = false
				break
			}
		}
		if isRoot {
			roots = append(roots, root)
		}
	}

	vertices, edges := g.countVerticesAndEdges()
	if len(roots) == 1 && g.countConnectedComponents() == 1 && vertices == edges+1 {
		return "Tree", nil
	}

	g.used = make(map[string]bool)
	for _, root := range roots {
		if !g.isWithoutCyclesBFS(root) {
			return "Not a tree and not a forest", nil
		}
	}
	return "Forest", nil
}
