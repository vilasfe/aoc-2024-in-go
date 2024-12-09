package day7

import (
  "fmt"
  "iter"
  //"math"
  "os"
  // "regexp"
  "slices"
  "strconv"
  "strings"

  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
  Use:   "day7",
  Short: "day7",
  Long:  `day7`,
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

func evalProdSum(col []int64) []int64 {
  r := []int64{}

  if len(col) == 1 {
    return []int64{col[0]}
  }

  pred := evalProdSum(col[:len(col)-1])

  for _, p := range pred {
    r = append(r, p + col[len(col)-1])
    r = append(r, p * col[len(col)-1])
  }

  return r
}

func part1(s string) int64 {
  total := int64(0)

  // Parse file line by line
  for _, line := range strings.Split(s, "\n") {
    if line == "" {
      continue
    }
    splitLine := strings.Split(line, ":")

    lineTotal, l_err := strconv.ParseInt(splitLine[0], 10, 64)
    if l_err != nil {
      panic(l_err)
    }

    vars := Map(strings.Fields(splitLine[1]), func(item string) int64 {
      val, err := strconv.ParseInt(item, 10, 64)
      if err != nil {
        panic(err)
      }
      return val
    })

    // fmt.Printf("Evaluating: %v\n", vars)

    rowTotal := evalProdSum(vars)

    // If feasible then add the result to the total
    if slices.Contains(rowTotal, lineTotal) {
      total += lineTotal
    }

  }

  return total
}

func part2(s string) int64 {
  total := int64(0)

  return total
}

