package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"
)

func PostList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemPost
	return pageList(query, &systemModel.AISystemPost{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{"status": "status"}, "sort ASC, id ASC")
}

func PostAccess() ([]systemModel.AISystemPost, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var posts []systemModel.AISystemPost
	err = softDelete(db).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&posts).Error
	return posts, err
}

func CreatePost(payload systemRequest.PostPayload) (*systemModel.AISystemPost, error) {
	return createRow[systemModel.AISystemPost]("ai_system_post", postPayloadData(payload))
}

func UpdatePost(id string, payload systemRequest.PostPayload) (*systemModel.AISystemPost, error) {
	return updateRow[systemModel.AISystemPost]("ai_system_post", id, postPayloadData(payload))
}

func DeletePost(id string) error {
	return deleteByID(&systemModel.AISystemPost{}, id)
}

func postPayloadData(payload systemRequest.PostPayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "name", payload.Name)
	setColumn(data, "code", payload.Code)
	setColumn(data, "sort", payload.Sort)
	setColumn(data, "status", payload.Status)
	setColumn(data, "remark", payload.Remark)
	return data
}
