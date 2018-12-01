package main

import (
	"fmt"
)

func main() {

	arr := []int{1, 2, 3, 4, 5, 6, 7}
	//ret := binarySearch2(arr, 0, len(arr), 5)
	ret := binarySearch(arr, len(arr), 5)
	fmt.Println("ret:", ret)

}

//二分查找--迭代版：采取分治策略。在顺序列表中，折半查找，一边大、一边小，不断迭代查找。
//循环相比递归少了很多内存分配和压栈的操作,开销会少很多，最好做法通过循环迭代实现。
// time:O(log^n)
func binarySearch(arr []int, n int, searchVal int) int {
	//定义查找区间为[l,r]
	l, r := 0, n-1
	//历史上著名的bug l与r都是int的最大值得情况下则l+r越界
	//mid := (l+r)/2

	for l <= r {
		mid := l + (r-l)/2

		if arr[mid] == searchVal {
			return mid
		}
		if arr[mid] < searchVal {
			l = mid + 1

		} else {
			r = mid - 1
		}
	}
	//不存在
	return -1
}

//二分查找--递归版
func binarySearch2(arr []int, l, r int, searchVal int) int {

	mid := l + (r-l)/2
	for l <= r {

		if arr[mid] == searchVal {
			return mid
		}
		if arr[mid] > searchVal {
			return binarySearch2(arr, l, mid-1, searchVal)
		}
		if arr[mid] < searchVal {
			return binarySearch2(arr, mid+1, r, searchVal)
		}
	}
	return -1

}
