package main

import (
	"context"
	"fmt"
	"log"

	bootstrap "github.com/chachidani/interview-coach-backend/Bootstrap"
	"github.com/chachidani/interview-coach-backend/Delivery/router"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	log.Println("Starting app...")
	env := bootstrap.NewEnv()

	mongoURI := env.DBUri
	if mongoURI == "" {
		mongoURI = "mongodb+srv://nardos_user:nardos_user@cluster0.k7yt2ba.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
		fmt.Println("DBUri environment variable is not set")
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
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
