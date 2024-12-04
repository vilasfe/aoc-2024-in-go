package day2

import (
  "fmt"
  "iter"
  "os"
  "slices"
  "strconv"
  "strings"

  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
  Use:   "day2",
  Short: "day2",
  Long:  `day2`,
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

func IsSafe(s []int64) bool {

  diff := []int64{}

  for l,r := range Zip(s[:len(s)-1], s[1:]) {
    diff = append(diff, l - r)
  }

  // check for zeros
  if slices.Contains(diff, 0) {
    return false
  }

  // check all negative or all positive
  if CountIf(diff, func(n int64) bool { return n < 0 }) < int64(len(diff)) && CountIf(diff, func(n int64) bool { return n > 0 }) < int64(len(diff)) {
      return false
  }

  // now check > 3, since we already checked for zeros above
  if slices.ContainsFunc(diff, func(n int64) bool { return AbsInt64(n) > 3 }) {
    return false
  }

  return true
}

func IsSafe2(s []int64) bool {
  // 610 is too low
  // 626 is the right answer but breaks my test
  if IsSafe(s) {
    // fmt.Printf("%v GOOD\n", s)
    return true
  } else {
    for i := 0; i < len(s) - 1; i++ {
      s2 := make([]int64, len(s))
      copy(s2, s)

      s2 = slices.Delete(s2, i, i+1)
      if IsSafe(s2) {
        fmt.Printf("%v GOOD\n", s2)
        return true
      } else {
        fmt.Printf("%v\n", s2)
      }
    }
  }

  fmt.Printf("%v NOT GOOD\n", s)
  return false
}

func part1(s string) int64 {

  total := int64(0)

  // Parse file line by line
  for _, line := range strings.Split(s, "\n") {
    if line == "" {
      continue
    }
    slice := Map(strings.Fields(line), func(item string) int64 {
      val, err := strconv.ParseInt(item, 10, 64)
      if err != nil {
        panic(err)
      }
      return val
    })
    if IsSafe(slice) {
      // fmt.Printf("%v GOOD PART1\n", slice)
      total++
    }
  }

  return total
}

func part2(s string) int64 {
  total := int64(0)

  // Parse file line by line
  for _, line := range strings.Split(s, "\n") {
    if line == "" {
      continue
    }
    slice := Map(strings.Fields(line), func(item string) int64 {
      val, err := strconv.ParseInt(item, 10, 64)
      if err != nil {
        panic(err)
      }
      return val
    })
    if IsSafe2(slice) {
      total++
    }
  }

  return total
}

