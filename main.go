package main

import (
	"encoding/json"
	"log"

	events "github.com/hld3/event-common-go/events"
	"github.com/hld3/event-core-processor-go/database"
	"github.com/hld3/event-core-processor-go/receiver"
	"github.com/hld3/event-core-processor-go/services"
	"github.com/streadway/amqp"
)

func main() {
	conn := database.Connect("testdatabase")

	rec, err := receiver.NewReceiver()
	if err != nil {
		log.Fatal("Failed to initialize rmq receiver:", err)
	}
	defer rec.Close()

	msgs, err := rec.StartReceiving()
	if err != nil {
		log.Fatal("Failed to start the receiver:", err)
	}

	// Keep app running... forever
	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			log.Println("Received a message:", string(msg.Body))
			headers := msg.Headers
			eventType, exists := headers["eventType"].(string)
			if !exists {
				log.Println("Missing eventType header, can not route event.")
				continue
			}
			eventRouter(eventType, msg, conn)
		}
	}()

	// event := services.UserDataEvent{NodeId: "1", UserId: "2", Username: "Name"}
	// services.ProcessUserDataEvent(event, conn)
	log.Println(" [*] Waiting for message. CTRL+C to exit.")
	<-forever
}

func eventRouter(eventType string, msg amqp.Delivery, conn *database.MongoConnection) {
	switch eventType {
	case "UserDataEvent":
		var message events.UserDataEvent
		err := json.Unmarshal([]byte(msg.Body), &message)
		if err != nil {
			log.Println("Could not unmarshal rmq message:", err)
		}
		services.ProcessUserDataEvent(message, conn)
	}
}
