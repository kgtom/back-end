## 学习大纲
* [一、高并发背景](#1)
* [二、解决方案1--简单粗博法](#2)
* [三、解决方案2--相对优雅法](#3)
* [四、解决方案3--真正有控制力法(工作池)](#4)
* [五、总结](#5)

## <span id="1">一、高并发背景</span>


## <span id="2">二、解决方案1--简单粗暴法</span>
~~~go
package main

import (
	"fmt"
)

type Job struct {
	Name string
}

var jobChan = make(chan string)

func (j *Job) Handle() {
	fmt.Println("正在处理..." + j.Name)
	jobChan <- j.Name
}

func main() {

	req := [3]string{
		"a",
		"b",
		"c",
	}

	//接收请求
	for _, name := range req {
		job := &Job{Name: name}
		go job.Handle()

	}

	for range req {
		fmt.Println(<-jobChan)
	}
	fmt.Printf("main end")

}

~~~

## <span id="2">三、解决方案1--相对优雅法</span>
~~~go
package main

import (
	"fmt"
)

type Job struct {
	Name string
}

var jobChan = make(chan *Job, 10)
var quit = make(chan string)

func (j *Job) Handle() {
	fmt.Println("正在处理..." + j.Name)

}

func main() {

	req := [3]string{
		"a",
		"b",
		"c",
	}

	//接收请求，发送给Job进行处理
	for _, name := range req {
		job := &Job{Name: name}
		jobChan <- job
	}

	//开启goroutine处理传过来的请求
	go func() {
		for {
			select {
			case job := <-jobChan:
				fmt.Println("接收到job...")
				job.Handle()

			case q := <-quit:
				fmt.Println("end", q)
				return

			}
		}
	}()

	quit <- "全部处理完成"
	fmt.Printf("main end")
}

~~~
## <span id="3">四、真正有控制力法(工作池)</span>


>reference
* [百万实践](https://blog.csdn.net/Jeanphorn/article/details/79018205)
* [实践案例](https://blog.csdn.net/artong0416/article/details/77530843#%E7%9C%9F%E6%AD%A3%E6%8E%A7%E5%88%B6%E5%8D%8F%E7%A8%8B%E6%95%B0%E9%87%8F%E5%B9%B6%E5%8F%91%E6%89%A7%E8%A1%8C%E7%9A%84%E4%BB%BB%E5%8A%A1%E6%95%B0)
