
## RDB(redis DataBase)
* RDB 用于保存和还原Redis服务器所有数据库中的所有键值对数据。
* 两种命令：
 - save:服务器进程执行，期间会堵塞客户端请求
 - bgsave:服务器进程会fork一个子进程去做，执行期间不能执行save，如果执行 bgrewriteAOF 则会排队稍后执行。bgsave和bgrewriteAOF都是对磁盘写操作，不能同时执行。
* RDB执行时机:任意一条件满足，则执行
 - save选项：例如：900秒内存，修改1次；300秒内修改10次，60秒内修改1000次
 - 计数器或者lastSave时间
## AOF(append-only file)

### AOF写入
 客户端请求命令先保存在AOF缓冲区，然后根据appendfsync的配置不同，写入到AOF文件。比如：
 ~~~
 10.1.183.247:6379[29]>sadd animails cat
(integer) 1
10.1.183.247:6379[29]> sadd animails dog tiger
(integer) 2
10.1.183.247:6379[29]> srem animails cat
(integer) 1
10.1.183.247:6379[29]> sadd animails cat lion
(integer) 2
10.1.183.247:6379[29]>
10.1.183.247:6379[29]> smembers animails
1) "lion"
2) "dog"
3) "cat"
4) "tiger"
10.1.183.247:6379[29]>
 ~~~
 * 缺点：缓冲区会写入三笔sadd 和一笔srem，这样的话，随着时间流逝， AOF文件体积变大。为了解决变大问题，采用AOF重写。
### AOF重写
* 原理：首先从数据库中读取现在的值，用一条命令记录键值对，代替之前记录多条命令。减少AOF文件大小，重新写入新文件，替换原来文件。
* 上面例子，重写如下：
~~~
10.1.183.247:6379[29]> sadd animails lion dog cat tiger
~~~
* 缺点：重写可以解决文件体积大小，但大量写操作，调用aof_rewrite函数线程会堵塞，redis单进程单线程，在重写期间服务器将无法处理客户端请求，为了解决这个问题，出现了AOF后台重写
### AOF后台重写(BGREWRITEAOF)
* 原理：为了解决单线程堵塞问题，重写功能交给子进程去做，服务器进程(父进程)可以继续处理命令请求。
* 子进程带有服务器进程副本，不使用子线程是为了避免使用锁的情况下，保证数据的安全性。
* 子进程完成文件重写后，给父进程发送信号，然后父进程用新的AOF文件替换旧的AOF文件，期间只有父进程处理信号函数时有堵塞。


## RDB和AOF区别
RDB占磁盘小，恢复快，容易丢数据；AOF占磁盘大，恢复慢，数据最终一致
> reference
《Redis设计与实现-第十一章》
