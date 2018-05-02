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

package orders

import (
	"github.com/jakeschurch/instruments"
)

// node is an element associated within a LinkedList.
type node struct {
	*instruments.Order
	next, prev *node
}

// Newnode returns a new reference to a Order node.
func NewNode(o *instruments.Order, prev, next *node) *node {
	var node = &node{
		Order: o, next: next, prev: prev,
	}
	if x := node.prev; x != nil {
		x.next = node
	}
	return node
}

// Next returns a reference to a node's next Holdingnode pointer.
func (node *node) Next() *node {
	return node.next
}

// Prev returns a reference to a node's prev Holdingnode pointer.
func (node *node) Prev() *node {
	return node.prev
}
