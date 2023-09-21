package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	db      *mongo.Database
	Client  *mongo.Client
	Context context.Context
	Cancel  context.CancelFunc
}

func Connect(dbName string) *MongoConnection {
	conn := &MongoConnection{}
	uri := "mongodb://127.0.0.1:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(conn.Context, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("There was an error connecting to the database:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("No connection to the database found:", err)
	}

	conn.Context = ctx
	conn.Cancel = cancel
	conn.Client = client
	conn.db = client.Database(dbName)

	log.Println("Successfully connected to MongoDB")

	return conn
}

func CloseConnection(connection *MongoConnection) {
	connection.Cancel()
	if err := connection.Client.Disconnect(connection.Context); err != nil {
		log.Fatal("Error disconnecting:", err)
	}
	log.Println("Connection to MongoDB closed")
}
