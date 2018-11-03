package routes

import "github.com/gin-gonic/gin"

// Upload doc
// @Summary Upload route
// @Description Upload a new file
// @Accept mpfd
// @Produce json
// @Success 200 {object} model.FileResponse
// @Failure 500 {object} routes.APIErrorMessage
// @Failure 400 {object} routes.APIErrorMessage
// @Router /data [put]
func (ctl *Controller) Upload(c *gin.Context) {
	
}
