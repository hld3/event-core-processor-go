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
	services.SetMongoConnection(conn)

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
			log.Println("Processing event of type:", eventType)
			eventRouter(eventType, msg)
		}
	}()

	log.Println(" [*] Waiting for message. CTRL+C to exit.")
	<-forever
}

func eventRouter(eventType string, msg amqp.Delivery) {
	var message events.BaseEvent
	unmarshalBaseEvent([]byte(msg.Body), &message)
	// Was an attempt at making the events more specific to use "bounded generics". Doesn't work but I am keeping for now as a reference.
	// convert the base event to the BaseEvent struct
	// err := json.Unmarshal([]byte(msg.Body), &message)
	// if err != nil {
	// 	log.Println("Could not unmarshal rmq message:", err)
	// }

	switch eventType {
	case "UserDataEvent":
		// Was an attempt at making the events more specific to use "bounded generics". Doesn't work but I am keeping for now as a reference.
		// take the payload from BaseEvent and convert it to bytes to the specific event type.
		// eventPayload := map[string]interface{}{}
		// if err := json.Unmarshal([]byte(msg.Body), &eventPayload); err == nil {
		// 	var payload events.UserDataEvent
		// 	convertPayload(eventPayload, &payload)
		// 	services.ProcessUserDataEvent(payload)
		// }
		rawPayload, ok := message.Payload.(RawPayload)
		if !ok {
			log.Println("Failed to cast Payload:")
			return
		}
		var payload events.UserDataEvent
		if err := json.Unmarshal(rawPayload.Data, &payload); err != nil {
			log.Println("Error unmarshalling UserDataEvent:", err)
		} else {
			services.ProcessUserDataEvent(payload)
		}
	case "GroupDataEvent":
		// Was an attempt at making the events more specific to use "bounded generics". Doesn't work but I am keeping for now as a reference.
		// take the payload from BaseEvent and convert it to bytes to the specific event type.
		// eventPayload := map[string]interface{}{}
		// if err := json.Unmarshal([]byte(msg.Body), &eventPayload); err == nil {
		// 	var payload events.GroupDataEvent
		// 	convertPayload(eventPayload, &payload)
		// 	services.ProcessGroupDataEvent(payload)
		// }
		rawPayload, ok := message.Payload.(RawPayload)
		if !ok {
			log.Println("Failed to cast Payload:")
			return
		}
		var payload events.GroupDataEvent
		if err := json.Unmarshal(rawPayload.Data, &payload); err != nil {
			log.Println("Error unmarshalling UserDataEvent:", err)
		} else {
			services.ProcessGroupDataEvent(payload)
		}
	}
}

// Was an attempt at making the events more specific to use "bounded generics". Doesn't work but I am keeping for now as a reference.
// func convertPayload[T events.Payload](payload interface{}, to *T) {
// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		log.Println("Error converting payload to bytes:", err)
// 	}
// 	// convert the payload bytes to the GroupDataEvent struct
// 	err = json.Unmarshal(payloadBytes, &to)
// 	if err != nil {
// 		log.Println("Error converting payload to struct:", err)
// 	}
// }

func unmarshalBaseEvent(data []byte, event *events.BaseEvent) error {
	type helper struct {
		MessageId string          `json:"messageId"`
		DateCode  string          `json:"dateCode"`
		Payload   json.RawMessage `json:"payload"`
	}

	var h helper
	if err := json.Unmarshal(data, &h); err != nil {
		return err
	}

	event.MessageId = h.MessageId
	event.DateCode = h.DateCode
	event.Payload = RawPayload{Data: h.Payload}

	return nil
}

type RawPayload struct {
	Data json.RawMessage
}

func (r RawPayload) IsEventPayload() {}
