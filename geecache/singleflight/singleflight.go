package singleflight

import (
	"sync"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

//	func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
//		if c, ok := g.m[key]; ok {
//			c.wg.Wait()
//			return c.val, c.err
//		}
//		c := new(call)
//		c.wg.Add(1)
//		c.val, c.err = fn()
//		g.m[key] = c
//		c.wg.Done()
//		//这个下面不会出现不安全的情况吗
//		delete(g.m, key)
//		return c.val, c.err
//	}
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
