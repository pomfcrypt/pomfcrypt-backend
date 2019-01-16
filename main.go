package main

import (
	"github.com/kataras/iris"
	"github.com/pomfcrypt/pomfcrypt-backend/routes"
	"github.com/spf13/viper"
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
	IrisEngine *iris.Application
	Controller *routes.Controller
}

func main() {
	// Initialize a Gin Engine (web server)
	engine := PomfEngine{IrisEngine: iris.Default(), Controller: routes.NewController(&routes.Settings{MaxSize: 256000000, FilenameLength: 4, UploadsDirectory: "uploads", Salt: "$1salt$!"})}

	// Initialize viper (configuration management)
	viper.SetConfigName("config")
	// Use yaml as the configuration language
	viper.SetConfigType("yaml")
	// Search for configuration in the following paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/pomfcrypt/")
	// Automatically bind environment variables to the configuration
	viper.AutomaticEnv()

	// Provide information about the project as index
	engine.IrisEngine.Get("/", func(c iris.Context) { c.JSON("https://github.com/pomfcrypt/pomfcrypt-backend") })

	// Create the API group ("party")
	api := engine.IrisEngine.Party("/api/v1/")

	api.Put("/file", engine.Controller.Upload)

	engine.IrisEngine.Run(iris.Addr("127.0.0.1:3000"), iris.WithCharset("UTF-8"))
}
