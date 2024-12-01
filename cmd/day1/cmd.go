package day1

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
  Use:   "day1",
  Short: "day1",
  Long:  `day1`,
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

func ParseFile(s string) ([]int64, []int64) {
  l := []int64{}
  r := []int64{}

  // Parse file line by line
  for _, line := range strings.Split(s, "\n") {
    if line == "" {
      continue
    }
    slice := strings.Fields(line)
    val, err := strconv.ParseInt(slice[0], 10, 64)
    if err != nil {
      panic(err)
    }
    l = append(l, val)

    val, err = strconv.ParseInt(slice[1], 10, 64)
    if err != nil {
      panic(err)
    }
    r = append(r, val)

    // fmt.Printf("%v\n", l)
    // fmt.Printf("%v\n", r)
  }

  return l, r

}

func part1(s string) int64 {

  l, r := ParseFile(s)

  l = slices.Sorted(Values(slices.All(l)))
  r = slices.Sorted(Values(slices.All(r)))

  // Make a channel for the map-reduce of summing the differences
  valDiff := make(chan int64)

  //for k, t := range Zip(slices.Sorted(Values(slices.All(l))), slices.Sorted(Values(slices.All(r)))) {
  for k := range len(l) {
    go func(i int) {
      if l[i] < r[i] {
        valDiff <- r[i] - l[i]
      } else {
        valDiff <- l[i] - r[i]
      }
    }(k)
  }

  // Wait for each goroutine to return so we can take the sum
  total := int64(0)
  for k := 0; k < len(l); k++ {
    total += <-valDiff
  }

  return total
}

func part2(s string) int64 {

  l, r := ParseFile(s)

  l = slices.Sorted(Values(slices.All(l)))
  r = slices.Sorted(Values(slices.All(r)))

  counted := make(chan int64)

  for k := range len(l) {
    go func(i int) {
      counted <- l[i] * CountIf(r, func(x int64) bool { return x == l[i] })
    }(k)
  }

  // Wait for each goroutine to return so we can take the sum
  total := int64(0)
  for k := 0; k < len(l); k++ {
    total += <-counted
  }

  return total
}

