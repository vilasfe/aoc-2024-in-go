package day10

import (
  "fmt"
  "iter"
  "math"
  "os"
  // "regexp"
  "slices"
  "strconv"
  "strings"

  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
  Use:   "day10",
  Short: "day10",
  Long:  `day10`,
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

func part1(s string) int64 {

  // fmt.Printf("Initial string has len: %d\n", len(s))

  // Make the input into a grid of ints
  grid := [][]int{}
  for _, line := range strings.Split(s, "\n") {
    row := Map(strings.Split(strings.TrimSpace(line), ""), func(item string) int {
      val, err := strconv.ParseInt(item, 10, 32)
      if err != nil {
        return -1
      }
      return int(val)
    })
    if len(row) > 0 {
      grid = append(grid, row)
    }
  }

  // fmt.Printf("Input Grid is %dx%d=%d\n", len(grid), len(grid[0]), len(grid) * len(grid[0]))
  // print2D(grid)

  // Setup an adjacency matrix where each nodeID is len(grid[0]) * row + col
  // Set all distances to math.MaxInt32 for infinity before we update the distance matrix
  graph := [][]int{}
  gridSize := len(grid) * len(grid[0])
  for i := 0; i < gridSize; i++ {
    row := slices.Repeat([]int{math.MaxInt32}, gridSize)
    graph = append(graph, row)
  }

  // Set the diagonal to zero and adjacent cells to 1 where appropriate
  for i := range len(grid) {
    rowIdx := i * len(grid[0])
    for j := range len(grid[0]) {
      if grid[i][j] != -1 {
        graph[rowIdx+j][rowIdx+j] = 0
        // fmt.Printf("Checking %d at (%d,%d) %d\n", grid[i][j], i, j, rowIdx+j)
        l := -1
        r := -1
        u := -1
        d := -1
        if j > 0 { l = grid[i][j-1] }
        if j < len(grid[0])-1 { r = grid[i][j+1] }
        if i > 0 { u = grid[i-1][j] }
        if i < len(grid)-1 { d = grid[i+1][j] }
        // fmt.Printf("l: %d, r: %d, u: %d, d: %d\n", l, r, u, d)

        // Check left
        if j > 0 && l - grid[i][j] == 1 {
          // fmt.Println("Adding l")
          graph[rowIdx+j][rowIdx+j-1] = 1
        }
        // Check right
        if j < len(grid[0])-1 && r - grid[i][j] == 1 {
          // fmt.Println("Adding r")
          graph[rowIdx+j][rowIdx+j+1] = 1
        }

        // Check up
        if i > 0 && u - grid[i][j] == 1 {
          // fmt.Println("Adding u")
          graph[rowIdx+j][rowIdx-len(grid[0])+j] = 1
        }

        // Check down
        if i < len(grid)-1 && d - grid[i][j] == 1 {
          // fmt.Println("Adding d")
          graph[rowIdx+j][rowIdx+len(grid[0])+j] = 1
        }
        // fmt.Printf("row: %v\n", graph[rowIdx+j])
      }
    }
  }

  // fmt.Println("Original adjacency matrix")
  // print2D(graph)

  // Find all zeroes as trailheads
  trailheads := find2D(grid, 0)
  // fmt.Printf("trailheads: %v\n", trailheads)

  // Find all nines as goal summits
  summits := find2D(grid, 9)
  // fmt.Printf("summits: %v\n", summits)

  // Get an APSP distance matrix
  dist := floydWarshall(graph)

  // Sum over all zeros how many nines can they get to
  total := int64(0)
  for _, trailhead := range trailheads {
    t := trailhead[0] * len(grid[0]) + trailhead[1]
    for _, summit := range summits {
      s := summit[0] * len(grid[0]) + summit[1]
      if dist[t][s] != math.MaxInt32 {
        total++
      }
    }
  }

  return total
}

func part2(s string) int64 {

  return int64(0)
}

