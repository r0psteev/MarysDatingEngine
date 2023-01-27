package main

import (
	"log"
	"fmt"
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	amqp "github.com/rabbitmq/amqp091-go"
)


func ProcessQueueObject(payload []byte, ch *amqp.Channel) (error) {
	var couple[2]string
	err := json.Unmarshal(payload, &couple)
	failOnError(err, fmt.Sprintf("Failed Unmarshalling %v", payload))
	log.Printf(fmt.Sprintf("[+] Processing queueobj %v", couple))
	err = pushRelation2Neo4j(couple[0], couple[1])
	return err
}

func pushRelation2Neo4j(person1 string, person2 string) (error) {
	// let's first find the alreay existing instance of person1
	query := fmt.Sprintf("MERGE (n: Male {name: '%s'})\n", person1)
	// Then let's first find the alreay existing instance of person2
	query += fmt.Sprintf("MERGE (m: Female {name: '%s'})\n", person2)
	// Then created relationship between them
	query += fmt.Sprintf("MERGE (n) -[:dated]- (m)")
	ctx := context.Background()
	uri := fmt.Sprintf(
		`neo4j://%s:%s`,
		neo4j_host,
		neo4j_port,
	)
	_, err := push2Neo4j(ctx, uri, neo4j_username, neo4j_password, query)
	return err
}

func push2Neo4j(
	ctx context.Context,
	uri,
	username,
	password,
	query string,
) (interface{}, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
			return "", err
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	greeting, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx, query, nil)
			if err != nil {
					return nil, err
			}

			if result.Next(ctx) {
					return result.Record().Values[0], nil
			}

			return nil, result.Err()
	})
	if err != nil {
			return "", err
	}
	return greeting, nil
}