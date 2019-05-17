### etcd 应用场景

1  服务发现（Service Discovery） 

2  消息发布与订阅 

3  负载均衡 

4  分布式通知与协调 

5  分布式锁、分布式队列 

6  集群监控与Leader竞选

### etcd 工作原理

#### 选主 

- 当集群初始化时候，每个节点都是Follower角色；

- 集群中存在至多1个有效的主节点，通过心跳与其他节点同步数据；

- 当Follower在一定时间内没有收到来自主节点的心跳，会将自己角色改变为Candidate，并发起一次选主投票；当收到包括自己在内超过半数节点赞成后，选举成功；当收到票数不足半数选举失败，或者选举超时。若本轮未选出主节点，将进行下一轮选举（出现这种情况，是由于多个节点同时选举，所有节点均为获得过半选票）。

- Candidate节点收到来自主节点的信息后，会立即终止选举过程，进入Follower角色。

  为了避免陷入选主失败循环，每个节点未收到心跳发起选举的时间是一定范围内的随机值，这样能够避免2个节点同时发起选主。

#### 日志复制

​	日志复制是指主节点将每次操作形成日志条目，并持久化到本地磁盘，然后通过网络IO发送给其他节点。其他节点根据日志的逻辑时钟(TERM)和日志编号(INDEX)来判断是否将该日志记录持久化到本地。当主节点收到包括自己在内超过半数节点成功返回，那么认为该日志是可提交的(committed），并将日志输入到状态机，将结果返回给客户端。

#### 安全性

​	选主以及日志复制并不能保证节点间数据一致。当一个某个节点挂掉了，一段时间后再次重启，并当选为主节点。在其挂掉这段时间内，集群会正常工作，那么会有日志提交。这些提交的日志无法传递给挂掉的节点。当挂掉的节点再次当选主节点，它将缺失部分已提交的日志。

​	Raft解决的办法是，在选主逻辑中，对能够成为主的节点加以限制，确保选出的节点已定包含了集群已经提交的所有日志。如果新选出的主节点已经包含了集群所有提交的日志，那就不需要从和其他节点比对数据了。

### ETCD网络层实现

​	在目前的实现中，ETCD通过HTTP协议对外提供服务，同样通过HTTP协议实现集群节点间数据交互。网络层的主要功能是实现了服务器与客户端(能发出HTTP请求的各种程序)消息交互，以及集群内部各节点之间的消息交互。

​	ETCD-SERVER 大体上可以分为网络层，Raft模块，复制状态机，存储模块，架构图如图所示。

![94f36270eb8fd5e75fb744c78b33e7dfa709bd92.png](https://github.com/golang-everyday/golang-everyday/blob/master/picture/94f36270eb8fd5e75fb744c78b33e7dfa709bd92.png?raw=true)

- 网络层：提供网络数据读写功能，监听服务端口，完成集群节点之间数据通信，收发客户端数据；

- Raft模块：完整实现了Raft协议；

- 存储模块：KV存储，WAL文件，SNAPSHOT管理

- 复制状态机：这个是一个抽象的模块，状态机的数据维护在内存中，定期持久化到磁盘，每次写请求会持久化到WAL文件，并根据写请求的内容修改状态机数据。

  

  节点之间网络拓扑结构 ETCD集群的各个节点之间需要通过HTTP协议来传递数据，表现在：

- Leader 向Follower发送心跳包, Follower向Leader回复消息；

- Leader向Follower发送日志追加信息；

- Leader向Follower发送Snapshot数据；

- Candidate节点发起选举，向其他节点发起投票请求；

- Follower将收的写操作转发给Leader;


​       各个节点在任何时候都有可能变成Leader, Follower, Candidate等角色，为了减少创建链接开销，ETCD节点在启动之初就创建了和集群其他节点之间的长链接。每个节点会向其他节点宣告自己监听的端口，该端口只接受来自其他节点创建链接的请求。因此，ETCD集群节点之间的网络拓扑是一个任意2个节点之间均有长链接相互连接的网状结构。

![de5daff21c3a6ef46d204494a1946bd380ee8296.png](https://github.com/golang-everyday/golang-everyday/blob/master/picture/de5daff21c3a6ef46d204494a1946bd380ee8296.png?raw=true)

### 节点之间消息交互

​	在ETCD实现中，消息采取了分类处理，抽象出了2种类型消息传输通道：Stream类型通道和Pipeline类型通道。这两种消息传输通道都使用HTTP协议传输数据，通过protocol buffer进行封装。

![e53bbd1a314675f6a777e3c0c41705eaf758f6c9.png](https://github.com/golang-everyday/golang-everyday/blob/master/picture/e53bbd1a314675f6a777e3c0c41705eaf758f6c9.png?raw=true)

- Stream类型通道：点到点之间维护HTTP长链接，主要用于传输数据量较小的消息，例如追加日志，心跳等。

  ​     使用了golang的http包实现Stream类型通道：

  1）被动发起方监听端口, 并在对应的url上挂载相应的handler（当前请求来领时，handler的ServeHTTP方法会被调用）

  2）主动发起方发送HTTP GET请求；

  3）监听方的Handler的ServeHTTP访问被调用(框架层传入http.ResponseWriter和http.Request对象），其中http.ResponseWriter对象作为参数传入Writter-Goroutine，该Goroutine的主循环就是将Raft模块传出的message写入到这个responseWriter对象里；http.Request的成员变量Body传入到Reader-Gorouting，该Gorutine的主循环就是不断读取Body上的数据，decode成message 通过Channel传给Raft模块。

- Pipeline类型通道：点到点之间不维护HTTP长链接，短链接传输数据，用完即关闭。用于传输数据量大的消息，例如snapshot数据

   使用了golang的http包实现Pipeline类型通道：

  1）根据参数配置，启动N个Goroutines；

  2）每一个Goroutines的主循环阻塞在消息Channel上，当收到消息后，通过POST请求发出数据，并等待回复。