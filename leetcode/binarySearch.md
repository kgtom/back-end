
### 测试用例

~~~go
arr := []int{1, 5, 3, 7, 2, 6}
	r := binarySerach.BinarySearh(arr, 3, 0, len(arr)-1)
	r1 := binarySerach.BinarySearchV2(arr, 0, len(arr)-1, 3)
	fmt.Println("arr:", r,r1)
~~~

### 源码


~~~go
//递归版
func BinarySearh(arr []int, searchKey, left, right int) int {

	for left > right {
		return 0
	}

	mid := (left + right) / 2
	if searchKey == arr[mid] {
		return mid
	} else if searchKey > arr[mid] {
		BinarySearh(arr, searchKey, mid+1, right)
	} else {
		BinarySearh(arr, searchKey, left, mid-1)
	}
	return -1
}

//迭代版
func BinarySearchV2(arr []int, l, r, searchKey int) int {
	if l <= r {
		mid := l + (r-l)/2
		if arr[mid] == searchKey {
			return mid
		}
		if arr[mid] < searchKey {
			//return BinarySearchV2(arr, mid+1, r, searchKey)
			l = mid + 1
		} else {
			//return BinarySearchV2(arr, l, mid-1, searchKey)
			r = mid - 1
		}
	}
	return -1
}


~~~
