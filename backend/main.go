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

func dfs(gridSize int, start, end Coordinate) []Coordinate {
	visited := make([][]bool, gridSize)
	for i := range visited {
		visited[i] = make([]bool, gridSize)
	}

	var path []Coordinate
	var found bool

	
	var dfsHelper func(x, y int) bool
	dfsHelper = func(x, y int) bool {
		if x == end.X && y == end.Y { 
			path = append(path, Coordinate{x, y})
			found = true
			return true
		}

		if !isValid(x, y, gridSize, visited) { 
			return false
		}

		visited[x][y] = true              
		path = append(path, Coordinate{x, y}) 

		for _, dir := range directions { 
			if dfsHelper(x+dir[0], y+dir[1]) {
				return true
			}
		}
		path = path[:len(path)-1]
		return false
	}

	dfsHelper(start.X, start.Y)

	if found {
		return path
	}
	return []Coordinate{} 
}


func findPathHandler(w http.ResponseWriter, r *http.Request) {
	var req PathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	path := dfs(20, req.Start, req.End)

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
