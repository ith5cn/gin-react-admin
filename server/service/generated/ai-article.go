package generated

import (
	commonResponse "server/model/common/response"
	generatedModel "server/model/generated"
	systemService "server/service/system"
)

func AiarticleList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []generatedModel.Aiarticle
	return systemService.PageList(query, &generatedModel.Aiarticle{}, &data, map[string]string{"title": "title"}, map[string]string{"status": "status"}, "id DESC")
}

func CreateAiarticle(data map[string]interface{}) (*generatedModel.Aiarticle, error) {
	return systemService.CreateRecord[generatedModel.Aiarticle]("ai_article", data)
}

func UpdateAiarticle(id string, data map[string]interface{}) (*generatedModel.Aiarticle, error) {
	return systemService.UpdateRecord[generatedModel.Aiarticle]("ai_article", id, data)
}

func DeleteAiarticle(id string) error {
	return systemService.DeleteRecord(&generatedModel.Aiarticle{}, id)
}
