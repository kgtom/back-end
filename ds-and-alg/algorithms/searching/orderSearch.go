package main

import (
	"fmt"
)

func main() {

	arr := []int{1, 3, 4, 2, 5}
	ret := orderSearch(arr, 4)
	fmt.Println("ret:", ret)
}

func orderSearch(arr []int, val int) int {
	idx := 0
	for i := 0; i < len(arr); i++ {
		if val == arr[i] {
			idx = i
			break
		}
	}
	return idx
}
