package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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

	if err := m.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %v", err)
	}

	log.Println("connected to mongodb: ptbot")
	return m, nil
}

func (m *Mongo) createIndexes(ctx context.Context) error {
	usersCol := m.DB.Collection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "tg_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := usersCol.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	log.Println("indexes created successfully")
	return nil
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
