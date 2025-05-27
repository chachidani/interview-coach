package handler


import (
	"context"
	"net/http"
	"time"

	bootstrap "github.com/chachidani/interview-coach-backend/Bootstrap"
	"github.com/chachidani/interview-coach-backend/Delivery/router"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	env    *bootstrap.Env
	client *mongo.Client
	db     *mongo.Database
)

func init() {
	env = bootstrap.NewEnv()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := env.DBUri
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	// Ping MongoDB to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	// Get database instance
	dbName := env.DBName
	if dbName == "" {
		dbName = "interview_coach"
	}
	db = client.Database(dbName)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	// Setup routes using your router package
	router.Setup(env, env.ContextTimeout, *db, engine)

	// Serve the request
	engine.ServeHTTP(w, r)
}
