package counter_test

import (
	"sort"
	"testing"

	"github.com/dmytro-vovk/shta/internal/counter"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	c := counter.New(10)

	for k, v := range map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
		"D": 4,
		"E": 5,
		"F": 6,
		"G": 7,
		"H": 8,
		"I": 9,
		"J": 10,
		"K": 11,
		"L": 12,
		"M": 13,
	} {
		for range v {
			c.Add(k)
		}
	}

	c.Add("X")

	top := c.Top()
	sort.Strings(top)

	assert.Equal(t, []string{"D", "E", "F", "G", "H", "I", "J", "K", "L", "M"}, top)

	c.Add("X")
	c.Add("X")
	c.Add("X")
	c.Add("X")

	top = c.Top()
	sort.Strings(top)

	assert.Equal(t, []string{"E", "F", "G", "H", "I", "J", "K", "L", "M", "X"}, top)

	assert.Equal(t, 9, c.Get("I"))
	assert.Equal(t, 5, c.Get("X"))
}
