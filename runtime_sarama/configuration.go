package runtime_sarama

import "github.com/Shopify/sarama"

func DefaultConfiguration() *sarama.Config {

	saramaconfig := sarama.NewConfig()

	saramaconfig.Net.MaxOpenRequests = 100

	saramaconfig.Consumer.IsolationLevel = sarama.ReadCommitted
	saramaconfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	saramaconfig.Consumer.Offsets.AutoCommit.Enable = true
	saramaconfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	saramaconfig.Producer.Idempotent = false
	saramaconfig.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	saramaconfig.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	saramaconfig.Producer.Return.Successes = true
	saramaconfig.Producer.Partitioner = sarama.NewReferenceHashPartitioner

	return saramaconfig
}
