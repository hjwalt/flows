package runtime_neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jConnection interface {
	Start() error
	Stop()
	Driver() neo4j.DriverWithContext
	Session(context.Context) neo4j.SessionWithContext
}
