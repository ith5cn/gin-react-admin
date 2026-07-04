package request

import (
	"bytes"
	"errors"
	"strconv"
)

// FlexInt16 是兼容字符串数字的 int16。
// 前端字典下拉（Ith5Select）提交的值是字符串（"1"），普通输入是数字（1），
// 这里在反序列化时统一成 int16，service 层不用再关心两种形态。
type FlexInt16 int16

// UnmarshalJSON 同时接受 1 和 "1" 两种 JSON 形态。
func (f *FlexInt16) UnmarshalJSON(data []byte) error {
	trimmed := bytes.Trim(bytes.TrimSpace(data), `"`)
	if len(trimmed) == 0 || string(trimmed) == "null" {
		return errors.New("invalid number")
	}
	value, err := strconv.ParseInt(string(trimmed), 10, 16)
	if err != nil {
		return err
	}
	*f = FlexInt16(value)
	return nil
}

// Int16Ptr 转回 *int16，方便直接写入 GORM 更新 map。
func (f *FlexInt16) Int16Ptr() *int16 {
	if f == nil {
		return nil
	}
	value := int16(*f)
	return &value
}

// CrontabPayload 是定时任务创建/更新入参，指针字段为 nil 表示不改动（部分更新）。
type CrontabPayload struct {
	Name      *string    `json:"name"`
	Type      *FlexInt16 `json:"type"`
	Target    *string    `json:"target"`
	Parameter *string    `json:"parameter"`
	TaskStyle *FlexInt16 `json:"taskStyle"`
	Rule      *string    `json:"rule"`
	Singleton *FlexInt16 `json:"singleton"`
	Status    *FlexInt16 `json:"status"`
	Remark    *string    `json:"remark"`
}
