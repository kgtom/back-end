## 一、redis数据结构

### 1.简单动态字符串SDS(simple dynamic string)
* 为了实现对字符串的高效操作，redis自己构建SDS抽象数据结构：
~~~
struct sdshdr{
//记录buf数组中已使用的字节数量，等于sds保存字符串的长度
int len;
//记录buf数组中未使用的字节数量
int free;
//字节数组，用于保存字符串
char buf[]
}
~~~

例如: set name =tom 
~~~
10.1.183.247:6379[29]> set name tom
OK
10.1.183.247:6379[29]> type name
string
10.1.183.247:6379[29]> object encoding name
"embstr"
10.1.183.247:6379[29]> strlen name
(integer) 3
~~~
len：3，free:0 buf:'t''o''m''\0',最后多一个空字符，为了使用重用C字符串函数库的函数

* 与C字符串优势：
 - 获取字符串长度O(1),因为可以直接获取len,C的话需要遍历O(N)
 - 避免缓冲区溢出，对SDS修改时，底层api会先check一下是否满足空间，不满足的话先自动分配再执行操作
 - 减少修改字符串时带来的内存分配次数，对于C字符串修改，必然会造成内存重新分配，而SDS通过空间预分配和惰性空间释放优化策略
 - 二进制安全，C字符串只能保存文本，SDS可以保存文本，也可以保存二进制数据
 - 兼容部分C字符串函数库,<string.H>

* 使用场景：redis五种对象都会使用的数据结构；AOF模块的缓冲区
### 2.linkedlist链表
* redis的链表通过链表节点带有prev和next组成了一个双端链表，获取某个节点的前置或者后置节点复杂度O(1),通过链表中len获取链表数量O(1)，只有在遍历查找遍历O(N)
* 使用场景：list对象底层实现，包括消息队列、发布与订阅

### 字典
* 字典又称为map，用于保存键值对的抽象结构体，redis数据库底层实现就是字典。

~~~
typedef struct dict{
//类型特定函数，包括计算hash值、复制、销毁等
dictType *type;

//私有数据，包括特定函数的可选参数
void *privdata;

//哈希表，一个平时用，另一个用于rehash时使用
dictht ht[2]

//当rehash不进行时，此值为1
int trehashidx;

}
~~~

* 哈希算法：通过哈希算法计算哈希值，通过哈希值或者sizemask计算出索引值idx，根据索引值放在哈希表数组为idx的哈希节点上。
* 解决哈希冲突，使用链地址法。将冲突的放在哈希节点的next指针，组成单向链表。
* 使用场景：hash对象、set 对象、redis数据库对象
### ziplist
* 压缩链表：redis节约内存而开发的顺序型数据结构,每个节点可以保存一个字节数组或者整数值。新增、获取单个值O(1),获取所有、删除需要遍历O(N)。
* 使用场景：list对象、hash对象、zset对象
### skiplist
* 跳跃表是一种有序数据结构，它通过每个节点维持着指向其它节点的指针达到快速访问的目的。查找O(logN)，最坏O(N)。由zskipList跳跃表信息(表头、表尾、长度) 和 zskiplistNode跳跃表节点组成。跳跃表节点按照score大小排序，如果score相同，则按照member排序
* zskiplistNode跳跃表节点包括：

~~~
typedef stuct zskiplistNode{
//后退指针：用于从表尾向表头访问节点
struct zskiplistNode *backward;


//doublel浮点型分值，用于排序
double score;
//成员对象：sds字符串对象
robj *obj;

//层高：数组包含多少个元素
struct zskiplistLevel{
//前进指针：用于从表头向表尾访问节点
struct zskiplistNode *forward;
//跨度：记录两个节点之间距离
unsigned int span;
}[]level;
}
~~~
* 使用场景：zset对象

### inset
* 整形结合是集合键set底层实现之一(当一个set只包含整数元素且数量不多时)。
* inset 底层整形数组，新增、查询某个值都是O(1)，删除或遍历O(N),检查某个值是否在集合O(logN)
## 二、redis五种对象

### string字符串对象
int 、raw、embstr,后两者实际 动态字符串(SDS)

### list对象
   * ziplist：字符串长度小于64字节且数量不超过512
   * linkedlist：双向链表：不满足上述条件的
### hash对象
   * ziplist：键和值长度小于64字节且键值对数量小于512
   * dict：实际hash，不满足上述条件的
### set对象
  * inset：整形整合，底层数组组成，保存元素是整数且元素数量不超过512
  * dict：实际hash,新增、删除、修改都是O(1)，不满足上述条件的

### zset对象
   * ziplist：压缩链表，两个紧挨着的节点，一个节点存储member，另一个节点存储score
   * skiplist：跳跃表,实际上在链表的基础上改造生成的，将链表数据进行升级，提升若干索引层，加上一层索引后，查找一个结点需要遍历的结点个数减少了，每层遍历最多3个结点即可，而跳表的高度为 h ，所以每次查找一个结点时，需要遍历的结点数为 3*跳表高度 ，所以忽略低阶项和系数后的时间复杂度就是O(logN) 也就是说查找效率提高了。


### redis 为什么使用skiplist 不使用avl平衡树或者hash呢？
 * skiplist 和avl 适合范围查找O(logN)，单个key查找使用hash O(1)
 * 范围查找：avl先找起点，然后中序再找，skiplist找到起点后，遍历链表即可
 * 实现复杂度：avl的插入和删除节点会引发树的调整，逻辑复杂，除了数据，还多了左右指针，以及叶子节点指针，skiplist只需要修改相邻指针就可以
 * 内存：avl 两个节点树 ，skiplist 修改节点指针，但是每个节点的指针小于<2,所以比B+树占用空间小
 * skiplist 底层既有dict(查询zset的score使用 o(1))，也有两个有序链表，方便根据范围查找查找，如 zrange zrevrange 等

## redis 为什么使用skiplist做索引而不使用B+树呢？

* 因为B+树的原理是：叶子节点存储数据，非叶子节点存储索引，B+树的每个节点可以存储多个关键字，它将节点大小设置为磁盘页的大小，充分利用了磁盘预读的功能。每次读取磁盘页时就会读取一整个节点,每个叶子节点还有指向前后节点的指针，为的是最大限度的降低磁盘的IO。B+树纯粹是为了mysql这种IO数据库准备的。B+树的每个节点的数量都是一个mysql分区页的大小
* redis 的数据在内存中读取耗费的时间是从磁盘的IO读取的百万分之一，从内存中读取数据而不涉及磁盘IO,索引使用skiplist。


### ziplist设计
* 特殊编码的，各个数据项在一个连续的内存空间，设计目标就是为了提高存储效率，O(1)复杂度push和pop，比如：zset中zscore,list中pop,hash中hget
* 缺点：当数据变动时，内存会重新分配；当数据多，ziplist转成dict 、linkedlist、skiplist，查找需要遍历,比如 list中lrange 


> reference
* 《redis设计与实现》
