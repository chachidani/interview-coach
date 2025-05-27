package router

import (
	"time"

	bootstrap "github.com/chachidani/interview-coach-backend/Bootstrap"
	"github.com/chachidani/interview-coach-backend/Delivery/controller"
	domain "github.com/chachidani/interview-coach-backend/Domain"
	"github.com/chachidani/interview-coach-backend/Infrastructure/middleware"
	repository "github.com/chachidani/interview-coach-backend/Repository"
	usecases "github.com/chachidani/interview-coach-backend/Usecases"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, r *gin.Engine) {
	// Initialize services
	jwtService := middleware.NewJWTService(env.AccessTokenSecret)
	middleware.SetJWTService(jwtService)
	passwordService := middleware.NewPasswordService()
	geminiRepository := repository.NewGeminiRepository()
	roomRepository := repository.NewRoomRepository(db, domain.CollectionRoom, geminiRepository)

	publicRouter := r.Group("/api/v1")
	NewSignUpRoutes(publicRouter, env, timeout, db, passwordService)
	NewLoginRoutes(publicRouter, env, timeout, db, passwordService, jwtService)

	protectedRouter := r.Group("/user/me")
	protectedRouter.Use(middleware.AuthMiddleware())
	NewRoomRoutes(protectedRouter, env, timeout, db, geminiRepository)

	NewOverallFeedbackRoutes(protectedRouter, env, timeout, &db, geminiRepository, roomRepository)
}

func NewSignUpRoutes(router *gin.RouterGroup, env *bootstrap.Env, timeout time.Duration, db mongo.Database, passwordService *middleware.PasswordService) {
	sr := repository.NewSignUpRepository(db, domain.CollectionUser, passwordService)
	sc := &controller.SignUpController{
		SignUpUsecase: usecases.NewSignUpUsecase(sr, timeout),
	}
	router.POST("/signup", sc.SignUp)
	router.GET("/users", sc.GetUser)
}

func NewLoginRoutes(router *gin.RouterGroup, env *bootstrap.Env, timeout time.Duration, db mongo.Database, passwordService *middleware.PasswordService, jwtService *middleware.JWTService) {
	lr := repository.NewLoginRepository(db, domain.CollectionUser, passwordService, jwtService)
	lc := &controller.LoginController{
		LoginUsecase: usecases.NewLoginUsecase(lr, timeout),
	}
	router.POST("/login", lc.Login)
}

func NewRoomRoutes(router *gin.RouterGroup, env *bootstrap.Env, timeout time.Duration, db mongo.Database, geminiRepository domain.GeminiRepository) {
	rr := repository.NewRoomRepository(db, domain.CollectionRoom, geminiRepository)
	rc := &controller.RoomController{
		RoomUsecase: usecases.NewRoomUsecase(rr, timeout),
	}
	router.POST("/rooms", rc.CreateRoom)
	router.GET("/rooms/:id", rc.GetRoom)
	router.GET("/rooms/user/:id", rc.GetRoomsWithUserID)
	router.PUT("/rooms/:id", rc.UpdateRoom)
	router.DELETE("/rooms/:id", rc.DeleteRoom)
	router.POST("/rooms/:id/messages", rc.AddMessageToRoom)
}

func NewOverallFeedbackRoutes(router *gin.RouterGroup, env *bootstrap.Env, timeout time.Duration, db *mongo.Database, geminiRepository domain.GeminiRepository, roomRepository domain.RoomRepository) {
	rr := repository.NewOverallFeedbackRepository(db, domain.CollectionOverallFeedback, geminiRepository, roomRepository)
	rc := &controller.OverallFeedbackController{
		OverallFeedbackUsecase: usecases.NewOverallFeedbackUsecase(rr, timeout),
	}
	router.POST("/overall-feedbacks", rc.CreateOverallFeedback)
	router.GET("/overall-feedbacks/user/:user_id", rc.GetOverallFeedback)
}
