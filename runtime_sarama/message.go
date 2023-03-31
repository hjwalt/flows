package runtime_sarama

import (
	"github.com/Shopify/sarama"
	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/runway/logger"
)

// mapping sarama consumer message to internal representation
func FromConsumerMessage(source *sarama.ConsumerMessage) (message.Message[message.Bytes, message.Bytes], error) {

	byteFormat := format.Bytes()

	// deserialise key
	key, err := byteFormat.Unmarshal(source.Key)
	if err != nil {
		logger.ErrorErr("consumer message map key deserialisation failure", err)
		return message.Message[message.Bytes, message.Bytes]{}, err
	}

	// deserialise value
	value, err := byteFormat.Unmarshal(source.Value)
	if err != nil {
		logger.ErrorErr("consumer message map value deserialisation failure", err)
		return message.Message[message.Bytes, message.Bytes]{}, err
	}

	// map headers
	headers := make(map[string][]message.Bytes)
	for _, header := range source.Headers {
		keyString := string(header.Key)
		if _, exists := headers[keyString]; !exists {
			headers[keyString] = make([]message.Bytes, 0)
		}
		headers[keyString] = append(headers[keyString], header.Value)
	}

	return message.Message[message.Bytes, message.Bytes]{
		Topic:     source.Topic,
		Partition: source.Partition,
		Offset:    source.Offset,
		Timestamp: source.Timestamp,
		Key:       key,
		Value:     value,
		Headers:   headers,
	}, nil
}

func FromConsumerMessages(sources []*sarama.ConsumerMessage) ([]message.Message[message.Bytes, message.Bytes], error) {

	mappedMessages := make([]message.Message[message.Bytes, message.Bytes], 0)
	for _, source := range sources {
		mapped, err := FromConsumerMessage(source)
		if err != nil {
			logger.ErrorErr("consumer messages map deserialisation failure", err)
			return make([]message.Message[message.Bytes, message.Bytes], 0), err
		}
		mappedMessages = append(mappedMessages, mapped)
	}
	return mappedMessages, nil
}

// mapping internal representation into sarama producer message
func ToProducerMessage(source message.Message[message.Bytes, message.Bytes]) (*sarama.ProducerMessage, error) {
	// map headers
	headers := make([]sarama.RecordHeader, 0)
	for k, vs := range source.Headers {
		for _, v := range vs {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: v,
			})
		}
	}

	return &sarama.ProducerMessage{
		Topic:   source.Topic,
		Headers: headers,
		Key:     sarama.ByteEncoder(source.Key),
		Value:   sarama.ByteEncoder(source.Value),
	}, nil
}

func ToProducerMessages(sources []message.Message[message.Bytes, message.Bytes]) ([]*sarama.ProducerMessage, error) {

	mappedMessages := make([]*sarama.ProducerMessage, 0)
	for _, source := range sources {
		mapped, err := ToProducerMessage(source)
		if err != nil {
			logger.ErrorErr("consumer messages map deserialisation failure", err)
			return make([]*sarama.ProducerMessage, 0), err
		}
		mappedMessages = append(mappedMessages, mapped)
	}
	return mappedMessages, nil
}
