### 关于内存分配的堆、栈
* 堆：一般来讲存放较大对象，程序员手动申请，
* 栈：通常变量、函数参数，一般使用完成就被释放


### 什么是内存逃逸
* 原本分配在栈的变量，在特殊情况下，分配到了堆上
 
### 为什么会发生内存逃逸
* GC压力增大，变量分配到了堆上
* 动态分配内存空间产生一些内存碎片

### 为什么要进行逃逸分析
* 最大好处，减少GC压力，尽量分配到栈上，函数返回时就释放，不需要GC标记

### 如何查看逃逸分析

* 通过gcflags查看编译过程
$ go build -gcflags '-m -l' main.go

* 通过反编译查看内存分配情况
$ go tool compile -S main.go

#### 分配到栈上
~~~go
package main

import (
	"fmt"
)

func main() {

	var a = 1
	fmt.Println("a", a)
}

~~~

分配情况如下，没有 runtime.newobject()，说明在栈上的分配
~~~
      ....
        0x00a2 00162 (<unknown line number>)    MOVQ    96(SP), BP
        0x00a7 00167 (<unknown line number>)    ADDQ    $104, SP
        0x00ab 00171 (<unknown line number>)    RET
        0x00ac 00172 (<unknown line number>)    NOP
        0x00ac 00172 (main.go:7)        PCDATA  $0, $-1
        0x00ac 00172 (main.go:7)        PCDATA  $2, $-1
        0x00ac 00172 (main.go:7)        CALL    runtime.morestack_noctxt(SB)
        0x00b1 00177 (main.go:7)        JMP     0
        0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 0f 86 99  
        
        ....

~~~


#### 分配到堆上

~~~go
package main

import (
	"fmt"
)

type User struct {
	ID   uint32
	Name string
}

func GetUserInfo() *User {

	return &User{
		ID:   1,
		Name: "tom",
	}
}
func main() {

	fmt.Println("userInfo", GetUserInfo())
}

~~~
~~~

tom@MacBook-Pro:~/go/src/demo$  go build -gcflags '-m -l' main.go
# command-line-arguments
./main.go:16:9: &User literal escapes to heap
./main.go:21:14: "userInfo" escapes to heap
./main.go:21:37: GetUserInfo() escapes to heap
./main.go:21:13: main ... argument does not escape

~~~
~~~
go tool compile -S main.go

 0x0022 00034 (main.go:21)       LEAQ    type."".User(SB), AX
        0x0029 00041 (main.go:16)       PCDATA  $2, $0
        0x0029 00041 (main.go:16)       MOVQ    AX, (SP)
        0x002d 00045 (main.go:16)       CALL    runtime.newobject(SB)
        

~~~
发现newobject(),说明分配到堆上。


### 关于竞争检测

~~~go
package main

import (
	"fmt"
	"sync"
	"time"
)

type User struct {
	ID   uint32
	Name string
}

func main() {
	

	arr := make([]*User, 1000)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(a []*User) {
			now := time.Now()
			handleOneUser(a)
			fmt.Println("time:", time.Since(now))
			wg.Done()
		}(arr)
	}
	wg.Wait()
}

func handleOneUser(user []*User) {
	user[0] = &User{
		ID:   1,
		Name: "tom",
	}
	//user[1] = &User{
	//	ID:   2,
	//	Name: "jerry",
	//}

}
~~~

检测竞争：go run -race main.go,发现对共享资源的抢占会导存在竞争

~~~

tom@MacBook-Pro:~/go/src/demo$ go run -race main.go
time: 13.218µs
==================
WARNING: DATA RACE
Write at 0x00c0000aa000 by goroutine 7:
  main.main.func1()
      /Users/tom/go/src/demo/main.go:34 +0xe9

Previous write at 0x00c0000aa000 by goroutine 6:
  main.main.func1()
      /Users/tom/go/src/demo/main.go:34 +0xe9

Goroutine 7 (running) created at:
  main.main()
      /Users/tom/go/src/demo/main.go:21 +0xcd

Goroutine 6 (running) created at:
  main.main()
      /Users/tom/go/src/demo/main.go:21 +0xcd
==================
time: 127.892µs
time: 6.655µs
time: 4.033µs
time: 3.888µs
time: 3.618µs
.....
~~~

### 解决竞争:互斥锁或者读写锁

* 互斥锁：写大于读的场景，尽量小范围使用，避免大量堵塞
* 读写锁：读多写少的场景，

~~~
package main

import (
	"fmt"
	"sync"
	"time"
)

type User struct {
	ID   uint32
	Name string
}

func main() {
	var lock sync.Mutex

	arr := make([]*User, 1000)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(a []*User) {
			now := time.Now()
			lock.Lock()
			handleOneUser(a)
			fmt.Println("time:", time.Since(now))
			lock.Unlock()
			wg.Done()
		}(arr)
	}
	wg.Wait()
}

func handleOneUser(user []*User) {
	user[0] = &User{
		ID:   1,
		Name: "tom",
	}
	//user[1] = &User{
	//	ID:   2,
	//	Name: "jerry",
	//}

}

~~~
### 总结
* 函数参数不一定都使用指针传递
* 注意并发情况下共享资源的处理



