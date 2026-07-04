package system

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	systemModel "server/model/system"
)

// UploadRoot 是本地上传文件的存储根目录，router 会把它挂载为 /uploads 静态服务。
const UploadRoot = "runtime/uploads"

const (
	maxImageSize = 10 << 20 // 图片上限 10MB
	maxFileSize  = 50 << 20 // 普通文件上限 50MB
)

// 扩展名白名单：宁可漏放也不错放。
// 注意不收 .svg（内嵌脚本有 XSS 风险）、不收任何可执行格式。
var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true, ".ico": true,
}

var fileExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true, ".ico": true,
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true,
	".txt": true, ".md": true, ".csv": true, ".json": true, ".xml": true,
	".zip": true, ".rar": true, ".7z": true, ".gz": true, ".tar": true,
	".mp3": true, ".wav": true, ".mp4": true, ".avi": true, ".mov": true, ".webm": true,
}

// SaveUploadedImage 保存上传图片：在通用校验之外，额外用文件内容嗅探确认真的是图片，
// 防止把 .php 改名 .jpg 这类伪装文件传上来。
func SaveUploadedImage(fileHeader *multipart.FileHeader, userID uint) (*systemModel.AISystemAttachment, error) {
	return saveUpload(fileHeader, userID, imageExts, maxImageSize, true)
}

// SaveUploadedFile 保存普通上传文件。
func SaveUploadedFile(fileHeader *multipart.FileHeader, userID uint) (*systemModel.AISystemAttachment, error) {
	return saveUpload(fileHeader, userID, fileExts, maxFileSize, false)
}

// saveUpload 是上传的通用实现：
// 校验扩展名和大小 → 读内容算 sha256 → 命中相同 hash 直接复用（秒传）→
// 落盘到按日期分目录的路径 → 写附件表。
// 存储文件名用内容 hash 而不是原文件名，天然避免路径穿越、重名覆盖和中文乱码问题。
func saveUpload(fileHeader *multipart.FileHeader, userID uint, allowedExts map[string]bool, maxSize int64, requireImage bool) (*systemModel.AISystemAttachment, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExts[ext] {
		return nil, NewBizError(fmt.Sprintf("不支持的文件类型 %s", ext))
	}
	if fileHeader.Size <= 0 {
		return nil, ErrUploadEmptyFile
	}
	if fileHeader.Size > maxSize {
		return nil, NewBizError(fmt.Sprintf("文件大小超过限制（最大 %dMB）", maxSize>>20))
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	// DetectContentType 按文件头部字节嗅探真实类型，不信任请求头声明的 Content-Type。
	mimeType := http.DetectContentType(content)
	if requireImage && !strings.HasPrefix(mimeType, "image/") {
		return nil, ErrUploadNotImage
	}

	sum := sha256.Sum256(content)
	hash := hex.EncodeToString(sum[:])

	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	// 相同内容已上传过则直接复用记录，不重复占磁盘（秒传）。
	var existing systemModel.AISystemAttachment
	if err := softDelete(db).Where("hash = ?", hash).First(&existing).Error; err == nil {
		return &existing, nil
	}

	relativeDir := filepath.Join(time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"))
	objectName := hash[:32] + ext
	if err := os.MkdirAll(filepath.Join(UploadRoot, relativeDir), 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(UploadRoot, relativeDir, objectName), content, 0644); err != nil {
		return nil, err
	}

	storagePath := filepath.ToSlash(filepath.Join(relativeDir, objectName))
	url := "/uploads/" + storagePath
	createdBy := int(userID)
	attachment := systemModel.AISystemAttachment{
		StorageMode: 1, // 1 = 本地存储
		OriginName:  ptrString(fileHeader.Filename),
		ObjectName:  ptrString(objectName),
		Hash:        ptrString(hash),
		MimeType:    ptrString(mimeType),
		StoragePath: ptrString(storagePath),
		Suffix:      ptrString(strings.TrimPrefix(ext, ".")),
		SizeByte:    &fileHeader.Size,
		SizeInfo:    ptrString(humanSize(fileHeader.Size)),
		URL:         ptrString(url),
		CreatedBy:   &createdBy,
	}
	if err := db.Create(&attachment).Error; err != nil {
		return nil, err
	}
	attachment.ResourceType = AttachmentResourceType(attachment.MimeType)
	return &attachment, nil
}

// humanSize 把字节数转成可读大小，如 1.5 MB。
func humanSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
