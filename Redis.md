

<h2 id="引言">引言</h2>
<h3 id="为什么写这篇文章">为什么写这篇文章?</h3>
<p>考虑到绝大部分写业务的程序员，在实际开发中使用redis的时候，只会setvalue和getvalue两个操作，对redis整体缺乏一个认知，所以对redis常见问题做一个总结，希望能够弥补知识盲点。</p>
<h3 id="复习要点">复习要点?</h3>
<p>本文围绕以下几点进行阐述<br>
1、<a href="#1">为什么使用redis</a><br>
2、<a href="#2">使用redis有什么缺点</a><br>
3、<a href="#3">单线程的redis为什么这么快</a><br>
4、<a href="#4">redis的数据类型，以及每种数据类型的使用场景</a><br>
5、<a href="#5">redis的过期策略以及内存淘汰机制</a><br>
6、<a href="#6">redis和数据库双写一致性问题</a><br>
7、<a href="#7">如何应对缓存穿透和缓存雪崩问题</a><br>
8、<a href="#8">如何解决redis的并发竞争问题</a><br>
9、<a href="#9">如何利用redis分布式锁实现控制并发</a><br>
10、<a href="#10">Redis与Memcached比较</a></p>
<h2 id="正文">正文</h2>
<h3 id="span-id11、为什么使用redisspan"><span id="1">1、为什么使用redis</span></h3>
<p><strong>分析</strong>:博主觉得在项目中使用redis，主要是从两个角度去考虑:<strong>性能</strong>和<strong>并发</strong>。当然，redis还具备可以做分布式锁等其他功能，但是如果只是为了分布式锁这些其他功能，完全还有其他中间件(如zookpeer等)代替，并不是非要使用redis。因此，这个问题主要从性能和并发两个角度去答。<br>
<strong>回答</strong>:如下所示，分为两点<br>
<strong>（一）性能</strong><br>
如下图所示，我们在碰到需要执行耗时特别久，且结果不频繁变动的SQL，就特别适合将运行结果放入缓存。这样，后面的请求就去缓存中读取，使得请求能够<strong>迅速响应</strong>。<br>
<img src="https://images.cnblogs.com/cnblogs_com/rjzheng/1202350/o_redis1.png" alt="image"><br>
<strong>题外话：<strong>忽然想聊一下这个</strong>迅速响应</strong>的标准。其实根据交互效果的不同，这个响应时间没有固定标准。不过曾经有人这么告诉我:"在理想状态下，我们的页面跳转需要在<strong>瞬间</strong>解决，对于页内操作则需要在<strong>刹那</strong>间解决。另外，超过<strong>一弹指</strong>的耗时操作要有进度提示，并且可以随时中止或取消，这样才能给用户最好的体验。"<br>
那么<strong>瞬间、刹那、一弹指</strong>具体是多少时间呢？<br>
根据《摩诃僧祗律》记载</p>
<pre><code>一刹那者为一念，二十念为一瞬，二十瞬为一弹指，二十弹指为一罗预，二十罗预为一须臾，一日一夜有三十须臾。
</code></pre>
<p>那么，经过周密的计算，一<strong>瞬间</strong>为0.36 秒,一<strong>刹那</strong>有 0.018 秒.一<strong>弹指</strong>长达 7.2 秒。<br>
<strong>（二）并发</strong><br>
如下图所示，在大并发的情况下，所有的请求直接访问数据库，数据库会出现连接异常。这个时候，就需要使用redis做一个缓冲操作，让请求先访问到redis，而不是直接访问数据库。<br>
<img src="https://images.cnblogs.com/cnblogs_com/rjzheng/1202350/o_redis2.png" alt="image"></p>
<h3 id="span-id22、使用redis有什么缺点span"><span id="2">2、使用redis有什么缺点<span></span></span></h3>
<p><strong>分析</strong>:大家用redis这么久，这个问题是必须要了解的，基本上使用redis都会碰到一些问题，常见的也就几个。<br>
<strong>回答</strong>:主要是四个问题<br>
(一)缓存和数据库双写一致性问题<br>
(二)缓存雪崩问题<br>
(三)缓存击穿问题<br>
(四)缓存的并发竞争问题<br>
这四个问题，我个人是觉得在项目中，比较常遇见的，具体解决方案，后文给出。</p>
<h3 id="span-id33、单线程的redis为什么这么快span"><span id="3">3、单线程的redis为什么这么快<span></span></span></h3>
<p><strong>分析</strong>:这个问题其实是对redis内部机制的一个考察。其实根据博主的面试经验，很多人其实都不知道redis是单线程工作模型。所以，这个问题还是应该要复习一下的。<br>
<strong>回答</strong>:主要是以下三点<br>
(一)纯内存操作<br>
(二)单线程操作，避免了频繁的上下文切换<br>
(三)采用了非阻塞<strong>I/O多路复用机制</strong></p>
<p><strong>题外话：<strong>我们现在要仔细的说一说I/O多路复用机制，因为这个说法实在是太通俗了，通俗到一般人都不懂是什么意思。博主打一个比方：小曲在S城开了一家快递店，负责同城快送服务。小曲因为资金限制，雇佣了</strong>一批</strong>快递员，然后小曲发现资金不够了，只够买<strong>一辆</strong>车送快递。<br>
<strong>经营方式一</strong><br>
客户每送来一份快递，小曲就让一个快递员盯着，然后快递员开车去送快递。慢慢的小曲就发现了这种经营方式存在下述问题</p>
<ul>
<li>几十个快递员基本上时间都花在了抢车上了，大部分快递员都处在闲置状态，谁抢到了车，谁就能去送快递</li>
<li>随着快递的增多，快递员也越来越多，小曲发现快递店里越来越挤，没办法雇佣新的快递员了</li>
<li>快递员之间的协调很花时间</li>
</ul>
<p>综合上述缺点，小曲痛定思痛，提出了下面的经营方式<br>
<strong>经营方式二</strong><br>
小曲只雇佣一个快递员。然后呢，客户送来的快递，小曲按<strong>送达地点</strong>标注好，然后<strong>依次</strong>放在一个地方。最后，那个快递员<strong>依次</strong>的去取快递，一次拿一个，然后开着车去送快递，送好了就回来拿下一个快递。</p>
<p><strong>对比</strong><br>
上述两种经营方式对比，是不是明显觉得第二种，效率更高，更好呢。在上述比喻中:</p>
<ul>
<li>每个快递员------------------&gt;每个线程</li>
<li>每个快递--------------------&gt;每个socket(I/O流)</li>
<li>快递的送达地点--------------&gt;socket的不同状态</li>
<li>客户送快递请求--------------&gt;来自客户端的请求</li>
<li>小曲的经营方式--------------&gt;服务端运行的代码</li>
<li>一辆车----------------------&gt;CPU的核数</li>
</ul>
<p>于是我们有如下结论<br>
1、经营方式一就是传统的并发模型，每个I/O流(快递)都有一个新的线程(快递员)管理。<br>
2、经营方式二就是I/O多路复用。只有单个线程(一个快递员)，通过跟踪每个I/O流的状态(每个快递的送达地点)，来管理多个I/O流。</p>
<p>下面类比到真实的redis线程模型，如图所示<br>
<img src="https://images.cnblogs.com/cnblogs_com/rjzheng/1202350/o_redis3.png" alt="image"><br>
参照上图，简单来说，就是。我们的redis-client在操作的时候，会产生具有不同事件类型的socket。在服务端，有一段I/0多路复用程序，将其置入队列之中。然后，文件事件分派器，依次去队列中取，转发到不同的事件处理器中。<br>
需要说明的是，这个I/O多路复用机制，redis还提供了select、epoll、evport、kqueue等多路复用函数库，大家可以自行去了解。</p>
<h3 id="span-id44、redis的数据类型，以及每种数据类型的使用场景span"><span id="4">4、redis的数据类型，以及每种数据类型的使用场景<span></span></span></h3>
<p><strong>分析</strong>：建议，在项目中用到后，再类比记忆，体会更深，不要硬记。基本上，一个合格的程序员，五种类型都会用到。<br>
<strong>回答</strong>：一共五种</p>
<h4 id="一string">(一)String</h4>
<p>最常规的set/get操作，value可以是String也可以是数字。<br>
<strong>场景</strong>：</p>
<ul>
<li><strong>缓存功能</strong>：(序列化用户信息、详情页面信息）</li>
<li>结合 <strong>incr命令</strong>做 <a href="http://www.redis.cn/commands/incr.html">计数、限流api功能</a></li>
</ul>
<h4 id="二hash">(二)hash</h4>
<p>特别适合用于存储对象，value存放结构化的对象，方便的就是操作其中的某个字段。</p>
<p><strong>场景</strong>：</p>
<ul>
<li>
<p>做<strong>用户登录</strong>的时候，就是用这种数据结构存储用户信息，以cookieId作为key，设置30分钟为缓存过期时间，能很好的模拟出类似session的效果。</p>
</li>
<li>
<p>每条微博都有点赞数、评论数、转发数和浏览数四条属性，这时用<code>hash</code>进行计数会更好，将该计数器的 key 设为<code>weibo:weibo_id</code>，<code>hash</code>的 field 为<code>like_number</code>、<code>comment_number</code>、<code>forward_number</code>和<code>view_number</code>，在对应操作后通过<strong>hincrby</strong>使<code>hash 中</code>的 field 自增</p>
</li>
</ul>
<h4 id="三list">(三)list</h4>
<p>简单的字符串列表，按照插入顺序排序</p>
<p><strong>场景</strong>：</p>
<ul>
<li>
<p><strong>做简单的消息队列的功能</strong>。</p>
</li>
<li>
<p><strong>做基于redis的分页功能</strong>，性能极佳，用户体验好。</p>
</li>
<li>
<p><code>list</code>作为双向链表，不光可以作为队列使用。如果将它用作栈便可以成为一个公用的时间轴。当用户发完微博后，都通过<code>lpush</code>将它存放在一个 key 为<code>LATEST_WEIBO</code>的<code>list</code>中，之后便可以通过<code>lrange</code>取出当前最新的微博。</p>
</li>
</ul>
<h4 id="四set">(四)set</h4>
<p>String 类型的无序集合。集合成员是唯一的，这就意味着集合中不能出现重复的数据。<br>
Redis 中集合是通过哈希表实现的，所以添加，删除，查找的复杂度都是 O(1)。</p>
<p><strong>场景</strong>：</p>
<ul>
<li>
<p>因为set堆放的是一堆不重复值的集合。所以可以做<strong>全局去重的功能</strong>。</p>
</li>
<li>
<p><strong>好友关系</strong> 使用交集、并集、差集等操作，可以<strong>计算共同喜好，全部的喜好，自己独有的喜好等功能</strong>。</p>
</li>
<li>
<p><strong>倒排索引</strong> 是构造搜索功能的最常见方式，在 Redis 中也可以通过<code>set</code>进行建立倒排索引，这里以简单的拼音 + 前缀搜索城市功能举例：</p>
</li>
</ul>
<p>假设一个城市<code>北京</code>，通过拼音词库将<code>北京</code>转为<code>beijing</code>，再通过前缀分词将这两个词分为若干个前缀索引，有：<code>北</code>、<code>北京</code>、<code>b</code>、<code>be</code>…<code>beijin</code>和<code>beijing</code>。将这些索引分别作为<code>set</code>的 key（例如:<code>index:北</code>）并存储<code>北京</code>的 id，倒排索引便建立好了。接下来只需要在搜索时通过关键词取出对应的<code>set</code>并得到其中的 id 即可。</p>
<h4 id="五sorted-set">(五)sorted set</h4>
<p>有序集合和集合一样也是string类型元素的集合,且不允许重复的成员。</p>
<p><strong>场景</strong>：</p>
<ul>
<li>
<p><strong>排行榜</strong>  如果应用有一个发帖排行榜的功能，便选择<code>sorted set</code>吧，将集合的 key 设为<code>POST_RANK</code>。当用户发帖后，使用<code>zincrby</code>将该用户 id 的 score 增长 1。<code>sorted set</code>会重新进行排序，用户所在排行榜的位置也就会得到实时的更新。</p>
</li>
<li>
<p>参照另一篇<a href="https://www.cnblogs.com/rjzheng/p/8972725.html">《分布式之延时任务方案解析》</a>，该文指出了sorted set可以用来做<strong>延时任务</strong>。</p>
</li>
</ul>
<h3 id="span-id55、redis的过期策略以及内存淘汰机制span"><span id="5">5、redis的过期策略以及内存淘汰机制</span></h3>
<p><strong>分析</strong>:这个问题其实相当重要，到底redis有没用到家，这个问题就可以看出来。比如你redis只能存5G数据，可是你写了10G，那会删5G的数据。怎么删的，这个问题思考过么？还有，你的数据已经设置了过期时间，但是时间到了，内存占用率还是比较高，有思考过原因么?<br>
<strong>回答</strong>:<br>
redis采用的是定期删除+惰性删除策略。<br>
<strong>为什么不用定时删除策略?</strong><br>
定时删除,用一个定时器来负责监视key,过期则自动删除。虽然内存及时释放，但是十分消耗CPU资源。在大并发请求下，CPU要将时间应用在处理请求，而不是删除key,因此没有采用这一策略.<br>
<strong>定期删除+惰性删除是如何工作的呢?</strong><br>
定期删除，redis默认每个100ms检查，是否有过期的key,有过期key则删除。需要说明的是，redis不是每个100ms将所有的key检查一次，而是随机抽取进行检查(如果每隔100ms,全部key进行检查，redis岂不是卡死)。因此，如果只采用定期删除策略，会导致很多key到时间没有删除。<br>
于是，惰性删除派上用场。也就是说在你获取某个key的时候，redis会检查一下，这个key如果设置了过期时间那么是否过期了？如果过期了此时就会删除。<br>
<strong>采用定期删除+惰性删除就没其他问题了么?</strong><br>
不是的，如果定期删除没删除key。然后你也没即时去请求key，也就是说惰性删除也没生效。这样，redis的内存会越来越高。那么就应该采用<strong>内存淘汰机制</strong>。<br>
在redis.conf中有一行配置</p>
<pre><code># maxmemory-policy volatile-lru
</code></pre>
<p>该配置就是配内存淘汰策略的(什么，你没配过？好好反省一下自己)<br>
1）noeviction：当内存不足以容纳新写入数据时，新写入操作会报错。<strong>应该没人用吧。</strong><br>
2）allkeys-lru：当内存不足以容纳新写入数据时，在键空间中，移除最近最少使用的key。<strong>推荐使用，目前项目在用这种。</strong><br>
3）allkeys-random：当内存不足以容纳新写入数据时，在键空间中，随机移除某个key。<strong>应该也没人用吧，你不删最少使用Key,去随机删。</strong><br>
4）volatile-lru：当内存不足以容纳新写入数据时，在设置了过期时间的键空间中，移除最近最少使用的key。<strong>这种情况一般是把redis既当缓存，又做持久化存储的时候才用。不推荐</strong><br>
5）volatile-random：当内存不足以容纳新写入数据时，在设置了过期时间的键空间中，随机移除某个key。<strong>依然不推荐</strong><br>
6）volatile-ttl：当内存不足以容纳新写入数据时，在设置了过期时间的键空间中，有更早过期时间的key优先移除。<strong>不推荐</strong><br>
ps：如果没有设置 expire 的key, 不满足先决条件(prerequisites); 那么 volatile-lru, volatile-random 和 volatile-ttl 策略的行为, 和 noeviction(不删除) 基本上一致。</p>
<h3 id="span-id66、redis和数据库双写一致性问题span"><span id="6">6、redis和数据库双写一致性问题</span></h3>
<p><strong>分析</strong>:一致性问题是分布式常见问题，还可以再分为最终一致性和强一致性。数据库和缓存双写，就必然会存在不一致的问题。答这个问题，先明白一个前提。就是<strong>如果对数据有强一致性要求，不能放缓存。<strong>我们所做的一切，只能保证最终一致性。另外，我们所做的方案其实从根本上来说，只能说</strong>降低不一致发生的概率</strong>，无法完全避免。因此，有强一致性要求的数据，不能放缓存。<br>
<strong>回答</strong>:</p>
<ul>
<li>
<p>第一种采取正确更新策略，先更新数据库，再删缓存。其次，因为可能存在删除缓存失败的问题，提供一个补偿措施即可，例如利用消息队列。</p>
</li>
<li>
<p>使用 mysql中 binlog获取数据异步操作缓存</p>
</li>
</ul>
<h3 id="span-id77、如何应对缓存穿透和缓存雪崩问题span"><span id="7">7、如何应对缓存穿透和缓存雪崩问题</span></h3>
<p><strong>分析</strong>:这两个问题，说句实在话，一般中小型传统软件企业，很难碰到这个问题。如果有大并发的项目，流量有几百万左右。这两个问题一定要深刻考虑。<br>
<strong>回答</strong>:如下所示<br>
<strong>缓存穿透</strong>，即黑客故意去请求缓存中不存在的数据，导致所有的请求都怼到数据库上，从而数据库连接异常。<br>
<strong>解决方案</strong>:<br>
(一)利用互斥锁，缓存失效的时候，先去获得锁，得到锁了，再去请求数据库。没得到锁，则休眠一段时间重试<br>
(二)采用异步更新策略，无论key是否取到值，都直接返回。value值中维护一个缓存失效时间，缓存如果过期，异步起一个线程去读数据库，更新缓存。需要做<strong>缓存预热</strong>(项目启动前，先加载缓存)操作。<br>
(三)提供一个能迅速判断请求是否有效的拦截机制，比如，利用布隆过滤器，内部维护一系列合法有效的key。迅速判断出，请求所携带的Key是否合法有效。如果不合法，则直接返回。<br>
<strong>缓存雪崩</strong>，即缓存同一时间大面积的失效，这个时候又来了一波请求，结果请求都怼到数据库上，从而导致数据库连接异常。<br>
<strong>解决方案</strong>:<br>
(一)给缓存的失效时间，加上一个随机值，避免集体失效。<br>
(二)使用互斥锁，但是该方案吞吐量明显下降了。<br>
(三)双缓存。我们有两个缓存，缓存A和缓存B。缓存A的失效时间为20分钟，缓存B不设失效时间。自己做缓存预热操作。然后细分以下几个小点</p>
<ul>
<li>I 从缓存A读数据库，有则直接返回</li>
<li>II A没有数据，直接从B读数据，直接返回，并且异步启动一个更新线程。</li>
<li>III 更新线程同时更新缓存A和缓存B。</li>
</ul>
<h3 id="span-id88、如何解决redis的并发竞争key问题span"><span id="8">8、如何解决redis的并发竞争key问题</span></h3>
<p><strong>分析</strong>:这个问题大致就是，同时有多个子系统去set一个key。这个时候要注意什么呢？大家思考过么。需要说明一下，博主提前百度了一下，发现答案基本都是推荐用redis事务机制。博主**不推荐使用redis的事务机制。**因为我们的生产环境，基本都是redis集群环境，做了数据分片操作。你一个事务中有涉及到多个key操作的时候，这多个key不一定都存储在同一个redis-server上。因此，<strong>redis的事务机制，十分鸡肋。</strong><br>
**回答:**如下所示<br>
(1)如果对这个key操作，<strong>不要求顺序</strong><br>
这种情况下，准备一个分布式锁，大家去抢锁，抢到锁就做set操作即可，比较简单。<br>
(2)如果对这个key操作，<strong>要求顺序</strong><br>
假设有一个key1,系统A需要将key1设置为valueA,系统B需要将key1设置为valueB,系统C需要将key1设置为valueC.<br>
期望按照key1的value值按照 valueA–&gt;valueB–&gt;valueC的顺序变化。这种时候我们在数据写入数据库的时候，需要保存一个时间戳。假设时间戳如下</p>
<pre><code>系统A key 1 {valueA  3:00}
系统B key 1 {valueB  3:05}
系统C key 1 {valueC  3:10}
</code></pre>
<p>那么，假设这会系统B先抢到锁，将key1设置为{valueB 3:05}。接下来系统A抢到锁，发现自己的valueA的时间戳早于缓存中的时间戳，那就不做set操作了。以此类推。</p>
<p>其他方法，比如利用队列，将set方法变成串行访问也可以。总之，灵活变通。</p>
<h3 id="span-id99、如何利用redis分布式锁实现控制并发span"><span id="9">9、如何利用redis分布式锁实现控制并发</span></h3>
<h4 id="redis命令解释">redis命令解释</h4>
<p>说道Redis的分布式锁都是通过setNx命令结合getset来实现的，在讲之前我们先了解下setNx和getset的意思，在redis官网是这样解释的<br>
注：redis的命令都是原子操作</p>
<h4 id="setnx-key-value">SETNX key value</h4>
<p>将 key 的值设为 value ，当且仅当 key 不存在。<br>
若给定的 key 已经存在，则 SETNX 不做任何动作。<br>
SETNX 是『SET if Not eXists』(如果不存在，则 SET)的简写。<br>
<strong>可用版本：</strong><br>
1.0.0+<br>
<strong>时间复杂度：</strong><br>
O(1)<br>
<strong>返回值：</strong><br>
设置成功，返回 1 。<br>
设置失败，返回 0 。</p>
<pre><code>redis&gt; EXISTS job                # job 不存在
(integer) 0
redis&gt; SETNX job "programmer"    # job 设置成功
(integer) 1
redis&gt; SETNX job "code-farmer"   # 尝试覆盖 job ，失败
(integer) 0
redis&gt; GET job                   # 没有被覆盖
"programmer"
</code></pre>
<h4 id="getset-key-value">GETSET key value</h4>
<p>将给定 key 的值设为 value ，并返回 key 的旧值(old value)。<br>
当 key 存在但不是字符串类型时，返回一个错误。<br>
<strong>可用版本：</strong><br>
1.0.0+<br>
<strong>时间复杂度：</strong><br>
O(1)<br>
<strong>返回值：</strong><br>
返回给定 key 的旧值。<br>
当 key 没有旧值时，也即是， key 不存在时，返回 nil 。</p>
<pre><code>redis&gt; GETSET db mongodb    # 没有旧值，返回 nil
(nil)
redis&gt; GET db
"mongodb"
redis&gt; GETSET db redis      # 返回旧值 mongodb
"mongodb"
redis&gt; GET db
"redis"
</code></pre>
<h4 id="思路">思路</h4>
<p>为了让分布式锁的算法更稳键些，持有锁的客户端在解锁之前应该再检查一次自己的锁是否已经超时，再去做DEL操作，因为可能客户端因为某个耗时的操作而挂起，操作完的时候锁因为超时已经被别人获得，这时就不必解锁了。</p>
<p>###<span id="10">Redis与Memcached比较</span></p>
<p>两者都是非关系型内存键值数据库。有以下主要不同：</p>
<blockquote>
<p>数据类型</p>
</blockquote>
<p>Memcached 仅支持字符串类型，而 Redis 支持五种不同种类的数据类型，使得它可以更灵活地解决问题。</p>
<blockquote>
<p>数据持久化</p>
</blockquote>
<p>Redis 支持两种持久化策略：RDB 快照和 AOF 日志，而 Memcached 不支持持久化。</p>
<blockquote>
<p>分布式</p>
</blockquote>
<p>Memcached 不支持分布式，只能通过在客户端使用一致性哈希这样的分布式算法来实现分布式存储，这种方式在存储和查询时都需要先在客户端计算一次数据所在的节点。</p>
<p>Redis Cluster 实现了分布式的支持。</p>
<blockquote>
<p>内存管理机制</p>
</blockquote>
<p>在 Redis 中，并不是所有数据都一直存储在内存中，可以将一些很久没用的 value 交换到磁盘。而 Memcached 的数据则会一直在内存中。</p>
<p>Memcached 将内存分割成特定长度的块来存储数据，以完全解决内存碎片的问题，但是这种方式会使得内存的利用率不高，例如块的大小为 128 bytes，只存储 100 bytes 的数据，那么剩下的 28 bytes 就浪费掉了。</p>
<h2 id="总结">总结</h2>
<p>本文对redis的常见问题做了一个总结。大部分是博主自己在工作中遇到，以及以前面试别人的时候，爱问的一些问题。另外，<strong>不推荐大家临时抱佛脚</strong>，真正碰到一些有经验的工程师，其实几下就能把你问懵。最后，希望大家有所收获吧。</p>
<blockquote>
<p>reference：</p>
</blockquote>
<p><a href="https://www.cnblogs.com/rjzheng/p/9096228.html">https://www.cnblogs.com/rjzheng/p/9096228.html</a><br>
<a href="https://blog.csdn.net/fuyuwei2015/article/details/72870131">https://blog.csdn.net/fuyuwei2015/article/details/72870131</a><br>
<a href="http://www.scienjus.com/redis-use-case/">http://www.scienjus.com/redis-use-case/</a><br>
<a href="https://github.com/CyC2018/Interview-Notebook/blob/master/notes/Redis.md#%E5%9B%9Bredis-%E4%B8%8E-memcached">https://github.com/CyC2018/Interview-Notebook/blob/master/notes/Redis.md#四redis-与-memcached</a></p>

