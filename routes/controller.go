package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type Settings struct {
	MaxSize          int64  `json:"max_size"`
	FilenameLength   int    `json:"filename_length"`
	Salt             string `json:"salt"`
	UploadsDirectory string `json:"uploads_directory"`
}

type Controller struct {
	Settings *Settings `json:"settings"`
}

func NewController(settings *Settings) *Controller { return &Controller{Settings: settings} }

type APIMessage struct{ Message string `json:"message"` }

func (e *APIMessage) Throw(c *gin.Context) { c.JSON(200, e) }

func NewAPIMessage(message string) *APIMessage { return &APIMessage{Message: message} }

type APIErrorMessage struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func NewAPIError(message string, error error) *APIErrorMessage { return &APIErrorMessage{Message: message, Error: error} }

func (e *APIErrorMessage) Throw(c *gin.Context, status int) { c.JSON(status, e) }

func (c *Controller) BuildPath(path string) string {
	absPath, err := filepath.Abs(c.Settings.UploadsDirectory + "/" + path)
	if err != nil {
		logrus.Fatal("Failed to open directory: ", err)
	}
	return absPath
}
