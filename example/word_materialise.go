package example

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/uptrace/bun"
)

// CREATE TABLE IF NOT EXISTS public.flows_materialised
// (
//	id                   VARCHAR(255)  NOT NULL,
//	key_content          VARCHAR(255)  NULL,
//	value_content        VARCHAR(255)  NULL,
//	timestamp_ms         BIGINT,
//	PRIMARY KEY(id)
// );

type FlowsMaterialised struct {
	bun.BaseModel `bun:"table:flows_materialised"`
	Id            string `bun:",pk"`
	KeyContent    string
	ValueContent  string
	TimestampMs   int64
}

func FlowsMaterialisedMap(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]FlowsMaterialised, error) {
	return []FlowsMaterialised{
		{
			Id:           string(m.Key),
			KeyContent:   string(m.Key),
			ValueContent: string(m.Value),
			TimestampMs:  m.Timestamp.UnixMilli(),
		},
	}, nil
}

func WordMaterialise() runtime.Runtime {
	materialiseConfiguration := flows.MaterialisePostgresqlFunctionConfiguration[FlowsMaterialised]{
		PostgresqlConfiguration: []runtime.Configuration[*runtime_bun.PostgresqlConnection]{
			runtime_bun.WithApplicationName("flows"),
			runtime_bun.WithConnectionString("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		},
		KafkaConsumerConfiguration: []runtime.Configuration[*runtime_sarama.Consumer]{
			runtime_sarama.WithConsumerSaramaConfig(runtime_sarama.DefaultConfiguration()),
			runtime_sarama.WithConsumerBroker("localhost:9092"),
			runtime_sarama.WithConsumerTopic("word-count"),
			runtime_sarama.WithConsumerGroupName("flows-word-materialise"),
		},
		MaterialiseMapFunction: FlowsMaterialisedMap,
	}

	return materialiseConfiguration.Runtime()
}
