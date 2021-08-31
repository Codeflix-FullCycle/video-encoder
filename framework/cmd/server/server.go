package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Codeflix-FullCycle/encoder/application/services"
	"github.com/Codeflix-FullCycle/encoder/framework/database"
	"github.com/Codeflix-FullCycle/encoder/framework/queue"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db database.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("Error parsing boolean env var")
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("Error parsing boolean env var")
	}

	db.AutoMigrateDb = autoMigrateDb
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")
}

func main() {
	messageChannel := make(chan amqp.Delivery)
	resultChannel := make(chan services.JobWorkerResult)

	dbConnect, err := db.Connect()

	if err != nil {
		log.Fatalln(err)
	}

	defer dbConnect.Close()

	rabbitmq := queue.NewRabbitMQ()
	ch := rabbitmq.Connect()

	rabbitmq.Consume(messageChannel)

	jobManager := services.NewJobManager(dbConnect, rabbitmq, messageChannel, resultChannel)

	jobManager.Start(ch)
}
