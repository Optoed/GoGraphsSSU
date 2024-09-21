package main

import "testing"

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
