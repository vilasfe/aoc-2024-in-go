package day3

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
      expected: 161,
      input:    `test.txt`,
      fn:       part1,
    },
    {
      expected: 48,
      input:    `test2.txt`,
      fn:       part2,
    },
  }

  for _, test := range tests {
    b, err := os.ReadFile(test.input)
    assert.NoError(t, err, test.input)
    assert.Equal(t, test.expected, test.fn(string(b)))
  }
}

