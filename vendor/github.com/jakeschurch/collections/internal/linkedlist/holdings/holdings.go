// Copyright (c) 2018 Jake Schurch
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR k PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package holdings

import (
	"errors"
	"sync"

	"github.com/jakeschurch/instruments"

	"github.com/jakeschurch/collections/internal/cache"
)

var ErrIndexRange = errors.New("index specified out of bounds")

// items is a container for multiple linked list instances.
type items struct {
	data []*list
	len  uint32
}

// newItems is a constructor for the items struct.
func newItems() *items {
	return &items{
		data: make([]*list, 0),
		len:  0,
	}
}

// get returns a linked list from data located at index i.
func (n *items) get(i uint32) (*list, error) {
	if i > n.len {
		return nil, ErrIndexRange
	}
	return n.data[i], nil
}

func (n *items) remove(i uint32) {
	n.data[i] = nil
}

func (n *items) insert(holding instruments.Holding, i uint32) {
	var node = NewNode(holding, nil, nil)

	if i >= n.len {
		n.grow(i)
	}
	if n.data[i] == nil {
		n.data[i] = NewList(*instruments.NewSummary(holding))
	}
	n.data[i].Push(node)
}

func (n *items) grow(i uint32) {
	for ; n.len <= i; n.len = (1 + n.len) * 2 {
	}
	n.data = append(n.data, make([]*list, n.len)...)
}

// ------------------------------------------------------------------

// Holdings is a collection that stores a cache and list.
type Holdings struct {
	sync.RWMutex
	cache *cache.Cache
	*items
}

// New returns a new Holdings instance.
func New() *Holdings {
	return &Holdings{
		cache: cache.New(),
		items: newItems(),
	}
}

func (h *Holdings) Keys() map[string]uint32 {
	return h.cache.Map()
}

// Get returns a list associated with a key from Holdings.list.
// If none are associated with specific key, return nil.
func (h *Holdings) Get(key string) (*list, error) {
	var list *list

	var i, err = h.cache.Get(key)
	if err != nil {
		return nil, err
	}

	h.Lock()
	if list, err = h.items.get(i); err != nil {
		return list, err
	}
	h.Unlock()
	return list, nil
}

// Insert a holding into a Holdings's items linked list.
func (h *Holdings) Insert(holding instruments.Holding) {
	var i uint32
	var err error

	h.Lock()
	if i, err = h.cache.Put(holding.Name); err != nil {
		i, _ = h.cache.Get(holding.Name)
	}
	h.items.insert(holding, i)
	h.Unlock()
}

// Remove a holding into a Holdings's items linked list.
// If nothing can be removed, return error.
func (h *Holdings) Remove(key string) (err error) {
	var i uint32

	h.Lock()
	if i, err = h.cache.Remove(key); err != nil {
		h.Unlock()
		return err
	}
	h.items.remove(i)
	h.Unlock()
	return nil
}

// Update takes a quote as input and updates
// the items list with new quoted data.
// If the cache does not contain an entry with the quote's name,
// method will return error; otherwise nil.
func (h *Holdings) Update(quote instruments.Quote) error {
	var i, err = h.cache.Get(quote.Name)
	if err != nil {
		return err
	}
	h.Lock()
	h.data[i].UpdateMetrics(quote.Bid.Price, quote.Ask.Price, quote.Timestamp)
	h.Unlock()
	return nil
}
func (h *Holdings) GetSlice(key string) ([]*instruments.Holding, error) {
	var i uint32
	var err error
	var nodes *list
	var holdings = make([]*instruments.Holding, 0)

	h.RLock()
	// Query cache to see if we hold Orders with same key.
	// If not return error.
	if i, err = h.cache.Get(key); err != nil {
		h.RUnlock()
		return holdings, err
	}
	// Query items list to return linked list of nodes.
	// Return error if i outside of last index.
	if nodes, err = h.items.get(i); err != nil {
		h.RUnlock()
		return holdings, err
	}
	for x := nodes.head.next; x != nil; x = x.next {
		if x.Holding != nil {
			holdings = append(holdings, x.Holding)
		}
	}
	h.RUnlock()
	return holdings, nil
}
