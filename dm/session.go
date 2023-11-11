package dm

import (
	"fmt"
	"sync"
)

type book struct {
	accs  map[string]Account
	mutex sync.RWMutex
}

func newSessionsBook() *book {
	return &book{
		accs: make(map[string]Account),
	}
}

func (b *book) new(acc Account) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.accs[acc.Name()]; ok {
		return fmt.Errorf("User %s already exists", acc.Name())
	}
	b.accs[acc.Name()] = acc
	return nil
}

func (b *book) get(name string) (Account, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if _, ok := b.accs[name]; !ok {
		return nil, fmt.Errorf("User %s not found", name)
	}
	return b.accs[name], nil
}

func (b *book) end(acc Account) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.accs[acc.Name()]; !ok {
		return fmt.Errorf("User %s not found", acc.Name())
	}
	delete(b.accs, acc.Name())
	return nil
}
