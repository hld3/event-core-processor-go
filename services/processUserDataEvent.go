package services

import (
	"log"

	"github.com/hld3/event-core-processor-go/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	USER_DATA_COLLECTION  = "userData"
	NODE_COUNT_COLLECTION = "nodeCount"
)

var conn *database.MongoConnection

type UserDataEvent struct {
	NodeId   string `json:"nodeId"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

type Executor interface {
	Execute(event UserDataEvent)
}

type UserProcessingService struct{}

func (s *UserProcessingService) Execute(event UserDataEvent) {
	log.Println("Running user processing service")
	userData := bson.M{"nodeId": event.NodeId, "username": event.Username, "userId": event.UserId}
	collection := conn.DB.Collection(USER_DATA_COLLECTION)
	_, err := collection.InsertOne(conn.Context, userData)
	if err != nil {
		log.Println("There was an error saving the user data:", err)
		return
	}
	log.Println("User successfully added to the database:", event.UserId)
}

type UserNodeCountService struct{}

func (s *UserNodeCountService) Execute(event UserDataEvent) {
	log.Println("Running user node count service")
	collection := conn.DB.Collection(NODE_COUNT_COLLECTION)

	filter := bson.M{"nodeId": event.NodeId}
	update := bson.M{"$inc": bson.M{"count": 1}}
	updateOptions := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(conn.Context, filter, update, updateOptions)
	if err != nil {
		log.Println("Error updating user node count:", err)
		return
	}
	log.Println("User count updated for node:", event.NodeId)
}

func ProcessUserDataEvent(event UserDataEvent, db *database.MongoConnection) {
	conn = db
	services := []Executor{&UserProcessingService{}, &UserNodeCountService{}}
	for _, service := range services {
		service.Execute(event)
	}
	log.Println("Finished processing user event data")
}
