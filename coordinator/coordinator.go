package coordinator

import (
	"github.com/aryann/db/storage"
)

type Coordinator struct {
	storage storage.Storage
	sync    chan struct{}
}

func NewCoordinator(storage storage.Storage) *Coordinator {
	return &Coordinator{
		storage: storage,
		sync:    make(chan struct{}, 1),
	}
}

func (c *Coordinator) run(fn func() error) error {
	c.sync <- struct{}{}
	defer func() { <-c.sync }()
	return fn()
}

func (c *Coordinator) Insert(key string, payload string) error {
	return c.run(func() error {
		return c.storage.Insert(key, payload)
	})
}
