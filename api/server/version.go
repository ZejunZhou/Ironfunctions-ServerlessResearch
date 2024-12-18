package server

import (
	"net/http"

	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/version"
	"github.com/gin-gonic/gin"
)

func handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": version.Version})
}
