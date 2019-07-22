package pool

import (
	"log"
	"sync"
	"time"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

// Pool maintains a pool of elasticsearch clients. They are retrieved and
// stored back on every request
type Pool struct {
	clients []*Client
	// currIndex contains position of the last retrieved client
	currIndex  int
	mux        *sync.Mutex
	numClients int
}

func newClient(id int, test bool) *Client {
	es, err := elasticsearch7.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
	return &Client{es, id, false}
}

// NewPool creates a new pool of clients
func NewPool(numClients int, test bool) *Pool {
	var clients []*Client
	for i := 0; i < numClients; i++ {
		cl := newClient(i, test)
		clients = append(clients, cl)
	}
	mutex := &sync.Mutex{}
	return &Pool{clients: clients, mux: mutex, numClients: numClients}
}

// Len returns the number of clients
func (p *Pool) Len() int {
	return p.numClients
}

// Client encapsulates an Elasticsearch client
type Client struct {
	*elasticsearch7.Client
	id     int
	locked bool
}

// Get returns the first available client
func (p *Pool) Get() *Client {
	var client *Client
	for {
		p.mux.Lock()
		client = p.clients[p.currIndex]
		if client.locked {
			p.currIndex = (p.currIndex + 1) % p.numClients
			p.mux.Unlock()
			// Leave some time for other goroutines to return
			// their clients
			time.Sleep(1 * time.Millisecond)
			continue
		}
		break
	}
	client.locked = true
	p.mux.Unlock()
	return client
}

// Return returns a new client to the pool
func (p *Pool) Return(cl *Client) {
	p.mux.Lock()
	// client to be available, set currIndex to its id position, so it can be
	// grabbed on next request
	p.currIndex = cl.id
	cl.locked = false
	// Add the client back to its position in the pool. The client might have
	// been renewed due to a connection error. Each subscriber is responsible to
	// maintain its state
	p.clients[cl.id] = cl
	p.mux.Unlock()
}
