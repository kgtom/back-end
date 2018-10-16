## 学习大纲
* [一、高并发背景](#1)
* [二、解决方案1--简单粗暴法](#2)
* [三、解决方案2--相对优雅法](#3)
* [四、解决方案3--真正有控制力法(工作池)](#4)
* [五、总结](#5)

## <span id="1">一、高并发背景</span>
  有大流量场景下，服务端一些耗时任务，例如定时批量同步用户信息、报表生成、通知类信息，这时候我们可以异步进行处理。具体有以下三种方式：

## <span id="2">二、解决方案1--简单粗暴法</span>
### 1.分析
简单粗暴，但无法控制创建goroutine数量，数量多了，内存也会暴涨，调度也会增多，从而影响性能。
所以我们需要把控以下 开启goroutine数量。继续看方案2
### 2.代码演示
~~~go
package main

import (
	"fmt"
)

type Job struct {
	Name string
}

var quit = make(chan string)

func (j *Job) Handle() {
	fmt.Println("正在处理:" + j.Name)
	quit <- j.Name
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
		//直接开启
		go job.Handle()

	}
	//time.Sleep(time.Second) 为了看到效果，也可以使用 quit chan

	for range req {
		fmt.Println("完成处理:", <-quit)
	}
	fmt.Printf("main end")

}


~~~

## <span id="3">三、解决方案2--相对优雅法</span>

### 1.分析
 使用了缓冲队列一定程度上了提高了并发，但也是治标不治本，大规模并发只是推迟了问题的发生时间。当请求速度远大于队列的处理速度时，缓冲区很快被打满，后面的请求一样被堵塞了。
### 2.代码演示
#### demo1:使用sync.WaitGroup保证每一个worker都处理完成，再退出主goroutine。
~~~go
package main

import (
	"fmt"
	"sync"
)

type Job struct {
	Name string
}

//限制同时工作的job数量(即请求数量) 1000
var jobChan = make(chan Job, 1000)

//若要等待Worker处理完成，我们就要使用sync.WaitGroup
var wg sync.WaitGroup

//从job 队列中获取一个job来处理
func Worker(jobChan <-chan Job) {
	defer wg.Done()
	for job := range jobChan {

		//handel()
		fmt.Println("正在处理：", job.Name)

	}

}

func main() {

	req := [3]string{
		"a",
		"b",
		"c",
	}

	//接收请求
	for _, name := range req {
		job := Job{Name: name}
		jobChan <- job

		go Worker(jobChan)
		wg.Add(1)

	}
	close(jobChan) //关闭chan,优雅方式通知worker不需要再进行工作了。如果注释掉，发生deadlock，因为for range 一直在接收。
	wg.Wait()

	fmt.Printf("main end")

}

~~~

或者 将main()使用下面代码,开启一个goroutine：
~~~go

func main() {

	req := [3]string{
		"a",
		"b",
		"c",
	}

	//接收请求
	for _, name := range req {
		job := Job{Name: name}
		//入队
		jobChan <- job

	}

	wg.Add(1)
	go Worker(jobChan)

	close(jobChan) //去掉注释，导致deadlock。因为range 一直在读。

	wg.Wait()

	fmt.Printf("main end")

}

~~~

#### demo2:缓冲区慢了，提示给调用方(判断队列是否慢：即限流目的)
~~~go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	Name string
}

//限制同时工作的job数量(即请求数量) 1000
var jobChan = make(chan Job, 5)

//若要等待Worker处理完成，我们就要使用sync.WaitGroup
var wg sync.WaitGroup

//从job 队列中获取一个job来处理
func Worker(jobChan <-chan Job) {
	defer wg.Done()
	for job := range jobChan {
		time.Sleep(time.Second * 1) //模拟业务处理1s
		//handel()
		fmt.Println("正在处理：", job.Name)

	}

}

func main() {

	req := [3]string{
		"a",
		"b",
		"c",
	}

	//接收请求
	for _, name := range req {
		job := Job{Name: name}

		//入队
		//jobChan <- job 这种方式如果缓冲区慢了，就堵塞与此，使用IsEnqueue 若返回false提示给调用方。

		r := IsEnqueue(job, jobChan)
		if !r {
			fmt.Println("IsEnqueue:", r)
			return
		}

		go Worker(jobChan)
		wg.Add(1)

	}
	close(jobChan) //关闭chan,如果注释掉，发生deadlock，因为for range 一直在接收。

	wg.Wait()

	fmt.Printf("main end")

}

//如果缓冲区慢了，告知调用方(实质：如地铁限流目的相同,实现非阻塞的生产者模式)
func IsEnqueue(job Job, jobChan chan Job) bool {
	select {
	case jobChan <- job:
		return true
	default:
		return false
	}
}


~~~

#### demo3:防止协程假死（不保证每一个worker都执行完成），使用sync.WaitGroup超时处理

如果不想要一直等sync.WaitGroup的完成，即不想等着每一个worker处理完成，可以设置一个超时时间，我们可以使用select实现。
如果 wg 先返回，那么close(ch)执行后，case<- ch:有效就会执行,返回true，否则执行超时分支，返回false。 

~~~go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	Name string
}

//限制同时工作的job数量(即请求数量) 1000
var jobChan = make(chan Job, 1000)

//若要等待Worker处理完成，我们就要使用sync.WaitGroup
var wg sync.WaitGroup

//从job 队列中获取一个job来处理
func Worker(jobChan <-chan Job) {
	defer wg.Done()
	for job := range jobChan {
		time.Sleep(time.Second * 1) //模拟业务处理1s
		//handel()
		fmt.Println("正在处理：", job.Name)

	}

}

func main() {

	req := [3]string{
		"a",
		"b",
		"c",
	}

	//接收请求
	for _, name := range req {
		job := Job{Name: name}
		jobChan <- job

		go Worker(jobChan)
		wg.Add(1)

	}
	close(jobChan) //如果注释掉，发生deadlock，因为for range 一直在接收。

	//wg.Wait()

	//设置3s超时时间，设置模拟业务5s,case 会走超时分支，提前结束wg
	//r := WaitTimeout(&wg, 2*time.Second)

	//设置3s超时时间，设置模拟业务1s,case 会走ch分支，正常结束wg
	r := WaitTimeout(&wg, 6*time.Second)
	if r {
        fmt.Println("执行worker完成退出")
    } else {
        fmt.Println("执行worker超时退出")
    }
	fmt.Printf("main end")

}

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

~~~
#### demo4:使用context，停止worker工作，不再继续执行
~~~go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	Name string
}

//限制同时工作的job数量(即请求数量) 1000
var jobChan = make(chan Job, 1000)

//若要等待Worker处理完成，我们就要使用sync.WaitGroup
var wg sync.WaitGroup

//从job 队列中获取一个job来处理
func Worker(ctx context.Context, jobChan <-chan Job) {

	for {
		select {
		case <-ctx.Done():
			return

		case job := <-jobChan:
			fmt.Println("正在处理：", job.Name)
			time.Sleep(1 * time.Second) //模拟业务耗时1s
		}
	}

}

func main() {

	req := [5]string{
		"a",
		"b",
		"c",
		"d",
		"e",
	}

	//接收请求
	for _, name := range req {
		job := Job{Name: name}
		//入队
		jobChan <- job

	}
	close(jobChan) //去掉注释，导致deadlock。因为range 一直在读。

	ctx, cancel := context.WithCancel(context.Background())

	go Worker(ctx, jobChan)
	time.Sleep(2 * time.Second) //worker中模拟耗时1s,这儿等待2后超时，执行cancel
	cancel()

	fmt.Printf("main end")

}

~~~
#### demo5:不使用context，停止worker工作，不再继续执行
~~~go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	Name string
}

//限制同时工作的job数量(即请求数量) 1000
var jobChan = make(chan Job, 1000)
var quit = make(chan struct{})

//若要等待Worker处理完成，我们就要使用sync.WaitGroup
var wg sync.WaitGroup

//从job 队列中获取一个job来处理
func Worker(jobChan <-chan Job, quit <-chan struct{}) {

	for {
		select {

		case job := <-jobChan:
			fmt.Println("正在处理：", job.Name)
			time.Sleep(1 * time.Second) //模拟业务耗时1s
		case val, ok := <-quit:
			fmt.Println("收到取消操作信号", val, ok)

			return
		}

	}

}

func main() {

	req := [5]string{
		"a",
		"b",
		"c",
		"d",
		"e",
	}

	//接收请求
	for _, name := range req {
		job := Job{Name: name}
		//入队
		jobChan <- job

	}
	close(jobChan) //去掉注释，导致deadlock。因为range 一直在读。

	go Worker(jobChan, quit)
	time.Sleep(2 * time.Second) //worker中模拟耗时1s,这儿等待2后超时，执行close(quit),此时worker中收到quit 信号
	close(quit)
	//quit <- struct{}{}
	time.Sleep(1 * time.Second)
	fmt.Println("main end")

}

~~~
## <span id="4">四、真正有控制力法(工作池job/worker模式)</span>
相对优雅的方法不能控制goroutine数量，只是延后了请求的爆发。真正有控制力方法就是 job/worker模式。既控制排队任务job，又控制goroutine数量。

~~~go
package main

import (
	"fmt"
	"strconv"
	"time"
)

//待执行的工作者
type Job struct {
	Name string
}

//待执行工作者的队列channal
var JobQueue chan Job

//最大worker线程数
var (
	MaxWorker = 8
)

//执行任务的工作者单元
type Worker struct {
	WorkerPool chan chan Job //工作池(每个元素是一个job的私有channal)
	JobChannel chan Job      //获取Job进行处理
	quit       chan bool     //退出信号
	no         int           //编号
}

//创建一个新worker
func NewWorker(workerPool chan chan Job, no int) Worker {

	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		no:         no,
	}
}

//循环  监听任务和结束信号
func (w Worker) Start() {
	go func() {
		for {
			//将JobChannel放入到工作池中.
			w.WorkerPool <- w.JobChannel
			fmt.Println("w.WorkerPool <- w.JobChannel", w.no)

			select {
			case job := <-w.JobChannel:

				// 收到任务,处理job
				fmt.Println("收到工作任务job,正在处理：", job.Name)

			case <-w.quit:
				// 收到退出信号
				return
			}
		}
	}()
}

// 停止信号
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

//调度中心
type Dispatcher struct {
	//工作者池
	WorkerPool chan chan Job
	//worker数量(开启线程数)
	MaxWorkers int
}

//创建调度中心
func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, MaxWorkers: maxWorkers}
}

//工作者池的初始化
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 1; i < d.MaxWorkers+1; i++ {
		worker := NewWorker(d.WorkerPool, i)
		worker.Start()
	}
	go d.dispatch()
}

//调度
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:

			go func(job Job) {
				//等待空闲worker (任务多的时候会阻塞这里)
				//从WorkerPool中获取一个jobChannel
				jobChannel := <-d.WorkerPool

				// 将任务放到上述woker的私有任务channal中
				// 向jobChannel中发送job,worker的Start()接收端会被唤醒,
				jobChannel <- job

			}(job)
		}
	}
}

func main() {
	JobQueue = make(chan Job, 5)
	dispatcher := NewDispatcher(MaxWorker)
	//开启调度中心
	dispatcher.Run()
	time.Sleep(1 * time.Second)
	//将job添加到队列
	go addQueue()

	time.Sleep(1000 * time.Second)
}

func addQueue() {
	for i := 0; i < 10; i++ {
		// 新建一个工作任务
		job := Job{Name: "tom" + strconv.Itoa(i)}
		// 任务放入任务队列channal
		JobQueue <- job

		time.Sleep(1 * time.Second)
	}
}

/*
一个任务的执行过程如下
JobQueue <- job  新工作任务入队
job := <-JobQueue: 调度中心收到任务
jobChannel := <-d.WorkerPool 从工作者池取到一个工作者
jobChannel <- job 任务给到工作者
job := <-w.JobChannel 工作者取出任务
{{1}} 执行任务
w.WorkerPool <- w.JobChannel 工作者在放回工作者池
*/

~~~

>reference
* [百万实践](https://blog.csdn.net/Jeanphorn/article/details/79018205)
* [实践案例](https://blog.csdn.net/artong0416/article/details/77530843#%E7%9C%9F%E6%AD%A3%E6%8E%A7%E5%88%B6%E5%8D%8F%E7%A8%8B%E6%95%B0%E9%87%8F%E5%B9%B6%E5%8F%91%E6%89%A7%E8%A1%8C%E7%9A%84%E4%BB%BB%E5%8A%A1%E6%95%B0)
* [Golang 任务队列策略](https://juejin.im/entry/5a1675315188254d28733457)
* [百万实践2](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)
