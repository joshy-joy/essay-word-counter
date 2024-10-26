package minheap

type Heap struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// MinHeap is a custom type that implements heap.Interface
type MinHeap []Heap

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Count < h[j].Count } // Min-Heap: smallest at top
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push adds an element to the heap
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Heap))
}

// Pop removes the smallest element from the heap
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
