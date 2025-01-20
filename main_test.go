package main

import (
	"reflect"
	"testing"
)

// TestDijkstra_Simple случай, где путь точно есть.
func TestDijkstra_Simple(t *testing.T) {
	maze := [][]int{
		{1, 2, 0},
		{2, 0, 1},
		{9, 1, 0},
	}
	start := [2]int{0, 0}
	finish := [2]int{2, 1}

	expected := [][2]int{
		{0, 0},
		{1, 0},
		{2, 0},
		{2, 1},
	}

	path := dijkstra(maze, start, finish)
	if path == nil {
		t.Fatalf("Результат работы программы nil")
	}

	if !reflect.DeepEqual(path, expected) {
		t.Errorf("Ожидаемый результат: %v, не совпадает с результатом работы программы: %v", expected, path)
	}
}

// TestDijkstra_NoPath случай, когда пути нет.
func TestDijkstra_NoPath(t *testing.T) {
	maze := [][]int{
		{1, 0, 1},
		{0, 0, 0},
		{1, 0, 1},
	}
	start := [2]int{0, 0}
	finish := [2]int{2, 2}

	path := dijkstra(maze, start, finish)
	if path != nil {
		t.Errorf("Ожидалось nil, но результат работы программы %v", path)
	}
}

// TestDijkstra_StartOrFinishWall случай, когда старт или финиш - стена.
func TestDijkstra_StartOrFinishWall(t *testing.T) {
	maze := [][]int{
		{0, 1},
		{1, 1},
	}
	start := [2]int{0, 0}
	finish := [2]int{1, 1}

	path := dijkstra(maze, start, finish)
	if path != nil {
		t.Errorf("Ожидалось nil, но результат работы программы %v", path)
	}

	maze2 := [][]int{
		{1, 1},
		{1, 0},
	}
	start2 := [2]int{0, 0}
	finish2 := [2]int{1, 1}

	path2 := dijkstra(maze2, start2, finish2)
	if path2 != nil {
		t.Errorf("Ожидалось nil, но результат работы программы %v", path2)
	}
}
