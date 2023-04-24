// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example demonstrates an integer heap built using the Heap type.
package heap_test

import (
	"fmt"

	"github.com/weiwenchen2022/container/heap"
)

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func Example_intHeap() {
	h := heap.New(func(i, j int) bool {
		return i < j
	},
		heap.WithData([]int{2, 1, 5}),
	)
	h.Init()
	h.Push(3)

	fmt.Printf("minimum: %d\n", h.Peek())

	for h.Len() > 0 {
		fmt.Printf("%d ", h.Pop())
	}
	fmt.Println()

	// Output:
	// minimum: 1
	// 1 2 3 5
}
