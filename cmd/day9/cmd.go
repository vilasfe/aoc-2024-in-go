package day9

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
  Use:   "day9",
  Short: "day9",
  Long:  `day9`,
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

func Checksum(list []int64) int64 {
  total := int64(0)

  for i, v := range list {
    if v >= int64(0) {
      total += int64(i) * v
    }
  }

  return total
}

// Create a brute-force disk map
func CreateDiskMap(l []int64) []int64 {
  m := []int64{}

  nextFile := int64(0)

  isFile := true
  for i := int64(0); i < int64(len(l)); i++ {
    // i = file ID
    // l[i] = number of file blocks or free space (file ID = -1)
    fileID := int64(-1)
    if isFile {
      fileID = nextFile
      nextFile++
    }
    isFile = !isFile
    m = append(m, slices.Repeat([]int64{fileID}, int(l[i]))...)
  }

  // fmt.Printf("Exploded disk map: %v\n", m)

  return m
}

// Just brute-force it for now
func Fragment(l []int64) []int64 {

  diskMap := CreateDiskMap(l)

  f := []int64{}

  for freePtr, stackPtr := 0, len(diskMap)-1; freePtr < stackPtr+1; {

    // copy the values until we find the next free blocks
    for freePtr < stackPtr+1 && diskMap[freePtr] != -1 {
      f = append(f, diskMap[freePtr])
      freePtr++
    }

    // Make sure the stackPtr points to something other than -1
    for stackPtr >= 0 && diskMap[stackPtr] == -1 {
      stackPtr--
    }

    // Now copy from the right until we find the next non-free blocks {
    for freePtr < stackPtr+1 && stackPtr >= 0 && diskMap[freePtr] == -1 && diskMap[stackPtr] != -1 {
      f = append(f, diskMap[stackPtr])
      freePtr++
      stackPtr--
    }
  }

  // fmt.Printf("Fragmented disk map: %v\n", f)

  return f
}


func Defragment(l []int64) []int64 {

  diskMap := CreateDiskMap(l)

  // for each file, try to move it left
  for fileID := slices.Max(slices.Sorted(slices.Values(diskMap))); fileID >= 0; fileID-- {
    // There's a better way from hanging onto the input from the caller as an index but
    // we can brute force things here
    fileLen := CountIf(diskMap, func(x int64) bool { return x == fileID })

    filePtr := slices.Index(diskMap, fileID)

    toInsert := IndexSeq(diskMap, slices.Repeat([]int64{-1}, int(fileLen)))

    // if found to the left of the original location
    if toInsert != -1 && filePtr != -1 && toInsert < filePtr {
      // Put the new elements in place
      diskMap = slices.Replace(diskMap, toInsert, toInsert + int(fileLen), slices.Repeat([]int64{fileID}, int(fileLen))...)
      // Now go find and replace it for the rest of the slice
      for i := toInsert + int(fileLen); i < len(diskMap); i++ {
        if diskMap[i] == fileID {
          diskMap[i] = -1
        }
      }
    }
  }

  // fmt.Printf("Defragmented disk map: %v\n", diskMap)

  return diskMap
}

func part1(s string) int64 {

  // convert string to []int64
  intList := Map(strings.Split(strings.TrimSpace(s), ""), func(item string) int64 {
    val, err := strconv.ParseInt(item, 10, 64)
    if err != nil {
      panic(err)
    }
    return val
  })

  // Fragment slice and return Checksum
  return Checksum(Fragment(intList))
}

func part2(s string) int64 {

  // convert string to []int64
  intList := Map(strings.Split(strings.TrimSpace(s), ""), func(item string) int64 {
    val, err := strconv.ParseInt(item, 10, 64)
    if err != nil {
      panic(err)
    }
    return val
  })

  // Defragment slice and return Checksum
  return Checksum(Defragment(intList))
}

