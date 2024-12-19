package day11

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
  Use:   "day11",
  Short: "day11",
  Long:  `day11`,
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

func splitStone(s string) (int64, int64) {
  s_list := []string{s[:len(s)/2], s[len(s)/2:]}

  r_list := Map(s_list, func(item string) int64 {
      val, err := strconv.ParseInt(item, 10, 64)
      if err != nil {
        panic(err)
      }
      return val
    })
  return r_list[0], r_list[1]
}

func blink(rocks []int64) []int64 {
  newRocks := []int64{}

  for _, r := range rocks {
    s := strconv.FormatInt(r, 10)

    switch {
    case r == 0:
      newRocks = append(newRocks, 1)
    case len(s)%2 == 0:
      s1, s2 := splitStone(s)
      newRocks = append(newRocks, s1, s2)
      // fmt.Printf("%d %d\n", s1, s2)
    default:
      newRocks = append(newRocks, r*2024)
    }
  }

  return newRocks
}

func part1(s string) int64 {

  slice := Map(strings.Fields(s), func(item string) int64 {
      val, err := strconv.ParseInt(item, 10, 64)
      if err != nil {
        panic(err)
      }
      return val
    })

  // fmt.Printf("%v\n", slice)

  for i := 0; i < 25; i++ {
    slice = blink(slice)
    // fmt.Printf("%v\n", slice)
  }

  return int64(len(slice))
}

func part2(s string) int64 {

  return int64(0)
}

