package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"
)

// NoticeList 分页查询通知公告，最新的在前。
func NoticeList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemNotice
	return pageList(query, &systemModel.AISystemNotice{}, &data,
		map[string]string{"title": "title"},
		map[string]string{"type": "type", "status": "status"},
		"id DESC")
}

// CreateNotice 新增通知公告，标题必填。
func CreateNotice(payload systemRequest.NoticePayload) (*systemModel.AISystemNotice, error) {
	if payload.Title == nil || *payload.Title == "" {
		return nil, ErrNoticeTitleRequired
	}
	return createRow[systemModel.AISystemNotice]("ai_system_notice", noticePayloadData(payload))
}

// UpdateNotice 更新通知公告（部分更新语义）。
func UpdateNotice(id string, payload systemRequest.NoticePayload) (*systemModel.AISystemNotice, error) {
	if payload.Title != nil && *payload.Title == "" {
		return nil, ErrNoticeTitleRequired
	}
	return updateRow[systemModel.AISystemNotice]("ai_system_notice", id, noticePayloadData(payload))
}

// DeleteNotice 按 ID 删除通知公告。
func DeleteNotice(id string) error {
	return deleteByID(&systemModel.AISystemNotice{}, id)
}

// noticePayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过。
func noticePayloadData(payload systemRequest.NoticePayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "title", payload.Title)
	setColumn(data, "type", payload.Type.Int16Ptr())
	setColumn(data, "content", payload.Content)
	setColumn(data, "status", payload.Status.Int16Ptr())
	setColumn(data, "remark", payload.Remark)
	return data
}
