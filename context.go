package main

import (
	"sync"
)

type Context struct {
	config *Config
	start  *sync.WaitGroup
	stop   chan bool

	l     *sync.RWMutex
	store map[string]interface{}
}

func NewContext(config *Config) *Context {
	start := &sync.WaitGroup{}
	start.Add(config.concurrency)
	stop := make(chan bool)
	return &Context{config, start, stop, &sync.RWMutex{}, make(map[string]interface{})}
}

func (c *Context) SetString(key string, value string) {
	c.l.Lock()
	defer c.l.Unlock()
	c.store[key] = value
}

func (c *Context) GetString(key string) string {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.store[key].(string)
}

func (c *Context) SetInt(key string, value int) {
	c.l.Lock()
	defer c.l.Unlock()
	c.store[key] = value
}

func (c *Context) GetInt(key string) int {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.store[key].(int)
}
