/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 全局对象容器
 */

package di

import "sync"

type container struct {
	mux sync.RWMutex
	m   map[string]interface{}
}

//add object
func (c *container) Add(name string, object interface{}) {
	c.mux.Lock()
	if c.m == nil {
		c.m = make(map[string]interface{})
	}
	c.m[name] = object
	c.mux.Unlock()
}

//Remove object
func (c *container) Remove(name string) {
	if _, ok := c.m[name]; ok {
		c.mux.Lock()
		delete(c.m, name)
		c.mux.Unlock()
	}
}

func (c *container) Get(name string) (interface{}, bool) {
	c.mux.RLock()
	object, ok := c.m[name]
	c.mux.RUnlock()
	return object, ok
}

var c *container

func init() {
	c = New()
}

//New
func New() *container {
	return &container{}
}

func Container() *container {
	return c
}

func Add(name string, object interface{}) {
	if name != "" {
		c.Add(name, object)
	} else {
		panic("name can not be null")
	}
}

func Remove(name string) {
	c.Remove(name)
}

func Get(name string) (interface{}, bool) {
	return c.Get(name)
}
