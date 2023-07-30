package example

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/runtime"
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

func FlowsMaterialisedMap(c context.Context, m message.Message[string, string]) ([]FlowsMaterialised, error) {
	return []FlowsMaterialised{
		{
			Id:           m.Key,
			KeyContent:   m.Key,
			ValueContent: m.Value,
			TimestampMs:  m.Timestamp.UnixMilli(),
		},
	}, nil
}

func WordMaterialise() runtime.Runtime {
	materialiseConfiguration := flows.MaterialisePostgresqlOneToOneFunctionConfiguration[FlowsMaterialised, string, string]{
		Name:                     "flows-word-materialise",
		InputTopic:               topic.String("word-count"),
		Function:                 FlowsMaterialisedMap,
		InputBroker:              "localhost:9092",
		OutputBroker:             "localhost:9092",
		HttpPort:                 8081,
		PostgresConnectionString: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	}

	return materialiseConfiguration.Runtime()
}
