package scheduler

type Heap []Entry

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i int, j int) bool {
	return h[i].Next.Before(h[j].Next)
}

func (h Heap) Swap(i int, j int) {
	h[i].Index, h[j].Index = h[j].Index, h[i].Index
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x any) {
	*h = append(*h, x.(Entry))
}

func (h *Heap) Pop() any {
	items := *h
	n := len(items)
	x := items[n-1]
	*h = items[0 : n-1]

	return x
}
