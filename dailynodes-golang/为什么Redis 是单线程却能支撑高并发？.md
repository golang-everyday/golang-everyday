## 为什么Redis 是单线程却能支撑高并发？



## **几种 I/O 模型**

为什么 Redis 中要使用 I/O 多路复用这种技术呢？

首先，Redis 是跑在单线程中的，所有的操作都是按照顺序线性执行的，但是由于读写操作等待用户输入或输出都是阻塞的，所以 I/O 操作在一般情况下往往不能直接返回，这会导致某一文件的 I/O 阻塞导致整个进程无法对其它客户提供服务，而 I/O 多路复用就是为了解决这个问题而出现的。

## **Blocking I/O**

先来看一下传统的阻塞 I/O 模型到底是如何工作的：当使用 read 或者 write 对某一个文件描述符（File Descriptor 以下简称 FD)进行读写时，如果当前 FD 不可读或不可写，整个 Redis 服务就不会对其它的操作作出响应，导致整个服务不可用。

这也就是传统意义上的，也就是我们在编程中使用最多的阻塞模型：

![img](https://mmbiz.qpic.cn/mmbiz_png/eQPyBffYbufXLqJSl9NGibsgbag6icanicPVzYGggto1cNBp5b3Tc2XzXoKvUE5ve5vjM1G0iclkY9PtdicWNHibiaeNA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

阻塞模型虽然开发中非常常见也非常易于理解，但是由于它会影响其他 FD 对应的服务，所以在需要处理多个客户端任务的时候，往往都不会使用阻塞模型。

## **I/O 多路复用**

虽然还有很多其它的 I/O 模型，但是在这里都不会具体介绍。

阻塞式的 I/O 模型并不能满足这里的需求，我们需要一种效率更高的 I/O 模型来支撑 Redis 的多个客户（redis-cli），这里涉及的就是 I/O 多路复用模型了：

![img](https://mmbiz.qpic.cn/mmbiz_png/eQPyBffYbufXLqJSl9NGibsgbag6icanicPNMdorNIz3r7J0ic305hMqOWQ8uukGMkcoFJsI5Siao2ZqfO1ZZnLcOtA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

在 I/O 多路复用模型中，最重要的函数调用就是 select，该方法的能够同时监控多个文件描述符的可读可写情况，当其中的某些文件描述符可读或者可写时，select 方法就会返回可读以及可写的文件描述符个数。

关于 select 的具体使用方法，在网络上资料很多，这里就不过多展开介绍了；

与此同时也有其它的 I/O 多路复用函数 epoll/kqueue/evport，它们相比 select 性能更优秀，同时也能支撑更多的服务。

## **Reactor 设计模式**

Redis 服务采用 Reactor 的方式来实现文件事件处理器（每一个网络连接其实都对应一个文件描述符）

![img](https://mmbiz.qpic.cn/mmbiz_png/eQPyBffYbufXLqJSl9NGibsgbag6icanicPUHkTHo8sByibr95Zv2arEiaU47ibZNhLicMiaffticWa1GLib4Q55gnl0ZBBg/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

文件事件处理器使用 I/O 多路复用模块同时监听多个 FD，当 accept、read、write 和 close 文件事件产生时，文件事件处理器就会回调 FD 绑定的事件处理器。

虽然整个文件事件处理器是在单线程上运行的，但是通过 I/O 多路复用模块的引入，实现了同时对多个 FD 读写的监控，提高了网络通信模型的性能，同时也可以保证整个 Redis 服务实现的简单。

## **I/O 多路复用模块**

I/O 多路复用模块封装了底层的 select、epoll、avport 以及 kqueue 这些 I/O 多路复用函数，为上层提供了相同的接口。![img](https://mmbiz.qpic.cn/mmbiz_jpg/eQPyBffYbufXLqJSl9NGibsgbag6icanicPqtibrzicDWbicKGoibVqp9e8M2DFmIqIYKsTB0VOFXAIKSRdwHbaU1dquA/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

在这里我们简单介绍 Redis 是如何包装 select 和 epoll 的，简要了解该模块的功能，整个 I/O 多路复用模块抹平了不同平台上 I/O 多路复用函数的差异性，提供了相同的接口：

- static int aeApiCreate(aeEventLoop *eventLoop)
- static int aeApiResize(aeEventLoop *eventLoop, int setsize)
- static void aeApiFree(aeEventLoop *eventLoop)
- static int aeApiAddEvent(aeEventLoop *eventLoop, int fd, int mask)
- static void aeApiDelEvent(aeEventLoop *eventLoop, int fd, int mask)
- static int aeApiPoll(aeEventLoop *eventLoop, struct timeval *tvp)

同时，因为各个函数所需要的参数不同，我们在每一个子模块内部通过一个 aeApiState 来存储需要的上下文信息：

```
// select
typedef struct aeApiState {
    fd_set rfds, wfds;
    fd_set _rfds, _wfds;
} aeApiState;

// epoll
typedef struct aeApiState {
    int epfd;
    struct epoll_event *events;
} aeApiState;
```

这些上下文信息会存储在 eventLoop 的 void *state 中，不会暴露到上层，只在当前子模块中使用。



