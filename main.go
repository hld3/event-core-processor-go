package main

import (
	"github.com/hld3/event-core-processor-go/database"
	"github.com/hld3/event-core-processor-go/services"
)

func main() {
    conn := database.Connect("testdatabase")
    event := services.UserDataEvent{NodeId: "1", UserId: "2", Username: "Name"}
    services.ProcessUserDataEvent(event, conn)
}
