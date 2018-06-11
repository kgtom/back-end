常见排序算法



 
1、[冒泡](#1) 
2、[使用redis有什么缺点](#2)
3、[单线程的redis为什么这么快](#3)
4、[redis的数据类型，以及每种数据类型的使用场景](#4)
5、[redis的过期策略以及内存淘汰机制](#5) 
6、[redis和数据库双写一致性问题](#6)
7、[如何应对缓存穿透和缓存雪崩问题](#7) 
8、[如何解决redis的并发竞争问题](#8)
9、[如何利用Redis分布式锁实现控制并发](#9)

<span id="1">1、冒泡</span>
//1.冒泡：相邻两个数，两两比较，交互位置，一次循环下来，最大值就在最后

//时间复杂度O(n^2) n：多次循环

func bubbleSort(arr []int) {

	for i := 0; i < len(arr)-1; i++ {

		for j := i + 1; j < len(arr); j++ {

			if arr[i] > arr[j] {

				arr[i], arr[j] = arr[j], arr[i]

			}

		}

	}

}

func bubbleSortBySortpkg(arr sort.Interface) {

	for i := 0; i < arr.Len(); i++ {

		for j := i + 1; j < arr.Len(); j++ {

			if arr.Less(i, j) {

			}

		}

	}

}
2.选择：选取特定索引值与数组其它元素比较

//时间复杂度O(n^2)

func selectionSort(arr []int) {

	for i := 0; i < len(arr); i++ {

		k := i

		for j := i + 1; j < len(arr); j++ {

			if arr[k] > arr[j] {

				k = j

			}

		}

		arr[i], arr[k] = arr[k], arr[i]

		//fmt.Println("a:", arr)

	}

}

func selectionSortBySortpkg(arr sort.Interface) {

	r := arr.Len()

	for i := 0; i < r; i++ {

		min := i

		for j := i + 1; j < r; j++ {

			if arr.Less(j, min) {

				min = j

			}

		}

		arr.Swap(i, min)

	}

}
3.插入：一条记录插入到已排好的有序表中(相邻两两交换,较慢，改进方案：希尔排序)
//3.插入：一条记录插入到已排好的有序表中(相邻两两交换,较慢，改进方案：希尔排序)

func insertionSort(arr []int) {

	var n = len(arr)

	for i := 1; i < n; i++ {

		j := i //第 j 元素是通过不断向左比较并交换

		for j > 0 {

			//fmt.Println("i:", i, "j:", j, "arr[j-1]:", arr[j-1], "arr[j]:", arr[j])

			if arr[j-1] > arr[j] {

				arr[j-1], arr[j] = arr[j], arr[j-1]

			}

			//fmt.Println("a", arr)

			j = j - 1

		}

	}

}

func insertionSortBypkg(arr sort.Interface) {

	var n = arr.Len()

	for i := 1; i < n; i++ {

		for j := i; j > 0 && arr.Less(j, j-1); j-- {

			arr.Swap(j, j-1)

		}

	}

}
4.希尔排序
// 希尔排序：选择合适的插入排序对间隔 h,然后交换不相邻元素
func ShellSort(a []int) {

	n := len(a)
	h := 1
	for h < n/3 { //寻找合适的间隔h
		h = 3*h + 1
	}
	for h >= 1 {
		//将数组变为间隔h个元素有序
		for i := h; i < n; i++ {
			//间隔h插入排序
			j := i
			fmt.Println("h:", h, "j:", j, "a[j]:", a[i], "a[j-h]:", a[i-h], a)
			for ; j >= h && a[j] < a[j-h]; j -= h {
				//swap(a, j, j-h)
				a[j], a[j-h] = a[j-h], a[j]
				fmt.Println("swap:", a)
			}
		}
		h /= 3
	}
	fmt.Println("h:", h)

}

func swap(slice []int, i int, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
5.归并排序
//4.归并排序的核心思想是将两个有序的数列合并成一个大的有序的序列。通过递归，层层合并，即为归并,利用二分法实现的排序算法，时间复杂度为nlogn.

func MergeSort(arr []int) {
	temp := make([]int, len(arr)) //辅助切片
	mergeSort2(arr, 0, len(arr)-1, temp)
}

func mergeSort2(arr []int, low, high int, temp []int) {

	if low < high {

		mid := (low + high) / 2
		fmt.Println("low:", low, "mid:", mid)
		mergeSort2(arr, low, mid, temp)    //分
		mergeSort2(arr, mid+1, high, temp) //分

		merge(arr, low, mid, high, temp) //合

	}

}

func merge(a []int, first, middle, end int, temp []int) {

	i, m, j, n, k := first, middle, middle+1, end, 0
	for i <= m && j <= n {
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
	for i <= m {
		temp[k] = a[i]
		k++
		i++
	}
	for j <= n {
		temp[k] = a[j]
		k++
		j++
	}
	for ii := 0; ii < k; ii++ {
		a[first+ii] = temp[ii]
	}
	fmt.Printf("sort: arr: %v\n", a)
}
6.快排
func quickSort(arr []int) {

	quickSort2(arr, 0, len(arr)-1)

	fmt.Println("quick :", arr)

}

func quickSort2(arr []int, left, right int) {
	fmt.Println("left:", left, "right:", right)
	if left >= right {
		//起始位置大于或等于终止位置，说明不再需要排序
		return
	}

	//i := partition(arr, left, right)

	i := left

	j := right

	key := arr[(left+right)/2]

	for {

		for arr[i] < key {

			i++
		}

		for arr[j] > key {

			j--
		}

		if i >= j {

			break
		}
		arr[i], arr[j] = arr[j], arr[i]
		//fmt.Println("i:", i,"j:",j,"arr:",arr)
	}

	quickSort2(arr, left, i-1)

	quickSort2(arr, i+1, right)
	//fmt.Println("end")

}

func partition(sortArray []int, left, right int) int {
	i := left

	j := right

	if left < right {

		key := sortArray[(left+right)/2]

		for {

			for sortArray[i] < key {

				i++

			}

			for sortArray[j] > key {

				j--

			}

			if i >= j {

				break

			}

			sortArray[i], sortArray[j] = sortArray[j], sortArray[i]

		}
		fmt.Println("i:", i)

	}
	return i
}
完整源码
总结：
快速排序是最快的通用排序算法，它的内循环的指令很少，而且它还能利用缓存，因为它总是顺序地访问数据。它的运行时间近似为 ~cNlogN，这里的 c 比其他线性对数级别的排序算法都要小
详见：https://play.golang.org/p/cAThuI0eobN

reference：

http://www.cnblogs.com/agui521/p/6918229.html

https://github.com/gaopeng527/go_Algorithm/blob/master/sort.go

https://blog.csdn.net/wangshubo1989/article/details/75135119

https://github.com/arnauddri/algorithms/
