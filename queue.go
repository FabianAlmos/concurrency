package main

import (
	"sync"
)

const (
	defaultBufferSize     = 10000
	defaultListenersCount = 1
)

type AsyncQueue struct {
	wg          sync.WaitGroup
	queue       chan func()
	stop        chan struct{}
	brokerCount int
}

func NewAsyncQueue() *AsyncQueue {
	return &AsyncQueue{
		queue:       make(chan func(), defaultBufferSize),
		stop:        make(chan struct{}),
		brokerCount: defaultListenersCount,
	}
}

func (s *AsyncQueue) WithBrokerCount(cnt int) *AsyncQueue {
	s.brokerCount = cnt
	return s
}

func (s *AsyncQueue) Append(job func()) {
	s.wg.Add(1)
	s.queue <- job
}

func (s *AsyncQueue) Start() {
	for i := 0; i < s.brokerCount; i++ {
		go s.listen()
	}
}

func (s *AsyncQueue) Shutdown() {
	s.wg.Wait()
	for i := 0; i < s.brokerCount; i++ {
		s.stop <- struct{}{}
	}
}

func (s *AsyncQueue) listen() {
	for {
		select {
		case job := <-s.queue:
			job()
			s.wg.Done()
		case <-s.stop:
			return
		}
	}
}
