package cmd

import (
  "aoc-2024-in-go/cmd/day1"
  "aoc-2024-in-go/cmd/day2"
  "aoc-2024-in-go/cmd/day3"
  "aoc-2024-in-go/cmd/day4"
  "aoc-2024-in-go/cmd/day5"
  "aoc-2024-in-go/cmd/day6"
  "aoc-2024-in-go/cmd/day7"
  "aoc-2024-in-go/cmd/day8"
  "aoc-2024-in-go/cmd/day9"
  "aoc-2024-in-go/cmd/day10"
  "aoc-2024-in-go/cmd/day11"
  "aoc-2024-in-go/cmd/day12"
  "aoc-2024-in-go/cmd/day13"
  "aoc-2024-in-go/cmd/day14"

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
  Cmd.AddCommand(day4.Cmd)
  Cmd.AddCommand(day5.Cmd)
  Cmd.AddCommand(day6.Cmd)
  Cmd.AddCommand(day7.Cmd)
  Cmd.AddCommand(day8.Cmd)
  Cmd.AddCommand(day9.Cmd)
  Cmd.AddCommand(day10.Cmd)
  Cmd.AddCommand(day11.Cmd)
  Cmd.AddCommand(day12.Cmd)
  Cmd.AddCommand(day13.Cmd)
  Cmd.AddCommand(day14.Cmd)
}

func Execute() {
  if err := Cmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}


