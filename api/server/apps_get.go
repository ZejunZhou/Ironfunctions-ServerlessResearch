package server

import (
	"context"
	"net/http"

	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleAppGet(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)

	appName := c.MustGet(api.AppName).(string)
	app, err := s.Datastore.GetApp(ctx, appName)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, appResponse{"Successfully loaded app", app})
}
