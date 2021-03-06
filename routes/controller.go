package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type Settings struct {
	MaxSize          int64  `json:"max_size"`
	FilenameLength   int    `json:"filename_length"`
	Salt             string `json:"salt"`
	UploadsDirectory string `json:"uploads_directory"`
	Debug            bool   `json:"debug"`
}

var debug = false

type Controller struct {
	Settings *Settings `json:"settings"`
}

func NewController(settings *Settings) *Controller {
	debug = settings.Debug
	return &Controller{Settings: settings}
}

type APIMessage struct{ Message string `json:"message"` }

func (e *APIMessage) Throw(c *gin.Context) { c.JSON(200, e) }

func NewAPIMessage(message string) *APIMessage { return &APIMessage{Message: message} }

type APIErrorMessage struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func NewAPIError(message string, error error) *APIErrorMessage { return &APIErrorMessage{Message: message, Error: error} }

func (e *APIErrorMessage) Throw(c iris.Context, status int) {
	c.StatusCode(status)
	if debug {
		// Throw whole error (marshal details)
		c.JSON(e)
	} else {
		// Disable verbose error output
		c.JSON(iris.Map{"message": e.Message})
	}
}

func (ctl *Controller) BuildPath(path string) string {
	absPath, err := filepath.Abs(ctl.Settings.UploadsDirectory + "/" + path)
	if err != nil {
		logrus.Fatal("Failed to open directory: ", err)
	}
	return absPath
}
