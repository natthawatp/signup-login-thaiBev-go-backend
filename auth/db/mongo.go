package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-auth-backend/auth/config"
)

var Client *mongo.Client
var Database *mongo.Database

func ConnectMongo(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.DBHost, cfg.DBPort))
	clientOptions.SetAuth(options.Credential{
		Username: cfg.DBUser,
		Password: cfg.DBPassword,
	})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ Cannot connect to MongoDB: ", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB ping failed: ", err)
	}

	fmt.Println("✅ Connected to MongoDB")
	Client = client
	Database = client.Database(cfg.MongoDB)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB ping failed: ", err)
	}

	fmt.Println("✅ Connected to MongoDB")

	Client = client
	Database = client.Database(cfg.MongoDB)
}
