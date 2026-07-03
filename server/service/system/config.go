package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"

	"gorm.io/gorm"
)

func ConfigGroupList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemConfigGroup
	return pageList(query, &systemModel.AISystemConfigGroup{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{}, "sort ASC, id ASC")
}

func CreateConfigGroup(data map[string]interface{}) (*systemModel.AISystemConfigGroup, error) {
	return createSimple[systemModel.AISystemConfigGroup]("ai_system_config_group", data, configGroupColumns())
}

func UpdateConfigGroup(id string, data map[string]interface{}) (*systemModel.AISystemConfigGroup, error) {
	return updateSimple[systemModel.AISystemConfigGroup]("ai_system_config_group", id, data, configGroupColumns())
}

func DeleteConfigGroup(id string) error {
	return deleteByID(&systemModel.AISystemConfigGroup{}, id)
}

func ConfigList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemConfig
	return pageList(query, &systemModel.AISystemConfig{}, &data, map[string]string{"name": "name", "key": "`key`"}, map[string]string{"groupId": "group_id"}, "sort ASC, id ASC")
}

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

func BatchUpdateConfig(groupID uint, configs []map[string]interface{}) error {
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, item := range configs {
			id := item["id"]
			payload := requestData(item, configColumns())
			payload["group_id"] = groupID
			setDefaultTimes(payload, false)
			if id == nil {
				setDefaultTimes(payload, true)
				if err := tx.Model(&systemModel.AISystemConfig{}).Create(payload).Error; err != nil {
					return err
				}
				continue
			}
			if err := tx.Model(&systemModel.AISystemConfig{}).Where("id = ?", id).Updates(payload).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func CreateConfig(data map[string]interface{}) (*systemModel.AISystemConfig, error) {
	return createSimple[systemModel.AISystemConfig]("ai_system_config", data, configColumns())
}

func UpdateConfig(id string, data map[string]interface{}) (*systemModel.AISystemConfig, error) {
	return updateSimple[systemModel.AISystemConfig]("ai_system_config", id, data, configColumns())
}

func DeleteConfig(id string) error {
	return deleteByID(&systemModel.AISystemConfig{}, id)
}

func configColumns() map[string]string {
	return map[string]string{"groupId": "group_id", "group_id": "group_id", "key": "key", "value": "value", "name": "name", "inputType": "input_type", "input_type": "input_type", "configSelectData": "config_select_data", "config_select_data": "config_select_data", "sort": "sort", "remark": "remark"}
}

func configGroupColumns() map[string]string {
	return map[string]string{"name": "name", "code": "code", "sort": "sort", "remark": "remark"}
}
