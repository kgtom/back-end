package main

import "fmt"

func main() {

	//整形切片
	arr := []int{3, 2, 1, 6}
	bubbleSort2(arr)
	fmt.Println("arr:", arr)

}

//相邻数比较，大的放后面，小的放前面。time O(n^2) space O(n)
func bubbleSort(arr []int) {

	//外层循环 控制循环轮次，内层循环 控制每轮的比较次数。
	for i := 0; i < len(arr); i++ {
		//每一轮都要进行两两比较 //len(arr)-1-i:每次循环，最后i位置已排好序，不用再遍历比较了
		for j := 0; j < len(arr)-1-i; j++ {

			if arr[j] > arr[j+1] {
				arr[j+1], arr[j] = arr[j], arr[j+1]
				fmt.Println(arr)
			}
		}
	}
}

//注意：如果在某次遍历中没有发生交换，那么就不必进行下次遍历，因为序列是有序的。
//此方法减少遍历次数
func bubbleSort2(arr []int) {
	var isSort = true
	for i := 0; i < len(arr) && isSort; i++ {
		isSort = false //len(arr)-1-i:每次循环，最后i位置已排好序，不用再遍历比较了
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j] > arr[j+1] {
				arr[j+1], arr[j] = arr[j], arr[j+1]
				isSort = true
			}
		}

	}
}
