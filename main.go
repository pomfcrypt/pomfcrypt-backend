package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type PomfEngine struct {
	GinEngine *gin.Engine
}

func main() {
	// Initialize a Gin Engine (web server)
	engine := PomfEngine{GinEngine: gin.Default()}

	// Initialize viper (configuration management)
	viper.SetConfigName("config")
	// Use yaml as the configuration language
	viper.SetConfigType("yaml")
	// Search for configuration in the following paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/pomfcrypt/")
	// Automatically bind environment variables to the configuration
	viper.AutomaticEnv()

	// Handle index request
	engine.GinEngine.GET("/", func(c *gin.Context) {
		// Respond with a link to the project meta repository
		NewError("Please refer to the repository to read the API documentation: https://github.com/pomfcrypt/pomfcrypt", 404).Throw(c)
	})

	// Create the API group
	api := engine.GinEngine.Group("/api")
	{
		api.POST("/files/list/", func(c *gin.Context) {
		})
	}

	// Run the server on the http port 3000
	engine.GinEngine.Run(":3000")
}
