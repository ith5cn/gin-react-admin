package system

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	systemModel "server/model/system"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// 用户 Excel 导入导出。表头在导出、模板、导入三处共用，保证格式一致：
// 用户用"导出的文件改一改再导回来"是最常见的使用方式，表头一旦对不上会非常恼人。

// userExcelHeaders 是用户 Excel 的固定表头（前 6 列参与导入，后面的列只导出展示）。
var userExcelHeaders = []string{"用户名", "昵称", "手机号", "邮箱", "部门ID", "状态(1正常 2停用)"}

// userImportDefaultPassword 是导入用户的初始密码，管理员导入后应通知用户尽快修改。
const userImportDefaultPassword = "123456"

// UserImportResult 是导入结果汇总，errors 里是逐行失败原因。
type UserImportResult struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}

// ExportUsersExcel 导出用户列表为 xlsx 字节流。
// 与 UserList 用同一套过滤条件和数据权限：导出的永远是操作者"看得到"的数据。
func ExportUsersExcel(operatorID uint, query map[string]string) ([]byte, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	scope, err := UserDataScope(operatorID)
	if err != nil {
		return nil, err
	}

	base := softDelete(db.Model(&systemModel.AISystemUser{}))
	base = applyUserDataScope(base, scope, operatorID)
	base = applyFilters(base, query,
		map[string]string{"username": "username", "nickname": "nickname", "phone": "phone", "email": "email"},
		map[string]string{"status": "status", "deptId": "dept_id"},
	)

	var users []systemModel.AISystemUser
	if err := base.Order("id ASC").Find(&users).Error; err != nil {
		return nil, err
	}

	file := excelize.NewFile()
	defer file.Close()
	sheet := file.GetSheetName(0)

	headers := append(append([]string{}, userExcelHeaders...), "创建时间")
	if err := file.SetSheetRow(sheet, "A1", &headers); err != nil {
		return nil, err
	}
	for i, user := range users {
		createTime := ""
		if user.CreateTime != nil {
			createTime = user.CreateTime.Format("2006-01-02 15:04:05")
		}
		row := []interface{}{
			user.Username,
			stringValue(user.Nickname),
			stringValue(user.Phone),
			stringValue(user.Email),
			uintPtrCell(user.DeptID),
			user.Status,
			createTime,
		}
		cell, err := excelize.CoordinatesToCellName(1, i+2)
		if err != nil {
			return nil, err
		}
		if err := file.SetSheetRow(sheet, cell, &row); err != nil {
			return nil, err
		}
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// UserImportTemplateExcel 生成导入模板：表头 + 一行示例数据。
func UserImportTemplateExcel() ([]byte, error) {
	file := excelize.NewFile()
	defer file.Close()
	sheet := file.GetSheetName(0)

	if err := file.SetSheetRow(sheet, "A1", &userExcelHeaders); err != nil {
		return nil, err
	}
	example := []interface{}{"zhangsan", "张三", "13800138000", "zhangsan@example.com", "", 1}
	if err := file.SetSheetRow(sheet, "A2", &example); err != nil {
		return nil, err
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// ImportUsersExcel 从 xlsx 导入用户，逐行处理：单行失败不中断整个导入，
// 最后汇总成功/失败数和失败原因。全部失败时返回业务错误（前端能直接看到原因）。
func ImportUsersExcel(reader io.Reader) (*UserImportResult, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, ErrImportNotExcel
	}
	defer file.Close()

	rows, err := file.GetRows(file.GetSheetName(0))
	if err != nil {
		return nil, err
	}
	// 第一行是表头，数据从第二行开始。
	if len(rows) <= 1 {
		return nil, ErrImportEmptyRows
	}

	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	// 初始密码整批只做一次 bcrypt（成本因子下逐行加密会让大批量导入慢到分钟级）。
	passwordHash, err := hashPassword(userImportDefaultPassword)
	if err != nil {
		return nil, err
	}

	result := &UserImportResult{Errors: []string{}}
	for index, row := range rows[1:] {
		lineNo := index + 2
		if err := importUserRow(db, row, passwordHash); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %s", lineNo, err.Error()))
			continue
		}
		result.Success++
	}

	if result.Success == 0 && result.Failed > 0 {
		return nil, NewBizError("导入失败：" + summarizeImportErrors(result.Errors))
	}
	return result, nil
}

// importUserRow 处理一行导入数据：列顺序与 userExcelHeaders 一致。
func importUserRow(db *gorm.DB, row []string, passwordHash string) error {
	username := strings.TrimSpace(cellAt(row, 0))
	if username == "" {
		return fmt.Errorf("用户名为空")
	}

	var count int64
	if err := db.Model(&systemModel.AISystemUser{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("用户名 %s 已存在", username)
	}

	data := map[string]interface{}{
		"username": username,
		"password": passwordHash,
		"status":   int16(1),
	}
	if nickname := strings.TrimSpace(cellAt(row, 1)); nickname != "" {
		data["nickname"] = nickname
	}
	if phone := strings.TrimSpace(cellAt(row, 2)); phone != "" {
		data["phone"] = phone
	}
	if email := strings.TrimSpace(cellAt(row, 3)); email != "" {
		data["email"] = email
	}
	if deptRaw := strings.TrimSpace(cellAt(row, 4)); deptRaw != "" {
		deptID, err := strconv.ParseUint(deptRaw, 10, 64)
		if err != nil {
			return fmt.Errorf("部门ID %q 不是数字", deptRaw)
		}
		data["dept_id"] = uint(deptID)
	}
	if statusRaw := strings.TrimSpace(cellAt(row, 5)); statusRaw != "" {
		status, err := strconv.ParseInt(statusRaw, 10, 16)
		if err != nil || (status != 1 && status != 2) {
			return fmt.Errorf("状态 %q 无效，只能是 1 或 2", statusRaw)
		}
		data["status"] = int16(status)
	}

	setDefaultTimes(data, true)
	return db.Model(&systemModel.AISystemUser{}).Create(data).Error
}

// cellAt 安全取单元格：excelize 返回的行会去掉行尾空单元格，越界按空串处理。
func cellAt(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return row[index]
}

// summarizeImportErrors 拼接失败原因，最多展示 5 条，避免弹窗爆炸。
func summarizeImportErrors(errors []string) string {
	const maxShown = 5
	if len(errors) <= maxShown {
		return strings.Join(errors, "；")
	}
	return strings.Join(errors[:maxShown], "；") + fmt.Sprintf("……（共 %d 条失败）", len(errors))
}

// uintPtrCell 把 *uint 转成 Excel 单元格值：nil 显示为空而不是 0。
func uintPtrCell(value *uint) interface{} {
	if value == nil {
		return ""
	}
	return *value
}
