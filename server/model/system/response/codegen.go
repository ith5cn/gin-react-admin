package response

type CodegenPreviewFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Group   string `json:"group"`
}
