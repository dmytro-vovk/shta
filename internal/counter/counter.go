package counter

import (
	"sync"
)

type Counter struct {
	m      sync.RWMutex
	values map[string]int
	top    map[string]int
	max    int
}

func New(max int) *Counter {
	return &Counter{
		values: make(map[string]int),
		top:    make(map[string]int),
		max:    max,
	}
}

// Get returns the number of time the key was added
func (c *Counter) Get(key string) int {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.values[key]
}

// Add increments counter for the key
func (c *Counter) Add(key string) {
	c.m.Lock()

	c.values[key]++
	c.top[key] = c.values[key]

	if len(c.top) > c.max {
		var (
			minKey string
			minVal int
		)

		for k, v := range c.top {
			if v <= c.values[key] {
				if minVal == 0 || minVal >= v {
					minVal = v
					minKey = k
				}
			}
		}

		delete(c.top, minKey)
	}

	c.m.Unlock()
}

// Top returns most frequently added keys
func (c *Counter) Top() []string {
	c.m.RLock()
	defer c.m.RUnlock()

	keys := make([]string, 0, len(c.top))

	for k := range c.top {
		keys = append(keys, k)
	}

	return keys
}
