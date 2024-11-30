package day1

import (
  "fmt"
  "os"
  // "strconv"
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

func part1(s string) int64 {

  // Parse file line by line
  for _, line := range strings.Split(s, "\n") {
    println(line)
  }

  return 0
}

func part2(s string) int64 {
  return 0
}

