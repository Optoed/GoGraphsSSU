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
		fmt.Println("5. add_edge <vertex1> <vertex2> <weight (for an unweighted graph weight = 1)>")
		fmt.Println("6. remove_vertex <name>")
		fmt.Println("7. remove_edge <vertex1> <vertex2>")
		fmt.Println("8. print_graph")
		fmt.Println("9. save_to_file <filename.json>")
		fmt.Println("10. load_from_file <filename.json>")
		fmt.Println("11. list_graphs")
		//task 2: adj list la 2, point 18:
		fmt.Println("12. vertices_from_u_and_not_from_v <vertex U> <vertex V>")
		//task 3: adj list la 3, point 6:
		fmt.Println("13. print_hanging_vertices")
		//task4: is g a subgraph of otherG
		fmt.Println("14. is_subgraph_of <ofOtherGraphID>")
		//task5: can get a tree if erase one vertex from graph (with dfs)
		fmt.Println("15. is_almost_tree")
		fmt.Println("18. exit")

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
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды.")
				continue
			}
			graphID, err := strconv.Atoi(parts[1])
			if err != nil || graphs[graphID] == nil {
				fmt.Println("Неверный ID графа.")
				continue
			}
			activeGraphID = graphID
			fmt.Println("Переключен на граф с ID:", graphID)

		case strings.HasPrefix(command, "copy_graph"):
			parts := strings.Split(command, " ")
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды.")
				continue
			}
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
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды.")
				continue
			}
			vertex := parts[1]
			graphs[activeGraphID].AddVertex(vertex)
			fmt.Println("Вершина добавлена:", vertex)

		case strings.HasPrefix(command, "add_edge"):
			parts := strings.Split(command, " ")
			if len(parts) != 4 {
				fmt.Println("Неверный формат команды.")
				continue
			}
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
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды.")
				continue
			}
			vertex := parts[1]
			graphs[activeGraphID].RemoveVertex(vertex)
			fmt.Println("Вершина удалена:", vertex)

		case strings.HasPrefix(command, "remove_edge"):
			parts := strings.Split(command, " ")
			if len(parts) != 3 {
				fmt.Println("Неверный формат команды.")
				continue
			}
			v1, v2 := parts[1], parts[2]
			graphs[activeGraphID].RemoveEdge(v1, v2)
			fmt.Println("Ребро удалено:", v1, v2)

		case strings.HasPrefix(command, "print_graph"):
			graphs[activeGraphID].PrintAdjList()

		case strings.HasPrefix(command, "save_to_file"):
			parts := strings.Split(command, " ")
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды.")
				continue
			}
			filename := parts[1]
			err := graphs[activeGraphID].SaveToFileJSON(filename)
			if err != nil {
				fmt.Println("Ошибка при сохранении файла:", err)
			} else {
				fmt.Println("Граф сохранён в файл", filename)
			}

		case strings.HasPrefix(command, "load_from_file"):
			parts := strings.Split(command, " ")
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды.")
				continue
			}
			filename := parts[1]
			newGraph, err := NewGraphFromFileJSON(filename)
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

		case strings.HasPrefix(command, "vertices_from_u_and_not_from_v"):
			parts := strings.Split(command, " ")
			if len(parts) != 3 {
				fmt.Println("Неверный формат команды.")
				continue
			}
			vertexU, vertexV := parts[1], parts[2]
			exist, verticesFromUAndNotFromV := graphs[activeGraphID].VerticesFromUAndNotFromV(vertexU, vertexV)
			if !exist {
				fmt.Printf("Нет ни одной такой вершины в графе с id = %d,"+
					" в которую есть дуга из вершины %s, но нет из %s\n",
					activeGraphID, vertexU, vertexV)
			} else {
				fmt.Printf("Найдены следующие вершины в графе с id = %d,"+
					" в которые есть дуга из вершины %s, но нет из %s:\n",
					activeGraphID, vertexU, vertexV)
				fmt.Println(verticesFromUAndNotFromV)
			}

		case strings.HasPrefix(command, "print_hanging_vertices"):
			exist, hangingVertices := graphs[activeGraphID].HangingVertices()
			if !exist {
				fmt.Println("В графе отсутствуют висячие вершины графа (степени 1).")
			} else {
				fmt.Println("Все висячие вершины графа (степени 1):")
				fmt.Println(hangingVertices)
			}

		case strings.HasPrefix(command, "is_subgraph_of"):
			parts := strings.Split(command, " ")
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды.")
				continue
			}
			otherGraphId, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Неверно указан ID другого графа")
				continue
			}
			if graphs[activeGraphID].isSubgraphOf(graphs[otherGraphId]) {
				fmt.Printf("Да, текущий (active) граф с ID = %d является подграфом другого графа с ID = %d"+
					" (то есть все вершины и ребра графа текущего графа с ID = %d присутствуют в графе c ID = %d\n",
					activeGraphID, otherGraphId, activeGraphID, otherGraphId)
			} else {
				fmt.Printf("Нет, текущий (active) граф с ID = %d НЕ является подграфом другого графа с ID = %d"+
					" (то есть НЕ все вершины и ребра графа текущего графа с ID = %d присутствуют в графе c ID = %d\n",
					activeGraphID, otherGraphId, activeGraphID, otherGraphId)
			}

		case strings.HasPrefix(command, "is_almost_graph"):
			res, err := graphs[activeGraphID].isAlmostTree()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if res {
				fmt.Println("Yes, we can delete one vertex and make a tree")
			} else {
				fmt.Println("No, we can't")
			}

		case command == "exit":
			fmt.Println("Завершение программы.")
			return

		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}
