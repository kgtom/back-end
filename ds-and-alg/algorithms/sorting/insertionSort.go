package main

import "fmt"

func main() {

	arr := []int{3, 2, 1, 6}
	insertionSort(arr)
	fmt.Println("arr:", arr)

}

//待排序序列：time O(n^2)，已排序序列：time O(n)
func insertionSort(arr []int) {
	//外层循环控制轮次，内层循环：比较大小、交换
	for i := 0; i < len(arr); i++ {
		//向前找元素
		for j := i; j > 0; j-- {
			//一一比较遇到小的交换
			if arr[j-1] > arr[j] {
				arr[j-1], arr[j] = arr[j], arr[j-1]
			}
		}
	}
}

