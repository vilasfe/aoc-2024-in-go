package day6

import (
  "fmt"
  "iter"
  //"math"
  "os"
  "regexp"
  // "slices"
  // "strconv"
  "strings"

  "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
  Use:   "day6",
  Short: "day6",
  Long:  `day6`,
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

func FindGuardStart(g []string) []int64 {
  r := regexp.MustCompile(`[\^<>v]`)

  for i, row := range g {
    // fmt.Printf("Searching: %v\n", row)
    j := r.FindIndex([]byte(row))
    if j != nil {
      pos := []int64{int64(i), int64(j[0])}
      return pos
    }
  }
  return []int64{}
}

func MoveGuard(g []string, pos []int64) ([]string, []int64) {

  newPos := make([]int64, len(pos))
  copy(newPos, pos)

  newG := make([]string, len(g))
  copy(newG, g)

  tmp := []rune(newG[pos[0]])
  tmp[pos[1]] = 'X'
  newG[pos[0]] = string(tmp)

  switch g[pos[0]][pos[1]] {
  case '^':
    newPos[0]--
    if g[newPos[0]][newPos[1]] == '#' {
      newPos = pos
      tmpRow := []rune(newG[pos[0]])
      tmpRow[pos[1]] = '>'
      newG[pos[0]] = string(tmpRow)
    } else {
      tmpRow := []rune(newG[newPos[0]])
      tmpRow[newPos[1]] = '^'
      newG[newPos[0]] = string(tmpRow)
    }
  case '>':
    newPos[1]++
    if g[newPos[0]][newPos[1]] == '#' {
      newPos = pos
      tmpRow := []rune(newG[pos[0]])
      tmpRow[pos[1]] = 'v'
      newG[pos[0]] = string(tmpRow)
    } else {
      tmpRow := []rune(newG[newPos[0]])
      tmpRow[newPos[1]] = '>'
      newG[newPos[0]] = string(tmpRow)
    }
  case 'v':
    newPos[0]++
    if g[newPos[0]][newPos[1]] == '#' {
      newPos = pos
      tmpRow := []rune(newG[pos[0]])
      tmpRow[pos[1]] = '<'
      newG[pos[0]] = string(tmpRow)
    } else {
      tmpRow := []rune(newG[newPos[0]])
      tmpRow[newPos[1]] = 'v'
      newG[newPos[0]] = string(tmpRow)
    }
  case '<':
    newPos[1]--
    if g[newPos[0]][newPos[1]] == '#' {
      newPos = pos
      tmpRow := []rune(newG[pos[0]])
      tmpRow[pos[1]] = '^'
      newG[pos[0]] = string(tmpRow)
    } else {
      tmpRow := []rune(newG[newPos[0]])
      tmpRow[newPos[1]] = '<'
      newG[newPos[0]] = string(tmpRow)
    }
  }

  return newG, newPos
}

func AtBoundary(g []string) bool {
  pos := FindGuardStart(g)

  // fmt.Printf("Checking %v for bounds of %d, %d\n", pos, len(g), len(g[0]))

  // outbound at top
  if pos[0] == 0 && g[pos[0]][pos[1]] == '^' {
    return true
  }

  // outbound at bottom
  if pos[0] == int64(len(g)-2) && g[pos[0]][pos[1]] == 'v' {
    return true
  }

  // outbound at left
  if pos[1] == 0 && g[pos[0]][pos[1]] == '<' {
    return true
  }

  // outbound at right
  if pos[1] == int64(len(g[0])-1) && g[pos[0]][pos[1]] == '>' {
    return true
  }

  return false
}

func TraverseGrid(g []string) []string {
  // Deep Copy the graph so we can modify it
  grid := make([]string, len(g))
  copy(grid, g)

  guard := FindGuardStart(g)

  // limit to the total discrete spaces available
  gridSize := int64(0)
  for _, row := range g {
    gridSize += int64(strings.Count(row, "."))
  }
  for i := int64(1); i < int64(gridSize); i++ {
    grid, guard = MoveGuard(grid, guard)
    if AtBoundary(grid) {
      // Add in the current position before returning
      tmpRow := []rune(grid[guard[0]])
      tmpRow[guard[1]] = 'X'
      grid[guard[0]] = string(tmpRow)

      return grid
    }
  }
  return grid
}

func part1(s string) int64 {

  grid := strings.Split(s, "\n")

  grid = TraverseGrid(grid)

  total := int64(0)
  for _, t := range grid {
    total += int64(strings.Count(t, "X"))
  }
  return total
}

func HasLoop(g []string) bool {

  // Deep Copy the graph so we can modify it
  grid := make([]string, len(g))
  copy(grid, g)

  guard := FindGuardStart(g)

  type move struct {
    row int64
    col int64
    dir rune
  }

  visited := map[move]struct{}{}

  // Add the starting location to the set
  currentMove := move{row: guard[0], col: guard[1], dir: rune(grid[guard[0]][guard[1]])}
  visited[currentMove] = struct{}{}
  // fmt.Printf("inserting: %v\n", currentMove)

  // limit to the total discrete spaces available
  gridSize := int64(0)
  for _, row := range g {
    gridSize += int64(strings.Count(row, "."))
  }
  for i := int64(1); i < int64(gridSize) * 2; i++ {
    grid, guard = MoveGuard(grid, guard)
    if AtBoundary(grid) {
      return false
    }
    currentMove = move{row: guard[0], col: guard[1], dir: rune(grid[guard[0]][guard[1]])}
    // check if we have visited before
    _, ok := visited[currentMove]
    if ok {
      // fmt.Printf("Previously visited: %v\n", currentMove)
      return true
    } else {
      // Add the current location to the set
      // fmt.Printf("inserting: %v\n", currentMove)
      visited[currentMove] = struct{}{}
    }
  }

  return false
}

func part2(s string) int64 {
  total := int64(0)

  grid := strings.Split(s, "\n")
  //gridSize := int64(strings.Count(s, "."))
  gridSize := int64((len(grid)-1) * len(grid[0]))

  guard := FindGuardStart(grid)

  // get the grid with the list of visited locations from part 1
  grid = TraverseGrid(grid)

  // Reset the starting point
  tmpRow := []rune(grid[guard[0]])
  tmpRow[guard[1]] = '^'
  grid[guard[0]] = string(tmpRow)

  // for each visited location, check to see if there is a loop
  loopChan := make(chan bool, gridSize)

  for i_iter := range len(grid) {
    go func(i int) {
      for j_iter := range len(grid[i]) {
        go func(j int) {
          if grid[i][j] != 'X' {
            loopChan <- false
          } else {
            // Deep Copy the graph so we can modify it
            newG := make([]string, len(grid))
            copy(newG, grid)

            tmp := []rune(newG[i])
            tmp[j] = '#'
            newG[i] = string(tmp)

            l := HasLoop(newG)
            // fmt.Printf("Tested %d,%d for %v\n", i, j, l)

            loopChan <- l
          }
        }(j_iter)
      }
    }(i_iter)
  }


  for k := int64(0); k < gridSize; k++ {
    if <-loopChan {
      total++
    }
  }

  return total
}

