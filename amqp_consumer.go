package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func startAMQPConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		"img_task_queue",   // queue
		"imaginary_worker", // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	failOnError(err, "Failed to register a consumer")

	// imgTaskQueueMSG := make(chan Task)
	// imgResultQueueMSG := make(chan ImageResultQueueMSG)
	forever := make(chan bool)

	go func() {
		for d := range msgs {

			task, error := readTask(d.Body)
			imgResult, _ := RunProcess(task)

			if error != nil {
				imgTaskMsg, _ := json.Marshal(task)
				SendAMQPMsg("img_task_queue", imgTaskMsg)
			}

			if task.ID == imgResult.ID {
				imgResultMsg, _ := json.Marshal(imgResult)
				SendAMQPMsg("img_result_queue", imgResultMsg)
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
