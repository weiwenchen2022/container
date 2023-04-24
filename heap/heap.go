// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap implements a min-heap.
//
// A heap is a tree with the property that each node is the
// minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, to provid the less function with the (negative) priority as the
// ordering, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete source.
package heap

// The Heap type implements a min-heap with the following invariants (established after
// Init has been called or if the data is empty or sorted):
//
//	!h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()
//
// To create a heap use heap.New.
type Heap[E any] struct {
	less func(ei, ej E) bool

	s        []E
	zero     E
	setIndex func(E, int)
}

type option[E any] func(*Heap[E])

// WithData sets heap using s as its initial contents.
// The heap takes ownership of s, and the caller should not use s after this call.
func WithData[E any](s []E) option[E] {
	return func(h *Heap[E]) {
		h.s = s
	}
}

// WithInitialCap sets heap's initial space to hold n heap elements.
func WithInitialCap[E any](n int) option[E] {
	return func(h *Heap[E]) {
		h.s = make([]E, 0, n)
	}
}

// WithSetIndex sets heap's setIndex field to function f.
// The function is called by the heap methods.
func WithSetIndex[E any](f func(E, int)) option[E] {
	return func(h *Heap[E]) {
		h.setIndex = f
	}
}

// New returns a min-heap according to the less function,
// with initial space to hold n elements.
func New[E any](less func(i, j E) bool, opts ...option[E]) *Heap[E] {
	h := &Heap[E]{less: less}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Len reports the number of elements in the heap.
func (h *Heap[E]) Len() int { return len(h.s) }

func (h *Heap[E]) swap(i, j int) {
	h.s[i], h.s[j] = h.s[j], h.s[i]

	if h.setIndex != nil {
		h.setIndex(h.s[i], i)
		h.setIndex(h.s[j], j)
	}
}

func (h *Heap[E]) push(x E) {
	if h.setIndex != nil {
		h.setIndex(x, len(h.s))
	}

	h.s = append(h.s, x)
}

func (h *Heap[E]) pop() (x E) {
	n := len(h.s)
	x, h.s[n-1] = h.s[n-1], h.zero // avoid memory leak
	if h.setIndex != nil {
		h.setIndex(x, -1) // for safety
	}
	h.s = h.s[:n-1]
	return x
}

// Init establishes the heap invariants.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
func (h *Heap[E]) Init() {
	// heapify
	n := len(h.s)
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Push(x E) {
	h.push(x)
	h.up(len(h.s) - 1)
}

// Pop removes and returns the minimum element (according to less function that provided to New) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func (h *Heap[E]) Pop() E {
	n := len(h.s) - 1
	h.swap(0, n)
	h.down(0, n)
	return h.pop()
}

// Peek returns the minimum element (according to less function that provided to New) from the heap.
// The complexity is O(1).
func (h *Heap[E]) Peek() E {
	return h.s[0]
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Remove(i int) E {
	n := len(h.s) - 1

	if n != i {
		h.swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}

	return h.pop()
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Fix(i int) {
	if !h.down(i, len(h.s)) {
		h.up(i)
	}
}

func (h *Heap[E]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if j == i || !h.less(h.s[j], h.s[i]) {
			break
		}

		h.swap(i, j)
		j = i
	}
}

func (h *Heap[E]) down(i0, n int) bool {
	i := i0

	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}

		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(h.s[j2], h.s[j1]) {
			j = j2 // = 2*i + 2 // right child
		}

		if !h.less(h.s[j], h.s[i]) {
			break
		}

		h.swap(i, j)
		i = j
	}

	return i > i0
}
