package main

import (
	"context"
	"flag"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "cfgpath", "./config.toml", "config file path")
}

func initClient(uri string) (*mongo.Client, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(uri),
		options.Client().SetReadPreference(readpref.SecondaryPreferred()),
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

type SendFuncType func() error

func sends(sendFunc SendFuncType, msgCnt int, pool chan int) {
	for i := 0; i < msgCnt; i++ {
		<-pool
		go func() {
			err := sendFunc()
			if err != nil {
				log.Printf("failed to send, %s", err)
			}
			pool <- 1
		}()
	}
}

func run(sendFunc SendFuncType, workerCnt int, genSecRateFunc func() int) {
	pool := make(chan int, workerCnt)
	for i := 0; i < workerCnt; i++ {
		pool <- 1
	}
	for range time.Tick(time.Second) {
		msgCnt := genSecRateFunc()
		sends(sendFunc, msgCnt, pool)
		log.Printf("send %d msgs\n", msgCnt)
	}
}

func main() {
	flag.Parse()

	// config
	config, err := initConfig(configPath)
	if err != nil {
		log.Fatalf("failed to init config, %s", err)
	}
	log.Printf("config: %+v\n", config)

	// client
	client, closeFunc, err := initClient(config.ConsumerUrl)
	if err != nil {
		log.Fatalf("failed to init client, %s", err)
	}
	defer closeFunc()

	// create data
	userRepo := NewUserRepo(client)
	err = userRepo.CreateUsers(100)
	if err != nil {
		log.Fatal(err)
	}

	// send read reqs
	sendFunc := func() error {
		_, err := userRepo.GetUsers(config.Timeout)
		if err != nil {
			log.Printf("failed to get users, %s", err)
		}
		return err
	}
	genSecRateFunc := createGenSecRateFunc(createGenMinRateFunc(config.Rates, config.Cnts))
	run(sendFunc, config.WorkerCnt, genSecRateFunc)
}
