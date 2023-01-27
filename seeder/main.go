/*
* Generate fake tasks you inject into rabbitmq for your Consumer counterpart
*/

package main

import (
	"log"
	"math/rand"
	"time"
	"os"
	"encoding/json"
	"fmt"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	amq_host = os.Getenv("AMQ_HOST")
	amq_port = os.Getenv("AMQ_PORT")
	amq_username = os.Getenv("AMQ_USERNAME")
	amq_password = os.Getenv("AMQ_PASSWORD")
	amq_task_queue = os.Getenv("AMQ_TASK_QUEUE")
	amq_output_queue = os.Getenv("AMQ_OUTPUT_QUEUE")
)

func main() {

	connString := fmt.Sprintf(
		`amqp://%s:%s@%s:%s`,
		amq_username,
		amq_password,
		amq_host,
		amq_port,
	)
	log.Println(connString)

	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to rabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// define a queue durable makes sure that even if the
	// rabbitmq process dies, the messages are still persisted to disk
	// this option should also be enforced on the consumer code.
	q, err := ch.QueueDeclare(
		amq_output_queue,
		true, // durable
		false,
		false,
		false,
		nil,
	)
	failOnError(err, fmt.Sprintf(`Could not Bind to %s`, q.Name))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {

		couple := dated()
		payload, err := json.Marshal(couple)
		failOnError(err, "unable to Marshal dated couple")
		log.Println(payload)
		log.Println(couple)

		var test [2]string;
		json.Unmarshal(payload, &test)
		log.Println(test)

		// Publish the people that dated on the Queue
		// the messages also need to be marked as persistence
		// to be eligible for storage on disk
		err = ch.PublishWithContext(ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Body: []byte(payload),
			})
		time.Sleep(5*time.Second)
	}

}

func dated() ([2]string) {
	var females = []string{
		"anna",
		"jessica",
		"susan",
		"chelsea",
		"rebecca",
		"bethanie",
		"chloe",
		"samantha",
		"claudia",
		"erica",
		"alima",
		"aissatou",
		"franca",
		"debora",
	}
	var males = []string{
		"mathew",
		"jeremy",
		"sam",
		"frank",
		"joshua",
		"stephane",
		"aristide",
		"mickey",
		"mike",
		"paul",
		"hugo",
		"peter",
		"claude",
		"leonard",
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	m := r.Intn(len(males))
	f := r.Intn(len(females))
	return [2]string{males[m], females[f]}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}