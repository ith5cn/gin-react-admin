package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
)

func LoginLogList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemLoginLog
	return pageList(query, &systemModel.AISystemLoginLog{}, &data, map[string]string{"username": "username"}, map[string]string{"status": "status", "ip": "ip"}, "id DESC")
}

func OperLogList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemOperLog
	return pageList(query, &systemModel.AISystemOperLog{}, &data, map[string]string{"username": "username", "serviceName": "service_name", "router": "router", "ip": "ip"}, map[string]string{}, "id DESC")
}
