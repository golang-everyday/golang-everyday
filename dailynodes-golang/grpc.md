gRPC 是在 HTTP/2 之上实现的 RPC 框架，HTTP/2 是第 7 层（应用层）协议，它运行在 TCP（第 4 层 - 传输层）协议之上，相比于传统的 REST/JSON 机制有诸多的优点：

1. 基于 HTTP/2 之上的二进制协议（Protobuf 序列化机制）；
2. 一个连接上可以多路复用，并发处理多个请求和响应；
3. 多种语言的类库实现；
4. 服务定义文件和自动代码生成（.proto 文件和 Protobuf 编译工具）。

### GRPC服务端创建流程

1    创建NettyServer，gRPC服务端的创建，初始化NettyServer，NettyServer负责监听socket地址，实现基于HTTP2协议的写入。

2    绑定 IDL定义服务接口实现类，gRPC和一些RPC框架的不同是，服务的接口实现类并不是同过反射实现的，而是通过proto工具生成的代码。服务启动后，将服务的接口实现类注册到grpc内部的服务注册中心上，请求消息来后，便可以通过服务名和方法名调用 ，直接调用启动的时候注册的服务实例，不需要反射进行调用，性能更优。

3    ServerImpl 负责整个 gRPC 服务端消息的调度和处理，创建 ServerImpl 实例过程中，会对服务端依赖的对象进行初始化，例如 Netty 的线程池资源 ， gRPC 的线程池 ，内部的服务注册类（InternalHandlerRegistry)

![img](https://upload-images.jianshu.io/upload_images/7706563-ee37da4e27423dee.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/670/format/webp)

### 1  gRpc消息接入流程

gRPC消息由netty /http/2 协议负责接入，通过grpc 注册的Http2Framelister将解码后的Http header和Http body 发送到gRPC的NettyServerHandler ，实现netty http/2的消息接入
 gRPC 请求消息接入流程如下:

![img](https://upload-images.jianshu.io/upload_images/7706563-695d455dd9439b99.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/692/format/webp)

### 2 gRPC消息响应模型

### 3 Netty Server  线程模型

gRPC 的线程模型遵循 Netty 的线程分工原则，即：协议层消息的接收和编解码由 Netty 的 I/O(NioEventLoop) 线程负责；后续应用层的处理由应用线程负责，防止由于应用处理耗时而阻塞 Netty 的 I/O 线程
 不过正是因为有了分工原则，grpc 之间会做频繁的线程切换，如果在一次grpc调用过程中，做了多次I/O线程到应用线程之间的切换，会导致性能的下降，这也是为什么grpc在一些私有协议支持不太友好的原因

## gRpc 的线程模型

### 1. BIO线程模型 ，例如tomcat的BIO线程模型

![img](https://upload-images.jianshu.io/upload_images/7706563-4aff067f11df2e0d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/639/format/webp)

缺点

改进:
优化后BIO线程模型采用了线程池的做法但是后端的应用处理线程仍然采用同步阻塞的模型，阻塞的时间取决对方I/O处理的速度和网络I/O传输的速度

### 2 异步非阻塞的线程模型

grpc的线程模型主要包含服务端线程模型，客户端线程模型

#### 服务端线程模型主要包括

1. 服务端的写入，客户端的接入线程（HTTP/2 Acceptor）
2. 网络I/O的读写线程
3. 服务接口调用线程

#### 客户端线程模型主要包含

1. 客户端的链接 （HTTP/2 Connector）
2. 网络I/O读写线程
3. 接口调用线程
4. 响应回调通知线程

### 2.1服务调度线程模型

I/O 通信线程模型
 gRPC的做法是服务监听线程和I/O线程分离Reactor多个线程模型 其工作原理如下:

![img](https://upload-images.jianshu.io/upload_images/7706563-14535640deeaa79f.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/591/format/webp)

服务调度线程模型

### 2.2客户端线程模型概述

gRPC 客户端线程模型工作原理如下图所示（同步阻塞调用为例)

![img](https://upload-images.jianshu.io/upload_images/7706563-b30682b98a14e819.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/651/format/webp)

#####  I/O 通信线程模型

相比于服务端，客户端的线程模型简单一些，它的工作原理如下：

![img](https://upload-images.jianshu.io/upload_images/7706563-d379bdaa6d8c3612.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/554/format/webp)

##### 客户端调用线程模型

![img](https://upload-images.jianshu.io/upload_images/7706563-ec610333a2d2fafa.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/744/format/webp)