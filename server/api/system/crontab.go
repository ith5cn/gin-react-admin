package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// CrontabList 定时任务分页列表。
func CrontabList(c *gin.Context) {
	data, err := systemService.CrontabList(queryMap(c))
	successOrFail(c, data, err)
}

// CreateCrontab 新增定时任务。
func CreateCrontab(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.CrontabPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateCrontab(payload)
	successOrFail(c, result, err)
}

// UpdateCrontab 更新定时任务。
func UpdateCrontab(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.CrontabPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateCrontab(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteCrontab 删除定时任务。
func DeleteCrontab(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteCrontab(c.Param("id")))
}

// RunCrontab 手动触发一次任务执行（异步），执行结果看任务日志。
func RunCrontab(c *gin.Context) {
	successOrFail(c, "执行已触发，结果请查看任务日志", systemService.RunCrontabOnce(c.Param("id")))
}

// CrontabLogList 任务执行日志分页列表。
func CrontabLogList(c *gin.Context) {
	data, err := systemService.CrontabLogList(queryMap(c))
	successOrFail(c, data, err)
}
