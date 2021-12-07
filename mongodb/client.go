package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitClient(uri string) (*mongo.Client, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// uuidCoder := NewUUIDCoder()
	// registry := bson.NewRegistryBuilder().
	// 	RegisterTypeEncoder(uuidCoder.TUUID, bsoncodec.ValueEncoderFunc(uuidCoder.UuidEncodeValue)).
	// 	RegisterTypeDecoder(uuidCoder.TUUID, bsoncodec.ValueDecoderFunc(uuidCoder.UuidDecodeValue)).
	// 	Build()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(uri),
		options.Client().SetReadPreference(readpref.SecondaryPreferred()),
		// options.Client().SetRegistry(registry),
	)
	if err != nil {
		return nil, nil, err
	}
	closeFunc := func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("failed to disconnect, %s", err)
		}
	}
	return client, closeFunc, nil
}

type AdminRepo struct {
	client *mongo.Client
}

func NewAdminRepo(client *mongo.Client) *AdminRepo {
	return &AdminRepo{client}
}

func (o *AdminRepo) EnableSharding(dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := bson.D{
		{Key: "enableSharding", Value: dbName},
	}
	err := o.client.Database("admin").RunCommand(ctx, cmd).Err()
	return err
}

func (o *AdminRepo) ShardCollection(dbName string, colName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := bson.D{
		{"shardCollection", fmt.Sprintf("%s.%s", dbName, colName)},
		{"key", bson.D{{"_id", "hashed"}}},
	}
	err := o.client.Database("admin").RunCommand(ctx, cmd).Err()
	return err
}
