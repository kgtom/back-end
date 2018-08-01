## 学习大纲
1. [一、os有调度器，为什么go 需要自己的调度器](#1)
2. [二、go调度器如何工作](#2) 
3. [三、go 调度器源码解析](#3) 


## <span id="1">一、os有调度器，为什么go 需要自己的调度器</span>
1.  os内核有自己的线程调度器，线程上下文切换非常耗时，当线程超过10w后，会发现线程切换是件麻烦的事情，关键是goroutine 不需要这些，反而是一种累赘
2.  os调度器不能做出服务goroutine调度。例如：gc时候 需要所有线程都要停止并等待其达到一定内存状态，若使用os调度器会是大量进程停止工作，此方案肯定不行，而go自己调度器只需要停止在cpu正在运行那个内核线程即可。
## <span id="2">二、go调度器如何工作</span>

### 1.了解三种线程模型

|模型	|含义	|特点|
| - | :- | :- | 
|N:1|	多个用户空间线程会运行在同一个OS线程上|上下文切换很迅速但是并不能利用多核系统的优势
|1:1|	一个线程对应一个OS线程|可以利用机器的所有CPU核心，但是上下文切换比较慢
|M:N|	多个用户空间线程会运行在多个OS线程上|上下文切换很迅速也能利用多核系统的优势，但调度器很复杂

 Go通过使用M:N模型综合前两种模型的优点。它会在多个OS线程上调度多个goroutines。
 ![Sched](/content/images/2018/08/Sched.png)
 
 [点击查看图解](https://github.com/kgtom/go-notes/blob/master/runtime2.md)
 
### 2.goroutine用户态线程
  goroutine就是Go语言提供的一种用户态线程，当我们创建了很多的goroutine，并且它们都是跑在M 内核线程之上的时候，至于在一个M上跑，还是多个M上跑，需要一个调度器来维护这些goroutine，确保所有的goroutine都最大化的使用cpu资源。
  
### 3.一图胜千言
![goroutine](/content/images/2018/07/goroutine.jpg)
* 地鼠(gopher)用小车运着一堆待加工的砖。M就可以看作图中的地鼠，P就是小车，G就是小车里装的砖。如果G 太多了，需要创建更多个M 去干活。没有P，M是不能运砖的。一个M坏了，runtime 将G 放到仓库中(全局队列中),再找新的M去运砖。
* 另外此图说明了[并发与并行的区别](http://www.aqee.net/docs/Concurrency-is-not-Parallelism/#slide-1)，如果只有一堆砖块，几个小地鼠有并发的去运输，并发分四个阶段：装、运输、卸载、送空车。图中两堆砖块，两个工作流程，说明并行去做。
        *  并发：同时(同一时间间隔)处理很多事情，交替做，重点组合。拿庆丰买包子为例，点餐、取餐分开执行，另一个例子：鼠标、键盘、显示器、硬盘同时工作。
        *  学术讲解并发：一种将一个程序分解成小片段独立执行的程序设计。golang中各个独立片段通过channel进行通信，符合csp模式。
        *  并行：同一时刻处理很多事情，同时，重点同时执行。大学餐厅为例，多个点餐处、多个取餐处
        

## <span id="3">三、go 调度器源码解</span>



[详情点击这里查看](https://github.com/kgtom/go-notes/blob/master/runtime.md)



>Reference：
* [github-kgtom](https://github.com/kgtom/go_case/blob/master/2018summary/goroutine%E7%90%86%E8%A7%A3)
* [aqee-Concurrency](http://www.aqee.net/docs/Concurrency-is-not-Parallelism/#slide-19)
* [github-runtime](https://github.com/kgtom/go-notes/blob/master/runtime.md)
* [github-runtime2](https://github.com/kgtom/go-notes/blob/master/runtime2.md)
