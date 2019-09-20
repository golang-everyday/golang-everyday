安装依赖库sarama
go get github.com/Shopify/sarama
该库要求kafka版本在0.8及以上，支持kafka定义的high-level API和low-level API，但不支持常用的consumer自动rebalance和offset追踪，所以一般得结合cluster版本使用。


安装依赖库sarama-cluste
go get github.com/bsm/sarama-cluster

需要kafka 0.9及以上版本



//生产者代码

package main

import (
​	"fmt"
​	"time"

	"github.com/Shopify/sarama"
)

type KafkaProducer struct {
​	RecvChan chan *sarama.ProducerMessage
​	QuitChan chan struct{}
​	Client   sarama.AsyncProducer
}

var KafkaPro *KafkaProducer

func init() {
​	KafkaPro = &KafkaProducer{
​		RecvChan: make(chan *sarama.ProducerMessage, 1),
​		QuitChan: make(chan struct{}, 1),
​	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	config.Version = sarama.V2_1_0_0
	
	//生产者
	client, err := sarama.NewAsyncProducer([]string{"192.168.4.17:9092"}, config)
	if err != nil {
		fmt.Println("producer close,err:", err)
		return
	}
	KafkaPro.Client = client
}

func (k *KafkaProducer) Start() {
​	go func(c sarama.AsyncProducer) {
​		errors := c.Errors()
​		success := c.Successes()
​		for {
​			select {
​			case err := <-errors:
​				if err != nil {
​					fmt.Println(err)
​				}
​			case <-success:
​			}
​		}
​	}(k.Client)

	for {
		select {
		case msg := <-k.RecvChan:
			k.Client.Input() <- msg
		case <-k.QuitChan:
			if err := k.Client.Close(); err != nil {
				fmt.Println(err)
			}
			return
		}
	}
}

func (k *KafkaProducer) SendMessage(msg *sarama.ProducerMessage) {
​	k.RecvChan <- msg
}



//消费者代码

package main

import (
​	"fmt"

	"github.com/Shopify/sarama"
)

type KafkaConsumer struct {
​	QuitChan chan struct{}
​	Consumer sarama.Consumer
}

var KafkaCon *KafkaConsumer

func init() {
​	KafkaCon = &KafkaConsumer{
​		QuitChan: make(chan struct{}, 1),
​	}
​	config := sarama.NewConfig()
​	config.Consumer.Return.Errors = true
​	config.Version = sarama.V2_1_0_0

	// consumer
	consumer, err := sarama.NewConsumer([]string{"192.168.4.17:9092"}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	KafkaCon.Consumer = consumer
}

func (k *KafkaConsumer) Start(topic string, partition int32, offset int64) {
​	partitionConsumer, err := k.Consumer.ConsumePartition(topic, partition, offset)
​	if err != nil {
​		fmt.Println(err)
​		return
​	}

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
				msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Printf("err :%s\n", err.Error())
		case <-k.QuitChan:
			if err := partitionConsumer.Close(); err != nil {
				fmt.Println(err)
			}
			return
		}
	}
}

func main() {
​	KafkaCon.Start("log", 0, sarama.OffsetOldest)
}
