package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
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

func CreatePost(data map[string]interface{}) (*systemModel.AISystemPost, error) {
	return createSimple[systemModel.AISystemPost]("ai_system_post", data, simpleColumns())
}

func UpdatePost(id string, data map[string]interface{}) (*systemModel.AISystemPost, error) {
	return updateSimple[systemModel.AISystemPost]("ai_system_post", id, data, simpleColumns())
}

func DeletePost(id string) error {
	return deleteByID(&systemModel.AISystemPost{}, id)
}

func simpleColumns() map[string]string {
	return map[string]string{"name": "name", "code": "code", "sort": "sort", "status": "status", "remark": "remark"}
}
