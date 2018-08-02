## 学习大纲
* [一、理解并发与并行](#1)
* [二、go并发编程(goroutine、channel、select)](#2)
* [三、小结](#3)


## 一、理解并发与并行

   并发编程的思想来自于多任务操作系统，允许终端用户同时运行多个程序。当一个线程不需要cup时，系统内核会把该线程挂起或者中断，让其它线程使用cpu。
 

并发是指两个或两个以上操作同时存在。(交替做事，如果是单核，两个或多个线程交替被cpu执行)
例子：买包子一个取号处---->一个取餐处
并行是指同时做多件事。（同时做事，必须是多核cpu，两个或多个线程分别分配到不同cpu上去执行）
例子：买包子多个取号处---->多个取餐处

并行是并发的一个子集。
 
## 二、go并发编程
  
### 1.Goroutine

Goroutine是由Go的运行时创建和管理，做为研发人员使用Go语句提交了一个并发任务，然后Go运行时系统并发的去执行。过程如下：
1.检查go函数及参数的合法性
2.本地p自由G列表或者调度器自由G列表获取可用的G
3.没有可用G就新建一个G,然后执行初始化操作对应G的状态GIdle/Grunnable
4.把这个G立刻放到P的runnext字段，等待运行，如果runnext字段有G，则会推到p可运行队列或者调度器可运行队列的队尾。
 
**注意：**
在执行上述4步之前，主goroutine其实默默的做了不少工作，设置栈空间、m0检测、初始化工作，如下：
设置每一个goroutine最大栈空间，32位机器 250M,64位机器 1G.
执行m0检测任务，与哪个p关联，寻找可运行G
创建特殊 defer语句，主goroutine退出时做一些善后工作
启用专门清理后台垃圾的goroutine，用于GC
执行 main 函数中init(), 然后在执行mian函数中的其它函数
 
Goroutine初始化时执行g0函数初始化线程栈2KB内存空间(线程1M)，理论上可用创建几十万goroutine.


* 暂停主main goroutine，让 go func()有机会得到cup时间片运行
~~~go
go func() {
 
        fmt.Println("hello world form goroutine")
    }()
    time.Sleep(1e9)
    //让main goroutine暂停一1s,让goroutine有机会执行
    //或者使用runtine.Gosched()暂停当前G,让其它G有机会运行
 ~~~
 

* 开启三个 goroutine，输出的结果不正确
~~~go
mySlice := []string{"tom", "lilei", "jim"}
 
    for_, name := range mySlice {
        go func() {
            fmt.Println("name:", name)
        }()
        //执行输出3条记录结果：name:jim，因为for语句执行太快，执行到最后一个循环name=jim,然后
        //go func()才得以执行，如何正确执行,解决见goroutine_2()
    }
 ~~~
 
* 解决方案：保证每一个goroutine都可以执行
~~~go
mySlice := []string{"tom", "lilei", "jim"}
    for_, name := range mySlice {
        go func(val string) {
            fmt.Println("name:", val)
        }(name)
        runtime.Gosched()
        //暂停一下当前G,让其它G得以执行，保证可以正确执行，但不能保证执行顺序
        //若要保证执行顺序，这涉及到同步范畴，使用channel,解决见goroutine_3()
    }
 
 ~~~
 
* 使用channel 同步解决顺序问题：

~~~go
mySlice := []string{"tom", "lilei", "jim"}
    ch := make(chan string)
    for_, name := range mySlice {
        go func(val string) {
 
            fmt.Println("syncname:", name)
            ch<- name
        }(name)
        <-ch
    }
 ~~~
 
### 2.Chnanel
 Channel多个goroutine之间同步和数据通信的一种并发安全管道。属于CSP(CommunicatingSequential Process)并发通信模型。
 
 **分为有缓冲和无缓冲**

* 无缓冲：同步作用，发送与接受同时进行。
~~~go
ch := make(chanbool)
    go func() {
 
        fmt.Println("dosomething")
        ch<- true
    }()
    <-ch
 ~~~
 
* 有缓冲：类似于队列，直到缓冲区容量满了，就无法继续写入。

~~~go
ch := make(chan int, 3)
    go func() {
        for i := 0; i < 5; i++ {
            ch<- i
        }
        close(ch) //不关闭的话，range一直处于接收状态 会造成deadlock
        //ch <- 8//panic:"send on closed channel"，对已关闭的channel再发送数据会panic
    }()
    for val := range ch {
        fmt.Println("i:", val)
    }
 ~~~
 
如果channel关闭的话，range就不会接收到值了，如果使用for接收的话，需要注意判断channel是否关闭。
~~~go
ch := make(chanint, 3)
    gofunc() {
        fori := 0; i < 5; i++ {
            ch<- i
        }
        close(ch) //不关闭的话，range 会造成deadlock
        //ch <- 8//panic:"send on closed channel"，对已关闭的channel再发送数据会panic
    }()
    
    for {
 
        ifval, ok := <-ch; ok{ //判断channel是否关闭
            fmt.Println("i:", val, "iClose:", ok)
        }else {
            fmt.Println("i:", val, "iClose:", ok)
//关闭的channel,可以接收，但不能写入
            return                
               //不加retrun的话，就无限循环下去了
        }
    }
 ~~~
 
### 3.Select:
Select 语句是一种用于通道发送和接受操作的专用语句。

~~~go 
package main
 
import (
    "fmt"
    "time"
)
 
func doWorker(c chanint, quit chanbool) {
    x, y := 0, 1
    for {
        select {
        case c <- x:
            x, y = y, x+y
        case <-quit:
            fmt.Println("it is timeto  quit")
            return
        }
    }
}
func main() {
    fmt.Println("mainroutine start")
    //demo 1 :主要练习 select case 的用法
    ch := make(chanint)
    timeout := make(chanbool, 1)
    go func() {
 
        //第一种方式：注释 time。Sleep(1e9)，两个case 都满足，go 随机选择一个执行
        //  //time.Sleep(1e9)/
        //  timeout<- true
        //  ch<- 3
 
        //第二种方式：ch <- 3，case 输出 ch
        // ch <- 3
        //time.Sleep(1e9)
        // timeout <-true
 
        //第三种方式：两个chan 没有send ，case 输出 default，这种情况下若没有default 分支，造成deadlock
        //ch <- 3
        time.Sleep(1e9)
        //timeout <-true
 
    }()
 
    select {
    caseval := <-ch:
        fmt.Println("ch...val:", val)
    case <-timeout:
        fmt.Println("timeout")
    default:
        fmt.Println("default")
    }
 
    //demo 2：主要练习开启的goroutine 完成操作后，发送信号通知main gorotine 结束。
 
    // c :=make(chan int)
    // quit :=make(chan bool)
    // go func() {
    //  for i:= 0; i < 10; i++ {
    //     fmt.Println(<-c)
    //  }
    //  //执行操作后，发生quit信号，通知main routine 结束
    //  quit<- true
    // }()
    // doWorker(c,quit)
    fmt.Println("mainroutine end")
}
 ~~~
 
* 如果有同时多个case接收数据,那么Go会伪随机的选择一个case处理(pseudo-random)。
* 如果没有case需要处理，则会选择default去处理。
* 如果没有default case，且有发送的chan 则select语句会阻塞，直到某个case需要处理。
 
 
## 小结：
* Goroutine: 让研发人员更加专注于业务逻辑，从os层面的逻辑抽离出来。
* Channel:简单、安全、高效的实现了多个goroutine之间的同步与信息传递。
* Select:可以处理多个channel。

> Reference
《郝大神go并发编程》
