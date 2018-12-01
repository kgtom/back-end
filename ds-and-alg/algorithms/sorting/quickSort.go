package main

import (
	"fmt"
)

func main() {

	arr := []int{6, 2, 7, 3, 8, 9}

	quickSort(arr, 0, len(arr)-1)
	fmt.Println("arr:", arr)

}

//快排 --递归版： 是对冒泡排序的一种改进，找到中轴数，通过一趟(先找大的，后找小的，双指针指向同一位置。)排序将数据分隔成两部分，一部分小 一部分大,
//对中轴数两边的数据，再分组分别进行上述的过程，
// 最佳分成规模相等的两部分 time:O(nlogn),最差：划分为一部分是大数组，另一部空数组，则 time:O(n^2)
func quickSort(arr []int, l, r int) {

	if l >= r {

		return
	}
	pivotIdx := partition(arr, l, r) //中轴数的下标
	fmt.Println("pivotIdx:", pivotIdx)
	quickSort(arr, l, pivotIdx-1) //左边递归
	quickSort(arr, pivotIdx+1, r) //右边递归

}

//找到中轴数
func partition(arr []int, left, right int) int {

	//默认以左边第一个arr[left]为中轴数
	pivotVal := arr[left]
	pivotIdx := left
	for left < right {

		//从右边开始，当小于pivotVal时 ，等待下面找个大的交换
		for arr[right] > pivotVal && left < right {
			right--
		}
		for arr[left] <= pivotVal && left < right {
			left++
		}
		//在lfet下标临界点，左右交换。
		arr[left], arr[right] = arr[right], arr[left]
		//fmt.Println(" arr-partition:", arr)

	}

	//将中轴数放在中间，左边的都比这个数小，右边的都比这个数大。一趟循环结束。
	arr[left], arr[pivotIdx] = arr[pivotIdx], arr[left]
	return left //返回临界点即中轴数
}

//快排--迭代版
func quickSort2(arr []int) {

}
