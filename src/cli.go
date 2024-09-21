package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func StartCLI() {
	graphs := make(map[int]*Graph)
	activeGraphID := 0
	graphs[activeGraphID] = NewEmptyGraph(false)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nТекущий граф ID:", activeGraphID)
		fmt.Println("Команды:")
		fmt.Println("1. create_graph <isDirected>")
		fmt.Println("2. switch_graph <graphID>")
		fmt.Println("3. copy_graph <fromGraphID>")
		fmt.Println("4. add_vertex <name>")
		fmt.Println("5. add_edge <vertex1> <vertex2> <weight>")
		fmt.Println("6. remove_vertex <name>")
		fmt.Println("7. remove_edge <vertex1> <vertex2>")
		fmt.Println("8. print_graph")
		fmt.Println("9. save_to_file <filename>")
		fmt.Println("10. load_from_file <filename> <isDirected>")
		fmt.Println("11. list_graphs")
		fmt.Println("12. exit")

		fmt.Print("\nВведите команду: ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch {
		case strings.HasPrefix(command, "create_graph"):
			parts := strings.Split(command, " ")
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды.")
				continue
			}
			isDirected := parts[1] == "true"
			graphID := len(graphs)
			graphs[graphID] = NewEmptyGraph(isDirected)
			activeGraphID = graphID
			fmt.Println("Граф создан с ID:", graphID)

		case strings.HasPrefix(command, "switch_graph"):
			parts := strings.Split(command, " ")
			graphID, err := strconv.Atoi(parts[1])
			if err != nil || graphs[graphID] == nil {
				fmt.Println("Неверный ID графа.")
				continue
			}
			activeGraphID = graphID
			fmt.Println("Переключен на граф с ID:", graphID)

		case strings.HasPrefix(command, "copy_graph"):
			parts := strings.Split(command, " ")
			fromGraphID, err := strconv.Atoi(parts[1])
			if err != nil || graphs[fromGraphID] == nil {
				fmt.Println("Неверный ID графа.")
				continue
			}
			newGraphID := len(graphs)
			graphs[newGraphID] = NewGraphCopy(*graphs[fromGraphID])
			activeGraphID = newGraphID
			fmt.Println("Граф скопирован с ID:", fromGraphID, "Новый граф ID:", newGraphID)

		case strings.HasPrefix(command, "add_vertex"):
			parts := strings.Split(command, " ")
			vertex := parts[1]
			graphs[activeGraphID].AddVertex(vertex)
			fmt.Println("Вершина добавлена:", vertex)

		case strings.HasPrefix(command, "add_edge"):
			parts := strings.Split(command, " ")
			v1, v2 := parts[1], parts[2]
			weight, err := strconv.Atoi(parts[3])
			if err != nil {
				fmt.Println("Неверный формат веса.")
				continue
			}
			graphs[activeGraphID].AddEdge(v1, v2, weight)
			fmt.Println("Ребро добавлено:", v1, v2, weight)

		case strings.HasPrefix(command, "remove_vertex"):
			parts := strings.Split(command, " ")
			vertex := parts[1]
			graphs[activeGraphID].RemoveVertex(vertex)
			fmt.Println("Вершина удалена:", vertex)

		case strings.HasPrefix(command, "remove_edge"):
			parts := strings.Split(command, " ")
			v1, v2 := parts[1], parts[2]
			graphs[activeGraphID].RemoveEdge(v1, v2)
			fmt.Println("Ребро удалено:", v1, v2)

		case strings.HasPrefix(command, "print_graph"):
			graphs[activeGraphID].PrintAdjList()

		case strings.HasPrefix(command, "save_to_file"):
			parts := strings.Split(command, " ")
			filename := parts[1]
			err := graphs[activeGraphID].SaveToFile(filename)
			if err != nil {
				fmt.Println("Ошибка при сохранении файла:", err)
			} else {
				fmt.Println("Граф сохранён в файл", filename)
			}

		case strings.HasPrefix(command, "load_from_file"):
			parts := strings.Split(command, " ")
			filename := parts[1]
			isDirected := parts[2] == "true"
			newGraph, err := NewGraphFromFile(filename, isDirected)
			if err != nil {
				fmt.Println("Ошибка при загрузке файла:", err)
			} else {
				graphID := len(graphs)
				graphs[graphID] = newGraph
				activeGraphID = graphID
				fmt.Println("Граф загружен из файла с ID:", graphID)
			}

		case strings.HasPrefix(command, "list_graphs"):
			for id := range graphs {
				fmt.Println("Граф ID:", id)
			}

		case command == "exit":
			fmt.Println("Завершение программы.")
			return

		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}
