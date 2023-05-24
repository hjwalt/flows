package runtime_sarama_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/test_helper"
	"github.com/hjwalt/runway/logger"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Container struct {
	Container testcontainers.Container
	Context   context.Context
	Endpoint  string
	T         *testing.T
}

func (c Container) Stop() {
	c.Container.Terminate(c.Context)
}

func CreateContainer(t *testing.T) Container {
	assert := assert.New(t)

	addrAdvertised := "localhost"
	// if config.GetEnvBool("CI", false) {
	// 	addrAdvertised = "docker"
	// } else {
	// 	addrAdvertised = "localhost"
	// }

	var containerError error
	containerContext := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "docker.redpanda.com/vectorized/redpanda:latest",
		ExposedPorts: []string{"29092:29092/tcp"},
		Env:          map[string]string{},
		Privileged:   true,
		WaitingFor:   wait.ForAll(wait.ForListeningPort("29092"), wait.ForLog("Successfully started Redpanda")),
		AutoRemove:   true,
		Cmd: []string{
			"redpanda",
			"start",
			"--overprovisioned",
			"--smp",
			"1",
			"--reserve-memory",
			"0M",
			"--node-id",
			"0",
			"--check=false",
			"--kafka-addr",
			"PLAINTEXT://0.0.0.0:29092",
			"--advertise-kafka-addr",
			"PLAINTEXT://" + addrAdvertised + ":29092",
		},
	}
	testContainer, containerError := testcontainers.GenericContainer(containerContext, testcontainers.GenericContainerRequest{
		ProviderType:     testcontainers.ProviderPodman,
		ContainerRequest: req,
		Started:          true,
		Logger:           test_helper.TestContainerLogging{},
	})
	assert.NoError(containerError)

	kafkaEndpoint, endpointErr := testContainer.PortEndpoint(containerContext, "29092", "")
	assert.NoError(endpointErr)

	logger.Info(kafkaEndpoint)

	// util.Log(t, testContainer)

	return Container{
		Container: testContainer,
		Context:   containerContext,
		Endpoint:  kafkaEndpoint,
		T:         t,
	}
}

// actual test with Kafka containers

func TestWithContainer(t *testing.T) {
	assert := assert.New(t)

	container := CreateContainer(t)
	defer container.Stop()

	table := []struct {
		name string
		test func(*testing.T)
	}{
		// Consumer test
		{
			name: "TestConsumerStartProperConfigurationShouldSucceed",
			test: func(t *testing.T) {
				consumer := runtime_sarama.NewConsumer(
					runtime_sarama.WithConsumerBroker(container.Endpoint),
					runtime_sarama.WithConsumerTopic("test-topic"),
					runtime_sarama.WithConsumerSaramaConfig(runtime_sarama.DefaultConfiguration()),
					runtime_sarama.WithConsumerRuntimeController(runtime.NewController()),
					runtime_sarama.WithConsumerLoop(ConsumerLoopForTest{}),
				)

				err := consumer.Start()

				assert.Nil(err)
				assert.Greater(len(consumer.GroupName), 0)

				consumer.Stop()
			},
		},
		{
			name: "TestConsumerStartWithConsumerGroupShouldUseNameAndSucceed",
			test: func(t *testing.T) {
				consumer := runtime_sarama.NewConsumer(
					runtime_sarama.WithConsumerBroker(container.Endpoint),
					runtime_sarama.WithConsumerTopic("test-topic"),
					runtime_sarama.WithConsumerSaramaConfig(runtime_sarama.DefaultConfiguration()),
					runtime_sarama.WithConsumerRuntimeController(runtime.NewController()),
					runtime_sarama.WithConsumerLoop(ConsumerLoopForTest{}),
					runtime_sarama.WithConsumerGroupName("consumer-group"),
				)

				err := consumer.Start()

				assert.Nil(err)
				assert.Equal(consumer.GroupName, "consumer-group")

				consumer.Stop()
			},
		},
		// Producer test
		{
			name: "TestProducerStartWithProducerShouldSucceed",
			test: func(t *testing.T) {
				producer := runtime_sarama.NewProducer(
					runtime_sarama.WithProducerBroker(container.Endpoint),
					runtime_sarama.WithProducerSaramaConfig(runtime_sarama.DefaultConfiguration()),
					runtime_sarama.WithProducerRuntimeController(runtime.NewController()),
				)

				err := producer.Start()

				assert.Nil(err)

				producer.Stop()
			},
		},
	}

	for _, tc := range table {
		t.Run(tc.name, tc.test)
	}
}

func TestCancelContextMultipleTimes(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())

	cancel()
	cancel()
	cancel()
	cancel()
	cancel()
	cancel()
	cancel()
}
