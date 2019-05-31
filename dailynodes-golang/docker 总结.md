

Namespace  :  做隔离pid，net，ipc，mnt，uts

control groups :  做资源限制

union file systems : Container和images的分层



##### 什么是docker images?

文件和metadata 的集合（root filesystem）

不同的imagse可以共享一个layer

images是只读的，负责app的存储和分发

##### 什么是docker container？

通过images创建，负责app的运行

##### cmd和entrypoint的区别：

cmd : 容器启动时默认执行，可以被docker run后面的命令覆盖

entrypoint : 容器启动后指定执行，不可以被docker run后面的命令覆盖



##### docker的network namespace

单机形式：Bridge network,Host network,None network

多机形式：Overlay network

查看本机所有的network  namespace : ip netns list

创建本机的network namespace :ip netns add test1

创建一对Veth pair，将两端接口放在test1和test2，这样test1和test2可以通信：

ip link add veth-test1 type veth peer name veth-test2

ip link set veth-test1 netns test1

ip link set veth-test2 netns test2

在test1和test2容器内为veth-test1、veth-test2分配ip地址

ip netns exec  test1 ip addr add 192.168.1.1/24 dev veth-test1

ip netns exec  test2 ip addr add 192.168.1.2/24 dev veth-test2

启动

ip netns exec test1 ip link set dev veth-test1 up

ip netns exec test2 ip link set dev veth-test2 up



创建bridge

docker network create  -d bridge my-mybridge

将docker连接到bridge上

docker network connect my-bridge test2 

多机docker通信：

创建overlay网络，然后指定overlay网络创建docker即可

使用blind mount映射同步，实时更新

##### 利用docker swarm 的service更新

docker service scale web = 2

docker service update --image xiaopeng163/python-flask-demo:2.0 web



### 1.6 Docker和k8s

#### 1.6.1 Docker

**Docker解决的问题**

- 保证程序运行环境的一致性；
- 降低配置开发环境、生产环境的复杂度和成本；
- 实现程序的快速部署和分发。

**Docker的核心技术**

镜像、容器、数据、网络

#### 1.6.2 K8s概念

**Cluster（集群）**

cluster是计算、存储和网络资源的集合，k8s利用这些资源运行各种基于容器的应用。

**Master（控制主节点）**

是Cluster的大脑，主要职责是调度，即决定应用放在哪里运行。为了实现高可用，可以运行多个Master。调度应用程序、维护应用程序的所需状态、扩展应用程序和滚动更新都是Master的主要工作。

**Node（节点）**

Node的职责是运行容器应用，由Master管理，负责监控并汇报容器的状态，同时根据Master管理容器的生命周期。

Node是Kubernetes集群中的工作机器，可以是物理机或虚拟机。每个节点都有一个kubelet，他是管理节点与Master节点进行通信的代理。

一个Kubernetes集群`至少有三个节点`。

**Pod（资源对象）**

Pod是K8s的最小工作单元，每个Pod含有一个或多个容器，Pod中的容器会作为一个整体被Master调度到一个Node上运行。

K8s引入Pod有两个目的：1）客观理性。有些容器天生需要紧密联系，一起工作。Pod提供了比容器更高层次的抽象，将他们封装到一个部署单元中。2）通信和资源共享。Pod所有的容器使用一个网络namespace，即相同的`IP地址和Port空间`，他们可以直接用localhost通信，同样这些容器可以`共享存储`。当k8s挂载volume到Pod，本质上是将volume挂载到Pod的每个容器。

**Controller（控制器）**

K8s通常不会直接创建Pod，而是通过Controller来管理Pod。Controller中定义了Pod的部署特性，比如有几个副本、在什么样的Node上运行等。K8s提供了多种Controller，包括Deployment、ReplcaSet、DaemonSet、StatefulSet、Job等。

Deployment是最常用的Controller，Deployment可以管理Pod的多个副本，并确保Pod按照期望的状态运行。

**Service（服务）**

Deploymen可以部署多个副本，每个Pod都有自己的IP，外界如何访问这些副本呢，通过IP吗？是Service。Kubernetes Service定义了外界访问的一组特定的Pod的方式。Service有自己的IP和端口，Service为Pod提供了负载均衡。

**Namespace（命名空间）**

Namespace 是对一组资源和对象的抽象集合，比如可以用来将系统内部的对象划分为不同的项目组或用户组。
如果有多个用户或项目组使用同一个Kubernetes Cluster，如何将他们创建的Controller、Pod等资源分开呢？ 答案就是Namespace。
Namespace可以将一个物理的Cluster逻辑上划分成多个虚拟Cluster，每个Cluster就是一个Namespace。不同Namespace里的资源是完全隔离的。

**Kubernetes默认创建了两个Namespace**

default：创建资源时如果不指定，将被放到这个Namespace中。
kube-system：Kubernetes自己创建的系统资源将放到这个Namespace中。