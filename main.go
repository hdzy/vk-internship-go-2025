package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Node — вершина графа
type Node struct {
	row, col int
	dist     int // накопленная стоимость достижения (приоритет)
	index    int // служебное поле для heap.Interface
}

// PriorityQueue — куча
type PriorityQueue []*Node

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].dist < (*pq)[j].dist
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := x.(*Node)
	n.index = len(*pq)
	*pq = append(*pq, n)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// readInputOneLineStartFinish читает данные из STDIN формата:
//
//	Первая строка: Размер лабиринта (два числа разделённых пробелом, означающих длину и ширину лабиринта).
//
//	Вторая строка: Структура лабиринта (список строк с числами от 0 до 9; числа отделяются пробелами;
//	строки - как в обычном тексте, отделяются переносами строк).
//
//	Третья строка: Две координаты клеток лабиринта - старт и финиш,
//	где каждая координата состоит из двух чисел - индекс строки и индекс столбца, разделенных пробелами.
func readInputOneLineStartFinish() (int, int, [][]int, [2]int, [2]int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	// [rows, cols]
	if !scanner.Scan() {
		return 0, 0, nil, [2]int{}, [2]int{}, errors.New("отсутствует первая строка")
	}
	dimsLine := strings.TrimSpace(scanner.Text())
	dims := strings.Fields(dimsLine)
	if len(dims) != 2 {
		return 0, 0, nil, [2]int{}, [2]int{}, errors.New("неверный формат первой строки")
	}
	rows, err1 := strconv.Atoi(dims[0])
	cols, err2 := strconv.Atoi(dims[1])
	if err1 != nil || err2 != nil || rows <= 0 || cols <= 0 {
		return 0, 0, nil, [2]int{}, [2]int{}, errors.New("некорректное значение первой строки")
	}

	// matrix
	maze := make([][]int, rows)
	for r := 0; r < rows; r++ {
		if !scanner.Scan() {
			return 0, 0, nil, [2]int{}, [2]int{}, fmt.Errorf("отсутствуется строка #%d", r+1)
		}
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) != cols {
			return 0, 0, nil, [2]int{}, [2]int{}, fmt.Errorf("неверный формат строки #%d требуется %d чисел", r+1, cols)
		}

		rowArr := make([]int, cols)
		for c := 0; c < cols; c++ {
			val, err := strconv.Atoi(parts[c])
			if err != nil {
				return 0, 0, nil, [2]int{}, [2]int{}, fmt.Errorf("ошибка преобразования %q в число", parts[c])
			}
			rowArr[c] = val
		}
		maze[r] = rowArr
	}

	// [rowStart colStart rowEnd colEnd]
	if !scanner.Scan() {
		return 0, 0, nil, [2]int{}, [2]int{}, errors.New("нет строки начальной и конечной точки")
	}
	lastLine := strings.TrimSpace(scanner.Text())
	coords := strings.Fields(lastLine)
	if len(coords) != 4 {
		return 0, 0, nil, [2]int{}, [2]int{}, errors.New("неверный формат строки начальной и конечной точки")
	}

	rowStart, errS1 := strconv.Atoi(coords[0])
	colStart, errS2 := strconv.Atoi(coords[1])
	rowEnd, errF1 := strconv.Atoi(coords[2])
	colEnd, errF2 := strconv.Atoi(coords[3])
	if errS1 != nil || errS2 != nil || errF1 != nil || errF2 != nil {
		return 0, 0, nil, [2]int{}, [2]int{}, errors.New("неверный формат строки старта/финиша")
	}

	start := [2]int{rowStart, colStart}
	finish := [2]int{rowEnd, colEnd}

	if rowStart < 0 || rowStart >= rows || colStart < 0 || colStart >= cols {
		return 0, 0, nil, [2]int{}, [2]int{}, errors.New("координаты старта выходят за пределы матрицы")
	}
	if rowEnd < 0 || rowEnd >= rows || colEnd < 0 || colEnd >= cols {
		return 0, 0, nil, [2]int{}, [2]int{}, errors.New("координаты финиша выходят за пределы матрицы")
	}

	return rows, cols, maze, start, finish, nil
}

// dijkstra ищет кратчайший путь по сумме весов,
// возвращает срез координат (row, col) от старта до финиша или nil, если путь не найден.
func dijkstra(maze [][]int, start, finish [2]int) [][2]int {
	rows := len(maze)
	if rows == 0 {
		return nil
	}
	cols := len(maze[0])

	if maze[start[0]][start[1]] == 0 || maze[finish[0]][finish[1]] == 0 {
		return nil
	}

	dist := make([][]int, rows)
	parent := make([][]*[2]int, rows)
	const INF = 1000000000
	for r := 0; r < rows; r++ {
		dist[r] = make([]int, cols)
		parent[r] = make([]*[2]int, cols)
		for c := 0; c < cols; c++ {
			dist[r][c] = INF
		}
	}

	sr, sc := start[0], start[1]
	dist[sr][sc] = maze[sr][sc]

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Node{row: sr, col: sc, dist: dist[sr][sc]})

	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for pq.Len() > 0 {
		node := heap.Pop(&pq).(*Node)
		r, c := node.row, node.col
		currDist := node.dist

		if currDist > dist[r][c] {
			continue
		}

		if r == finish[0] && c == finish[1] {
			break
		}

		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr < 0 || nr >= rows || nc < 0 || nc >= cols {
				continue
			}
			if maze[nr][nc] == 0 {
				continue
			}
			newDist := dist[r][c] + maze[nr][nc]
			if newDist < dist[nr][nc] {
				dist[nr][nc] = newDist
				parent[nr][nc] = &[2]int{r, c}
				heap.Push(&pq, &Node{row: nr, col: nc, dist: newDist})
			}
		}
	}

	fr, fc := finish[0], finish[1]
	if dist[fr][fc] == INF {
		return nil
	}

	path := make([][2]int, 0)
	cur := finish
	for {
		path = append(path, cur)
		p := parent[cur[0]][cur[1]]
		if p == nil {
			break
		}
		cur = [2]int{p[0], p[1]}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func main() {
	_, _, maze, start, finish, err := readInputOneLineStartFinish()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	path := dijkstra(maze, start, finish)
	if path == nil {
		fmt.Fprintln(os.Stderr, "Путь не найден")
		os.Exit(1)
	}

	fmt.Println("\nРезультат:")
	for _, coord := range path {
		fmt.Printf("%d %d\n", coord[0], coord[1])
	}
	fmt.Println(".")
}
