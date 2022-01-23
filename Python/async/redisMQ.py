# -*- encoding: utf-8 -*-
import redis
from redis import Redis

"""
exists
get mget keys randomkey
set setnx mset msetnx getset append del rename
"""

"""
System
    info
    dbsize
    monitor
    flushdb     删除当前选择数据库的数据
    flushall    删除所有数据库内的数据
    config (get set resetstat rewrite) // requirepass xxx
    slowlog (get 10)

    save        阻塞的RDB备份
    bgsave      非阻塞的RDB备份,fork进程后开始备份,fork过程中是阻塞的
    slaveof     开启主从复制


Key
    exists
    keys (*)
    randomkey
    del (a b c)
    expire
    move
    ttl         // 判断一个Key是否失效 0:过期 -1:永不过期 -2:失效


Type
    type


Transaction
    multi
    exec
    discard
    watch


String:
    存储原理: https://www.cnblogs.com/mingbai/p/11469294.htm 和 go语言 的切片特别相似
    get         // 不存在,过期,失效, 返回 nil
    set
    append      // 不存在,则创建
    strlen
    INCR        // 原子操作
    DECR
    INCRBY      // 原子操作
    DECRBY
    getset      // 原子操作,返回并赋值,不存在则创建
    setEX       // 设置一个Key并指定过期时间
    setNX       // 原子操作,不存在则set,否则不做操作
    MSET        // 原子操作 批量设置
    MGET        // 原子操作 批量获取,不存在返回 nil
    MSETNX      // 原子操作 所有Key均不存在,则批量设置,否则不做操作,返回0


List
    元素类型为 String, 左侧为头,右侧为尾,顶层实现就是一个链表

    SORT        // 排序

    LPUSH       // 从头插入一批数据,如果不存在则创建
    LPUSHX      // 从头插入一批数据,如果不存在则返回失败 0
    LRANGE      // 头的索引为0, 尾的索引可以视为-1, 首尾均包含在内. 和 python 的 list 特别相似
    LPOP        // 从头弹出一个值,是一个栈的操作
    LLEN        // 返回列表长度
    LREM        // 从头开始remove 指定个数个指定值的元素,并返回实际移除元素个数, 个数为-1 为全部移除
    LSET        // 给指定索引处的元素重新设值,正值索引为从头开始,负值索引为从尾开始
    LINDEX      // 指定元素索引获取值
    LTRIM       // 段截取 保留从开始索引到索引结束的元素,首尾均包含在内
    LINSERT     // 在从头到尾的第一个指定的元素之前或者之后插入值


    RPUSH       // 从尾部插入一批数据,从头到位的顺序为传参顺序, 如果不存在则创建
    RPUSHX      // 从尾部插入一批数据,如果不存在则返回失败 0, 返回所有元素数量
    RPOP        // 从尾部弹出,并返回该元素
    RPOPLPUSH   // 从l1尾部弹出,插入l2头部,可以是同一个列表

    BRPOP       // 从多个列表中阻塞读,任一信道有消息即可唤醒,消息到达立刻返回 (是如何保持连接的)
                   在消息队列中客户端不用空轮询(Polling)
    BLPOP       // 阻塞过久服务端会断开连接,客户端此时需要重试


Hash
    可以用来存储对象,一般业务中会以类名+ID来命名
    HSET       // 给Hash批量设置键值对,hash不存在则创建, 返回实际新增键的个数, 实际时间复杂度为 O(1)
    HGET       // 获取单个键的值
    HDEL       // 批量删除hash的键,返回实际删除的个数
    HEXISTS    // 判断单个键是否存在
    HLEN       // 返回hash的个数
    HSETNX     // 新增不存在的指定键值,如果存在,则返回0
    HINCRBY    // 原子操作, 指定键值 加 指定值,返回结果
    HGETALL    // 返回所有key和val,按照顺序 1key1val 排列
    HKEYS      // 返回所有key
    HVALS      // 返回所有value
    HMGET      // 批量get, 不存在返回 nil
    HMSET      // 批量设置, 返回 ok. 实际时间复杂度为 O(n)  PS: 注意 和 HSET 的区别
                  根据Redis 4.0.0，HMSET被视为已弃用。请在新代码中使用HSET。


Set
    元素类型为string类型，元素具有唯一性， 不允许存在重复的成员。多个集合类型之间可以进行并集、交集和差集运算
    场景1:判断某个Key是否存在

    SADD            // 批量向集合中加入元素,返回实际插入元素个数
    SMEMBERS        // 返回集合所有元素,返回元素顺序和插入顺序不一致
    SCARD           // 返回集合的元素个数
    SISMEMBER       // 判断元素是否存在
    SPOP            // 随机的移除并返回Set中的某一成员
    SREM            // 移除指定元素,并返回实际移除元素个数
    SRANDMEMBER     // 随机返回一位成员
    SMOVE           // 将set1中的一个元素移动到另一个set2中, s2不存在则会创建,返回 0 或者 1


Zset
    内部实现是 "跳跃表"
    每个元素都会关联–个double类型的分数score(表示权重)，可以通过权重的大小排序，元素的score可以相同
    大型在线游戏的积分排行榜, ZADD ZRANGE ZRANK
    构建索引数据

    ZADD                    // 批量向集合中加入绑定分数的元素,返回实际插入元素个数,参数 INCR ==> ZINCRBY
    ZCARD                   // 返回元素个数
    ZCOUNT                  // 返回指定分数范围内的元素个数
    ZREM                    // 移除指定元素,并返回实际移除元素个数
    ZSCORE                  // 返回指定元素的分数,返回类型是 String
    ZINCRBY                 // 增加指定元素指定的分数,如果不存在则创建,并将增加分数设置为其分数,返回最终分数
    ZRANDMEMBER             // 随机返回指定个数个元素

    ZRANK                   // 返回指定元素在 zs 中按分数从小到大排序后所在的位置,起始位置为0
    ZREVRANK                // 以分数从高到低的方式获取并返回此区间内的成员。

    ZRANGE                  // 返回指定范围的元素, 可以增加参数 withscores 一并返回分数
                               默认按照排序后的等级(索引)排序,可指定参数按照score,lex返回,返回顺序默认都是升序的
                               不同的返回方式指定 min max 的参数特点不一致.
                               很多命令可以通过该命令参数实现

                               zrange zs 0 -1
                               ZRANGE zs (3 +inf  byscore withscores
                               zrange zs - + bylex     // 先按照分数排序,在按照字典序

    ZREMRANGEBYRANK         // key start stop 返回实际移除的元素个数 -inf 表示最小的rank +inf 表示最大的rank
    ZREMRANGEBYSCORE        // key min max 返回实际移除的元素个数 -inf 表示最小的score +inf 表示最大的score

    BZPOPMIN


TODO  了解各种类型底层实现原理

NOTE  API 有很多,记是记不住的,根据实际应用场景去记忆 "那么几个" 就行,其他的根据实际需求去查阅即可.
"""


class MQ(object):
    """
    1- LPUSH BRPOP
        无ACK(需要额外维护)
        不可重复消费
        不能广播(支持P/S模型)
        不支持分组消费

    2- PUBLISH SUBSCRIBE UNSUBSCRIBE
        生产者往信道内发送消息, 中间件将消息复制到多个消息队列,每个消息队列由对应的消费组消费

        广播
        多信道订阅
        消息即时发送,消息不用等待消费者读取,消费者会自动接收到信道发布的消息

        客户端下线,重新订阅不能接收历史消息,而且意味着必须先开始消费,后生产才不会丢失消息
        不能保证消费者接收消息的时间是一致的
        在消息的生产远大于消费速度时,若消费者客户端出现消息积压，到一定程度，会被强制断开，导致消息意外丢失。
        Pub/Sub 模式不适合做消息存储，消息积压类的业务，而是擅长处理广播，即时通讯，即时反馈的业务。

    3- SORTED SET
        zADD ZPOPMAX
        生产者维护消息顺序,生产者给定一个排序ID确定消息顺序,
        当然ID作为score需要是单调递增的,通常可以使用时间戳+序号的方案

        就是可以自定义消息ID，在消息ID有意义时，比较重要。

        不允许重复消息（因为是集合），同时消息ID确定有错误会导致消息的顺序出错。

    4- XADD

        5.0之后新增的数据结构, 支持多播的可持久化消息队列，实现借鉴了Kafka设计。

        持久化的
        消息链表(Stream)
        消息唯一ID
        消费组(名称唯一)
        消费者(组内名称唯一)
        游标(last_delivered_id)
        ACK 消息列表(Pending Entries List),消费者(客户端)已经被读取到的消息,但是还没有ACK
            用来确保客户端至少消费了消息一次，而不会在网络传输的中途丢失了没处理

        每个消费组(Consumer Group)的状态都是独立的, 同一份Stream内部的消息会被每个消费组都消费到。
        同一个消费组(Consumer Group)可以挂接多个消费者(Consumer)，这些消费者之间是竞争关系

        XADD        ID 可以自己指定 但是必须是从 0-0开始
        XDEL        逻辑删除
        XRANGE      获取指定范围内的消息(会过滤已经删除的消息)
        XLEN
        XGROUP      指定起始消息ID(初始化last_delivered_id, 指定消费组起始消费的消息),创建消费组
        XTRIM       物理删除

        独立消费
            使用xread脱离消费组(不用创建消费组)直接消费,(可视为只有一个成员的消费组??),
            可视为将Stream当成普通的消息队列来使用,在没有消息时,也可以阻塞等待.

        消费组消费
            XREADGROUP GROUP group consumer [COUNT count] [BLOCK milliseconds] [NOACK] STREAMS key [key ...] ID [ID ...]
            XACK key group ID [ID ...]

        消息积累太多
            消息逻辑删除,消息会积累,使用XADD 指定stream最大长度maxlen，就将老的消息干掉，确保最多不超过指定长度

        忘记ACK
            导致 Pending entry list 越来越大,占用内存越来越多

        PEL避免消息丢失
            在客户端消费者读取Stream消息时，Redis服务器将消息回复给客户端的过程中，客户端突然断开了连接，消息就丢失了。
            但是PEL里已经保存了发出去的消息ID。待客户端重新连上之后，可以再次收到PEL中的消息ID列表。
            不过此时xreadgroup的起始消息必须是任意有效的消息ID，一般将参数设为0-0，
            表示读取所有的PEL消息以及自last_delivered_id之后的新消息。

        分区Partition
            Stream的消费模型借鉴了kafka的消费分组的概念，它弥补了Redis Pub/Sub不能持久化消息的缺陷。
            但是它又不同于kafka，kafka的消息可以分partition，而Stream不行。
            如果非要分parition的话，得在客户端做，提供不同的Stream名称，对消息进行hash取模来选择往哪个Stream里塞.


    0-0

    XREADGROUP GROUP gp cs STREAMS st > id [id...]  // id 可选,如果给定id则是从PEL中读取
    XPENDING
    XACK


    Redis做消息队列的缺点:
        Redis 本身可能会丢数据
        面对消息积压，Redis 内存资源紧张
        如果PEL中消息没有被ACK,那么消费者一直从消息队列里消费会导致PEL会变大,
        而且消费者在后续也不能确定该消息是否被成功消费,

        如果要再次消费,且消费PEL中的消息,则需要指定ID ,此时 group的last_deliver_id是不变的.
    """


class DL(object):
    """
    锁要具有以下两个身份信息: 1-某事中某个资源的锁 2-被某人所持有

    原子性要求:
        加锁: set run_job_<id> user_token EX 10 NX
        解锁: GETDEL[值必须是一个字符串].  或者使用 Lua脚本先get后del

        锁加唯一标识的原因:
            当资源操作过程中锁失效,那么去释放锁的话,会把其它拿到锁的进程的锁给释放掉

    过期时间不好评估:
        加锁时，先设置一个过期时间，然后我们开启一个「守护线程」，定时去检测这个锁的失效时间，
        如果锁快要过期了，操作共享资源还未完成，那么就自动对锁进行「续期」，重新设置过期时间。
        这确实一种比较好的方案:
            如果你是 Java 技术栈，幸运的是，已经有一个库把这些工作都封装好了：Redisson。
        采用了「自动续期」的方案来避免锁过期，这个守护线程我们一般也把它叫做「看门狗」线程.

        但是当STW 时,续期会失效,导致锁被提前释放,这个时候只能由底层资源来兜低, 所谓无绝对安全
        但是,底层资源既然是互持,那要分布式锁搞毛 ????

        状态机(如下单),最后还是要用CP模型(CP模型之间的比较)来卡
        数据库兜底(共享资源互持) TODO 增强容错性????,既然兜底那redis岂不是冗余了?
        原子关播(ZK的ZAB, etcd的raft)

    Redisson
        可重入锁
        乐观锁
        公平锁
        读写锁
        Redlock（红锁，下面会详细讲）

    [主从集群+哨兵] 下故障转移导致锁的问题
        主库异常宕机时，哨兵可以实现「故障自动切换」，把从库提升为主库，继续提供服务，以此保证可用性。
        主从复制是异步的, 当主SET 未同步到从库上时,就完成切换,会导致锁丢失

        RedLock
            1-客户端先获取「当前时间戳T1」
            2-客户端依次向这 5 个 Redis 实例发起加锁请求（用前面讲到的 SET 命令），且每个请求会设置超时时间
             （毫秒级，要远小于锁的有效时间），如果某一个实例加锁失败（包括网络超时、锁被其它人持有等各种异常情况），
             就立即向下一个 Redis 实例申请加锁
            3- 如果客户端从 3 个（大多数）以上 Redis 实例加锁成功，则再次获取「当前时间戳T2」，
               如果 T2 - T1 < 锁的过期时间，此时，认为客户端加锁成功，否则认为加锁失败
            4- 加锁成功，去操作共享资源（例如修改 MySQL 某一行，或发起一个 API 请求）
            5- 加锁失败，向「全部节点」发起释放锁请求（前面讲到的 Lua 脚本释放锁）

        分布式系统会遇到的三座大山：NPC。
            N：Network Delay，网络延迟
            P：Process Pause，进程暂停（GC）
            C：Clock Drift，时钟漂移

        Martin 表示 红锁(锁的分布式) 不能兼顾效率和正确性
            客户端1在获取锁时, 触发了GC,进入STW状态,等唤醒后锁已经失效,客户端2此时也可以拿到锁,产生冲突
                客户端 1 请求锁定节点 A、B、C、D、E
                客户端 1 的拿到锁后，进入 GC（时间比较久）
                所有 Redis 节点上的锁都过期了
                客户端 2 获取到了 A、B、C、D、E 上的锁
                客户端 1 GC 结束，认为成功获取锁
                客户端 2 也认为获取到了锁，发生「冲突」

            系统时间发生跳跃,即各服务器之间在不同的时间内时间不同步,导致锁的提前释放
            Redlock 必须「强依赖」多个节点的时钟是保持同步的，一旦有节点时钟发生错误，那这个算法模型就失效了。
                客户端 1 获取节点 A、B、C 上的锁，但由于网络问题，无法访问 D 和 E
                节点 C 上的时钟「向前跳跃」，导致锁到期,或者故障重启,C上的数据没了
                客户端 2 获取节点 C、D、E 上的锁，由于网络问题，无法访问 A 和 B
                客户端 1 和 2 现在都相信它们持有了锁（冲突）

            Redlock 的算法是建立在「同步模型」基础上的，同步模型的假设，在分布式系统中是有问题的。
            在混乱的分布式系统的中，你不能假设系统时钟就是对的，所以，你必须非常小心你的假设。

            fencing token 的方案，保证分布式锁的正确性.
                TODO 共享资源怎么拒绝后来者的请求 ??? 这个是安全的吗?
                TODO mysql的乐观锁 ???


            一个好的分布式锁，无论 NPC 怎么发生，可以不在规定时间内给出结果，但并不会给出一个错误的结果。
            也就是只会影响到锁的「性能」（或称之为活性），而不会影响它的「正确性」。

            Martin 的结论：
                1、Redlock 不伦不类：它对于效率来讲，Redlock 比较重，没必要这么做，
                   而对于正确性来说，Redlock 是不够安全的。
                2、时钟假设不合理：该算法对系统时钟做出了危险的假设（假设多个节点机器时钟都是一致的），
                   如果不满足这些假设，锁就会失效。
                3、无法保证正确性：Redlock 不能提供类似 fencing token 的方案，所以解决不了正确性的问题。
                   为了正确性，请使用有「共识系统」的软件，例如 Zookeeper。

        Redis 作者 Antirez 的反驳
            手动修改时钟：不要这么做就好了，否则你直接修改 Raft 日志，那 Raft 也会无法工作...
            时钟跳跃：通过「恰当的运维」，保证机器时钟不会大幅度跳跃（每次通过微小的调整来完成），实际上是可以做到的

            T2-T1>EXPIRE 可以保证,客户端在确认拿到锁之前,判断redis中的kEY是否过期

            客户端确认拿到了锁，去操作共享资源的途中发生了问题，导致锁失效，那这不止是 Redlock 的问题，
            任何其它锁服务例如 Zookeeper，都有类似的问题，这不在讨论范畴内

            质疑 fencing token 机制:
                客户端要带着这个 token 去改 MySQL 的某一行，这就需要利用 MySQL 的「事务隔离性」来做。
                UPDATE table T SET val = $new_val WHERE id = $id AND current_token < $token
                如果操作的不是 MySQL 呢？例如向磁盘上写一个文件，或发起一个 HTTP 请求，
                那这个方案就无能为力了，这对要操作的资源服务器，提出了更高的要求。

                ??? 拿锁之后，在到数据库层面互持？A写了，B在去写不等，那Ａ最后不去擦掉呢？不就造成死锁了？
                UPDATE table T SET val = $new_val WHERE id = $id AND current_token = $redlock_value

                ??? 既然资源服务器都有了「互斥」能力，那还要分布式锁干什么　？？？

                分布式锁的本质，是为了「互斥」，只要能保证两个客户端在并发时，
                一个成功，一个失败就好了，不需要关心「顺序性」。

        Zookeeper:
            客户端通过定时心跳与服务端维护一个　Session,如果长时间收不到客户端的心跳，服务端则认为会话失效，
            然后主动删除锁．当然客户端也可以主动释放锁．

            客户端１拿锁，进入GC，长时间服务端删除节点；客户端２拿锁，客户端１ＧＣ结束，仍然认为自己有锁．
            ??? 同理也可以增加一个守护线程检测锁是否失效，但是谁来守护守护线程是否失效　？？？

            Zookeeper 的优缺点：
                不需要考虑锁的过期时间
                watch 机制，加锁失败，可以 watch 等待锁释放，实现乐观锁

                性能不如 Redis
                部署和运维成本高
                客户端与 Zookeeper 的长时间失联，锁被释放问题

        作者总结：
            保证时钟不偏移，可以使用Redlock,但是时钟偏移是时有发生的，不建议使用．
            使用类　fencing token 机制，在操作资源上加锁，NOTE 但是操作资源上的锁的释放，又会引来问题？？？

        多个系统操作共享数据:
            1-查出共享数据,更新时间戳,并缓存到redis
            2-对共享资源加锁,加锁失败等待或者返回
            3-如果遇到锁失效,则当前要更新的数据开始时间是不是大于缓存中的时间戳
            4-是则继续,否则返回
    """


class HA(object):
    """
    主从模式
        # 配置主节点的ip和端口
        slaveof 192.168.1.10 6379
        # 从redis2.6开始，从节点默认是只读的
        slave-read-only yes
        # 假设主节点有登录密码，是123456
        masterauth 123456
        slaveof no one

        手动实现主从切换,自主决定主从角色

    哨兵模式
        在主从模式下，redis 同时提供了哨兵命令redis-sentinel，哨兵是一个独立的进程，作为进程，它会独立运行。
        其原理是哨兵进程向所有的 redis 机器发送命令，等待 Redis 服务器响应，从而监控运行的多个 Redis 实例。


    集群模式

    数据同步问题
    """
    pass


def code_test():
    a = "大萨达"  # 字面量默认使用 UTF-8 编码
    ab = a.encode(encoding="UTF-8")
    print(a, ab, ab.decode("GBK"))


if __name__ == '__main__':
    connection_pool = redis.ConnectionPool(host='127.0.0.1', db=0, password="holyzing")
    rd: Redis = redis.Redis(connection_pool=connection_pool)
    rd.setnx("a", 1)
    r = rd.incr("a", 1)


# 已经实现的轮子
# Redlock-py (Python implementation).
# Pottery (Python implementation).
# Aioredlock (Asyncio Python implementation).

# redis的数据结构设计的挺巧妙的, redis的数据结构设计 ？

# actor akka
