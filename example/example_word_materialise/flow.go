package example_word_materialise

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/logger"
	"github.com/uptrace/bun"
)

const (
	Instance = "flows-word-materialise"
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

func FlowsMaterialisedMap(c context.Context, m flow.Message[string, string]) ([]FlowsMaterialised, error) {
	logger.Info("materialising")
	return []FlowsMaterialised{
		{
			Id:           m.Key,
			KeyContent:   m.Key,
			ValueContent: m.Value,
			TimestampMs:  m.Timestamp.UnixMilli(),
		},
	}, nil
}

func Registrar(ci inverse.Container) flows.Prebuilt {
	return flows.MaterialisePostgresqlOneToOneFunctionConfiguration[FlowsMaterialised, string, string]{
		Name:                     Instance,
		InputTopic:               flow.StringTopic("word-count"),
		Function:                 FlowsMaterialisedMap,
		InputBroker:              "localhost:9092",
		OutputBroker:             "localhost:9092",
		HttpPort:                 8081,
		PostgresConnectionString: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	}
}

func Register(m flows.Main) {
	m.Prebuilt(Instance, Registrar)
}
