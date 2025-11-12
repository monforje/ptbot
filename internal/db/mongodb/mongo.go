package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Mongo struct {
	client *mongo.Client
	DB     *mongo.Database
}

func New(dbURL string) (*Mongo, error) {
	m := &Mongo{}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbURL).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	m.client = client

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	m.DB = client.Database("ptbot")

	log.Println("connected to mongodb: ptbot")
	return m, nil
}

func (m *Mongo) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.client.Disconnect(ctx)
	if err != nil {
		return err
	}

	log.Println("mongodb stopped")
	return nil
}
