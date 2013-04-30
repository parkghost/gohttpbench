package main

import (
	"sync"
)

type Context struct {
	config *Config
	start  *sync.WaitGroup
	stop   chan struct{}
	rwm    *sync.RWMutex
	store  map[string]interface{}
}

func NewContext(config *Config) *Context {
	start := &sync.WaitGroup{}
	start.Add(config.concurrency)
	return &Context{config, start, make(chan struct{}), &sync.RWMutex{}, make(map[string]interface{})}
}

func (c *Context) SetString(key string, value string) {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	c.store[key] = value
}

func (c *Context) GetString(key string) string {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	return c.store[key].(string)
}

func (c *Context) SetInt(key string, value int) {
	c.rwm.Lock()
	defer c.rwm.Unlock()
	c.store[key] = value
}

func (c *Context) GetInt(key string) int {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	return c.store[key].(int)
}
