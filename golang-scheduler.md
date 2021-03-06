## 学习大纲
* [一、os有调度器，为什么go 需要自己的调度器](#1)
* [二、go调度器如何工作](#2) 
* [三、go 调度器源码解析](#3) 


## <span id="1">一、os有调度器，为什么go 需要自己的调度器</span>
#### 1.  os内核有自己的线程调度器，线程主要的消耗在上下文切换的延迟
   os有一个正在运行线程的列表，os给每一个线程分配一个时间片。保证其公平使用时间片。
   当线程超过10w后，会发现线程切换非常耗时，线程切换平均10us(微妙)，1s内最多10w个，如果是跨核切换，缓存失效更加耗时，关键是**goroutine 不需要这些**，反而是一种累赘;另外虽然有线程池，但线程池过多，频繁请求cpu，必然造成性能下降。
#### 2.  os调度器不能做出服务goroutine调度。
   例如：gc时候 需要所有线程都要停止并等待其达到一定内存状态，若使用os调度器会是大量进程停止工作，此方案肯定不行，而go自己调度器只需要停止正在cpu运行那个内核线程即可。
   
#### 3.另外说一下 os 线程与goroutine两点区别
   * os 线程，分配固定内存栈1M，理论上1G内存创建1000个,goroutine 动态扩容和收缩的栈空间 4kb，理论上1G内存创建250wgoroutine。
   * os 线程被调用时，需要cpu从用户态切换到内核态，调用完成后，再切换回来，若线程多则耗时大，如果是跨核切换，缓存失效更加耗时；goroutine属于两级线程模型，通过线程库将其与内核线程关联起来，避免了cpu用户态与内核态的切换， 比如：一个goroutine调用了time.Sleep或者被channel调用或者mutex操作阻塞时，调度器会使其进入休眠状态并且开始执行另一个goroutine，直到阻塞结束再去唤醒第一个goroutine。因为这种调度方式不需要进入内核的上下文，所以go runtine运行时调度goroutine比内核调度一个线程代价要低得多。
   
## <span id="2">二、go调度器如何工作</span>

### 1.了解三种线程模型

|模型	|含义	|特点|
| - | :- | :- | 
|N:1|	多个用户空间线程会运行在同一个OS线程上|上下文切换很迅速但是只在一个core并不能利用多核系统的优势
|1:1|	一个线程对应一个OS线程|可以利用机器的所有CPU核心，但是上下文切换比较慢
|M:N|	多个用户级线程会运行在多个OS线程上|上下文切换很迅速也能利用多核系统的优势，但调度器很复杂，例子：go的调度器runtine
 Go通过使用M:N模型综合前两种模型的优点。它会在多个OS线程上调度多个goroutines。
  线程库通过内核创建内核线程，线程库将应用程序线程(goroutine用户级线程)与内核线程关联起来，这种两级线程模型，使得内核资源消耗少了，也提高了用户级线程的管理。
 
 ### 1.三者关系
 ![sched](https://github.com/kgtom/back-end/blob/master/pic/sched.png)
 
 
 ### 2.G状态变化
 ![G_state](https://github.com/kgtom/back-end/blob/master/pic/G_state.jpg)

 ### 3.P状态变化
![P_state](https://github.com/kgtom/back-end/blob/master/pic/P_state.jpg)
 [点击查看图解](https://github.com/kgtom/go-notes/blob/master/runtime2.md)
 
 [点击查看图解](https://github.com/kgtom/go-notes/blob/master/runtime2.md)
 
### 2.goroutine用户态线程
 * goroutine就是Go语言提供的一种用户态线程，当我们创建了很多的goroutine，并且它们都是跑在M 内核线程之上的时候，至于在一个M上跑，还是多个M上跑，需要一个调度器来维护这些goroutine，确保所有的goroutine都最大化的使用cpu资源。
  * 当一个goroutine堵塞时，所在M线程会堵塞，但其他goroutine不会堵塞，M会释放P,P会被调度器转移到其他M1线程上。M1的来源：可能是新创建的，也可能是睡眠中的线程。(一般m的数量大于p的数量，多出来的m类似于线程池，随时待命。)
  * G的来源：从各处搜索可运行的G,例如：本地P的可运行G队列、调度器的可运行G队列、其它P的可运行G队列
  * goroutine：让研发人员更加专注于业务逻辑，从os层面的逻辑抽离出来。
  * G的数量比M多，go用少量的内核线程M支撑了大量goroutine的并发，多个goroutine共享一个内核线程，对于操作系统而言没有线程切换的开销
  * 为了充分利用线程和CPU资源，调度器采取以下两种调度策略：
    - 任务窃取(work-stealing):为了充分利用计算机资源，不能让p有的忙，有得闲，窃取其它地方的G(全局G列表、本地p队列，其它p上)
    - 减少堵塞：（正在运行G塞M）
        - 场景1：由于互斥量、channel堵塞，调度器将G调度出去，执行p上面的其它G
        - 场景2：网络请求和IO操作将当前G先放在网络轮询器去执行，执行完毕后，再放回本地p等待M的执行
        - 场景3：系统调用堵塞M，将当前G调度给其它M1，M等之后复用
        - 场景4：执行sleep,M回去执行其它G,当前G退下来，等待时机再次执行
  
  
### 3.一图胜千言
![goroutine](https://github.com/kgtom/back-end/blob/master/pic/goroutine.jpg)
* 地鼠(gopher)用小车运着一堆待加工的砖。M就可以看作图中的地鼠，P就是小车，G就是小车里装的砖。如果G 太多了，需要创建更多个M 去干活。没有P，M是不能运砖的。一个M坏了，runtime 将G 放到仓库中(全局队列中),再找新的M去运砖。
* 另外此图说明了[并发与并行的区别](http://www.aqee.net/docs/Concurrency-is-not-Parallelism/#slide-1)，如果只有一堆砖块，几个小地鼠有并发的去运输，并发分四个阶段：装、运输、卸载、送空车。图中两堆砖块，两个工作流程，说明并行去做。
        *  并发：同时(同一时间间隔)处理很多事情，交替做，重点组合。拿庆丰买包子为例，点餐、取餐分开执行，另一个例子：鼠标、键盘、显示器、硬盘同时工作。
        *  学术讲解并发：一种将一个程序分解成小片段独立执行的程序设计。golang中各个独立片段通过channel进行通信，符合csp模式。
        *  并行：同一时刻处理很多事情，同时，重点同时执行。大学餐厅为例，多个点餐处、多个取餐处
        *  并行一定是多核cpu,单核只有并发。

## <span id="3">三、go 调度器源码解</span>


### go scheduler 
`源码文件地址：/src/runtime/proc.go:15行`
~~~
 The scheduler's job is to distribute ready-to-run goroutines over worker threads.  
//调度器任务：将goroutne 分配到可用的工作线程上。
 The main concepts are:  
 G - goroutine.  
 M - worker thread, or machine.  
 P - processor, a resource that is required to execute Go code.  
M must have an associated P to execute Go code, however it can be blocked or in a syscall w/o an associated P.
//M 必须关联一个P才能执行go代码。
~~~


### bootstrap 启动顺序
~~~

// The bootstrap sequence is:  
//  
// call osinit  //os初始化
// call schedinit  //调度器初始化
// make & queue new G  //生成G，放到G队列，等待运行
// call runtime·mstart //调用mstart方法，启动M
~~~
1.  osinit ：ncpu的数量也就是P最大数量，也就是可运行G的队列的数量，实际上对并发运行G规模的一种限制。
2.  schedinit ：M最大数量初始设置1000，从环境变量GOMAXPROCS中获取P
3. make & queue new G 即：runtime·newproc ：生成G，放到队列中，等待运行
4. runtime·mstart ：生成G之后就需要启动M,并且将P关联到M上，没有P，M也无法启动，关联上P之后进入schedule（）开始干活。

`源码地址：src/runtime/proc.go中 mstart1() 1289行`
~~~go
    else if _g_.m != &m0 {
		acquirep(_g_.m.nextp.ptr())
		_g_.m.nextp = 0
	}
	schedule()
~~~

`源码地址：src/runtime/proc.go schedule()  2557行`

~~~go
if gp == nil {
	//从全局可运行G队列查找G
		if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {
			lock(&sched.lock)
			gp = globrunqget(_g_.m.p.ptr(), 1)
			unlock(&sched.lock)
		}
	}
	if gp == nil {
    //从本地P可运行G队列找G并判断是否是自旋状态
		gp, inheritTime = runqget(_g_.m.p.ptr())
		if gp != nil && _g_.m.spinning {
			throw("schedule: spinning with local work")
		}
	}
	if gp == nil {
    //从调度器的可运行G队列中查找G
		gp, inheritTime = findrunnable() // blocks until work is available
	}

	// This thread is going to run a goroutine and is not spinning anymore,
	// so if it was marked as spinning we need to reset it now and potentially
	// start a new spinning M.
    //重置自旋状态，然后启动M
	if _g_.m.spinning {
		resetspinning()
	}
  //G 开始运行了
	execute(gp, inheritTime)

~~~



### 上下文切换

* 保存现场：G不在占用cpu时间片运行的时候(gopark/gosched)，将相关寄存器的值给保存到内存中；
`源码地址：src/runtime/proc.go Gosched()  265行`
~~~go
func Gosched() {
	checkTimeouts()
	mcall(gosched_m)
}
~~~

`源码地址：src/runtime/proc.go gopark()  285行`
~~~go
gopark(){
	// can't do anything that might move the G between Ms here.
	mcall(park_m)
}
~~~
* 恢复现场:在G重新获得cpu时间片的时候(execute)，需要从内存把之前的寄存器信息全部放回到相应寄存器中去。




[详情点击这里查看](https://github.com/kgtom/go-notes/blob/master/runtime.md)



>Reference：
* [github-kgtom](https://github.com/kgtom/go_case/blob/master/2018summary/goroutine%E7%90%86%E8%A7%A3)
* [aqee-Concurrency](http://www.aqee.net/docs/Concurrency-is-not-Parallelism/#slide-19)
* [github-runtime](https://github.com/kgtom/go-notes/blob/master/runtime.md)
* [github-runtime2](https://github.com/kgtom/go-notes/blob/master/runtime2.md)
* [github-GoHackers](https://github.com/GoHackers/talks/blob/master/20150110/3.Go%E5%B9%B6%E5%8F%91%E7%BC%96%E7%A8%8B-%E9%83%9D%E6%9E%97.pdf)
