启动zookeeper  
zkServer.sh start

停止zookeeper
zkServer.sh stop

查看zookeeper状态
zkServer.sh status

----------------------------------------------------------------
启动kafka集群服务

bin/kafka-server-start.sh config/server.properties 

bin/kafka-server-start.sh -daemon config/server.properties (以守护进程启动)

停止kafka集群服务

bin/kafka-server-stop.sh

------------------------------------------------------------------


创建topic

bin/kafka-topics.sh --create --zookeeper 192.168.54.130:2181,192.168.54.128:2181,192.168.54.129:2181 --replication-factor 3 --partitions 3 --topic topic1


查看已经创建的topic信息

bin/kafka-topics.sh --list --zookeeper 192.168.54.130:2181


product生产命令

bin/kafka-console-producer.sh --broker-list 192.168.54.130:9092 --topic topic1


consumer消费命令

bin/kafka-console-consumer.sh --zookeeper 192.168.54.130:2181 --from-beginning --topic topic1

-----------------------------------------------------------------------------------------------------------

SASL权限 product 生产命令
./kafka-console-producer.sh --broker-list 192.168.54.130:9092 --topic topic1 --producer.config= /home/sun/apps/kafka/config/producer.properties


SASL权限 consumer 消费命令
./kafka-console-consumer.sh --bootstrap-server 192.168.54.130:9092 --topic topic1 --from-beginning
 --consumer.config=/home/sun/apps/kafka/config/consumer.properties

------------------------------------------------------------------------------------------------------------

##### SASL服务器增加权限
./kafka-acls.sh --authorizer-properties zookeeper.connect=192.168.54.130:2181 --allow-principal User:alice --operation=Read --operation=Write --topic=* 

./kafka-acls.sh --authorizer-properties zookeeper.connect=192.168.54.130:2181 --allow-principal User:alice --consumer --topic=* --group=* --add


##### SASL查询权限
./kafka-acls.sh --authorizer-properties zookeeper.connect=192.168.54.130:2181 --list --topic topic1


##### SASL服务器删除权限
./kafka-acls.sh --authorizer-properties zookeeper.connect=192.168.54.130:2181 --remove --allow-principal User:alice --allow-principal User:Alice  --operation Read --operation Write --topic topic1

------------------------------------------------------------------------------------------------------------

坑点：

1  使用docker配置zookeeper时， 若使用云服务器。云服务器采用虚拟化的技术，监听的网卡是属于物理网关的网卡，而虚拟化机内部自然没有这个网卡。需要让服务器进程监听0.0.0.0的ip地址，也就是监听所有网卡。
并设置 quorumListenOnAllIPs=true


2  程序报错：Could not find a 'KafkaServer' or 'sasl_plaintext.KafkaServer' entry in the JAAS configuration. System property 'java.security.auth.login.config' is not set   

解决办法：

export KAFKA_OPTS="-Djava.security.auth.login.config=/home/sun/apps/kafka/config/kafka_server_jaas.conf"

3  kafka的每个broker有默认的consumer group分组：test-consumer-group。在SASL消费前要给用户（alice）分配该group组的权限才能进行消费。