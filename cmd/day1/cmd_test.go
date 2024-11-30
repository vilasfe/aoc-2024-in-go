package day1

import (
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestParts(t *testing.T) {
  tests += []struct {
    expected int64
    input    string
    fn       func(string) int64
  }{
    {
      expected: 999,
      input:    `test.txt`,
      fn:       part1,
    },
  }

  for _, test := range tests {
    b, err := os.ReadFile(test.input)
    assert.NoError(t, err, test.input)
    assert.Equal(t, test.exptected, test.fn(string(b)))
  }
}

