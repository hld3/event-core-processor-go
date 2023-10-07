package services

import (
	"github.com/hld3/event-common-go/events"
	"github.com/hld3/event-core-processor-go/database"
)

var Conn *database.MongoConnection

type Executor interface {
	Execute(event events.UserDataEvent)
}

func SetMongoConnection(conn *database.MongoConnection) {
    Conn = conn
}
