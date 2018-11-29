package main

import "fmt"

func main() {

	//整形切片
	arr := []int{3, 2, 1, 6}
	selectionSort(arr)
	fmt.Println("arr:", arr)

}

//选取特定索引值(默认第0个)与其他元素比较，将最小的放在前面。time O(n^2)
func selectionSort(arr []int) {

	for i := 0; i < len(arr); i++ {

		min := i

		//每次找到最小的
		for j := i + 1; j < len(arr); j++ {

			if arr[j] < arr[min] {
				min = j
			}
		}
		//将最小的放在最前面，如果min==i，则不需交换值
		if min != i {
			arr[min], arr[i] = arr[i], arr[min]

		}
	}
}
