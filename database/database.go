package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongoDB(mongoURL string) (*mongo.Database, *mongo.Client, context.Context, context.CancelFunc) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("No se pudo hacer ping al servidor MongoDB: %v", err)
	}

	db := client.Database("todoer")

	return db, client, ctx, cancel
}

func CloseConnection(client *mongo.Client, context context.Context, cancel context.CancelFunc) {
	defer func() {
		cancel()
		if err := client.Disconnect(context); err != nil {
			panic(err)
		}
	}()
}
