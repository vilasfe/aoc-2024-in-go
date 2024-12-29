package day16

import (
  "container/heap"
  "fmt"
  "iter"
  "math"
  "os"
  // "regexp"
  "slices"
  // "strconv"
  "strings"

  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
  Use:   "day16",
  Short: "day16",
  Long:  `day16`,
  Run:   func(cmd *cobra.Command, args []string) {
    execute(cmd.Parent().Name(), cmd.Name())
  },
}

func execute(parent, command string) {
  b, err := os.ReadFile(fmt.Sprintf(`cmd/%s/1.txt`, command))

  if err != nil {
    logrus.Fatal(err)
  }

  logrus.Infof("score part1: %d", part1(string(b)))
  logrus.Infof("score part2: %d", part2(string(b)))
}

// Make a zip iterator for convenience
func Zip[T, U any](t []T, u []U) iter.Seq2[T, U] {
  return func(yield func(T, U) bool) {
    for i := range min(len(t), len(u)) {
      if !yield(t[i], u[i]) {
        return
      }
    }
  }
}

// Seq2 to Seq adapter
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
  return func(yield func(V) bool) {
    for _, v := range seq {
      if !yield(v) {
        return
      }
    }
  }
}

// Map function to apply a function to each element in a slice
func Map[T, V any](ts []T, fn func(T) V) []V {
  result := make([]V, len(ts))
  for i, t := range ts {
    result[i] = fn(t)
  }
  return result
}

// CountIf function
func CountIf[T any](slice []T, f func(T) bool) int64 {
  count := int64(0)
  for _, s := range slice {
    if f(s) {
      count++
    }
  }
  return count
}

func AbsInt64(x int64) int64 {
  if x < 0 {
    return -x
  }
  return x
}

func Unique[T comparable](slice []T) []T {
  keys := map[T]struct{}{}
  result := []T{}

  for _, v := range slice {
    if _, ok := keys[v]; !ok {
      keys[v] = struct{}{}
      result = append(result, v)
    }
  }
  return result
}

func UniqueFunc[T any](slice []T, cmp func (l, r T) int) (result []T) {
  for _, v := range slice {
    if !slices.ContainsFunc(result, func (r T) bool { return cmp(v, r) == 0 }) {
      result = append(result, v)
    }
  }
  return
}

// returns the first index of finding the subsequence e in s, or -1 if not present
func IndexSeq[S ~[]E, E comparable](s S, e S) int {

  if len(e) > len(s) {
    return -1
  }

  for i := range len(s) - len(e) {
    if slices.Equal(s[i:i+len(e)], e) {
      return i
    }
  }

  return -1
}

type PqItem struct {
  value int
  priority int
  index int
}

type PriorityQueue []*PqItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
  // Use greater if you want a max-heap, and less if you want a min-heap
  return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
  pq[i], pq[j] = pq[j], pq[i]
  pq[i].index = i
  pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
  n := len(*pq)
  item := x.(*PqItem)
  item.index = n
  *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
  old := *pq
  n := len(old)
  item := old[n-1]
  old[n-1] = nil // let the GC do its thing
  item.index = -1 // invalidate the index just in case
  *pq = old[:n-1]
  return item
}

func (pq *PriorityQueue) Contains(value int) bool {
  for _, i := range *pq {
    if i.value == value {
      return true
    }
  }
  return false
}

func (pq *PriorityQueue) FindByValue(value int) *PqItem {
  for _, i := range *pq {
    if i.value == value {
      return i
    }
  }
  return nil
}

// Update the priority and value of a PqItem in the queue
func (pq *PriorityQueue) Update(item *PqItem, value int, priority int) {
  item.value = value
  item.priority = priority
  heap.Fix(pq, item.index)
}

func dijkstra(graph SparseAdjMatrix, startIdx int) (cost []int, pred []int) {

  graphSize := graph.MaxKey() + 1

  // Init the data structures
  vQueue := make(PriorityQueue, graphSize)

  cost = slices.Repeat([]int{math.MaxInt32}, graphSize)
  pred = slices.Repeat([]int{-1}, graphSize)

  for i := range graphSize {
    vQueue[i] = &PqItem{ value: i, priority: math.MaxInt32, index: i }
  }
  heap.Init(&vQueue)

  // Set the start node to 0 distance
  vQueue.Update(vQueue[startIdx], startIdx, 0)
  cost[startIdx] = 0

  //fmt.Printf("Initial queue: \n")
  //for _, v := range vQueue {
  //  fmt.Printf("%v ",v)
  //}
  //fmt.Println("Done with initial queue")

  for vQueue.Len() > 0 {
    u := heap.Pop(&vQueue).(*PqItem)
    // fmt.Printf("Processing: %v\n", u)

    // For each neighbor still in the queue
    // fmt.Printf("Neighbors of %d: %v\n", u.value, graph[u.value])
    for  v, _ := range graph[u.value] {
      if vQueue.Contains(v) {
        alt := cost[u.value] + graph[u.value][v]
        // fmt.Printf("Updating neighbor: %d-%d with cost %d\n", u.value, v, alt)
        if alt < cost[v] {
          cost[v] = alt
          pred[v] = u.value
          i := vQueue.FindByValue(v)
          vQueue.Update(i, v, alt)
        }
      }
    }
  }

  return
}

func floydWarshall(graph [][]int) [][]int {
  // Deep copy the adjacency matrix as the distance matrix
  n := len(graph)
  dist := make([][]int, n)
  for i := range dist {
    dist[i] = make([]int, n)
    copy(dist[i], graph[i])
  }

  for k := range n {
    for i := range n {
      for j := range n {
        if dist[i][k] != math.MaxInt32 && dist[k][j] != math.MaxInt32 && dist[i][k]+dist[k][j] < dist[i][j] {
          dist[i][j] = dist[i][k] + dist[k][j]
        }
      }
    }
  }

  return dist
}

func print2D[S ~[]E, E any](s []S) {
  for _, r := range s {
    fmt.Printf("%v\n", r)
  }
}

// returns a list of row/col coordinates and an empty list if not found
func find2D[S ~[]E, E comparable](s []S, e E) [][]int {
  retVal := [][]int{}

  for i, v := range s {
    for j, c := range v {
      if c == e {
        retVal = append(retVal, []int{i, j})
      }
    }
  }
  return retVal
}

func find2DString(s []string, e rune) [][]int {
  retVal := [][]int{}

  for i, v := range s {
    for j, c := range v {
      if c == e {
        retVal = append(retVal, []int{i, j})
      }
    }
  }
  return retVal
}

// this assumes a DAG since it has no "seen check"
func depthFirstSearch[S ~[][]E, E int](g S, start int, visit func(int) ) {
  queue := []int{}
  queue = append(queue, start)

  for len(queue) > 0 {
    v := queue[0]
    visit(v)
    if len(queue) > 1 {
      queue = queue[1:]
    } else {
      queue = []int{}
    }
    for i, val := range g[v] {
      if val > 0 && val < math.MaxInt32 {
        queue = append(queue, i)
      }
    }
  }
}

// GCD using Euclid's algo
func GCD(a, b int) int {
  for b != 0 {
    a, b = b, a%b
  }
  return a
}

// LCM using GCD
func LCM(a, b int) int {
  return (a*b) / GCD(a, b)
}

type pos struct {
  row int64
  col int64
}

func Compare(l, r pos) int {
  if l.row == r.row {
    return int(l.col - r.col)
  }
  return int(l.row - r.row)
}

func Add(l, r pos) pos {
  return pos{ row: l.row + r.row, col: l.col + r.col }
}

func gridFromInput(s string) []string {
  grid := []string{}
  for _, line := range strings.Split(s, "\n") {
    if len(line) > 0 {
      grid = append(grid, line)
    }
  }
  return grid
}

func bfs2D(g []string, start pos, seen [][]bool) []pos {

  // fmt.Printf("BFS starting from %v %s %t\n", start, string(g[start.row][start.col]), seen[start.row][start.col])

  cluster := []pos{}

  neighbors := []pos{start}

  for len(neighbors) > 0 {
    v := neighbors[0]
    // fmt.Printf("Handling %v from %v\n", v, neighbors)
    neighbors = neighbors[1:]

    if g[v.row][v.col] == g[start.row][start.col] && !seen[v.row][v.col] {
      cluster = append(cluster, v)
      if v.row > 0 {
        neighbors = append(neighbors, Add(v, pos{row: -1, col: 0}))
      }
      if v.row < int64(len(g) - 1) {
        neighbors = append(neighbors, Add(v, pos{row: 1, col: 0}))
      }
      if v.col > 0 {
        neighbors = append(neighbors, Add(v, pos{row: 0, col: -1}))
      }
      if v.col < int64(len(g[0]) - 1) {
        neighbors = append(neighbors, Add(v, pos{row: 0, col: 1}))
      }
      seen[v.row][v.col] = true
    }
  }

  return cluster
}

// Return the grid and the move string with newlines removed
func ParseFile(s string) (g []string) {

  for _, r := range strings.Split(s, "\n") {
    if r != "" {
      g = append(g, r)
    }
  }
  return
}

type SparseAdjMatrix map[int]map[int]int

func (sm *SparseAdjMatrix) MaxKey() int {
  maxKey := -1
  for k, _ := range *sm {
    if k > maxKey { maxKey = k }
  }
  return maxKey
}

func (sm *SparseAdjMatrix) At(row, col int) (int, bool) {
  if c, ok := (*sm)[row]; ok {
    if val, ok2 := c[col]; ok2 {
      return val, true
    }
  }
  return -1, false
}

func (sm *SparseAdjMatrix) Set(row, col int, val int) {
  if _, ok := (*sm)[row]; !ok {
    (*sm)[row] = make(map[int] int)
  }
  (*sm)[row][col] = val
}

// Return an adjacency matrix for the grid
func gridToGraph(grid []string) (graph SparseAdjMatrix) {
  graph = make(SparseAdjMatrix)
  // Adjacency matrix is in the form of i = 2*(r * len(grid[0] + c) + (vert = 0, horiz = 1)
  // It costs 1 unit to go one step in the direction of travel
  // And 1000 units to turn from v to h or h to v
  // And 0 to stay in place
  for r := range len(grid) {
    rowIdx := r * len(grid[0])
    for c := range len(grid[0]) {
      if grid[r][c] != '#' {
        cellIdx := 2 * (rowIdx + c)
        graph.Set(cellIdx, cellIdx+1, 1000)
        graph.Set(cellIdx+1, cellIdx, 1000)
        // check vertical movement
        if r > 0 && grid[r-1][c] != '#' {
          graph.Set(cellIdx, 2*(rowIdx-len(grid[0])+c), 1)
        }
        if r < len(grid)-1 && grid[r+1][c] != '#' {
          graph.Set(cellIdx, 2*(rowIdx+len(grid[0])+c), 1)
        }
        // Check horizontal movement
        if c > 0 && grid[r][c-1] != '#' {
          graph.Set(cellIdx+1, 2*(rowIdx+c-1)+1, 1)
        }
        if c < len(grid[0])-1 && grid[r][c+1] != '#' {
          graph.Set(cellIdx+1, 2*(rowIdx+c+1)+1, 1)
        }
      }
    }
  }
  // fmt.Printf("Map: %v\n", graph)
  return
}

func part1(s string) int64 {

  grid := ParseFile(s)

  //for _, s := range grid {
  //  fmt.Println(s)
  //}

  start := find2DString(grid, 'S')[0]
  end := find2DString(grid, 'E')[0]

  startIdx := 2*(start[0] * len(grid[0]) + start[1])
  endIdx := 2*(end[0] * len(grid[0]) + end[1])

  graph := gridToGraph(grid)

  // Ignore vert since it always starts facing East
  //cost1, _ := dijkstra(graph, startIdx)
  cost2, _ := dijkstra(graph, startIdx+1)

  // return int64(slices.Min([]int{cost1[endIdx], cost1[endIdx+1], cost2[endIdx], cost2[endIdx+1]}))
  return int64(slices.Min([]int{cost2[endIdx], cost2[endIdx+1]}))
}

func part2(s string) int64 {

  return int64(0)
}

