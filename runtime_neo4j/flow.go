package runtime_neo4j

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/logger"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InsertConstraint(conn Neo4jConnection) {
	ctx := context.Background()

	session := conn.Session(ctx)

	logger.Info("inserting")

	session.ExecuteWrite(
		ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			transaction.Run(
				ctx,
				"CREATE CONSTRAINT FOR (n:Topic) REQUIRE n.name IS UNIQUE",
				map[string]any{},
			)
			return nil, nil
		},
	)

	session.ExecuteWrite(
		ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			transaction.Run(
				ctx,
				"CREATE CONSTRAINT FOR (n:Table) REQUIRE n.name IS UNIQUE",
				map[string]any{},
			)
			return nil, nil
		},
	)
}

func InsertTopic[IK any, IV any](conn Neo4jConnection, topic flow.Topic[IK, IV]) {
	ctx := context.Background()

	session := conn.Session(ctx)

	logger.Info("inserting topic")

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		createTopic, err := transaction.Run(
			ctx,
			"MERGE (topic:Topic {name: $topic}) RETURN topic.name",
			map[string]any{
				"topic": topic.Name(),
			},
		)
		if err != nil {
			logger.ErrorErr("error creating input topic", err)
		}
		if createTopic.Next(ctx) {
			logger.Info(createTopic.Record().Values[0].(string))
		}

		return nil, nil
	})

	if err != nil {
		logger.ErrorErr("failed to insert topic", err)
	}
}

func InsertRelationStateless[IK any, IV any, OK any, OV any](conn Neo4jConnection, name string, input flow.Topic[IK, IV], output flow.Topic[OK, OV]) {
	ctx := context.Background()

	session := conn.Session(ctx)

	logger.Info("inserting stateless relation")

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		createRelation, err := transaction.Run(
			ctx,
			"MATCH (input:Topic {name: $input}), (output:Topic {name: $output}) MERGE (input)-[r:STATELESS {name: $name}]->(output) RETURN r.name",
			map[string]any{
				"input":  input.Name(),
				"output": output.Name(),
				"name":   name,
			},
		)
		if err != nil {
			logger.ErrorErr("error creating relationship", err)
		}
		if createRelation.Next(ctx) {
			logger.Info(createRelation.Record().Values[0].(string))
		}

		return nil, nil
	})

	if err != nil {
		logger.ErrorErr("failed to submit relation information", err)
	}
}

func InsertRelationStateful[IK any, IV any, OK any, OV any](conn Neo4jConnection, name string, input flow.Topic[IK, IV], output flow.Topic[OK, OV]) {
	ctx := context.Background()

	session := conn.Session(ctx)

	logger.Info("inserting stateless relation")

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		createRelation, err := transaction.Run(
			ctx,
			"MATCH (input:Topic {name: $input}), (output:Topic {name: $output}) MERGE (input)-[r:STATEFUL {name: $name}]->(output) RETURN r.name",
			map[string]any{
				"input":  input.Name(),
				"output": output.Name(),
				"name":   name,
			},
		)
		if err != nil {
			logger.ErrorErr("error creating relationship", err)
		}
		if createRelation.Next(ctx) {
			logger.Info(createRelation.Record().Values[0].(string))
		}

		return nil, nil
	})

	if err != nil {
		logger.ErrorErr("failed to submit relation information", err)
	}
}

func InsertTable(conn Neo4jConnection, table string) {
	ctx := context.Background()

	session := conn.Session(ctx)

	logger.Info("inserting topic")

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		createTopic, err := transaction.Run(
			ctx,
			"MERGE (table:Table {name: $table}) RETURN table.name",
			map[string]any{
				"table": table,
			},
		)
		if err != nil {
			logger.ErrorErr("error creating table entry", err)
		}
		if createTopic.Next(ctx) {
			logger.Info(createTopic.Record().Values[0].(string))
		}

		return nil, nil
	})

	if err != nil {
		logger.ErrorErr("failed to insert topic", err)
	}
}
func InsertRelationMaterialise[IK any, IV any](conn Neo4jConnection, name string, input flow.Topic[IK, IV], table string) {
	ctx := context.Background()

	session := conn.Session(ctx)

	logger.Info("inserting stateless relation")

	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		createRelation, err := transaction.Run(
			ctx,
			"MATCH (input:Topic {name: $input}), (table:Table {name: $table}) MERGE (input)-[r:MATERIALISE {name: $name}]->(table) RETURN r.name",
			map[string]any{
				"input": input.Name(),
				"table": table,
				"name":  name,
			},
		)
		if err != nil {
			logger.ErrorErr("error creating relationship", err)
		}
		if createRelation.Next(ctx) {
			logger.Info(createRelation.Record().Values[0].(string))
		}

		return nil, nil
	})

	if err != nil {
		logger.ErrorErr("failed to submit relation information", err)
	}
}
