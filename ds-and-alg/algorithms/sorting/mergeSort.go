package main

import "fmt"

func main() {

	arr := []int{8, 4, 5, 7, 1, 3, 6, 2}
	mergeSort(arr)
	fmt.Println("arr:", arr)

}

//归并排序使用分治策略，将大序列分解成小序列，然后再合并起来。
// 递归拆分序列：time O(log^n),每次合并操作 time: O(n)
// 总的time:n(log^n)

func mergeSort(arr []int) {
	temp := make([]int, len(arr)) //辅助序列，避免递归时重复开辟空间
	mergeSort2(arr, 0, len(arr)-1, temp)
}

//分
// [8457] [3162]
// [84] [57] [31][62]
// [8][4][5][7][3][1][6][2]
//合并
// [48] [57] [13] [26]
// [4578] [1236]
//[12345678]
func mergeSort2(arr []int, low, high int, temp []int) {

	if low < high {

		mid := (low + high) / 2

		mergeSort2(arr, low, mid, temp) //左边归并

		mergeSort2(arr, mid+1, high, temp) //右边归并

		merge(arr, low, mid, high, temp) //左右合并

	}

}
func merge(a []int, first, mid, end int, temp []int) {

	//fmt.Println("l:", first, "mid:", mid, "h:", end)
	i, j := first, mid+1 //左边 、右边
	k := 0               //临时序列索引
	for i <= mid && j <= end {
		if a[i] <= a[j] {
			temp[k] = a[i]
			k++
			i++
		} else {
			temp[k] = a[j]
			k++
			j++
		}
	}
	for i <= mid { //将左边剩余元素填充到temp
		temp[k] = a[i]
		k++
		i++
	}
	for j <= end { //将右边剩余元素填充到temp
		temp[k] = a[j]
		k++
		j++
	}
	//每次合并操作 time: O(n)
	for m := 0; m < k; m++ {
		//将temp 中全部元素拷贝到原数组中
		a[first+m] = temp[m]

	}
	//fmt.Println("arr:", a, temp)

}
