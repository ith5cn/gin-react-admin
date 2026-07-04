package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"
)

// PostList 分页查询岗位列表。
func PostList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemPost
	return pageList(query, &systemModel.AISystemPost{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{"status": "status"}, "sort ASC, id ASC")
}

// PostAccess 返回启用状态的岗位列表，用于下拉选择。
func PostAccess() ([]systemModel.AISystemPost, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var posts []systemModel.AISystemPost
	err = softDelete(db).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&posts).Error
	return posts, err
}

// CreatePost 新增岗位。
func CreatePost(payload systemRequest.PostPayload) (*systemModel.AISystemPost, error) {
	return createRow[systemModel.AISystemPost]("ai_system_post", postPayloadData(payload))
}

// UpdatePost 更新岗位（部分更新：nil 字段不动）。
func UpdatePost(id string, payload systemRequest.PostPayload) (*systemModel.AISystemPost, error) {
	return updateRow[systemModel.AISystemPost]("ai_system_post", id, postPayloadData(payload))
}

// DeletePost 按 ID 删除岗位。
func DeletePost(id string) error {
	return deleteByID(&systemModel.AISystemPost{}, id)
}

// postPayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过。
func postPayloadData(payload systemRequest.PostPayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "name", payload.Name)
	setColumn(data, "code", payload.Code)
	setColumn(data, "sort", payload.Sort)
	setColumn(data, "status", payload.Status)
	setColumn(data, "remark", payload.Remark)
	return data
}
