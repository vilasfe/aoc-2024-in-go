package day12

import (
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestParts(t *testing.T) {
  tests := []struct {
    expected int64
    input    string
    fn       func(string) int64
  }{
    {
      expected: 772,
      input:    `test.txt`,
      fn:       part1,
    },
    {
      expected: 140,
      input:    `test1.txt`,
      fn:       part1,
    },
    {
      expected: 1930,
      input:    `test2.txt`,
      fn:       part1,
    },

    {
      expected: 436,
      input:    `test.txt`,
      fn:       part2,
    },
    {
      expected: 80,
      input:    `test1.txt`,
      fn:       part2,
    },
    {
      expected: 1206,
      input:    `test2.txt`,
      fn:       part2,
    },
    {
      expected: 236,
      input:    `test3.txt`,
      fn:       part2,
    },
    {
      expected: 368,
      input:    `test4.txt`,
      fn:       part2,
    },
  }

  for _, test := range tests {
    b, err := os.ReadFile(test.input)
    assert.NoError(t, err, test.input)
    assert.Equal(t, test.expected, test.fn(string(b)))
  }
}

