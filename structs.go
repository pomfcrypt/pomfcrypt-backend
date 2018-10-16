package main

import (
	"github.com/gin-gonic/gin"
)

// This file contains common structs used by PomfCrypt.

// Error struct containing user-relevant information about the error
type Error struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}

// Initialization function for the Error struct
func NewError(message string, code int) Error {
	return Error{Message: message, Code: code}
}

// Throw an error
func (e Error) Throw(context *gin.Context) {
	// Throw the error with the given code and message
	context.JSON(e.Code, e)
}
