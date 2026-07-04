package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"

	"gorm.io/gorm"
)

// ConfigGroupList 分页查询配置分组。
func ConfigGroupList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemConfigGroup
	return pageList(query, &systemModel.AISystemConfigGroup{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{}, "sort ASC, id ASC")
}

// CreateConfigGroup 新增配置分组。
func CreateConfigGroup(payload systemRequest.ConfigGroupPayload) (*systemModel.AISystemConfigGroup, error) {
	return createRow[systemModel.AISystemConfigGroup]("ai_system_config_group", configGroupPayloadData(payload))
}

// UpdateConfigGroup 更新配置分组。
func UpdateConfigGroup(id string, payload systemRequest.ConfigGroupPayload) (*systemModel.AISystemConfigGroup, error) {
	return updateRow[systemModel.AISystemConfigGroup]("ai_system_config_group", id, configGroupPayloadData(payload))
}

// DeleteConfigGroup 删除配置分组。
func DeleteConfigGroup(id string) error {
	return deleteByID(&systemModel.AISystemConfigGroup{}, id)
}

// ConfigList 分页查询配置项，支持按 groupId 过滤。
func ConfigList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemConfig
	return pageList(query, &systemModel.AISystemConfig{}, &data, map[string]string{"name": "name", "key": "`key`"}, map[string]string{"groupId": "group_id"}, "sort ASC, id ASC")
}

// ConfigInfo 按分组编码或配置 key 查询配置，返回 key → value 的映射。
// key 是保留字，SQL 里必须用反引号包住，否则语法错误。
func ConfigInfo(code string) (map[string]string, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var configs []systemModel.AISystemConfig
	if err := db.Table("ai_system_config AS c").
		Select("c.*").
		Joins("JOIN ai_system_config_group AS g ON g.id = c.group_id").
		Where("c.delete_time IS NULL AND (g.code = ? OR c.`key` = ?)", code, code).
		Order("c.sort ASC, c.id ASC").
		Find(&configs).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, item := range configs {
		if item.Value == nil {
			result[item.Key] = ""
			continue
		}
		result[item.Key] = *item.Value
	}
	return result, nil
}

// BatchUpdateConfig 批量保存一组配置项：有 id 的更新、没 id 的新建，
// 全部放在一个事务里，避免保存到一半失败产生脏数据。
func BatchUpdateConfig(payload systemRequest.BatchUpdateConfigPayload) error {
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, item := range payload.Config {
			data := configPayloadData(item)
			data["group_id"] = payload.GroupID
			setDefaultTimes(data, false)
			if item.ID == nil {
				setDefaultTimes(data, true)
				if err := tx.Model(&systemModel.AISystemConfig{}).Create(data).Error; err != nil {
					return err
				}
				continue
			}
			if err := tx.Model(&systemModel.AISystemConfig{}).Where("id = ?", *item.ID).Updates(data).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CreateConfig 新增配置项。
func CreateConfig(payload systemRequest.ConfigPayload) (*systemModel.AISystemConfig, error) {
	return createRow[systemModel.AISystemConfig]("ai_system_config", configPayloadData(payload))
}

// UpdateConfig 更新配置项。
func UpdateConfig(id string, payload systemRequest.ConfigPayload) (*systemModel.AISystemConfig, error) {
	return updateRow[systemModel.AISystemConfig]("ai_system_config", id, configPayloadData(payload))
}

// DeleteConfig 删除配置项。
func DeleteConfig(id string) error {
	return deleteByID(&systemModel.AISystemConfig{}, id)
}

// configPayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过。
func configPayloadData(payload systemRequest.ConfigPayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "group_id", payload.GroupID)
	setColumn(data, "key", payload.Key)
	setColumn(data, "name", payload.Name)
	setColumn(data, "input_type", payload.InputType)
	setColumn(data, "config_select_data", payload.ConfigSelectData)
	setColumn(data, "sort", payload.Sort)
	setColumn(data, "remark", payload.Remark)
	// value 可能是字符串或数字（radio/select 组件），统一走 normalizeValue 归一化。
	if payload.Value != nil {
		data["value"] = normalizeValue(payload.Value)
	}
	return data
}

// configGroupPayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过。
func configGroupPayloadData(payload systemRequest.ConfigGroupPayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "name", payload.Name)
	setColumn(data, "code", payload.Code)
	setColumn(data, "sort", payload.Sort)
	setColumn(data, "remark", payload.Remark)
	return data
}
