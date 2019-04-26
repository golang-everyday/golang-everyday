一、安装java环境

二、安装zookeeper

1  解压zookeeper到指定文件下 

tar -zxvf zookeeper-3.4.10.tar.gz -C apps/

 

2  修改zoo.cfg 文件配置       

cd  /zookeeper/conf/

mv  zoo_sample.cfg  zoo.cfg

vi  zoo.cfg

 

```
 
```

dataDir=/home/sun/data/zkdata/

clientPort=2181

 

dataLogDir=/home/sun/log/zklog/

 

server.1=192.168.54.130:2888:3888

server.2=192.168.54.128:2888:3888

server.3=192.168.54.129:2888:3888

 

 

3  将配置文件发送到其他集群服务器

scp -r zookeeper-3.4.10/  192.168.123.103:/home/sun/apps/zookeeper-3.xx

scp -r zookeeper-3.4.10/  192.168.123.104:/home/sun/apps/zookeeper-3.xx

 

4  在每个zooKeeper 服务器节点，新建目录 /home/sun/data/zkdata，建好之后，在里面新建一个文件，文件名叫 myid（唯一标识）,里面存放的内容就是每个服务器的id（ID从1开始自增）

​    Mkdir  /home/hadoop/data/zkdata

​       cd  data/zkdata/

​       echo  1  >  myid

 

```
5  配置环境变量
```

​       vi .bashrc

 

\#Zookeeper

export ZOOKEEPER_HOME=/home/sun/apps/zookeeper-3.4.10

export PATH=$PATH:$ZOOKEEPER_HOME/bin

 

​       source .bashrc

 

```
6.每台服务器配置完成后，最后在各服务器节点，启动zookeeper
```

启动：zkServer.sh start
 停止：zkServer.sh stop
 查看状态：zkServer.sh status

 

三、安装kafka集群

```
1 解压kafka到指定文件夹 
```

tar        -zxvf kafka_2.11-0.8.2.0.tgz -C apps

cd  apps/

mv        kafka_2.11-0.8.2.0/ kafka

```
2 修改配置文件（每台kafka服务器）
```

cd apps/kafka/config/

vim  server.properties

 

 

//当前机器在集群中的唯一标识，和zookeeper的myid性质一样

broker.id=0

//当前kafka对外提供服务的端口默认是9092

port=9092

//本机ip地址

host.name=192.168.123.102

//log目录

log.dirs=/home/sun/log/kafka-logs

//设置zookeeper的连接端口

zookeeper.connect=192.168.123.102:2181,192.168.123.103:2181,192.168.123.104:2181

3 修改product.properties和consumer.properties配置文件

vim  producer.properties

metadata.broker.list=192.168.123.102:9092,192.168.123.103:9092,192.168.123.104:9092

 

vim  consumer.properties

zookeeper.connect=192.168.123.102:2181,192.168.123.103:2181,192.168.123.104:2181

4  修改环境变量

```
               vi .bashrc
 
               #Kafka
               export KAFKA_HOME=/home/sun/apps/kafka
               export PATH=$PATH:$KAFKA_HOME/bin
 
               source ~/.bashrc
 
5 启动kafka集群
 
               bin/kafka-server-start.sh config/server.properties
```

 

四．Kafka 权限模块的配置（SASL）

1  在 Kafka 安装目录下的 config/server.properties文件中配置以下信息

 

listeners=SASL_PLAINTEXT://192.168.123.102:9092

security.inter.broker.protocol=SASL_PLAINTEXT

sasl.mechanism.inter.broker.protocol=PLAIN

sasl.enabled.mechanisms=PLAIN

authorizer.class.name = kafka.security.auth.SimpleAclAuthorizer

super.users=User:admin

 

2  配置一个名 kafka_server_jaas.conf的配置文件，将配置文件放置在conf目录下

KafkaServer {

org.apache.kafka.common.security.plain.PlainLoginModule required

username="admin"

password="admin"

user_admin="admin"

user_alice="alice";

};

 

3  为Kafka 在 bin/kafka-server-start.sh中添加以下内容

export KAFKA_SASL_OPTS='-Djava.security.auth.login.config=/home/sun/apps/kafka/config/kafka_server_jaas.conf'

4  解决报错问题

export KAFKA_OPTS="-Djava.security.auth.login.config=/home/sun/apps/kafka/config/kafka_server_jaas.conf"

 

五．Kafka的客户端配置

**1** **如果是在程序里** 

 

（1）新建 kafka_client_jaas.conf ，内容是：

KafkaClient {

org.apache.kafka.common.security.plain.PlainLoginModule required

username="alice"

password="alice";

};

 

（2）程序里面的kafka配置信息里（producer 和 consumer）添加如下配置

 

```
config := sarama.NewConfig()
config.Net.SASL.Enable = true
config.Net.SASL.User = "alice"
config.Net.SASL.Password = "alice"
 
client, err := sarama.NewClient(strings.Split(hosts, ","), config)
if err != nil {
    log.Fatalf("unable to create kafka client: %q", err)
}
// consumer api
consumer, err := sarama.NewConsumerFromClient(client)
defer consumer.Close()
// producer api
producer, err := sarama.NewSyncProducerFromClient(client)
defer producer.Close()
```

 

**2**    **如果是在命令行**

 

（1）  在config目录下创建kafka_client_jaas.conf

KafkaClient {

​        org.apache.kafka.common.security.plain.PlainLoginModule required

​        username="alice"

​        password="alice";

};

（2） 在config下的producer.properties和consumer.properties添加下面配置

sasl.jaas.config=org.apache.kafka.common.security.plain.PlainLoginModule required \

username="alice" \

password="alice";

security.protocol=SASL_PLAINTEXT

sasl.mechanism=PLAIN

(3 )在bin下的kafka-console-producer.sh和kafka-console-consumer.sh下添加下面配置

if [ "x$KAFKA_OPTS" ]; then

export KAFKA_OPTS="-Djava.security.auth.login.config=/home/sun/apps/kafka/config/kafka_client_jaas.conf"

fi

 

( 4 )给alice用户设置sal权限

./kafka-acls.sh --authorizer-properties zookeeper.connect=192.168.54.130:2181 --allow-principal User:alice --operation=Read --operation=Write --topic=*  --add

 

 

./kafka-acls.sh --authorizer-properties zookeeper.connect=192.168.54.130:2181 --allow-principal User:alice --consumer --topic=*  --group=*  --add

( 5 ) 第一次握手验证权限时，需要的时间比较长，需要修改默认等待时间

vim server.properties

zookeeper.connection.timeout.ms=10000

 

​        vim producer.properties

request.timeout.ms=10000

 

vim consumer.properties

consumer.timeout.ms=10000

 

（6）生产者生产(sasl)命令 ：

./kafka-console-producer.sh --broker-list  192.168.123.102:9092 --topic topic1

--producer.config=/home/sun/apps/kafka/config/producer.properties

 

（7）消费者消费(sasl)命令 ：

./kafka-console-consumer.sh --bootstrap-server 192.168.123.102:9092 --topic test  --from-beginning

 --consumer.config=/home/sun/apps/kafka/config/consumer.properties

（8）使用docker配置zookeeper时， 由于云服务器采用虚拟化的技术，监听的网卡是属于物理网关的网卡，而虚拟机内部自然没有这个网卡。需要让服务器进程监听0.0.0.0的ip地址，也就是监听所有网卡。在server.properties中

设置 quorumListenOnAllIPs=true

 

 

（9）kafka 的每个broker有默认的consumer group分组：test-consumer-group。在SASL    消费前要给用户（alice）分配该group组的权限才能进行消费。