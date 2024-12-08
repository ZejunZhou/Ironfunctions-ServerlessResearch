package server

import (
	"context"
	"net/http"

	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api"
	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/models"
	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/runner/common"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleAppDelete(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	log := common.Logger(ctx)

	app := &models.App{Name: c.MustGet(api.AppName).(string)}

	routes, err := s.Datastore.GetRoutesByApp(ctx, app.Name, &models.RouteFilter{})
	if err != nil {
		log.WithError(err).Error(models.ErrAppsRemoving)
		c.JSON(http.StatusInternalServerError, simpleError(ErrInternalServerError))
		return
	}
	//TODO allow this? #528
	if len(routes) > 0 {
		log.WithError(err).Debug(models.ErrDeleteAppsWithRoutes)
		c.JSON(http.StatusBadRequest, simpleError(models.ErrDeleteAppsWithRoutes))
		return
	}

	err = s.FireBeforeAppDelete(ctx, app)
	if err != nil {
		log.WithError(err).Error(models.ErrAppsRemoving)
		c.JSON(http.StatusInternalServerError, simpleError(ErrInternalServerError))
		return
	}

	err = s.Datastore.RemoveApp(ctx, app.Name)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	err = s.FireAfterAppDelete(ctx, app)
	if err != nil {
		log.WithError(err).Error(models.ErrAppsRemoving)
		c.JSON(http.StatusInternalServerError, simpleError(ErrInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "App deleted"})
}
