

package main

import (
	"fmt"
	"time"
	"math/rand"
)

func quickSort(nums []int, chNums chan int) {
	if len(nums) == 0 {
		close(chNums)
		return
	}
	if len(nums) == 1 {
		chNums <- nums[0]
		close(chNums)
		return
	}
	leftNums, rightNums := []int{}, []int{}
	pivotVal := nums[0]
	nums = nums[1:]
	for _, v := range nums {
		if v >= pivotVal {
			rightNums = append(rightNums, v)
		} else {
			leftNums = append(leftNums, v)
		}
	}
	leftCh, rightCh := make(chan int, len(leftNums)), make(chan int, len(rightNums))

	// 
	go quickSort(leftNums, leftCh)
	go quickSort(rightNums, rightCh)

  //小
	for v := range leftCh {
		chNums <- v
	}
  //基准
	chNums <- pivotVal
  //大
	for v := range rightCh {
		chNums <- v
	}
  
	close(chNums)
	return
}

func main() {
	numsLength := 10
	nums := []int{}
	chNums := make(chan int, numsLength) //缓冲channel，存放排序后的数据

	for i := 0; i < numsLength; i++ {
		nums = append(nums, rand.Int()) //随机数 初始化数组
	}
	fmt.Println("print nums:", nums)
	now := time.Now()
	quickSort(nums, chNums)
	fmt.Println("nums len:", numsLength, "用时:", time.Since(now))
	//遍历ch
	for v := range chNums {
		fmt.Println(v)
	}
}



//**总结：**使用channel 和goroutine执行快排操作，利用channel存储排序后的数据
