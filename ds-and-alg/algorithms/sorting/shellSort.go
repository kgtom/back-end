package main

import "fmt"

func main() {

	arr := []int{3, 8, 1, 6, 4, 2}
	shellSort(arr)
	fmt.Println("arr:", arr)

}

//希尔排序是插入排序的高级版，插入排序每次只能一个一个的比较，而希尔排序可以选择间距(增量)。
//即：用增量来将数组进行分隔，直到增量为1。然后对分组后序列执行插入排序。

//最差的的情况 time O(n^2)，最好情况，已排序序列：time O(n)
func shellSort(arr []int) {
	increment := len(arr) / 2 //默认将每一趟排序分隔两段
	for increment >= 1 {
		for i := increment; i < len(arr); i++ {

			//这部分仍然是插入排序，只是间隔增量为increment，而不是1(j-- 被j-increment代替)
			for j := i; j >= increment; {
				if arr[j-increment] > arr[j] {
					arr[j], arr[j-increment] = arr[j-increment], arr[j]
				}

				j = j - increment
			}
		}

		increment = increment / 2 //不断缩小间隔增量
		//increment = int(increment * 5 / 11) //根据公式得出的理想值。

	}
}
