package main

import (
	"fmt"
)

func main() {

	//arr := []int{6, 2, 7, 3, 9, 8, 1}

	arr := []int{2, 1, 4}
	fmt.Println("arr:", arr)
	countSort(arr)

}

//计数排序：依靠一个计数数组来实现，不基于比较，time: O(n),space:O(2n)。这是一种用空间换时间的方法。
//重点：找到元素位置，则需要知道前面有几个元素，这样就知道了顺序，然后遍历待排序元素，根据位置输出即可。
func countSort(arr []int) {

	//1.获取最大值，开辟计数数组大小
	max := getMax(arr)
	countArr := make([]int, max+1)  // 因为从0开始，max+1 是为了防止 countArr[i-1] 计数时索引超出
	retArr := make([]int, len(arr)) //返回已排序的数组
	// 2.统计每一个元素次数
	for _, v := range arr {
		countArr[v]++
	}
	fmt.Println("countArr1:", countArr)

	// 3.统计当前 小于或等于countArr[i] 元素的总和，得出每一个元素的位置。
	// 例如有一个[1,3,4]，countArr[i]=3，因为小于或等于只有 1和3大，则排到在第2位。
	// countArr 数组类似于水桶。
	for i := 1; i <= max; i++ {
		countArr[i] += countArr[i-1]
	}

	fmt.Println("countArr2:", countArr)

	// 4.让 arr 中每个元素找到其位置，输出到新数组
	for _, v := range arr {
		//countArr[v]-1:从0开始
		retArr[countArr[v]-1] = v

		countArr[v]-- //遇到重复元素，将 countArr 减1
	}

	fmt.Println("retArr:", retArr)
	//fmt.Println("arr:", arr)
}

func getMax(arr []int) (max int) {
	max = arr[0]
	for _, v := range arr {
		if max < v {
			max = v
		}
	}
	return
}
