package config

import (
	"social-network-api/internal/controllers"
	"social-network-api/internal/db"
	"social-network-api/internal/models"
	"social-network-api/internal/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupApi(config models.Config) *gin.Engine {
	r := gin.Default()

	// config.AllowOrigins = []string{"http://localhost:3000"} // Reemplaza esto con el origen real de tu aplicación cliente
	// config.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, "OPTIONS"}
	// config.AllowHeaders = []string{"Authorization", "Content-Type"}

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Authorization", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))

	db := db.ConnectDatabase(config)

	jwtService := services.NewJwtService(config.SecretKey)
	authService := services.NewAuthService(db, jwtService)
	postService := services.NewPostService(db)
	profileService := services.NewProfileService(db)

	authRoutes(r, authService, jwtService)
	postRoutes(r, postService, jwtService)
	profileRoutes(r, profileService, jwtService)

	return r
}

func authRoutes(r *gin.Engine, authService *services.AuthService, jwtService *services.JWTService) {
	controller := controllers.NewUserController(authService, jwtService)
	auth := r.Group("/auth")

	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Register)
	auth.POST("/renew", controller.Renew)
}

func postRoutes(r *gin.Engine, postService *services.PostService, jwtService *services.JWTService) {
	controller := controllers.NewPostController(postService, jwtService)
	post := r.Group("/post")

	post.POST("", controller.NewPost)
	post.GET("", controller.GetPost)
	post.GET("/user/:id", controller.GetPostsByUser)
	post.POST("/message/:id", controller.MessagePost)
	post.POST("/reaction/:id", controller.ReactionPost)
}

func profileRoutes(r *gin.Engine, profileService *services.ProfileService, jwtService *services.JWTService) {
	controller := controllers.NewProfileController(profileService, jwtService)
	post := r.Group("/profile")

	post.GET("", controller.GetProfiles)
	post.GET("/:id", controller.GetByID)
	post.GET("/followed/:id", controller.GetByFollowed)
	post.GET("/follower/:id", controller.GetByFollower)
	post.PUT("", controller.EditProfile)
	post.POST("/followUnFollow/:id", controller.FollowUnFollow)
}
