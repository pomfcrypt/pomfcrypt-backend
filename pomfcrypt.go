//go:generate swagger generate spec

package main

import (
	"github.com/kataras/iris"
	"github.com/pomfcrypt/pomfcrypt-backend/routes"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

// PomfEngine is a container struct which holds the API routing engine and the API route controller
type PomfEngine struct {
	App        *iris.Application
	Controller *routes.Controller
}

// Kingpin CLI flag definitions
var (
	verbose          = kingpin.Flag("verbose", "Enable verbose output").Envar("POMF_VERBOSE").Short('v').Bool()
	maxSize          = kingpin.Flag("max-size", "Set maximum file size in bytes").Envar("POMF_MAX_SIZE").Default("256000000").Int64()
	filenameLength   = kingpin.Flag("filename-length", "Set random filename length").Envar("POMF_LEN_FILENAME").Default("4").Int()
	uploadsDirectory = kingpin.Flag("directory", "Upload directory").Short('d').Envar("POMF_DIR").Default("uploads").ExistingDir()
	salt             = kingpin.Flag("salt", "Set salt for encryption").Short('s').Envar("POMF_SALT").Default("salt").String()
	debug            = kingpin.Flag("debug", "Enable debug mode").Envar("POMF_DEBUG").Bool()
)

func main() {
	// Parse the CLI parameters given
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	if *debug {
		logrus.Warn("Debug mode is activated!")
	}

	// Filter to verbose level if --verbose is provided as console flag
	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Initialize a Gin Engine (web server)
	engine := PomfEngine{App: iris.Default(), Controller: routes.NewController(&routes.Settings{MaxSize: *maxSize, FilenameLength: *filenameLength, UploadsDirectory: *uploadsDirectory, Salt: "$1" + *salt + "$!", Debug: *debug})}
	logrus.Debug("Initialized web framework and controller")

	// Provide information about the project as index
	engine.App.Get("/", func(c iris.Context) { c.JSON("https://github.com/pomfcrypt/pomfcrypt-backend") })

	// Create the API group ("party")
	api := engine.App.Party("/api/v1/")

	// /api/v1/file API route (upload)
	api.Put("/file", engine.Controller.Upload)
	// /api/v1/file/:name API route (retrieve)
	api.Post("/file/{name:string}", engine.Controller.Retrieve)

	logrus.Debug("Attempting to run server")
	if err := engine.App.Run(iris.Addr("127.0.0.1:3000"), iris.WithCharset("UTF-8")); err != nil {
		logrus.Fatal("Failed to run server: ", err)
	}
}
