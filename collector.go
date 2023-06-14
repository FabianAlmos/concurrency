package main

import (
	"fmt"
	"sync"
	"time"
)

const maxCollectionCount = 500

type Collector struct {
	pipe       chan any
	collection []any
	mutex      sync.Mutex
	ticker     *time.Ticker
	stop       chan struct{}
	done       chan struct{}
}

func NewCollector() *Collector {
	return &Collector{
		pipe:       make(chan any),
		collection: make([]any, 0),
		stop:       make(chan struct{}),
		done:       make(chan struct{}),
	}
}

func (c *Collector) Append(item any) {
	c.pipe <- item
}

func (c *Collector) addToCollection(item any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.collection = append(c.collection, item)
}

func (c *Collector) count() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return len(c.collection)
}

func (c *Collector) clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.collection) == 0 {
		return
	}
	c.collection = make([]any, 0)
}

func (c *Collector) Start() {
	c.ticker = time.NewTicker(time.Minute)
	defer c.ticker.Stop()
	for {
		select {
		case <-c.ticker.C:
			c.release()
			c.clear()
		case item, ok := <-c.pipe:
			if !ok {
				return
			}
			c.addToCollection(item)
			if c.count() < maxCollectionCount {
				continue
			}
			c.release()
			c.clear()
		case <-c.stop:
			c.ticker.Stop()
			close(c.pipe)
			c.release()
			c.clear()
			// unblock code which stops collector until release is finished
			c.done <- struct{}{}
		}
	}
}

// Stop send message to stop channel and wait until collection be released
func (c *Collector) Stop() {
	c.stop <- struct{}{}
	<-c.done
}

func (c *Collector) release() {
	if c.count() == 0 {
		return
	}
	// Do some actions with collection data
	for _, item := range c.collection {
		fmt.Println(item)
	}
}
