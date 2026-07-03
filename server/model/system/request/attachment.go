package request

// AttachmentDeletePayload 是附件批量删除入参。
// RemoveSource 预留：是否删除真实文件，当前后端不处理。
type AttachmentDeletePayload struct {
	IDs          []uint `json:"ids"`
	RemoveSource bool   `json:"removeSource"`
}
