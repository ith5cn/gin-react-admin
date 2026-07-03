package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
)

func DictTypeList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemDictType
	return pageList(query, &systemModel.AISystemDictType{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{"status": "status"}, "id ASC")
}

func CreateDictType(data map[string]interface{}) (*systemModel.AISystemDictType, error) {
	return createSimple[systemModel.AISystemDictType]("ai_system_dict_type", data, dictTypeColumns())
}

func UpdateDictType(id string, data map[string]interface{}) (*systemModel.AISystemDictType, error) {
	return updateSimple[systemModel.AISystemDictType]("ai_system_dict_type", id, data, dictTypeColumns())
}

func DeleteDictType(id string) error {
	return deleteByID(&systemModel.AISystemDictType{}, id)
}

func DictDataList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemDictData
	return pageList(query, &systemModel.AISystemDictData{}, &data, map[string]string{"label": "label", "code": "code"}, map[string]string{"status": "status", "typeId": "type_id"}, "sort ASC, id ASC")
}

func DictAll() (map[string][]systemModel.AISystemDictData, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var dicts []systemModel.AISystemDictData
	if err := softDelete(db).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&dicts).Error; err != nil {
		return nil, err
	}
	result := make(map[string][]systemModel.AISystemDictData)
	for _, item := range dicts {
		if item.Code == nil || *item.Code == "" {
			continue
		}
		result[*item.Code] = append(result[*item.Code], item)
	}
	return result, nil
}

func CreateDictData(data map[string]interface{}) (*systemModel.AISystemDictData, error) {
	return createSimple[systemModel.AISystemDictData]("ai_system_dict_data", data, dictDataColumns())
}

func UpdateDictData(id string, data map[string]interface{}) (*systemModel.AISystemDictData, error) {
	return updateSimple[systemModel.AISystemDictData]("ai_system_dict_data", id, data, dictDataColumns())
}

func DeleteDictData(id string) error {
	return deleteByID(&systemModel.AISystemDictData{}, id)
}

func dictTypeColumns() map[string]string {
	return map[string]string{"name": "name", "code": "code", "status": "status", "remark": "remark"}
}

func dictDataColumns() map[string]string {
	return map[string]string{"typeId": "type_id", "type_id": "type_id", "label": "label", "value": "value", "color": "color", "code": "code", "sort": "sort", "status": "status", "remark": "remark"}
}
