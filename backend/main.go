package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PathRequest struct {
	Start Coordinate `json:"start"`
	End   Coordinate `json:"end"`
}

type PathResponse struct {
	Path []Coordinate `json:"path"`
}

var directions = [4][2]int{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1},
}

func isValid(x, y int, gridSize int, visited [][]bool) bool {
	return x >= 0 && y >= 0 && x < gridSize && y < gridSize && !visited[x][y]
}

func bfs(gridSize int, start, end Coordinate) []Coordinate {
	visited := make([][]bool, gridSize)
	for i := range visited {
		visited[i] = make([]bool, gridSize)
	}

	var queue []Coordinate
	var parent = make(map[Coordinate]Coordinate)

	visited[start.X][start.Y] = true
	queue = append(queue, start)

	found := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.X == end.X && current.Y == end.Y {
			found = true
			break
		}

		for _, dir := range directions {
			nextX := current.X + dir[0]
			nextY := current.Y + dir[1]
			nextCoord := Coordinate{nextX, nextY}

			if isValid(nextX, nextY, gridSize, visited) {
				visited[nextX][nextY] = true
				queue = append(queue, nextCoord)
				parent[nextCoord] = current
			}
		}
	}

	if found {
		var path []Coordinate
		for at := end; at != start; at = parent[at] {
			path = append([]Coordinate{at}, path...) 
		}
		path = append([]Coordinate{start}, path...)
		return path
	}

	return []Coordinate{} 
}

// func dfs(gridSize int, start, end Coordinate) []Coordinate {
// 	visited := make([][]bool, gridSize)
// 	for i := range visited {
// 		visited[i] = make([]bool, gridSize)
// 	}
// 	shortestPath := []Coordinate{}
// 	var path []Coordinate
	
// 	var dfsHelper func(x, y int) 
// 	dfsHelper = func(x, y int) {
// 		if x == end.X && y == end.Y { 
// 			// path = append(path, Coordinate{x, y})
// 			if len(path) < len(shortestPath) || len(shortestPath) == 0 {
// 					shortestPath = append([]Coordinate{}, path...)
// 			} 
// 			fmt.Println("Shortest path", shortestPath)
// 			return
// 		}

// 		if !isValid(x, y, gridSize, visited) { 
// 			return
// 		}

// 		visited[x][y] = true              
// 		path = append(path, Coordinate{x, y})
// 		if len(path) > len(shortestPath) && len(shortestPath) != 0 {
// 			return
// 		} 

// 		for _, dir := range directions { 
// 			dfsHelper(x+dir[0], y+dir[1]) 
// 		}
		
// 		path = path[:len(path)-1]
// 		fmt.Println(path)
// 	}

// 	dfsHelper(start.X, start.Y)
// 	return append(shortestPath, end)
// }

func findPathHandler(w http.ResponseWriter, r *http.Request) {
	var req PathRequest
	fmt.Println("Request received", req.Start, req.End)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	path := bfs(20, req.Start, req.End)

	resp := PathResponse{
		Path: path,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/find-path", findPathHandler)
	handler := cors.Default().Handler(mux)

	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
