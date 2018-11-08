package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/pomfcrypt/pomfcrypt-backend/docs"
	"github.com/pomfcrypt/pomfcrypt-backend/routes"
	"github.com/spf13/viper"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// PomfCrypt Backend Open API specification information
// @title PomfCrypt Backend
// @version 0.1
// @description PomfCrypt is a service which offers simple encrypted file uploading
// @termsOfService https://github.com/pomfcrypt/pomfcrypt
// @contact.name Daniel Malik
// @contact.url https://github.com/pomfcrypt/pomfcrypt
// @contact.email mail@fronbasal.de
// @license.name MIT License
// @host localhost:3000
// @BasePath /api/v1

// PomfEngine is a container struct which holds the API routing engine and the API route controller
type PomfEngine struct {
	GinEngine  *gin.Engine
	Controller *routes.Controller
}

func main() {
	// Initialize a Gin Engine (web server)
	engine := PomfEngine{GinEngine: gin.Default(), Controller: routes.NewController(&routes.Settings{MaxSize: 256000000, FilenameLength: 4, UploadsDirectory: "uploads", Salt: "$1salt$!"})}

	// Initialize viper (configuration management)
	viper.SetConfigName("config")
	// Use yaml as the configuration language
	viper.SetConfigType("yaml")
	// Search for configuration in the following paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/pomfcrypt/")
	// Automatically bind environment variables to the configuration
	viper.AutomaticEnv()

	// Respond with the generated swagger documentation
	engine.GinEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Provide information about the project as index
	engine.GinEngine.GET("/", func(c *gin.Context) { c.String(200, "https://github.com/pomfcrypt/pomfcrypt-backend") })

	// Create the API group
	v1 := engine.GinEngine.Group("/api/v1")
	{
		dataApi := v1.Group("/data")
		{
			dataApi.PUT("", engine.Controller.Upload)
		}
	}

	// Run the server on the http port 3000
	engine.GinEngine.Run(":3000")
}
