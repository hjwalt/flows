package metric

type Retry interface {
	RetryCount(topic string, partition int32, trycount int64)
}

type Produce interface {
	MessagesProducedIncrement(topic string, partition int32, count int64)
}

type Consume interface {
	MessagesProcessedIncrement(topic string, partition int32, count int64)
}
