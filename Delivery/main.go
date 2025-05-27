package main

import (
	"context"
	"fmt"
	"log"
	"time"

	bootstrap "github.com/chachidani/interview-coach-backend/Bootstrap"
	"github.com/chachidani/interview-coach-backend/Delivery/router"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	env := bootstrap.NewEnv()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := env.DBUri

	if mongoURI == "" {
		mongoURI = "mongodb+srv://nardos_user:Dear@1234@cluster0.k7yt2ba.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
		fmt.Println("DBUri environment variable is not set")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Ping MongoDB to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB")
	// Get database instance with default name if not set
	dbName := env.DBName
	if dbName == "" {
		dbName = "interview_coach"
	}
	db := client.Database(dbName)

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Setup routes
	router.Setup(env, env.ContextTimeout, *db, r)

	// Start server
	serverPort := env.ServerPort
	if serverPort == "" {
		serverPort = "8080"
	}

	fmt.Printf("Server is running on port %s\n", serverPort)
	if err := r.Run(":" + serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}


