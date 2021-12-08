package main

import (
	"flag"
	"log"
	"time"

	"github.com/cwza/simple_mongodb/mongodb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "cfgpath", "./config.toml", "config file path")
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

	// create mongodb client
	client, closeFunc, err := mongodb.InitClient(config.ConsumerUrl)
	if err != nil {
		log.Fatalf("failed to init client, %s", err)
	}
	defer closeFunc()

	// drop users
	userRepo := mongodb.NewUserRepo(client)
	err = userRepo.DropUsers()
	if err != nil {
		log.Fatalf("failed to drop users, %s", err)
	}

	// shard testdb.user
	if config.Mode == Shard {
		adminRepo := mongodb.NewAdminRepo(client)
		err = adminRepo.EnableSharding("testdb")
		if err != nil {
			log.Fatalf("failed to enable sharding on testdb, %s", err)
		}
		err = adminRepo.ShardCollection("testdb", "user")
		if err != nil {
			log.Fatalf("failed to shard on testdb.user, %s", err)
		}
	}

	// create users
	err = userRepo.CreateUsers(100)
	if err != nil {
		log.Fatalf("failed to create create users, %s", err)
	}

	// get users
	users, err := userRepo.GetUsers(30)
	if err != nil {
		log.Fatalf("failed to get users, %s", err)
	}
	log.Printf("users: %v", users)

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
