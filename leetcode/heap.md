### 测试

~~~go

 fmt.Println("heap ....")
	h := heap.NewMaxHeap()

	h.Insert(heap.Int(8))
	h.Insert(heap.Int(7))
	h.Insert(heap.Int(6))
	h.Insert(heap.Int(3))
	h.Insert(heap.Int(1))
	h.Insert(heap.Int(0))
	h.Insert(heap.Int(2))
	h.Insert(heap.Int(4))
	h.Insert(heap.Int(9))
	h.Insert(heap.Int(5))

	arr := make([]heap.Int, 0)
	for h.Len() > 0 {
		arr = append(arr, h.Extract().(heap.Int))
	}
	fmt.Println("arr:", arr)

~~~

### 代码

~~~go
package heap

import (
	"sync"
)

type Item interface {
	Less(i Item) bool
}

type Heap struct {
	sync.Mutex
	data []Item
	min  bool
}

func NewHeap() *Heap {

	return &Heap{data: make([]Item, 0)}
}
func NewMaxHeap() *Heap {
	return &Heap{
		data: make([]Item, 0),
		min:  false,
	}
}
func NewMinHeap() *Heap {
	return &Heap{
		data: make([]Item, 0),
		min:  true,
	}
}

func (h *Heap) IsEmpty() bool {
	return len(h.data) == 0
}

func (h *Heap) Len() int {
	return len(h.data)
}

func (h *Heap) Get(n int) Item {
	return h.data[n]
}

//尾部插入
func (h *Heap) Insert(i Item) {
	h.Lock()
	defer h.Unlock()
	h.data = append(h.data, i)
	h.shiftUp()
	return
}

func (h *Heap) Less(a, b Item) bool {
	if h.min {
		return a.Less(b)

	} else {

		return b.Less(a)
	}
}

//下沉:当前节点的左右节点比较，哪个小与它交换
//2k:=i<<1
//2k+1 := i<<1 + 1

//伪代码思路
func ShiftDown(i, n int, arr []int) {
	for i*2 <= n {
		T := i * 2
		//左右比
		if T+1 <= n && arr[T+1] < arr[T] {
			T++
		}
		//左右中间小的与当前i 比较，并交换
		if arr[i] < arr[T] {
			//swap( arr[ i ] , arr[ T ] );
			arr[T], arr[i] = arr[i], arr[T]
			i = T
		} else {
			break
		}
	}
}

func (h *Heap) shiftDown() {

	for i, child := 0, 1; i < h.Len() && i<<1+1 < h.Len(); i = child {
		child = i << 1

		if child+1 <= h.Len()-1 && h.Less(h.Get(child+1), h.Get(child)) {
			child++
		}

		if h.Less(h.Get(i), h.Get(child)) {
			break
		}

		h.data[i], h.data[child] = h.data[child], h.data[i]
	}
}

//上浮:从当前节点开始，和它的父节点比较，如果比父节点小则交互。

//伪代码思路
func ShiftUp(i int, arr []int) {

	for i/2 >= 1 {
		if arr[i] < arr[i/2] {

			arr[i], arr[i/2] = arr[i/2], arr[i]
			i = i / 2
		} else {
			break
		}
	}
}

func (h *Heap) shiftUp() {
	for i, parent := h.Len()-1, h.Len()-1; i > 0; i = parent {
		parent = i >> 1
		if h.Less(h.Get(i), h.Get(parent)) {
			h.data[parent], h.data[i] = h.data[i], h.data[parent]
		} else {
			break
		}
	}
}
func (h *Heap) Extract() (el Item) {
	h.Lock()
	defer h.Unlock()
	if h.Len() == 0 {
		return
	}

	el = h.data[0]
	last := h.data[h.Len()-1]
	if h.Len() == 1 {
		h.data = nil
		return
	}

	h.data = append([]Item{last}, h.data[1:h.Len()-1]...)
	h.shiftDown()

	return
}

type Int int

func (x Int) Less(than Item) bool {
	return x < than.(Int)
}

~~~
