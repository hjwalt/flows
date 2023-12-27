package runtime_neo4j

import (
	"context"

	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

// constructor
func NewBasicAuth(configurations ...runtime.Configuration[*Neo4JConnectionBasicAuth]) runtime.Runtime {
	connection := &Neo4JConnectionBasicAuth{}
	for _, configuration := range configurations {
		connection = configuration(connection)
	}
	return connection
}

func WithUrl(url string) runtime.Configuration[*Neo4JConnectionBasicAuth] {
	return func(c *Neo4JConnectionBasicAuth) *Neo4JConnectionBasicAuth {
		c.url = url
		return c
	}
}

func WithUser(user string) runtime.Configuration[*Neo4JConnectionBasicAuth] {
	return func(c *Neo4JConnectionBasicAuth) *Neo4JConnectionBasicAuth {
		c.user = user
		return c
	}
}

func WithPass(pass string) runtime.Configuration[*Neo4JConnectionBasicAuth] {
	return func(c *Neo4JConnectionBasicAuth) *Neo4JConnectionBasicAuth {
		c.pass = pass
		return c
	}
}

func WithRealm(realm string) runtime.Configuration[*Neo4JConnectionBasicAuth] {
	return func(c *Neo4JConnectionBasicAuth) *Neo4JConnectionBasicAuth {
		c.realm = realm
		return c
	}
}

type Neo4JConnectionBasicAuth struct {
	url    string
	user   string
	pass   string
	realm  string
	driver neo4j.DriverWithContext
}

func (conn *Neo4JConnectionBasicAuth) Start() error {
	logger.Info("starting neo4j")

	driver, err := neo4j.NewDriverWithContext(conn.url, neo4j.BasicAuth(conn.user, conn.pass, ""))
	if err != nil {
		logger.ErrorErr("failed to start neo4j", err)
		return err
	}

	info, err := driver.GetServerInfo(context.Background())
	if err != nil {
		logger.ErrorErr("failed to start neo4j", err)
		return err
	}

	logger.Info("neo4j info", zap.Any("info", info))

	conn.driver = driver
	return nil
}

func (conn *Neo4JConnectionBasicAuth) Stop() {
	logger.Info("stopping neo4j")
	conn.driver.Close(context.Background())
}

func (conn *Neo4JConnectionBasicAuth) Driver() neo4j.DriverWithContext {
	return conn.driver
}

func (conn *Neo4JConnectionBasicAuth) Session(ctx context.Context) neo4j.SessionWithContext {
	return conn.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
}
