package services

import (
	"log"

	"github.com/hld3/event-common-go/events"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	GROUP_DATA_COLLECTION = "currentGroupState"
)

// Change to use Execute interface when I add more services.
func ProcessGroupDataEvent(event events.GroupDataEvent) {
	log.Println("Processing the group data event:", event)
	groupData := bson.M{
		"name":             event.Name,
		"code":             event.Code,
		"groupId":          event.GroupId,
		"ownerId":          event.OwnerId,
		"knownLanguage":    event.KnownLanguage,
		"learningLanguage": event.LearningLanguage,
	}

	collection := Conn.DB.Collection(GROUP_DATA_COLLECTION)
	_, err := collection.InsertOne(Conn.Context, groupData)
	if err != nil {
		log.Println("Error saving the group data:", err)
		return
	}
	log.Println("Group data successfully added to the database:", event.Name, event.GroupId)
}
