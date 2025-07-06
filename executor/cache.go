package executor

import (
	"sync"
	"time"
)

type Cache struct {
	Nodes []*Node
	Exp   time.Time
}

type CacheMap map[string]*Cache

func (c CacheMap) Get(key string) []*Node {
	cache, exists := c[key]
	if !exists || time.Now().After(cache.Exp) {
		return nil
	}

	return cache.Nodes
}

func (c CacheMap) Set(key string, nodes []*Node, duration time.Duration) {
	c[key] = &Cache{
		Nodes: nodes,
		Exp:   time.Now().Add(duration),
	}
}

func NewCache(node []*Node, duration time.Duration) *Cache {
	return &Cache{
		Nodes: node,
		Exp:   time.Now().Add(duration),
	}
}

func NewPool() *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return CacheMap{}
		},
	}
}
