package boot

import (
	"log"
	"sync"
)

type container struct {
	items      sync.Map
	shutdownFn []shutdownFn
	once       sync.Once
}

type shutdownFn struct {
	name string
	fn   func()
}

func (c *container) Set(name string, item any, fn func()) *container {
	log.Printf("Setting up %s", name)
	c.items.Store(name, item)

	if fn != nil {
		c.shutdownFn = append(c.shutdownFn, shutdownFn{
			name: name,
			fn:   fn,
		})
	}

	return c
}

func (c *container) Get(name string) any {
	if it, ok := c.items.Load(name); ok {
		return it
	}

	return nil
}

func (c *container) shutdown() {
	c.once.Do(func() {
		for i := len(c.shutdownFn) - 1; i >= 0; i-- {
			log.Printf("Shutting down %s...", c.shutdownFn[i].name)
			c.shutdownFn[i].fn()
			log.Printf("Shutting down %s complete", c.shutdownFn[i].name)
		}
	})
}
