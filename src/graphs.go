package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jupp0r/go-priority-queue"
	"math"
	"os"
	"sort"
)

const resourcesDir = "C:\\Users\\User\\Desktop\\SSUhomework\\GraphsGoSSU\\resources\\"
const INF = math.MaxInt32
const NEGATIVE_INF = math.MinInt32

type Graph struct {
	adjList    map[string]map[string]int
	isDirected bool
	used       map[string]bool
	minEdge    map[string]int
}

func NewEmptyGraph(directed bool) *Graph {
	return &Graph{
		adjList:    make(map[string]map[string]int),
		isDirected: directed,
		used:       make(map[string]bool),
		minEdge:    make(map[string]int),
	}
}

func (g *Graph) UsedClear() {
	g.used = make(map[string]bool)
}

func (g *Graph) MinEdgeClear() {
	g.minEdge = make(map[string]int)
}

func (g *Graph) minEdgeMakeDistInfinity() {
	for v := range g.adjList {
		g.minEdge[v] = INF
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
	for v, neighbours := range g.adjList {
		cntVertices++
		for u := range neighbours {
			if g.isDirected || v <= u {
				cntEdges++
			}
		}
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
		vertices, edges := gCopy.countVerticesAndEdges()
		if vertices == edges+1 && gCopy.countConnectedComponents() == 1 {
			return true, nil
		}
	}
	return false, nil
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

	cntIn := make(map[string]int)

	for _, neighbours := range g.adjList {
		for u := range neighbours {
			cntIn[u]++
		}
	}

	for _, cnt := range cntIn {
		if cnt > 1 {
			return "Not a tree and not a forest", nil
		}
	}

	isTree := false

	for v, _ := range g.adjList {
		g.UsedClear()
		isWithoutCycle := g.isWithoutCyclesBFS(v)
		if !isWithoutCycle {
			return "Not a tree and not a forest", nil
		}
		if len(g.used) == len(g.adjList) {
			isTree = true
		}
	}

	if isTree {
		return "Tree", nil
	} else {
		return "Forest", nil
	}
}

type Item struct {
	previous, current string
	weight            int
}

// MSTPrime => task 7: Каркас III
// Дан взвешенный неориентированный граф из N вершни и M ребер. Найти MST с помощью алгоритма Прима
func (g *Graph) MSTPrime() (*Graph, int, error) {
	if g.isDirected {
		return nil, -1, errors.New("graph is directed")
	}

	g.UsedClear()
	if g.countConnectedComponents() > 1 {
		return nil, -1, errors.New("graph is not connected")
	}

	g.UsedClear()
	g.MinEdgeClear()
	g.minEdgeMakeDistInfinity()
	var start string
	for v := range g.adjList {
		start = v
		g.minEdge[start] = 0
		break
	}

	totalWeight := 0
	MST := NewEmptyGraph(false)
	priorityQueue := pq.New()
	priorityQueue.Insert(Item{previous: "", current: start, weight: 0}, 0)

	for priorityQueue.Len() > 0 {
		el, _ := priorityQueue.Pop()
		item := el.(Item)
		if g.used[item.current] {
			continue
		}
		g.used[item.current] = true
		totalWeight += item.weight
		if item.previous != "" {
			MST.AddEdge(item.previous, item.current, item.weight)
		}

		for to, w := range g.adjList[item.current] {
			if !g.used[to] && w < g.minEdge[to] {
				g.minEdge[to] = min(g.minEdge[to], w)
				priorityQueue.Insert(Item{previous: item.current, current: to, weight: w}, float64(g.minEdge[to]))
			}
		}
	}

	return MST, totalWeight, nil
}

// Алгоритм Флойда поиска расстояний между всеми парами вершин, для task 10 №9
func (g *Graph) FloydWarshall() (map[string]map[string]int, bool) {
	n := len(g.adjList)
	dist := make(map[string]map[string]int, n)
	for v := range g.adjList {
		dist[v] = make(map[string]int, n)
		for u := range g.adjList {
			if w, exist := g.adjList[v][u]; exist {
				dist[v][u] = w
			} else if v == u {
				dist[v][u] = 0
			} else {
				dist[v][u] = INF
			}
		}
	}

	//k - промежуточная вершина на пути между v и u
	for k := range g.adjList {
		for v := range g.adjList {
			for u := range g.adjList {
				if dist[v][k] < INF && dist[k][u] < INF {
					dist[v][u] = min(dist[v][u], dist[v][k]+dist[k][u])
				}
			}
		}
	}

	//проверка на отрицательный цикл
	for v := range g.adjList {
		if dist[v][v] < 0 {
			return dist, true
		}
	}

	return dist, false
}

// task10 9.Вывести длины кратчайших путей для всех пар вершин.
func (g *Graph) PrintFloydWarshall() {
	dist, hasNegativeCycles := g.FloydWarshall()
	if hasNegativeCycles {
		fmt.Println("graph has a negative cycle")
		return
	}
	fmt.Print("from\t")
	vertices := make([]string, 0)
	for el := range dist {
		vertices = append(vertices, el)
	}
	sort.Strings(vertices)
	for _, el := range vertices {
		fmt.Printf("%s\t", el)
	}
	fmt.Println()
	for _, v := range vertices {
		fmt.Printf("%s\t\t", v)
		for _, u := range vertices {
			fmt.Printf("%d\t", dist[v][u])
		}
		fmt.Println()
	}
}

type EdgeStruct struct {
	from string
	to   string
	dist int
}

func (g *Graph) BellmanFord(u, v string) ([]string, int, error) {
	n := len(g.adjList)
	dist := make(map[string]int, n)
	for vert := range g.adjList {
		dist[vert] = INF
	}
	dist[u] = 0

	path := make(map[string]string)

	edges := make([]EdgeStruct, 0)
	for v1, neighbours := range g.adjList {
		for v2, d := range neighbours {
			edges = append(edges, EdgeStruct{from: v1, to: v2, dist: d})
		}
	}

	for i := 0; i < n-1; i++ {
		for _, edge := range edges {
			if dist[edge.from] != INF && dist[edge.from]+edge.dist < dist[edge.to] {
				dist[edge.to] = dist[edge.from] + edge.dist
				path[edge.to] = edge.from
			}
		}
	}

	if dist[v] == INF {
		return nil, INF, errors.New("there is no path from u to v")
	}

	answerPath := make([]string, 0)
	curVertex := v
	for path[curVertex] != "" {
		answerPath = append(answerPath, curVertex)
		curVertex = path[curVertex]
		if path[curVertex] == "" {
			answerPath = append(answerPath, curVertex)
		}
	}

	for i, j := 0, len(answerPath)-1; i < j; i, j = i+1, j-1 {
		answerPath[i], answerPath[j] = answerPath[j], answerPath[i]
	}

	return answerPath, dist[v], nil
}

type VertexWeight struct {
	weight int
	vertex string
}

// Алгоритм Дейкстры для task 8: 10
func (g *Graph) Dijkstra(u string) map[string]int {
	d := make(map[string]int, len(g.adjList))
	for vertex := range g.adjList {
		d[vertex] = INF
	}
	d[u] = 0

	path := make(map[string]string)

	priorityQueue := pq.New()
	priorityQueue.Insert(VertexWeight{0, u}, 0)

	for priorityQueue.Len() > 0 {
		topElement, _ := priorityQueue.Pop()
		top := topElement.(VertexWeight)
		if top.weight > d[top.vertex] {
			continue
		}
		for to := range g.adjList[top.vertex] {
			if d[top.vertex]+g.adjList[top.vertex][to] < d[to] {
				d[to] = d[top.vertex] + g.adjList[top.vertex][to]
				path[to] = top.vertex
				priorityQueue.Insert(VertexWeight{d[to], to}, float64(d[to]))
			}
		}
	}

	return d
}

// GetRadius task 8: "Вес IV a"
// 10. Эксцентриситет вершины — максимальное расстояние из всех минимальных расстояний
// от других вершин до данной вершины.
// Найти радиус графа — минимальный из эксцентриситетов его вершин

func (g *Graph) GetRadius() (string, string, int) {
	dist := make(map[string]map[string]int)
	for v := range g.adjList {
		dist[v] = g.Dijkstra(v)
	}

	eccentricities := make(map[string]EdgeStruct, len(dist))
	for v, neighbours := range dist {
		maxDist := EdgeStruct{from: v, to: v, dist: NEGATIVE_INF}
		for u, curDist := range neighbours {
			if curDist > maxDist.dist {
				maxDist.to = u
				maxDist.dist = curDist
			}
		}
		eccentricities[v] = maxDist
	}

	minEccentricity := EdgeStruct{dist: INF}
	for _, struc := range eccentricities {
		if struc.dist < minEccentricity.dist {
			minEccentricity = struc
		}
	}

	return minEccentricity.from, minEccentricity.to, minEccentricity.dist
}
