package clientpool

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Pool struct {
	connections map[int]*ethclient.Client
	cursor      int
	mu          *sync.Mutex
}

func NewPool(clients []*ethclient.Client) *Pool {
	connections := make(map[int]*ethclient.Client)

	for i := 0; i < len(clients); i++ {
		connections[i] = clients[i]
	}

	return &Pool{
		connections: connections,
		cursor:      0,
		mu:          &sync.Mutex{},
	}
}

func (p *Pool) GetClient() (*ethclient.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	c, ok := p.connections[p.cursor]
	if !ok {
		return nil, errors.New("invalid cursor for clients map")
	}

	fmt.Printf("getting client %d: %+v\n", p.cursor, c)

	if p.cursor == len(p.connections)-1 {
		p.cursor = 0
	} else {
		p.cursor++
	}

	return c, nil
}
