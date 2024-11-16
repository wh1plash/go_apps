package main

import (
	"gin/controllers"
	"gin/initializers"
	"gin/middleware"
	"github.com/gin-gonic/gin"
	"log"
	//"gorm.io/driver/sqlite"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.InitDirs()
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/upload", middleware.RequireAuth, controllers.UploadHandler)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
