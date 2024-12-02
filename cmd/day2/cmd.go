package day2

import (
  "fmt"
  "iter"
  "os"
  //"slices"
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

func IsSafe(s []string) bool {
  val, err := strconv.ParseInt(s[0], 10, 64)
  if err != nil {
    panic(err)
  }

  nextVal, nextErr := strconv.ParseInt(s[1], 10, 64)
  if err != nil {
    panic(err)
  }

  diff := val - nextVal
  nextDiff := diff

  for i := 1; i < len(s); i++ {
    nextVal, nextErr = strconv.ParseInt(s[i], 10, 64)
    if nextErr != nil {
      panic(nextErr)
    }

    if diff > 0 {
      nextDiff = val - nextVal
    } else {
      nextDiff = nextVal - val
    }

    if nextDiff < 1 || nextDiff > 3 {
      return false
    }

    val = nextVal
  }

  return true
}

func part1(s string) int64 {

  total := int64(0)

  // Parse file line by line
  for _, line := range strings.Split(s, "\n") {
    if line == "" {
      continue
    }
    slice := strings.Fields(line)
    if IsSafe(slice) {
      total++
    }
  }

  return total
}

func part2(s string) int64 {
  total := int64(0)

  return total
}

