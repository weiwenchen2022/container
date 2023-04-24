// This example demonstrates a priority queue built using the Heap type.
package heap_test

import (
	"fmt"

	"github.com/weiwenchen2022/container/heap"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.

	// The index is updated and maintained by the heap methods by provide SetIndex method.
	index int // The index of the item in the heap.
}

func (i *Item) Less(j *Item) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return i.priority > j.priority
}

func (i *Item) SetIndex(n int) {
	i.index = n
}

// A PriorityQueue holds Items.
type PriorityQueue struct {
	*heap.Heap[*Item]
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	pq.Fix(item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func Example_priorityQueue() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	xs := make([]*Item, len(items))
	i := 0
	for value, priority := range items {
		xs[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := &PriorityQueue{heap.New((*Item).Less,
		heap.WithData(xs),
		heap.WithSetIndex((*Item).SetIndex)),
	}
	pq.Init()

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	pq.Push(item)
	pq.update(item, item.value, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := pq.Pop()
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
	fmt.Println()

	// Output:
	// 05:orange 04:pear 03:banana 02:apple
}
