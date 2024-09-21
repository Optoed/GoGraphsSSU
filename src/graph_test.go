package main

import "testing"

func TestGraph_AddVertex(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddVertex("A")

	if _, exists := graph.adjList["A"]; !exists {
		t.Errorf("Expected vertex A to be added, but it wasn't")
	}
}

func TestGraph_AddEdge(t *testing.T) {
	graph := NewEmptyGraph(true)
	graph.AddEdge("A", "B", 10)

	if _, exists := graph.adjList["A"]["B"]; !exists {
		t.Errorf("Expected edge between A and B, but it wasn't added")
	}
	if _, exists := graph.adjList["B"]["A"]; !exists && !graph.isDirected {
		t.Errorf("Expected edge between B and A, but it wasn't added")
	}
}

func TestGraph_RemoveVertex(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddVertex("A")
	graph.RemoveVertex("A")

	if _, exists := graph.adjList["A"]; exists {
		t.Errorf("Expected vertex A to be removed, but it still exists")
	}
}

func TestGraph_RemoveEdge(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)
	graph.RemoveEdge("A", "B")

	if _, exists := graph.adjList["A"]["B"]; exists {
		t.Errorf("Expected edge between A and B to be removed, but it still exists")
	}
	if _, exists := graph.adjList["B"]["A"]; !graph.isDirected && exists {
		t.Errorf("Expected edge between B and A to be removed, but it still exists")
	}
}

func TestGraph_SaveToFile(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)
	err := graph.SaveToFileJSON("test_graph.json")
	if err != nil {
		t.Errorf("Expected no error while saving graph to file, got %v", err)
	}
}

func TestGraph_LoadFromFile(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)
	err := graph.SaveToFileJSON("test_graph.json")
	if err != nil {
		t.Errorf("Expected no error while saving graph to file, got %v", err)
	}

	loadedGraph, err := NewGraphFromFileJSON("test_graph.json")
	if err != nil {
		t.Errorf("Expected no error while loading graph from file, got %v", err)
	}

	if _, exists := loadedGraph.adjList["A"]["B"]; !exists {
		t.Errorf("Expected edge between A and B in loaded graph, but it wasn't found")
	}
}

func TestGraph_CopyGraph(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)

	copiedGraph := NewGraphCopy(*graph)

	if _, exists := copiedGraph.adjList["A"]["B"]; !exists {
		t.Errorf("Expected edge between A and B in copied graph, but it wasn't found")
	}

	graph.RemoveEdge("A", "B")

	if _, exists := copiedGraph.adjList["A"]["B"]; !exists {
		t.Errorf("Expected edge between A and B in copied graph, but it wasn't found." +
			"Probably because of not deep copying")
	}
}

func TestGraph_RemoveEdgeNonExistent(t *testing.T) {
	graph := NewEmptyGraph(false)
	graph.AddEdge("A", "B", 10)
	graph.RemoveEdge("A", "C") // Remove non-existent edge

	if len(graph.adjList["A"]) != 1 {
		t.Errorf("Expected edge count for A to remain 1, got %d", len(graph.adjList["A"]))
	}
}
