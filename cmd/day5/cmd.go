package day5

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
  Use:   "day5",
  Short: "day5",
  Long:  `day5`,
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

func IsValid(s []int64, allowable map[int64][]int64) bool {

  // get the collection of pairs
  pairs := make(map[int64][]int64)

  for i := int64(0); i < int64(len(s) - 1); i++ {
    for j := int64(i+1); j < int64(len(s)); j++ {
      pairs[s[i]] = append(pairs[s[i]], s[j])
    }
  }

  // Find all pairs in the allowable list
  // If any do not exist then short-circuit return false
  for k, v_seq := range pairs {
    for _, v := range v_seq {
      if !slices.Contains(allowable[k], v) {
        // fmt.Printf("Did not find %d|%d\n", k, v)
        return false
      }
    }
  }

  // fmt.Printf("GOOD: %v\n", s)
  return true
}

func part1(s string) int64 {

  total := int64(0)

  // create a multimap of pairs for later lookup
  // implemented as map of slices
  allowable := make(map[int64][]int64)

  // Parse file line by line
  for _, line := range strings.Split(s, "\n") {
    if line == "" {
      continue
    }

    if strings.Contains(line, "|") {
      // split into left|right
      pair := Map(strings.Split(line, "|"), func(item string) int64 {
        val, err := strconv.ParseInt(item, 10, 64)
        if err != nil {
          panic(err)
        }
        return val
      })

      // fmt.Printf("Adding %d|%d\n", pair[0], pair[1])
      allowable[pair[0]] = append(allowable[pair[0]], pair[1])

    } else if strings.Contains(line, ",") {
      order := Map(strings.Split(line, ","), func(item string) int64 {
        val, err := strconv.ParseInt(item, 10, 64)
        if err != nil {
          panic(err)
        }
        return val
      })

      if IsValid(order, allowable) {
        // fmt.Printf("Summing: %d\n", order[len(order)/2 + 1])
        total += order[len(order)/2]
      }
    }

  }

  return total
}

func part2(s string) int64 {
  total := int64(0)

  return total
}

