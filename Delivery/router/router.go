package router

import (
	"time"

	bootstrap "github.com/chachidani/interview-coach-backend/Bootstrap"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, r *gin.Engine) {
	
}
