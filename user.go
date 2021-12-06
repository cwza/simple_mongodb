package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Name string
	Age  int
}

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(client *mongo.Client) *UserRepo {
	return &UserRepo{client.Database("testdb").Collection("user")}
}

func (o *UserRepo) CreateUsers(cnt int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := o.collection.Drop(ctx)
	if err != nil {
		return err
	}

	docs := make([]interface{}, cnt)
	for i := 0; i < cnt; i++ {
		user := User{Name: fmt.Sprintf("user%d", i), Age: 28}
		docs[i] = user
	}
	_, err = o.collection.InsertMany(ctx, docs)
	return err
}

func (o *UserRepo) GetUsers(timeout int) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cur, err := o.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	users := make([]User, 0)
	for cur.Next(ctx) {
		result := User{}
		// var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		users = append(users, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
