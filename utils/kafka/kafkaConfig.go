package kafka

import "sync"

var (
	kfLock        sync.Mutex
	KafkaProducer *Producer
	KafkaConsumer *Consumer
)

type KafkaConfig struct {
	Addr    []string `json:"addr"`
	GroupId string   `json:"group_id,optional"`
}

func (k *KafkaConfig) InitProducer() *Producer {
	if KafkaProducer != nil {
		return KafkaProducer
	}
	kfLock.Lock()
	defer kfLock.Unlock()
	if len(k.Addr) > 0 {
		KafkaProducer = NewProducer(k.Addr)
	}
	return KafkaProducer
}

func (k *KafkaConfig) InitConsumer() *Consumer {
	if KafkaConsumer != nil {
		return KafkaConsumer
	}
	kfLock.Lock()
	defer kfLock.Unlock()
	if len(k.Addr) > 0 && k.GroupId != "" {
		KafkaConsumer = NewConsumerGrp(k.Addr, k.GroupId)
	}
	return KafkaConsumer
}
