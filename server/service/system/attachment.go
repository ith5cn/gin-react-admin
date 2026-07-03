package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

var documentMimeTypes = []string{
	"application/pdf",
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.ms-excel",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
}

// AttachmentList 查询附件分页列表。
// 前端会传 resourceType=all/image/document/audio/video/application，这里按 MIME 类型做分类过滤。
func AttachmentList(query map[string]string) (*commonResponse.PageResult, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	page := parsePage(query)
	base := softDelete(db.Model(&systemModel.AISystemAttachment{}))
	base = applyAttachmentFilters(base, query)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}

	var rows []systemModel.AISystemAttachment
	if err := base.Order("id DESC").Offset((page.Page - 1) * page.Size).Limit(page.Size).Find(&rows).Error; err != nil {
		return nil, err
	}

	for i := range rows {
		rows[i].ResourceType = AttachmentResourceType(rows[i].MimeType)
	}

	return &commonResponse.PageResult{List: rows, Total: total}, nil
}

// DeleteAttachments 批量软删除附件记录。
// removeSource 目前只预留参数，不主动删除真实文件，避免本地/云存储误删。
func DeleteAttachments(ids []uint) error {
	if len(ids) == 0 {
		return ErrAttachmentIDsEmpty
	}

	db, err := systemDB()
	if err != nil {
		return err
	}

	return db.Model(&systemModel.AISystemAttachment{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{"delete_time": gorm.Expr("NOW()"), "update_time": gorm.Expr("NOW()")}).
		Error
}

func applyAttachmentFilters(db *gorm.DB, query map[string]string) *gorm.DB {
	if query["originName"] != "" {
		db = db.Where("origin_name LIKE ?", "%"+query["originName"]+"%")
	}
	if query["origin_name"] != "" {
		db = db.Where("origin_name LIKE ?", "%"+query["origin_name"]+"%")
	}
	if query["mimeType"] != "" {
		db = db.Where("mime_type LIKE ?", "%"+query["mimeType"]+"%")
	}
	if query["mime_type"] != "" {
		db = db.Where("mime_type LIKE ?", "%"+query["mime_type"]+"%")
	}
	if query["storageMode"] != "" {
		db = db.Where("storage_mode = ?", query["storageMode"])
	}
	if query["storage_mode"] != "" {
		db = db.Where("storage_mode = ?", query["storage_mode"])
	}

	db = applyDateRange(db, query)
	return applyResourceTypeFilter(db, query["resourceType"])
}

func applyDateRange(db *gorm.DB, query map[string]string) *gorm.DB {
	start := firstNonEmpty(query["startDate"], query["start_date"])
	end := firstNonEmpty(query["endDate"], query["end_date"])
	if startTime, ok := parseDate(start, false); ok {
		db = db.Where("create_time >= ?", startTime)
	}
	if endTime, ok := parseDate(end, true); ok {
		db = db.Where("create_time <= ?", endTime)
	}
	return db
}

func applyResourceTypeFilter(db *gorm.DB, resourceType string) *gorm.DB {
	switch strings.ToLower(strings.TrimSpace(resourceType)) {
	case "", "all":
		return db
	case "image":
		return db.Where("mime_type LIKE ?", "image/%")
	case "audio":
		return db.Where("mime_type LIKE ?", "audio/%")
	case "video":
		return db.Where("mime_type LIKE ?", "video/%")
	case "document":
		return db.Where("(mime_type LIKE ? OR mime_type IN ?)", "text/%", documentMimeTypes)
	case "application":
		return db.Where("(mime_type LIKE ? AND mime_type NOT IN ?)", "application/%", documentMimeTypes)
	default:
		return db
	}
}

// AttachmentResourceType 根据 MIME 类型返回前端展示分类。
func AttachmentResourceType(mimeType *string) string {
	normalized := ""
	if mimeType != nil {
		normalized = strings.ToLower(strings.TrimSpace(*mimeType))
	}
	if strings.HasPrefix(normalized, "image/") {
		return "image"
	}
	if strings.HasPrefix(normalized, "audio/") {
		return "audio"
	}
	if strings.HasPrefix(normalized, "video/") {
		return "video"
	}
	if strings.HasPrefix(normalized, "text/") || stringInSlice(normalized, documentMimeTypes) {
		return "document"
	}
	return "application"
}

func parseDate(value string, endOfDay bool) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	layouts := []string{"2006-01-02", "2006-01-02 15:04:05", time.RFC3339}
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err != nil {
			continue
		}
		if endOfDay && layout == "2006-01-02" {
			parsed = parsed.Add(24*time.Hour - time.Nanosecond)
		}
		return parsed, true
	}
	if timestamp, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Unix(timestamp, 0), true
	}
	return time.Time{}, false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func stringInSlice(value string, items []string) bool {
	for _, item := range items {
		if value == item {
			return true
		}
	}
	return false
}
