// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

import (
	stdheap "container/heap"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

type myHeap struct {
	*Heap[int]
}

func less(i, j int) bool {
	return i < j
}

func (h *myHeap) PushX(x int) {
	h.Push(x)
}

func (h *myHeap) PopX() int {
	return h.Pop()
}

func (h myHeap) verify(t *testing.T, i int) {
	t.Helper()

	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2

	if j1 < n {
		if h.less(h.s[j1], h.s[i]) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d",
				i, h.s[i], j1, h.s[j1])
			return
		}

		h.verify(t, j1)
	}

	if j2 < n {
		if h.less(h.s[j2], h.s[i]) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d",
				i, h.s[i], j2, h.s[j2])
			return
		}

		h.verify(t, j2)
	}
}

func TestInit0(t *testing.T) {
	t.Parallel()

	h := &myHeap{New(less)}
	for i := 20; i > 0; i-- {
		h.push(0) // all elements are the same
	}
	h.Init()
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := h.Pop()
		h.verify(t, 0)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestInit1(t *testing.T) {
	t.Parallel()

	h := &myHeap{New(less)}
	for i := 20; i > 0; i-- {
		h.push(i) // all elements are different
	}
	h.Init()
	h.verify(t, 0)

	for i := 1; h.Len() > 0; i++ {
		x := h.Pop()
		h.verify(t, 0)
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func Test(t *testing.T) {
	t.Parallel()

	h := &myHeap{New(less)}
	h.verify(t, 0)

	for i := 20; i > 10; i-- {
		h.push(i)
	}
	h.Init()
	h.verify(t, 0)

	for i := 10; i > 0; i-- {
		h.Push(i)
		h.verify(t, 0)
	}

	for i := 1; h.Len() > 0; i++ {
		x := h.Pop()
		if i < 20 {
			h.Push(20 + i)
		}
		h.verify(t, 0)

		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	t.Parallel()

	h := &myHeap{New(less)}
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	h.verify(t, 0)

	for h.Len() > 0 {
		i := h.Len() - 1
		x := h.Remove(i)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}

		h.verify(t, 0)
	}
}

func TestRemove1(t *testing.T) {
	t.Parallel()

	h := &myHeap{New(less)}
	for i := 0; i < 10; i++ {
		h.push(i)
	}
	h.verify(t, 0)

	for i := 0; h.Len() > 0; i++ {
		x := h.Remove(0)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}

		h.verify(t, 0)
	}
}

func TestRemove2(t *testing.T) {
	t.Parallel()

	const N = 10

	h := &myHeap{New(less)}
	for i := 0; i < N; i++ {
		h.push(i)
	}
	h.verify(t, 0)

	m := make(map[int]bool)
	for h.Len() > 0 {
		m[h.Remove((h.Len()-1)/2)] = true
		h.verify(t, 0)
	}

	if N != len(m) {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}

	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

type heapInterface interface {
	Len() int
	PushX(int)
	PopX() int
}

type stdHeap []int

func (h stdHeap) Len() int           { return len(h) }
func (h stdHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h stdHeap) Less(i, j int) bool { return h[i] < h[j] }

func (h *stdHeap) Push(v any) {
	*h = append(*h, v.(int))
}

func (h *stdHeap) Pop() (v any) {
	v, *h = (*h)[h.Len()-1], (*h)[:h.Len()-1]
	return
}

func (h *stdHeap) PushX(x int) {
	stdheap.Push(h, x)
}

func (h *stdHeap) PopX() int {
	return stdheap.Pop(h).(int)
}

type bench struct {
	setup func(*testing.B, heapInterface)
	perG  func(*testing.B, heapInterface)
}

func benchHeap(b *testing.B, bench bench) {
	for _, h := range [...]heapInterface{&stdHeap{}, &myHeap{}} {
		b.Run(fmt.Sprintf("%T", h), func(b *testing.B) {
			h = reflect.New(reflect.TypeOf(h).Elem()).Interface().(heapInterface)
			if bench.setup != nil {
				bench.setup(b, h)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bench.perG(b, h)
			}
		})
	}
}

func BenchmarkDup(b *testing.B) {
	const n = 10000

	benchHeap(b, bench{
		setup: func(b *testing.B, hi heapInterface) {
			switch h := hi.(type) {
			case *stdHeap:
				*h = make(stdHeap, 0, n)
			case *myHeap:
				*h = myHeap{New(less, WithInitialCap[int](n))}
			}
		},

		perG: func(b *testing.B, h heapInterface) {
			for i := 0; i < n; i++ {
				h.PushX(0) // all elements are the same
			}

			for h.Len() > 0 {
				h.PopX()
			}
		},
	})
}

func TestFix(t *testing.T) {
	t.Parallel()

	h := &myHeap{New(less)}
	h.verify(t, 0)

	for i := 200; i > 0; i -= 10 {
		h.Push(i)
	}
	h.verify(t, 0)

	if h.s[0] != 10 {
		t.Fatalf("Expected head to be 10, was %d", h.s[0])
	}

	h.s[0] = 210
	h.Fix(0)
	h.verify(t, 0)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 100; i > 0; i-- {
		elem := r.Intn(h.Len())
		if i&1 == 0 {
			h.s[elem] *= 2
		} else {
			h.s[elem] /= 2
		}
		h.Fix(elem)
		h.verify(t, 0)
	}
}
