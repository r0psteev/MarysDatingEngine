/*
* Consumes couples that dated from it's task queue
* and plots neo4j relations graphs between these couples.
*/

package main

import (
	"log"
	"os"
	"fmt"
	"time"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (

	// how to talk to rabbitmq
	amq_host = os.Getenv("AMQ_HOST")
	amq_port = os.Getenv("AMQ_PORT")
	amq_username = os.Getenv("AMQ_USERNAME")
	amq_password = os.Getenv("AMQ_PASSWORD")
	amq_task_queue = os.Getenv("AMQ_TASK_QUEUE")
	amq_output_queue = os.Getenv("AMQ_OUTPUT_QUEUE")

	/* how to talk to neo4j */
	neo4j_host= os.Getenv("NEO4J_HOST")
	neo4j_port = os.Getenv("NEO4J_PORT")
	neo4j_username = os.Getenv("NEO4J_USERNAME")
	neo4j_password = os.Getenv("NEO4J_PASSWORD")
)

func main() {

	// fetch messages from rabbitmq
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
		amq_task_queue,
		true, // durable
		false,
		false,
		false,
		nil,
	)
	failOnError(err, fmt.Sprintf(`Could not Bind to %s`, q.Name))

	err = ch.Qos(
		1, // prefetch count
		0, // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	// msgs has type chan Delivery
	msgs, err := ch.Consume(
			q.Name, // queue name
			"",     // consumer
			false,   // auto-ack: No auto-ack, i manually ack this time
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// So this channel acts as blocker which keeps
	// this main thread busy ?
	var forever chan struct{}

	go func() {
			for d := range msgs {
					//log.Printf("Received a message: %s", d.Body)
					//log.Printf("Message has content type: %s", d.ContentType)
					err = ProcessQueueObject(d.Body, ch)
					failOnError(err, fmt.Sprintf(`Failed to process queue object %s`, d.Body))
					time.Sleep(2 * time.Second)
					// this is acking one by one, but
					// batch acking is also possible
					d.Ack(false)
			}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}