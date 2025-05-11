package logs

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap/zapcore"
	"os"
)

// kafka生产者结构
type KafkaProducer struct {
	topic        string
	syncProducer sarama.SyncProducer
}

// kafka生产者
var kProducer *KafkaProducer

// kafka传输
func (kp *KafkaProducer) Write(p []byte) (n int, err error) {

	_, _, err = kp.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.ByteEncoder(p),
	})
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// 关闭kafka
func (kp *KafkaProducer) Close() error {

	var err error
	if kp.syncProducer != nil {
		err = kp.syncProducer.Close()
	}

	return err
}

// kafka写入对象
func getKafkaWriter(config LoggerConfig) zapcore.WriteSyncer {

	var err error
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.NoResponse
	saramaConfig.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(config.KafkaInfo.Address, saramaConfig)
	if err != nil {
		return nil
	}

	kProducer = &KafkaProducer{
		topic:        config.KafkaInfo.Topic,
		syncProducer: syncProducer,
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(kProducer),
		zapcore.AddSync(os.Stdout))
}

// 关闭kafka
func closeKafka() error {

	var err error
	if kProducer != nil {
		err = kProducer.Close()
	}

	return err
}
