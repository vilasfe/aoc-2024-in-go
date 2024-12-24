package day14

import (
  "fmt"
  "iter"
  "math"
  "os"
  "regexp"
  "slices"
  "strconv"
  "strings"

  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
  Use:   "day14",
  Short: "day14",
  Long:  `day14`,
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

type robot struct {
  position pos
  vel pos
}

func Move(r robot, bounds pos, steps int64) robot {
  v := pos{row: r.vel.row * steps, col: r.vel.col * steps}

  r.position = Add(r.position, v)

  // Wrap the row
  if r.position.row < 0 {
    r.position.row = bounds.row - (-1 * r.position.row) % bounds.row
  }
  if r.position.row >= bounds.row {
    r.position.row %= bounds.row
  }

  // Wrap the col
  if r.position.col < 0 {
    r.position.col = bounds.col - (-1 * r.position.col) % bounds.col
  }
  if r.position.col >= bounds.col {
    r.position.col %= bounds.col
  }

  return r
}

func ParseFile(s string) []robot {
  robots := []robot{}

  re := regexp.MustCompile(`p=([+-]?\d+),([+-]?\d+) v=([+-]?\d+),([+-]?\d+)`)
  for _, line := range strings.Split(s, "\n") {
    if line == "" {
      continue
    }
    m1 := re.FindAllStringSubmatch(line, -1)
    // fmt.Printf("%v\n", m1)
    m := m1[0]
    x, x_err := strconv.ParseInt(m[1], 10, 64)
    if x_err != nil {
      panic(x_err)
    }
    y, y_err := strconv.ParseInt(m[2], 10, 64)
    if y_err != nil {
      panic(y_err)
    }
    vx, vx_err := strconv.ParseInt(m[3], 10, 64)
    if vx_err != nil {
      panic(vx_err)
    }
    vy, vy_err := strconv.ParseInt(m[4], 10, 64)
    if vy_err != nil {
      panic(vy_err)
    }

    robots = append(robots, robot{position: pos{row: y, col: x}, vel: pos{row: vy, col: vx}})
  }

  return robots
}

func safetyFactor(robots []robot, bounds pos) int64 {
  halfCol := bounds.col / 2
  halfRow := bounds.row / 2

  q1 := CountIf(robots, func(r robot) bool {
    return r.position.row < halfRow && r.position.col < halfCol
  })

  q2 := CountIf(robots, func(r robot) bool {
    return r.position.row < halfRow && r.position.col > halfCol
  })

  q3 := CountIf(robots, func(r robot) bool {
    return r.position.row > halfRow && r.position.col > halfCol
  })

  q4 := CountIf(robots, func(r robot) bool {
    return r.position.row > halfRow && r.position.col < halfCol
  })

  // fmt.Printf("%d %d %d %d\n", q1, q2, q3, q4)
  return q1 * q2 * q3 * q4
}

func part1(s string) int64 {

  robots := ParseFile(s)

  // fmt.Printf("robots: %v\n", robots)

  bounds := pos{row: 7, col: 11}

  // Use the real bounds instead of the test bounds if we are on the real inputs
  if len(robots) > 14 {
    bounds = pos{row: 103, col: 101}
  }

  // fmt.Printf("Bounds: %v\n", bounds)

  for r, _ := range robots {
    robots[r] = Move(robots[r], bounds, 100)
  }

  // fmt.Printf("robots: %v\n", robots)

  return safetyFactor(robots, bounds)
}

func part2(s string) int64 {

  robots := ParseFile(s)

  // fmt.Printf("robots: %v\n", robots)

  bounds := pos{row: 7, col: 11}

  // Use the real bounds instead of the test bounds if we are on the real inputs
  if len(robots) > 14 {
    bounds = pos{row: 103, col: 101}
  }

  // fmt.Printf("Bounds: %v\n", bounds)

  minSafety := int64(math.MaxInt32)
  minIter := int64(0)

  for s := range bounds.row * bounds.col {
    for r, _ := range robots {
      robots[r] = Move(robots[r], bounds, 1)
    }
    safety := safetyFactor(robots, bounds)
    if safety < minSafety {
      minSafety = safety
      minIter = s
    }
  }

  // fmt.Printf("robots: %v\n", robots)

  return int64(minIter + 1)
}

