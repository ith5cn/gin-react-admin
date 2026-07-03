package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"
)

func DictTypeList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemDictType
	return pageList(query, &systemModel.AISystemDictType{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{"status": "status"}, "id ASC")
}

func CreateDictType(payload systemRequest.DictTypePayload) (*systemModel.AISystemDictType, error) {
	return createRow[systemModel.AISystemDictType]("ai_system_dict_type", dictTypePayloadData(payload))
}

func UpdateDictType(id string, payload systemRequest.DictTypePayload) (*systemModel.AISystemDictType, error) {
	return updateRow[systemModel.AISystemDictType]("ai_system_dict_type", id, dictTypePayloadData(payload))
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

func CreateDictData(payload systemRequest.DictDataPayload) (*systemModel.AISystemDictData, error) {
	return createRow[systemModel.AISystemDictData]("ai_system_dict_data", dictDataPayloadData(payload))
}

func UpdateDictData(id string, payload systemRequest.DictDataPayload) (*systemModel.AISystemDictData, error) {
	return updateRow[systemModel.AISystemDictData]("ai_system_dict_data", id, dictDataPayloadData(payload))
}

func DeleteDictData(id string) error {
	return deleteByID(&systemModel.AISystemDictData{}, id)
}

func dictTypePayloadData(payload systemRequest.DictTypePayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "name", payload.Name)
	setColumn(data, "code", payload.Code)
	setColumn(data, "status", payload.Status)
	setColumn(data, "remark", payload.Remark)
	return data
}

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
