## 本节大纲
* [第一章 mysql体系机构与存储引擎](#1)
* [第二章 Innodb存储引擎](#2)
* [第三章 文件](#3)
* [第四章 表](#4)
* [第五章 索引与算法](#5)
* [第六章 锁](#6)
* [第七章 事务](#7)
* [总结](#8)

## 前言
* mysql 使用广，其设计思路很有启发性，比如缓冲池与磁盘数据同步设计思想，与日常redis cache设计异曲同工。
* 多读经典，扩展思路
## <span id="1">第一章 mysql体系机构与存储引擎</span>

### 1.区分数据库与实例
* 数据库：文件集合
* 实例：是程序，对数据库操作都是在实例下进行的

### 2.体系结构
* server层：连接池组件、sql接口、缓存组件、查询分析器、优化分析器、物理文件
* 引擎层： 负责数据存储与读取 **插件式存储引擎**
* mysql数据库区别于其他数据库一个特点：插件式的表存储引擎，用户根据mysql存储引擎接口可以自定义自己存储引擎
### 3.存储引擎分类
* Innodb:支持事务，设计目标OLTP(在线事务处理online trans process),行级锁、支持外键、高并发性(通过mvcc多版本并发控制)、next-key-locking 策略避免幻读、聚集索引、insert buffer、double wrinte、自适应hash索引等
* MyISAM:不支持事务、表锁定并发性差、支持全文索引、面向OLAP(在线分析处理online analyze process)设计,适用于查询多场景，
* Memory:将表数据存放内存中，适用于临时数据的临时表以及数据仓库的维度，默认支持hash索引，不支持test、varchar按照定常字段进行，所以浪费空间。
* Archive:只支持insert和select操作，适用于高速插入和压缩功能

### 4.连接mysql常用方式
 连接mysql操作是一个连接进程和mysql数据库实例进行通信，从程序设计角度理解：进程之间的通信。
* tcp/ip套接字连接，如：mysql -h 192.168.0.10 -u tom -p,这种方式首先验证ip是否有权限连接
* 命名管道和共享内存
* unix域套接字，如：
 ~~~
 mysql -u tom -S /var/lib/mysql/mysql.sock
 ~~~


## <span id="2">第二章 Innodb存储引擎</span>

### 1.innodb体系架构
* 多个内存块，组成内存池，包括缓冲池innodb_buffer_pool、重做日志缓冲区redo log buffer、额外的内存池 additional memory pool
* 后台进程：刷新缓冲池中的数据，保证缓冲池和磁盘数据一致；db异常时，及时恢复

### 2.内存池
* 缓冲池：解决cpu速度与磁盘速度不匹配问题，提高数据库整体性能。先将读取页的数据放到缓冲池,再次请求的时候先检查pool中有没有，没有的话再读取磁盘上的页。
* 重做日志缓冲：缓存重做日志信息，然后一定策略(master thread、事务提交、pool空间小于1/2)刷新到磁盘的重做日志文件
* 额外的内存池：用来缓存LRU链表、等待、锁等数据结构。
* 缓冲池的管理：最近最少使用算法LRU，频繁使用的放在LRU List前端、当LRU为空的时候从磁盘获取数据先放在Free List 和 当LRUlist数据有更新时(脏页)放在Flush List，通过checkPoint刷新到磁盘
### 3.后台进程
* master thread :负责将缓冲池数据异步刷新到磁盘，保存数据一致性，
* IO thread：处理异步IO请求，分别用来处理insert buffer、重做日志、读写请求的IO回调  
* purge thread:用来回收undo页
* page cleaner thread：用来刷新脏页，减钱master thread 工作，提高引擎性能
### 4.checkPoint检查点技术
* 两种情况下不需要checkPoint技术：a.内存足够大，缓冲池可以缓存所有数据 b.重做日志可以无限增大
* 需要checkPoint技术：当数据库宕机时，重做日志时间非常久，恢复需要付出很大代价，因此checkPoint技术解决以下问题：即：redo log不可用，刷新缓存脏页，解决内存瓶颈，缩短恢复时间
 - 缩短数据库的恢复时间：checkPoint之前不用恢复，只需要恢复checkPoint之后的数据
 - 缓冲池不够用时，将脏页刷新到磁盘
 - 重组日志不可用时，将缓冲池页强制刷新到当前重做日志的位置

* checkPoint具体做：a.master thead 的定时刷新；b.LRUlist没有足够空间(脏页太多);c.redolog 不可用，d.数据库关闭
* checkPoint两种工作模式：默认fuzzy checkpoint 刷新局部脏页，另一种sharp checkpoint刷新所有页

### 5.innodb关键特性
* 插入缓冲Insert Buffer:对非聚集索引(非唯一)的插入或更新，不是每一次都插入到索引页，而是先判断是否在缓冲池，若在则直接插入，若不在先放入inset buffer,再一定频率将其与辅助索引页节点合并，提高插入的性能
* 两次写Double write:背景数据库宕机了，页损坏了，重做日志无法恢复了，所以重做日志之前，先写副本，如果页失效后，先从副本还原，再进行重做，提高数据页的可靠性
* 自适应哈希索引AdaptiveHash Index:数据库自优化，根据访问频率100次及查询条件访问的页超过1/16，等值查询做优化O(1)
* 异步IO Async IO:多个IO合并为一个，提高磁盘操作性能
* 刷新邻接页 Flush Neighbor Page:相邻页如果存在脏页一起刷新。


## <span id="3">第三 章文件</span>

### 1.参数文件
* 配置参数的文件，通过 mysql --help | grep my.cnf来查找
### 2.日志文件
* 错误文件error log:针对mysql启动运行关闭过程的记录；定位文件路径位置如下：
~~~
mysql> show variables like 'log_error';
+---------------+---------------------+
| Variable_name | Value               |
+---------------+---------------------+
| log_error     | /var/log/mysqld.log |
+---------------+---------------------+
1 row in set (0.00 sec)

~~~
* 二进制文件binlog：对数据库所有更改的操作，用于：数据库恢复、复制、审计(是否有攻击)，保证数据一致性
* 慢查询日志slow query log:是否开启此功能；设置long_query_time 阈值；通过mysqldumpslow命令分析文件，比如top10:mysqldumpslow -s -al -n 10 xxx.log
* 查询日志 log：记录所有对mysql数据库请求的信息
### 3.套接字文件
* unix系统下，本地连接mysql可以使用套接字方式，这样就需要一个套接字socket文件,查找文件位置：
~~~
mysql> show variables like 'socket';
+---------------+---------------------------+
| Variable_name | Value                     |
+---------------+---------------------------+
| socket        | /var/lib/mysql/mysql.sock |
+---------------+---------------------------+
1 row in set (0.00 sec)
~~~
### 4.pid文件
* mysql 启动，将自己进程写入一个文件中，文件路径

~~~
mysql> show variables like 'pid_file';
+---------------+----------------------------+
| Variable_name | Value                      |
+---------------+----------------------------+
| pid_file      | /var/run/mysqld/mysqld.pid |
+---------------+----------------------------+
1 row in set (0.00 sec)
~~~
### 5.表结构文件
* mysql数据存储是根据表进行的，每个表都会有与之对应的文件，frm为后缀名的文件，记录表的结构定义。
### 6.innodb存储引擎文件
* 重做日志文件：默认两个ib_logfile0,ib_logfile1,存储事务日志，如果宕机，使用该文件恢复到宕机前的时刻，保证数据的完整性，重做日志文件与二进制文件区别如下：
 - 记录范围：前者只记录innodb引擎的事务日志，后者记录mysql所有引擎的记录
 - 记录内容：前者记录页的更改的情况，后者记录一个事务具体操作内容
 - 写入时间：前者只要事务进行，就会写入，后者只能事务提交前提交，写磁盘一次

* 重做日志文件写入时机(第二章写过)，根据 innodb_flush_log_at_trx_commit 参数具体再写一下：
  - 当参数0：等待主线程每秒刷新，将redo log buffer 写入磁盘文件，如果宕机了，前一秒数据会丢失
  - 当参数1：事务提交时，将redo log buffer同步写入到磁盘
  - 当参数2：先写到redo log buffer，但buffer空间小于1/2时 ,异步写入到磁盘
  - 总结：为了保证事务ACID的持久性，必须参数设置为1或者2
 
* 表空间文件：管理innodb存储引擎的存储，分为共享表空间和独立表空间(设置innodb_file_pre_table之后)

## <span id="4">第四章 表</span>

### 1.索引组织表
* 表都是根据主键顺序组织存放的，这种方式称为表的索引组织表
* 如果没有主键，选择非空唯一索引，如果没有，则自动创建6字节_rowid
### 2.逻辑存储结构
* innodb存储逻辑结构最高层：表空间(第三章写过)，表空间由段segment、区extent、页page 组成，详细如下：
* 段：数据段b+树叶子节点，索引段b+树非叶子节点
* 区：一个区默认1M，每页默认16k,一个区有64个页
* 页：磁盘管理最小单位，每一页最多7992行记录，对于blog、varchar 可能会行溢出，数据页之外
### 3.约束
* 约束机制：保证数据完整性，包括：主键约束、外键约束、唯一约束、默认为空约束及触发器
* 约束与索引区别：
 - 前者保证数据的完整性的机制，后者存储在数据页上的数据结构，便于查询，一种存储方式
### 4.分区
* 数据库的应用分为两类：OLTP(在线事务处理)和OLAP(在线分析处理)
* OLTP:不建议分区，因为大部分场景通过索引返回几条记录。因为B+数，1000w表3次io，分区的话，反而增加会增加io
* OLAP:建议分区,提高查询性能。例如1亿行的表，查询用户某一年数据，可以按照时间戳分区，只扫描响应分区即可


## <span id="5">第五章 索引与算法</span>

### 1.数据结构与算法
* 二分查找：折半查找 O(logN)
* 二叉查找树-->平衡二叉树--->B+树:O(logN)
### 2.B+树索引
* 聚集索引:按照表的主键构造B+树，叶子节点存放数据，各叶子节点指针相连，即每个数据页都是通过一个双向链表相连。通常查询分析器采用聚集索引，通过主键查先查找到叶子节点，然后直接在叶子节点上根据范围查找或者排序查找数据

ps：B+树 二分查找数据所在的页，将页放在内存中，提高了查询速度
~~~

-- 主键范围查找,索引组织表通过非叶子节点，二分查找法定位数据在哪一页？然后再去页上二分查找找到记录
select id,name from user where id>1 and id<5
-- 主键排序查找
select id，name from user order by id limit 10 

-- 分页，尽量将offset换成id,重发利于b+树的特性。找到max_id的地方，取10条
select id,name from user where col=xxx and limit 100,10
select id,name from user where id>max_id limit 10
~~~
* 非聚集索引:叶子节点不仅存放数据，也存放键值对(索引键和主键),通过非聚集索引key，找到主键，然后根据主键找到行记录。如果一个非聚集索引数高度3，则需要3次io找到主键，如果聚集索引高度也是3，则同样需要3次io找到数据，共6次io。所以查询效率取决于树的高度。

~~~
-- 二次查询(回表查询)：例如非主键索引 id_card，先根据 id_card 找到id为1，再根据聚集索引找id=1的记录，找到name。即：从索引树找到 id_card 对应的聚集索引id,然后回到聚集索引中找到行记录的name
select id,name from user where id_card="110xxxx"

-- 不需要二次查询(回表)，例如非主键索引idCard，直接在索引树上，根据idCard找到记录
select id,id_card from user where id_card="110xxxx"

-- 结合复合索引或者覆盖索引，避免二次回表查询,建立复合索引(idcard,name),注意最左匹配原则
select id,name,id_card from  user where id_card="110xxxx"
~~~
### 3.哈希索引
* 哈希索引是自适应索引，根据查询频率和条件自动生成索引

### 4.全文索引
* 有两个列，一个word字段，另一个ilist字段，并且在word字段上有设立索引，ilist存放位置信息，故可以进行 相邻节点查找proximity search
* FTS Index cache 全文检索索引缓存是一个红黑树结构，其根据word、ilist进行排序。

### 5.覆盖检索
* 覆盖索引covering index：从非聚集索引中可以查询到记录，不需要到聚集索引中查询。

~~~
-- 非聚集索引，由于包含主键信息，所以叶子节点存放(id主键、idx_key非主键索引),下面sql用覆盖索引就可以查询到
select id，key from user where key=xxx
select key from user where key=xxx
select count(*) from user //数据量小的时候，优化器使用idx_key索引，而没有使用id聚集索引全表扫描,层级少，减少io，建议不使用count(*),使用count(idx_key)或者count(id)
~~~
### 6.联合索引
* 建立多列的索引，注意结合最左匹配原则，适用于查询或者排序
* 例如建立(a,b)，可以等同于 (a)、(ab)两个个索引，遇到范围查询(>、<、between、like)之后就用不上了
* 例如：建立(a,b,c)索引，等同于a,ab,abc三个索引
~~~
-- 用到索引，遵循最左匹配
where a=1 and b=2 and c=3
where a>1
where a=1 and b in(1,2) 
where a=1 and b=3 and c>5 

--不能用，不以a开头查询条件
where b=1 and c=2
where b>2
where a=1 or b=1 //应该单独建立各自索引

-- 部分用到
where a=1 and b>3 and c=4 //只用到a=1 and b>3，c的时候在前两个条件的数据里面，再找出c
where a>1 and b=3 //只用到a索引

-- 排序必须相同顺序
where a=1 order  by b desc  
where a=1 order by a asc  b desc  //不能用，顺序不同

~~~

### 7.全文检索
* 语法：match ...against()
* 倒排索引：存储单词与单词直接映射，通常使用关联数组实现。

~~~
-- 不能使用b+树
select id,name from user where name like '%tom%'

-- 使用全文检索
select id,name from user where match(name) against('tom')
~~~


## <span id="6">第六章 锁</span>

### 1.lock锁的分类
* 数据库使用锁目的：支持对共享资源并发访问时，提供数据的一致性和完整性
* lock 与latch区别：前者锁的对象是事务，包括行锁、表锁、意向锁，通过wait-for graph和timeout进行死锁检查与处理；后者锁的对象是临界资源，包括 互斥锁mutex和读写锁rwlock，没有死锁检查机制，全凭程序代码控制。
* lock 分类
 - 行级锁:共享锁 S Lock 和排它锁 X Lock
 - 表级锁:意向共享锁 IS Lock和意向排它锁 IX Lock,假如给某一行上X锁，则该行对应的表和页都上IX锁。如果加锁之前该表已有S锁，则需等待其释放后再加锁，即：IX与S锁不兼容，需等待其释放再用
 - 一致性非锁定锁：如果读取的行正在做update、delete操作，此时读取不需要等待行上的X锁，读取该行之前版本的快照数据。一个行的快照版本较多，则称为行的多版本技术。不同事务隔离级别读取版本不同，在READ COMIMTED 读取最新快照；在默认REPEATABLE READ 读取事务开始时的行数据版本。
 - 一致性锁定读：包括 select ...for update 和select ...lock in share mode 
 - 自增长锁：自增长值，通过 AUTO-INC locking，锁表，不是事务结束释放而是执行插入sql后立刻释放
 ~~~
 select max(auto_incr_col) from t for update
 ~~~
### 2.行锁的三种算法
* Record lock 单个行记录上的锁：锁定索引记录，若没有索引，则隐式锁住主键
~~~
-- id 主键,此时如果有另一个事务操作id=3的记录，则需要等待正在进行事务commit或者rollback
start transaction
select * from t where id=3 for update 
update t set xxx=xxxx where id=3
~~~

* Gap lock 间隙锁：锁定一个范围，但不包括记录本身，该锁母的：防止同一个事务的两次读取的数据不同，出现幻读的情况。两种方式关闭gap lock:innodb_locks_unsafe_for_binlog=1和隔离级别READ COMMITED

* Next-key lock ：record lock+gap lock，解决phantom problem幻读的问题，注意左开右闭区间。
~~~
--  当查询的列是主键或者唯一索引的时候，next-key lock 降级为 Record lock
 id 为t 表主键，则执行以下sql不需要等待锁释放
 inset into t select id=1
 inset into t select id=2
 
 -- 当查询的列是非聚集索引(不是主键或唯一索引)时候，next-key lock会锁定两个范围，一个（-min，xx）和(xxx,+max)
 a 是t表主键，b为t表普通索引
 b在t中有 1，2，3，4，5
 正在执行：select * from t where b=3 for update ，
 则下面sql 不执行了,因为锁定区间（1,3）和(3,5)
 insert into t (8,2)
 insert into t (8,5)
 
 -- ？补充区间的案例
~~~
### 3.锁问题
通过锁机制可以实现事务的隔离性，使事务可以并发工作。锁提高了并发，但带来了三种潜在问题：
* 脏读 Dirty Read:不同事务之间，当前事务读取到另一个事务未提交的数据(脏数据)，解决 REPEATABLE READ 级别
* 不可重复的读即幻读：不同事务之间，当前事务读取到了另一个事务已提交的数据，解决方案next-key lock
* 丢失更新：一个事务的更新操作会被另一个事务更新操作覆盖，解决方案：将事务处理串行，不要并行做，即加排它锁X。
### 4.堵塞与死锁
* 堵塞：等待另一个事务释放锁，再执行该事务。innodb_lock_wait_timeout 等待时间，默认50s
* 死锁：两个或两个以上事务争夺资源造成相互等待现象。解决方案 超时机制和wait-for graph 等待图

## <span id="7">第七章 事务 </span>

### 1.事务的分类
* 扁平事务：最简单的事务，要么成功，要么失败，不能提交或回滚某一部分
* 带保存节点事务：可以回滚到任意保存的节点
* 链事务：只能回滚到最近的一个保存点
* 分布式事务：通过事务管理器管理不同事务，遵循ACID,要么全部成功，要么全部回滚。
### 2.事务的实现
* redo log重做日志:提交事务修改的页操作，保证事务原子性和持久性。还包括 redo log buffer ，两者之间同步详见 第三章--innodb存储文件
* undo log:保证事务的一致性，当事务执行失败或者回滚操作时，回滚行记录到某个特定版本。

总结：innodb对数据库的修改，不仅产生redo log ，还要产生undo log。
### 3.事务的隔离级别
* Serializable（串行化）：串行化读，每次读都需要获取表级锁，读写会堵塞；
* Repeatable read（可重复读）：可避免脏读、不可重复读(幻读)的发生；
* Read committed（读已提交）：可避免脏读的发生，不能避免幻读即能读取到其它事务已提交的数据；
* Read uncommitted（读未提交）：允许脏读；级别最低
隔离级别由高到低，级别越高，执行效率越低

### 4.分布式事务
* 分布式事务的实现是：应用通过一个事务管理器实现对多个相同或不同的数据库实例的事务管理。
* 分布式事务与本地事务的区别是多了一个prepare的阶段，待收到所有节点的同意信息后再commit或rollback。
  ### 5.内部XA分布式事务
内部XA最常见的是binlog和innodb存储引擎。事务提交时会先写binlog再写redo log，在写完binlog宕机的情况下，mysql重启会先检查prepare的事务链表中的uxid事务是否已经提交，若没有则存储引擎层再做一次提交。

### 6.不好的事务习惯





## <span id="8">总结</span>

### 1.一条sql执行过程

* server层：一条SQL进入MySQL服务器，会依次经过连接池组件（进行鉴权，生成线程），查询缓存组件（是否被缓存过），SQL接口模块（简单的语法校验），查询解析模块(语法检查，sql是做什么的，生成语法树)，优化器模块（sql怎么做，索引的选择、表之间join），然后再进入innodb存储引擎。
* 引擎层：进入innodb后，首先会判断该SQL涉及到的页是否存在于缓存池中，如果不存在则从磁盘读取相应索引及数据页加载至缓存池。
 - 如果是select语句，读取数据(使用一致性非锁定读)，并将查询结果返回至服务器层
 - 如果是DML语句，读取到相关页，先试图给这个SQL涉及到的记录加行级排它锁X锁。加锁成功后，先写undo 页，逻辑地记录这些记录修改前的状态。然后再修改相关记录，这些操作会同步物理地记录至redo log buffer。
 - 如果涉及及非唯一辅助索引的更新，还需要使用insert buffer。
 - 事务提交时，会启用内部分布式事务，先将SQL语句记录到binlog中，再根据系统设置刷新redo log buffer至redo log，保证binlog与redo log的一致性。
 - 提交后，事务会释放对这些记录所加的锁，并将这些修改的记录所在的页放入innodb的flush list中，等待被page cleaner thread刷新到磁盘。
 - 这个事务产生的undo page如果没有被其它事务引用(insert的undo page不会被其它事务引用)，就会被放入history list中，等待被purge线程回收。

~~~
　　需要注意的是：
　　a.脏页的刷新采用的是checkpoint机制
　　b.DML语句不同undo页的格式也会不同。insert类型的undo log只记录了主键及对应的主键值，而update、delete则记录了主键及所有变更的字段值
　　c.一条设计不好的SQL，可能会导致大量的离散读、加载很多冗余的数据页至缓存中
  
  ~~~
  

### B+树高度3能存储多少行记录
* 前提条件：主键id,bigint 8个字节，一条行记录大小1k，1个页16k,也就说1页上有16条行记录。即：innodb_page_size:16K*1024=16384字节，每个键包括主键8个字节+指针6个字节=14个字节
~~~
mysql> show variables like 'innodb_page_size';
+------------------+-------+
| Variable_name    | Value |
+------------------+-------+
| innodb_page_size | 16384 |
+------------------+-------+
~~~
* 第1层树存储索引键个数：16384/14=1170个索引键，第1层就1页，即：1页上最多1170键值对。
* 第2层也是只存放索引键，从第1层眼神1170个键对应1170个页，每一个页上最多1170键值对，所以共有：页数 * 每页键值对=1170 * 1170=1368900键值对；假如存放数据的话，行记录数=节点索引键个数*单个叶子节点行记录数即：1170 * 16条行记录=18720行记录。
* 第3层：存储数据，从第二层共有1170*1170键值对*16=2千万左右行。

**总结**：数的高度为3时，可以存储2千万条记录

### 数据库三范式
* 每个列都**不可以再拆分**
* 非主键列完全依赖于主键,而不能是依赖于主键的一部分，**没有部分依赖**，保证一张表只描述一件事情，存在重复数据，需要拆分
* 非主键列只依赖于主键,不依赖于其他非主键,**消除传递依赖**，保证每一列与主键直接相关,a-->b-->c不能存在传递关系


~~~
student 表（学号，姓名，年龄，性别，院校名称，院校地址，院校电话）
存在：学号--->院校名称--->院校地址和电话
拆分：
student表:学号，姓名，年龄，性别，院校编号
Academy 院校表:院校编号,院校名称，院校地址，院校电话
~~~

总结：三范式是设计数据库理念，尽量遵守三范式,有时候处于性能考虑，接受冗余数据，方便查询。

### MaxCompute与关系数据库的区别
* 应用场景：前者适合海量存储和大数据分析，不适合做在线事务处理OLTP，后者适合OLTP
* DML：前者不支持update、delete，只有drop,后者都可以
* 前者不支持 索引、约束、主键，后者都可以

### varchar 和char 区别
* 定长和变长，比如"tom",存储char(10)，剩余7个空格字符，varchar只占3个字符
* 存储容量，varchar 最大65535个字符，char 255个字符
* 比如 用户身份证号，手机号 固定长度的字符串应该使用char而不是varchar来存储,这样可以节省空间且提高检索效率.

### 主从同步原理分析
* 1.Master 数据库只要发生变化，立马记录到binlog 日志文件中
* 2.Slave数据库启动一个I/O thread连接Master数据库，请求Master变化的binlog
* 3.Slave I/O获取到的binlog，保存到自己的Relay log 日志文件中。
* 4.Slave 有一个 SQL thread定时检查Realy log是否变化，变化那么就更新数据


### InnoDB 聚集索引与非聚集索引区别？

 **场景分析：** 通过 user 表,id 是主键索引，age是非聚集索引
* 1.select * from user where id >50 and id <100;
* 2.select * from user where age > 15 and age < 20';

**分析：**
* 第一条SQL语句根据id进行范围查询，因为(50, 100)范围内的记录在磁盘上按顺序存储，磁盘的顺序读取非常快、

* 第二条SQL语句查询 age （15, 20）范围内的记录，对应的主键id分布可能是离散的10，13，14，26等等；根据主键离散的读取数据，磁盘的离散读比如顺序读快，另外 非聚集索引比聚集索引查询增加了磁盘IO开销。
**小结：**
 * 聚集索引是按照主键顺序，构造B+树，所有数据只存放在叶子节点且有顺，叶子节点之间是一个双向链表，(使用二分查找快速定位到叶子节点),所以根据主键范围查找或者排序查找是非常快速的
 * 非聚集索引：叶子节点并不包含行记录的所有数据，叶子节点处理包括键值对外，还包含一个指向主键索引的bookmark书签(或者说指向数据存储的指针)，即叶子节点索引值和数据值不在一起，增加磁盘查询io，比聚集慢
>reference
* 《mysql技术内幕innodb存储引擎》
* [cnblogs](https://www.cnblogs.com/janehoo/p/7717041.html)

