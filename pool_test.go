package pool

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPoolWithInitialNumberofClients(t *testing.T) {
	p := NewPool(5, false)
	assert.Equal(t, 5, p.Len())
}

func indexer(p *Pool, wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		cl := p.Get()
		p.Return(cl)
	}
	wg.Done()
}

func TestConcurrentPool(t *testing.T) {
	p := NewPool(5, false)
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go indexer(p, &wg)
	}
	wg.Wait()
}

func faultyIndexer(p *Pool, wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		cl := p.Get()
		cl2 := newClient(cl.id, false)
		p.Return(cl2)
	}
	wg.Done()
}

func TestConcurrentPoolIndexerCreatesNewClient(t *testing.T) {
	p := NewPool(5, false)
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go faultyIndexer(p, &wg)
	}
	wg.Wait()
}
