package generated

import (
	systemAPI "server/api/system"
	generatedService "server/service/generated"

	"github.com/gin-gonic/gin"
)

func AiarticleList(c *gin.Context) {
	result, err := generatedService.AiarticleList(systemAPI.QueryMap(c))
	systemAPI.SuccessOrFail(c, result, err)
}

func CreateAiarticle(c *gin.Context) {
	data, ok := systemAPI.BindJSONMap(c)
	if !ok {
		return
	}
	result, err := generatedService.CreateAiarticle(data)
	systemAPI.SuccessOrFail(c, result, err)
}

func UpdateAiarticle(c *gin.Context) {
	data, ok := systemAPI.BindJSONMap(c)
	if !ok {
		return
	}
	result, err := generatedService.UpdateAiarticle(c.Param("id"), data)
	systemAPI.SuccessOrFail(c, result, err)
}

func DeleteAiarticle(c *gin.Context) {
	systemAPI.SuccessOrFail(c, map[string]interface{}{}, generatedService.DeleteAiarticle(c.Param("id")))
}
