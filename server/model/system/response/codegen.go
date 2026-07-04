package response

// CodegenPreviewFile 是代码生成预览的单个文件：路径、内容和前后端分组。
type CodegenPreviewFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Group   string `json:"group"`
}
