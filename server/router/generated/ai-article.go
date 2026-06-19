package generated

import (
	generatedAPI "server/api/generated"

	"github.com/gin-gonic/gin"
)

func RegisterAiarticleRoutes(group *gin.RouterGroup) {
	group.GET("/system/ai-article/index", generatedAPI.AiarticleList)
	group.POST("/system/ai-article", generatedAPI.CreateAiarticle)
	group.PUT("/system/ai-article/:id", generatedAPI.UpdateAiarticle)
	group.DELETE("/system/ai-article/:id", generatedAPI.DeleteAiarticle)
}
