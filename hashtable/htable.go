package hashtable

import (
	"Quantos/protocol"
	"fmt"
	"sync"
	"github.com/davecgh/go-spew/spew"
)

type Key interface{}
type Value interface{}

type HashTable struct {
	items map[int]Value
	lock sync.RWMutex
}


func _hash(k Key) int {
	key := fmt.Sprintf("%s", k)
	h := 0
	for i := 0; i < len(key); i++ {
		h = 31*h + int(key[i])
	}
	return h
}

func (ht *HashTable) Put(k Key, v Value) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	i := _hash(k)
	if ht.items == nil {
		ht.items = make(map[int]Value)
	}
	ht.items[i] = v

}

func (ht *HashTable) Get(k Key) Value {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	i := _hash(k)
	return ht.items[i]
}

func (ht *HashTable) Remove(k Key) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	i := _hash(k)
	delete(ht.items, i)
}

func (ht *HashTable) Size() int {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	return len(ht.items)
}

func (ht *HashTable) PrintHashTable() {
	spew.Dump(ht)
}

func (ht *HashTable) ToBytes() ([]byte, error) {
	return protocol.Marshal(ht.Items())
}

func (ht *HashTable) Items() map[int]interface{} {
	I := make(map[int]interface{}, ht.Size())
	for k, v := range ht.items {
		I[k] = v
	}
	return I
}
