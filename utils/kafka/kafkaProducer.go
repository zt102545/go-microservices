package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"go-microservices/utils/gos"
	"go-microservices/utils/logs"
	"time"
)

type Producer struct {
	producer sarama.AsyncProducer
}

func NewProducer(kafkaAddr []string) *Producer {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal     // ack=1
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages
	// Flush批量发送参数
	config.Producer.Flush.Bytes = 16384                       //>16kb 触发批量发送
	config.Producer.Flush.Messages = 10000                    //>10000 触发批量发送
	config.Producer.Flush.Frequency = 1000 * time.Millisecond //>1s 触发批量发送
	config.Producer.Flush.MaxMessages = 50000                 //>50000 上限

	config.Producer.Return.Errors = true
	//config.Producer.Return.Successes = true
	client, err := sarama.NewClient(kafkaAddr, config)
	if err != nil {
		logs.Err(context.Background(), "create kafka client fail,err:%v\n", err.Error(), logs.Flag("kafka"))
		return nil
	}
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		logs.Err(context.Background(), "create kafka producer fail,err:%v\n", err.Error(), logs.Flag("kafka"))
		return nil
	}
	instance := &Producer{
		producer: producer,
	}
	gos.GoSafe(func() {
		instance.Return()
	})
	return instance
}

func (p *Producer) SendMessage(topic, key, value string) {
	kValue := sarama.StringEncoder(value)
	eventKey := sarama.StringEncoder(key)
	timestamp := time.Now()
	message := &sarama.ProducerMessage{Topic: topic, Value: kValue, Timestamp: timestamp, Key: eventKey}
	logs.Info(context.Background(), "send message to kafka,topic:%v,key:%v,value:%v\n", topic, key, value, logs.Flag("kafka"))
	p.producer.Input() <- message
}

func (p *Producer) Return() {
	for {
		select {
		//case msg := <-p.producer.Successes():
		//	SuccessesNum++
		//	fmt.Printf("Successes,num:%v,msg:%v\n", SuccessesNum, msg)
		case err := <-p.producer.Errors():
			logs.Err(context.Background(), "send message to kafka is fail,err:%v\n", err, logs.Flag("kafka"))
		}
	}
}

func (p *Producer) Close() {
	p.producer.AsyncClose()
}
