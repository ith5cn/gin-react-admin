package system

import (
	systemModel "server/model/system"
	systemRequest "server/model/system/request"
)

// DeptList 查询部门并组装成树形结构。
func DeptList(query map[string]string) ([]*systemModel.AISystemDept, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var depts []systemModel.AISystemDept
	q := softDelete(db.Model(&systemModel.AISystemDept{}))
	q = applyFilters(q, query, map[string]string{"name": "name"}, map[string]string{"parentId": "parent_id", "status": "status"})
	if err := q.Order("sort ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	return BuildDeptTree(depts), nil
}

// CreateDept 新增部门，level 层级路径由 createWithLevel 自动维护。
func CreateDept(payload systemRequest.DeptPayload) (*systemModel.AISystemDept, error) {
	return createWithLevel[systemModel.AISystemDept]("ai_system_dept", deptPayloadData(payload))
}

// UpdateDept 更新部门，父级变化时同步重算 level。
func UpdateDept(id string, payload systemRequest.DeptPayload) (*systemModel.AISystemDept, error) {
	return updateWithLevel[systemModel.AISystemDept]("ai_system_dept", id, deptPayloadData(payload))
}

// DeleteDept 删除部门；存在子部门时拒绝删除。
func DeleteDept(id string) error {
	has, err := hasChildren("ai_system_dept", id)
	if err != nil {
		return err
	}
	if has {
		return ErrDeptHasChildren
	}
	return deleteByID(&systemModel.AISystemDept{}, id)
}

// DeptAccess 返回启用状态的部门，tree=true 时返回树形结构（下拉树用），
// 否则返回扁平列表。
func DeptAccess(tree bool) (interface{}, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var depts []systemModel.AISystemDept
	if err := softDelete(db).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	if tree {
		return BuildDeptTree(depts), nil
	}
	return depts, nil
}

// BuildDeptTree 把扁平部门列表组装成树，算法同 BuildMenuTree。
func BuildDeptTree(depts []systemModel.AISystemDept) []*systemModel.AISystemDept {
	nodeMap := make(map[uint]*systemModel.AISystemDept, len(depts))
	roots := make([]*systemModel.AISystemDept, 0)
	for i := range depts {
		dept := depts[i]
		dept.Children = []*systemModel.AISystemDept{}
		nodeMap[dept.ID] = &dept
	}
	for _, dept := range nodeMap {
		if isZeroParent(dept.ParentID) {
			roots = append(roots, dept)
			continue
		}
		if parent, ok := nodeMap[*dept.ParentID]; ok {
			parent.Children = append(parent.Children, dept)
		} else {
			roots = append(roots, dept)
		}
	}
	sortTreeChildren(roots, func(n *systemModel.AISystemDept) []*systemModel.AISystemDept { return n.Children }, func(a, b *systemModel.AISystemDept) bool {
		if a.Sort == b.Sort {
			return a.ID < b.ID
		}
		return a.Sort < b.Sort
	})
	return roots
}

// deptPayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过（部分更新语义）。
func deptPayloadData(payload systemRequest.DeptPayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "parent_id", payload.ParentID)
	setColumn(data, "name", payload.Name)
	setColumn(data, "status", payload.Status)
	setColumn(data, "sort", payload.Sort)
	setColumn(data, "remark", payload.Remark)
	return data
}
