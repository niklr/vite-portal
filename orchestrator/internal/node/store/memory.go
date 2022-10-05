package store

import (
	"errors"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/collections"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
)

type MemoryStore struct {
	db        collections.NameObjectCollectionI[types.Node]
	addresses mapset.Set[string]
	lock      sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		addresses: mapset.NewSet[string](),
		lock:      sync.RWMutex{},
	}
	s.Clear()
	return s
}

// ---
// Implement "StoreI" interface

func (s *MemoryStore) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.db = collections.NewNameObjectCollection[types.Node]()
	s.addresses.Clear()
}

func (s *MemoryStore) Close() {

}

func (s *MemoryStore) Count() int {
	return s.db.Count()
}

func (s *MemoryStore) GetById(id string) (n types.Node, found bool) {
	node := s.db.Get(id)
	if commonutil.IsEmpty(node) {
		return *new(types.Node), false
	}

	return node, true
}

func (s *MemoryStore) GetByIndex(index int) (n types.Node, found bool) {
	node := s.db.GetByIndex(index)
	if commonutil.IsEmpty(node) {
		return *new(types.Node), false
	}

	return node, true
}

func (s *MemoryStore) GetEnumerator() collections.EnumeratorI[types.Node] {
	return s.db.GetEnumerator()
}

func (s *MemoryStore) Add(n types.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := n.Validate()
	if err != nil {
		return err
	}

	if _, found := s.GetById(n.Id); found {
		return errors.New("a node with the same id already exists")
	}

	addr := n.ClientIp
	if s.addresses.Contains(addr) {
		return errors.New("a node with the same ip address already exists")
	}

	s.db.Add(n.Id, n)
	s.addresses.Add(addr)

	return nil
}

func (s *MemoryStore) Remove(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	existing, found := s.GetById(id)
	if !found {
		return nil
	}

	s.db.Remove(id)
	s.addresses.Remove(existing.ClientIp)

	return nil
}
