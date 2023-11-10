package dm

import (
	"fmt"
	"sync"
)

type book struct {
	users map[string]Session
	mutex sync.RWMutex
}

type Session interface {
	Receive(msg []byte, sender string) (int, error)
}

func newSessionsBook() *book {
	return &book{
		users: make(map[string]Session),
	}
}

func (b *book) new(acc Sender) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.users[acc.Name()]; ok {
		return fmt.Errorf("User %s already exists\n", acc.Name())
	}
	b.users[acc.Name()] = acc
	return nil
}

func (b *book) get(name string) (Session, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if _, ok := b.users[name]; !ok {
		return nil, fmt.Errorf("User %s not found\n", name)
	}
	return b.users[name], nil
}

func (b *book) end(acc Sender) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.users[acc.Name()]; !ok {
		return fmt.Errorf("User %s not found\n", acc.Name())
	}
	delete(b.users, acc.Name())
	return nil
}
