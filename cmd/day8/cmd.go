package day8

import (
  "fmt"
  "iter"
  //"math"
  "os"
  // "regexp"
  "slices"
  // "strconv"
  "strings"

  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
  Use:   "day8",
  Short: "day8",
  Long:  `day8`,
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

type pos struct {
  x int64
  y int64
}

func Compare(l, r pos) int {
  if l.x == r.x {
    return int(l.y - r.y)
  }
  return int(l.x - r.x)
}

func (p pos) Add(r pos) pos {
  return pos{ x: p.x + r.x, y: p.y + r.y }
}

func (p pos) Subtract(r pos) pos {
  return pos{ x: p.x - r.x, y: p.y - r.y }
}

func GetLocations(g []string, ant rune) []pos {
  p := []pos{}

  for i := int64(0); i < int64(len(g)); i++ {
    for j := int64(0); j < int64(len(g[i])); j++ {
      if rune(g[i][j]) == ant {
        p = append(p, pos{x: i, y: j})
      }
    }
  }

  return p
}

func CalcInterference(p []pos) []pos {
  i := []pos{}

  if len(p) < 2 {
    return i
  }

  // for each pair, find delta-x and delta-y
  for l := 0; l < len(p); l++ {
    for r := l+1; r < len(p); r++ {
      dx := p[l].x - p[r].x
      dy := p[l].y - p[r].y
      // Find places where 2dx and 2dy work
      i = append(i, pos{x: p[l].x + dx, y: p[l].y + dy})
      i = append(i, pos{x: p[r].x - dx, y: p[r].y - dy})
    }
  }

  return i
}

func CalcInterference2(p []pos, bounds pos) []pos {
  i := []pos{}

  if len(p) < 2 {
    return i
  }

  // for each pair, find delta-x and delta-y
  for l := 0; l < len(p); l++ {
    i = append(i, p[l])
    for r := l+1; r < len(p); r++ {
      diff := pos{x: p[l].x - p[r].x, y: p[l].y - p[r].y}

      // Find places where l & r make a pattern to extend to 0 or the size of the array
      for newPos := p[l].Add(diff); newPos.x >= 0 && newPos.y >= 0 && newPos.x < bounds.x && newPos.y < bounds.y; newPos = newPos.Add(diff) {
        i = append(i, newPos)
      }

      for newPos := p[r].Subtract(diff); newPos.x >= 0 && newPos.y >= 0 && newPos.x < bounds.x && newPos.y < bounds.y; newPos = newPos.Subtract(diff) {
        i = append(i, newPos)
      }
    }
  }

  return i
}

func part1(s string) int64 {

  grid := strings.Split(s, "\n")

  interference := []pos{}

  for c := 'a'; c < 'z'; c++ {
    loc := GetLocations(grid, c)
    if len(loc) > 1 {
      // fmt.Printf("%x has locations %v\n", c, loc)
      interference = append(interference, CalcInterference(loc)...)
    }
  }
  for c := 'A'; c < 'Z'; c++ {
    loc := GetLocations(grid, c)
    if len(loc) > 1 {
      // fmt.Printf("%s has locations %v\n", string(c), loc)
      localInterference := CalcInterference(loc)
      // fmt.Printf("%s has interference %v\n", string(c), localInterference)
      interference = append(interference, localInterference...)
    }
  }
  // Not a typo, including '9' makes it fail
  for c := '0'; c < '9'; c++ {
    loc := GetLocations(grid, c)
    if len(loc) > 1 {
      // fmt.Printf("%s has locations %v\n", string(c), loc)
      localInterference := CalcInterference(loc)
      // fmt.Printf("%s has interference %v\n", string(c), localInterference)
      interference = append(interference, localInterference...)
    }
  }

  // fmt.Printf("Interference locations: %v\n", interference)

  // filter interference based on size of graph
  interference = slices.DeleteFunc(interference, func(n pos) bool {
    return n.x < 0 || n.y < 0 || n.x >= int64(len(grid)) || n.y >= int64(len(grid[0]))
  })

  // fmt.Printf("Pruned interference locations: %v\n", interference)

  interference = Unique(interference)

  // sort the result for readability
  slices.SortFunc(interference, Compare)

  // fmt.Printf("Unique interference locations: %v\n", interference)

  return int64(len(interference))
}

func part2(s string) int64 {
  grid := strings.Split(s, "\n")

  interference := []pos{}

  maxSize := pos{x: int64(len(grid)), y: int64(len(grid[0]))}


  for c := 'a'; c < 'z'; c++ {
    loc := GetLocations(grid, c)
    if len(loc) > 1 {
      // fmt.Printf("%x has locations %v\n", c, loc)
      interference = append(interference, CalcInterference2(loc, maxSize)...)
    }
  }
  for c := 'A'; c < 'Z'; c++ {
    loc := GetLocations(grid, c)
    if len(loc) > 1 {
      // fmt.Printf("%s has locations %v\n", string(c), loc)
      localInterference := CalcInterference2(loc, maxSize)
      // fmt.Printf("%s has interference %v\n", string(c), localInterference)
      interference = append(interference, localInterference...)
    }
  }
  // Not a typo, including '9' makes it fail
  for c := '0'; c <= '9'; c++ {
    loc := GetLocations(grid, c)
    if len(loc) > 1 {
      // fmt.Printf("%s has locations %v\n", string(c), loc)
      localInterference := CalcInterference2(loc, maxSize)
      // fmt.Printf("%s has interference %v\n", string(c), localInterference)
      interference = append(interference, localInterference...)
    }
  }

  // filter interference based on size of graph
  interference = slices.DeleteFunc(interference, func(n pos) bool {
    return n.x < 0 || n.y < 0 || n.x >= int64(len(grid))-1 || n.y >= int64(len(grid[0]))
  })

  interference = Unique(interference)

  // sort the result for readability
  slices.SortFunc(interference, Compare)

  fmt.Printf("Unique interference locations: %v\n", interference)

  return int64(len(interference))
}

