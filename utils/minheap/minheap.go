package minheap

import (
	"container/heap"
)

// Heap represents an item in the MinHeap with a word and a count.
type Heap struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// MinHeap is a custom type that implements heap.Interface
type MinHeap struct {
	elements []Heap
	indexMap map[string]int // Tracks the index of each word in the heap
}

// NewMinHeap initializes a new MinHeap.
func NewMinHeap() *MinHeap {
	return &MinHeap{
		elements: []Heap{},
		indexMap: make(map[string]int),
	}
}

func (h MinHeap) Len() int           { return len(h.elements) }
func (h MinHeap) Less(i, j int) bool { return h.elements[i].Count < h.elements[j].Count } // Min-Heap: smallest at top
func (h MinHeap) Swap(i, j int) {
	// Swap the elements
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
	// Update the map with the new indices
	h.indexMap[h.elements[i].Word] = i
	h.indexMap[h.elements[j].Word] = j
}

// Push adds an element to the heap
func (h *MinHeap) Push(x interface{}) {
	item := x.(Heap)

	// Check if the word already exists using the map
	if idx, exists := h.indexMap[item.Word]; exists {
		// If it exists, update the count and fix the heap property
		h.elements[idx].Count = item.Count
		heap.Fix(h, idx)
	} else {
		// If it does not exist, add it to the heap
		h.elements = append(h.elements, item)
		h.indexMap[item.Word] = len(h.elements) - 1
		heap.Push(h, item)
	}
}

// Pop removes the smallest element from the heap
func (h *MinHeap) Pop() interface{} {
	old := h.elements
	n := len(old)
	item := old[n-1]

	// Update the heap elements and the index map
	h.elements = old[0 : n-1]
	delete(h.indexMap, item.Word)

	return item
}
