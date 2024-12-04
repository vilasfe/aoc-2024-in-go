package day4

import (
  "fmt"
  "iter"
  "os"
  // "regexp"
  // "slices"
  // "strconv"
  "strings"

  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
  Use:   "day4",
  Short: "day4",
  Long:  `day4`,
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

func part1(s string) int64 {

  total := int64(0)

  // Parse file line by line
  grid := strings.Split(s, "\n")
  for i, line := range grid {
    if line == "" {
      continue
    }

    // Count forward XMAS or backward SAMX
    total += int64(strings.Count(line, "XMAS"))
    total += int64(strings.Count(line, "SAMX"))

    if i > 2 {
      for j, c := range line {
        // Count vertical XMAS going up or SAMX going up
        if c == 'X' && grid[i-1][j] == 'M' && grid[i-2][j] == 'A' && grid[i-3][j] == 'S' {
          total++
        }
        if c == 'S' && grid[i-1][j] == 'A' && grid[i-2][j] == 'M' && grid[i-3][j] == 'X' {
          total++
        }
        // Count diagonals that go left and up or down
        if j > 2 {
          if c == 'X' && grid[i-1][j-1] == 'M' && grid[i-2][j-2] == 'A' && grid[i-3][j-3] == 'S' {
            total++
          }
          if c == 'S' && grid[i-1][j-1] == 'A' && grid[i-2][j-2] == 'M' && grid[i-3][j-3] == 'X' {
            total++
          }
        }
        // Count diagonals that go right and up or down
        if j < len(line) - 3 {
          if c == 'X' && grid[i-1][j+1] == 'M' && grid[i-2][j+2] == 'A' && grid[i-3][j+3] == 'S' {
            total++
          }
          if c == 'S' && grid[i-1][j+1] == 'A' && grid[i-2][j+2] == 'M' && grid[i-3][j+3] == 'X' {
            total++
          }
        }
      }
    }

  }

  return total
}

func part2(s string) int64 {
  total := int64(0)

  return total
}

