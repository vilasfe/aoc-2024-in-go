package cmd

import (
  "aoc-2024-in-go/cmd/day1"
  "aoc-2024-in-go/cmd/day2"
  "aoc-2024-in-go/cmd/day3"

  "fmt"
  "os"

  "github.com/spf13/cobra"
)

var Cmd = &cobra.Command {
  Use: "2024",
  Short: "2024",
  Long: `2024 is a command line utility for AOC`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff here
  },
}

func init() {
  Cmd.AddCommand(day1.Cmd)
  Cmd.AddCommand(day2.Cmd)
  Cmd.AddCommand(day3.Cmd)
}

func Execute() {
  if err := Cmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}


