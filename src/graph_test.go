package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGraph_AddVertex(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddVertex("A")

	if _, exists := graph.adjList["A"]; !exists {
		t.Errorf("Ожидалась добавление вершины A, но она не была добавлена")
	}
}

func TestGraph_AddEdge(t *testing.T) {
	graph := NewEmptyGraph(true)
	graph.AddEdge("A", "B", 10)

	if _, exists := graph.adjList["A"]["B"]; !exists {
		t.Errorf("Ожидалось наличие ребра между A и B, но оно не было добавлено")
	}
	if _, exists := graph.adjList["B"]["A"]; !exists && !graph.isDirected {
		t.Errorf("Ожидалось наличие ребра между B и A, но оно не было добавлено")
	}
}

func TestGraph_RemoveVertex(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddVertex("A")
	graph.RemoveVertex("A")

	if _, exists := graph.adjList["A"]; exists {
		t.Errorf("Ожидалось удаление вершины A, но она все еще существует")
	}
}

func TestGraph_RemoveEdge(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)
	graph.RemoveEdge("A", "B")

	if _, exists := graph.adjList["A"]["B"]; exists {
		t.Errorf("Ожидалось удаление ребра между A и B, но оно все еще существует")
	}
	if _, exists := graph.adjList["B"]["A"]; !graph.isDirected && exists {
		t.Errorf("Ожидалось удаление ребра между B и A, но оно все еще существует")
	}
}

func TestGraph_SaveToFile(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)
	err := graph.SaveToFileJSON("test_graph.json")
	if err != nil {
		t.Errorf("Ожидалось отсутствие ошибки при сохранении графа в файл, получена %v", err)
	}
}

func TestGraph_LoadFromFile(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)
	err := graph.SaveToFileJSON("test_graph.json")
	if err != nil {
		t.Errorf("Ожидалось отсутствие ошибки при сохранении графа в файл, получена %v", err)
	}

	loadedGraph, err := NewGraphFromFileJSON("test_graph.json")
	if err != nil {
		t.Errorf("Ожидалось отсутствие ошибки при загрузке графа из файла, получена %v", err)
	}

	if _, exists := loadedGraph.adjList["A"]["B"]; !exists {
		t.Errorf("Ожидалось наличие ребра между A и B в загруженном графе, но оно не было найдено")
	}
}

func TestGraph_CopyGraph(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)

	copiedGraph := NewGraphCopy(*graph)

	if _, exists := copiedGraph.adjList["A"]["B"]; !exists {
		t.Errorf("Ожидалось наличие ребра между A и B в скопированном графе, но оно не было найдено")
	}

	graph.RemoveEdge("A", "B")

	if _, exists := copiedGraph.adjList["A"]["B"]; !exists {
		t.Errorf("Ожидалось наличие ребра между A и B в скопированном графе, но оно не было найдено." +
			"Вероятно, из-за отсутствия глубокого копирования")
	}
}

func TestGraph_RemoveEdgeNonExistent(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)
	graph.RemoveEdge("A", "C") // Удаление несуществующего ребра

	if len(graph.adjList["A"]) != 1 {
		t.Errorf("Ожидалось, что количество рёбер для A останется 1, получено %d", len(graph.adjList["A"]))
	}
}

// task 5: удалить вершину из неориентированного графа и получить дерево
func TestGraph_isAlmostTree1(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("1", "2", 1)
	graph.AddEdge("2", "3", 1)
	graph.AddEdge("1", "3", 1)

	ok, _ := graph.isAlmostTree()
	if ok == false {
		t.Error("Полученный ответ нет, хотя можем удалить вершину 1 и получить дерево.")
	}
}

func TestGraph_isAlmostTree2(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("1", "2", 1)
	graph.AddEdge("3", "2", 1)
	graph.AddEdge("3", "4", 1)
	graph.AddEdge("1", "3", 1)
	graph.AddEdge("1", "4", 1)

	ok, _ := graph.isAlmostTree()
	if ok != true {
		t.Error("Полученный ответ нет, хотя можем удалить вершину 3 и получить дерево.")
	}
}

func TestGraph_isAlmostTree3(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("1", "2", 1)
	graph.AddEdge("3", "2", 1)
	graph.AddEdge("3", "4", 1)
	graph.AddEdge("1", "3", 1)
	graph.AddEdge("1", "4", 1)
	graph.AddEdge("2", "4", 1)

	ok, _ := graph.isAlmostTree()
	if ok == true {
		t.Error("Полученный ответ да, хотя никак не можем получить дерево удалением вершины")
	}
}

// task 6: проверить, является ли орграф деревом или лесом
func TestGraph_isDirectedGraphTheTreeOrForest1(t *testing.T) {
	dirgraph := NewEmptyGraph(true)
	dirgraph.AddEdge("1", "2", 1)
	dirgraph.AddEdge("2", "4", 1)
	dirgraph.AddEdge("1", "3", 1)
	dirgraph.AddEdge("2", "5", 1)

	answer, err := dirgraph.isDirectedGraphTheTreeOrForest()
	if err != nil {
		t.Error(err)
	}
	if answer != "Tree" {
		t.Errorf("this is a tree, but answer is %s\n", answer)
	} else {
		t.Log(answer)
	}
}

func TestGraph_isDirectedGraphTheTreeOrForest2(t *testing.T) {
	dirgraph := NewEmptyGraph(true)
	dirgraph.AddEdge("1", "2", 1)
	dirgraph.AddEdge("2", "3", 1)
	dirgraph.AddEdge("3", "1", 1)

	answer, err := dirgraph.isDirectedGraphTheTreeOrForest()
	if err != nil {
		t.Error(err)
	}
	if answer != "Not a tree and not a forest" {
		t.Errorf("True answer: Not a tree and not a forest. But we get: %s\n", answer)
	} else {
		t.Log(answer)
	}
}

func TestGraph_isDirectedGraphTheTreeOrForest3(t *testing.T) {
	dirgraph := NewEmptyGraph(true)
	dirgraph.AddEdge("1", "2", 1)
	dirgraph.AddEdge("3", "4", 1)
	dirgraph.AddVertex("5")

	answer, err := dirgraph.isDirectedGraphTheTreeOrForest()
	if err != nil {
		t.Error(err)
	}
	if answer != "Forest" {
		t.Errorf("Forest, but answer is %s\n", answer)
	}
	t.Log(answer)
}

func TestGraph_isDirectedGraphTheTreeOrForest4(t *testing.T) {
	dirgraph := NewEmptyGraph(true)
	dirgraph.AddEdge("2", "1", 1)
	dirgraph.AddEdge("3", "1", 1)

	answer, err := dirgraph.isDirectedGraphTheTreeOrForest()
	if err != nil {
		t.Error(err)
	}
	if answer != "Not a tree and not a forest" {
		t.Errorf("Not a tree and not a forest, but answer is %s\n", answer)
	}
	t.Log(answer)
}

// Task 7: MST (Prime)
func TestGraph_MSTPrime1(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("a", "b", 1)
	graph.AddEdge("a", "c", 2)
	graph.AddEdge("b", "d", 1)
	graph.AddEdge("c", "d", 1)
	graph.AddEdge("d", "e", 1)
	graph.AddEdge("d", "f", 50)
	graph.AddEdge("e", "f", 1)
	graph.AddEdge("f", "g", 1)
	graph.AddEdge("e", "g", 3)

	mst, totalWeight, err := graph.MSTPrime()
	if err != nil {
		t.Error(err)
	}
	t.Log("totalWeight = ", totalWeight)
	t.Log("MST:\n")
	mst.PrintAdjList()
	if totalWeight != (1 + 1 + 1 + 1 + 1 + 1) {
		t.Error("wrong weight: true weight is 6, but the answer is ", totalWeight)
	}
}

// Task 8: radius
func TestGraph_GetRadius(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("a", "b", 1)
	graph.AddEdge("a", "c", 2)
	graph.AddEdge("b", "d", 1)
	graph.AddEdge("c", "d", 1)
	graph.AddEdge("d", "e", 1)
	graph.AddEdge("d", "f", 50)
	graph.AddEdge("e", "f", 1)
	graph.AddEdge("f", "g", 1)
	graph.AddEdge("e", "g", 3)

	from, to, radius := graph.GetRadius()
	if !(radius == 3 && (from == "d" && to == "g" || from == "e" && to == "a")) {
		t.Errorf("the result is: radius = %d, from = %s, to = %s,"+
			" but true answer is: radius = 3, from = d and to = e or from = e and to = a", radius, from, to)
	} else {
		fmt.Printf("the result is right and equals: radius = %d, from = %s, to = %s\n", radius, from, to)
	}
}

// task 9: FordBellman
func TestGraph1_FordBellman(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("a", "b", 1)
	graph.AddEdge("a", "c", 2)
	graph.AddEdge("b", "d", 1)
	graph.AddEdge("c", "d", 1)
	graph.AddEdge("d", "e", 1)
	graph.AddEdge("d", "f", 50)
	graph.AddEdge("e", "f", 1)
	graph.AddEdge("f", "g", 1)
	graph.AddEdge("e", "g", 3)

	path, dist, err := graph.BellmanFord("g", "a")
	if err != nil {
		t.Error("Expected path = {g, f, e, d, b, a}, dist = 5, but catch error:" + err.Error())
	}

	ansString := "Expected path = {g, f, e, d, b, a}, dist = 5, we got: path = {"
	for i, el := range path {
		ansString += el
		if i != len(path)-1 {
			ansString += ", "
		}
	}
	ansString += "}, "
	ansString += "dist = " + strconv.Itoa(dist)

	checkPath := true
	truePath := []string{"g", "f", "e", "d", "b", "a"}
	if len(path) == len(truePath) {
		for i := range len(path) {
			if truePath[i] != path[i] {
				checkPath = false
				break
			}
		}
	} else {
		checkPath = false
	}

	if !checkPath {
		t.Error(ansString)
	}

	fmt.Println(ansString)
}

// task 10 (dist between all pairs)
func TestGraph_PrintFloydWarshall(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("a", "b", 1)
	graph.AddEdge("a", "c", 2)
	graph.AddEdge("b", "d", 1)
	graph.AddEdge("c", "d", 1)
	graph.AddEdge("d", "e", 1)
	graph.AddEdge("d", "f", 50)
	graph.AddEdge("e", "f", 1)
	graph.AddEdge("f", "g", 1)
	graph.AddEdge("e", "g", 3)

	graph.PrintFloydWarshall()
}

// task 11 (Алгоритм Форда-Фалкерсона поиска максимального потока
func TestGraph_FordFulkerson(t *testing.T) {
	// Пример использования
	n := 6 // количество вершин
	graph := NewMaxFlowFordFulkerson(n)

	// Добавление рёбер с их ёмкостями
	graph.AddEdge(1, 2, 30)
	graph.AddEdge(1, 3, 40)
	graph.AddEdge(1, 4, 20)
	graph.AddEdge(2, 3, 50)
	graph.AddEdge(2, 5, 40)
	graph.AddEdge(3, 4, 20)
	graph.AddEdge(3, 5, 30)
	graph.AddEdge(4, 5, 30)

	s := 1  // исток
	to := 5 // сток

	fmt.Printf("Максимальный поток: %d\n", graph.MaxFlow(s, to))
}
