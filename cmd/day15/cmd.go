package day15

import (
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
  Use:   "day15",
  Short: "day15",
  Long:  `day15`,
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

func swap[S any](x, y *S) {
  *x, *y = *y, *x
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
func ParseFile(s string) ([]string, string) {
  file := strings.Split(s, "\n\n")

  return strings.Split(file[0], "\n"), strings.ReplaceAll(file[1], "\n", "")
}

func Move(g []string, r pos, m rune) ([]string, pos) {

  dir := pos{}

  switch m {
    case '<': {
      dir = pos{row: 0, col: -1}
    }
    case '^': {
      dir = pos{row: -1, col: 0}
    }
    case '>': {
      dir = pos{row: 0, col: 1}
    }
    case 'v': {
      dir = pos{row: 1, col: 0}
    }
  }

  train := []pos{r}
  shouldMove := false

  for t := r; g[t.row][t.col] != '#' && g[t.row][t.col] != '.'; {
    t = Add(t, dir)

    if g[t.row][t.col] == 'O' {
      train = append(train, t)
    } else if g[t.row][t.col] == '.' {
      train = append(train, t)
      shouldMove = true
    }
  }

  if shouldMove {
    // update the grid by moving the train around
    slices.Reverse(train)
    // fmt.Printf("%s %v\n", string(m), train)
    for t, _ := range train[:len(train)-1] {
      temp := rune(g[train[t].row][train[t].col])
      tr := []rune(g[train[t].row])
      tr[train[t].col] = rune(g[train[t+1].row][train[t+1].col])
      g[train[t].row] = string(tr)

      tr2 := []rune(g[train[t+1].row])
      tr2[train[t+1].col] = temp
      g[train[t+1].row] = string(tr2)
    }
    return g, Add(r, dir)
  } else {
    // fmt.Printf("%s Not Moving\n", string(m))
    return g, r
  }
}

func ScoreGPS(g []string) int64 {
  total := int64(0)

  for _, i := range find2DString(g, 'O') {
    total += int64(100 * i[0] + i[1])
  }

  return total
}

func part1(s string) int64 {

  grid, moves := ParseFile(s)

  start := find2DString(grid, '@')[0]

  robot := pos{row: int64(start[0]), col: int64(start[1])}

  for _, m := range moves {
    grid, robot = Move(grid, robot, m)
  }

  return ScoreGPS(grid)
}

func part2(s string) int64 {

  return int64(0)
}

