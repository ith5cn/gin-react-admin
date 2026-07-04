package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"
)

// DictTypeList 分页查询字典类型。
func DictTypeList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemDictType
	return pageList(query, &systemModel.AISystemDictType{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{"status": "status"}, "id ASC")
}

// CreateDictType 新增字典类型。
func CreateDictType(payload systemRequest.DictTypePayload) (*systemModel.AISystemDictType, error) {
	return createRow[systemModel.AISystemDictType]("ai_system_dict_type", dictTypePayloadData(payload))
}

// UpdateDictType 更新字典类型。
func UpdateDictType(id string, payload systemRequest.DictTypePayload) (*systemModel.AISystemDictType, error) {
	return updateRow[systemModel.AISystemDictType]("ai_system_dict_type", id, dictTypePayloadData(payload))
}

// DeleteDictType 删除字典类型。
func DeleteDictType(id string) error {
	return deleteByID(&systemModel.AISystemDictType{}, id)
}

// DictDataList 分页查询字典数据，支持按 typeId 过滤。
func DictDataList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemDictData
	return pageList(query, &systemModel.AISystemDictData{}, &data, map[string]string{"label": "label", "code": "code"}, map[string]string{"status": "status", "typeId": "type_id"}, "sort ASC, id ASC")
}

// DictAll 一次返回全部启用的字典数据，按字典编码分组。
// 前端启动时调用一次并缓存到 store，避免每个下拉框都发请求。
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

// CreateDictData 新增字典数据。
func CreateDictData(payload systemRequest.DictDataPayload) (*systemModel.AISystemDictData, error) {
	return createRow[systemModel.AISystemDictData]("ai_system_dict_data", dictDataPayloadData(payload))
}

// UpdateDictData 更新字典数据。
func UpdateDictData(id string, payload systemRequest.DictDataPayload) (*systemModel.AISystemDictData, error) {
	return updateRow[systemModel.AISystemDictData]("ai_system_dict_data", id, dictDataPayloadData(payload))
}

// DeleteDictData 删除字典数据。
func DeleteDictData(id string) error {
	return deleteByID(&systemModel.AISystemDictData{}, id)
}

// dictTypePayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过。
func dictTypePayloadData(payload systemRequest.DictTypePayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "name", payload.Name)
	setColumn(data, "code", payload.Code)
	setColumn(data, "status", payload.Status)
	setColumn(data, "remark", payload.Remark)
	return data
}

// dictDataPayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过。
func dictDataPayloadData(payload systemRequest.DictDataPayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "type_id", payload.TypeID)
	setColumn(data, "label", payload.Label)
	setColumn(data, "value", payload.Value)
	setColumn(data, "color", payload.Color)
	setColumn(data, "code", payload.Code)
	setColumn(data, "sort", payload.Sort)
	setColumn(data, "status", payload.Status)
	setColumn(data, "remark", payload.Remark)
	return data
}
