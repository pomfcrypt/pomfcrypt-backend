package routes

type Controller struct{}

func NewController() *Controller { return &Controller{} }

type APIMessage struct{ Message string `json:"message"` }

func NewAPIMessage(message string) *APIMessage { return &APIMessage{Message: message} }

type APIErrorMessage struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func NewAPIError(message string, error error) *APIErrorMessage { return &APIErrorMessage{Message: message, Error: error} }
