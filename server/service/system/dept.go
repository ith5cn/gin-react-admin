package system

import (
	"errors"
	systemModel "server/model/system"
)

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

func CreateDept(data map[string]interface{}) (*systemModel.AISystemDept, error) {
	payload := requestData(data, deptColumns())
	return createWithLevel[systemModel.AISystemDept]("ai_system_dept", payload)
}

func UpdateDept(id string, data map[string]interface{}) (*systemModel.AISystemDept, error) {
	payload := requestData(data, deptColumns())
	return updateWithLevel[systemModel.AISystemDept]("ai_system_dept", id, payload)
}

func DeleteDept(id string) error {
	has, err := hasChildren("ai_system_dept", id)
	if err != nil {
		return err
	}
	if has {
		return errors.New("部门下存在子部门，无法删除")
	}
	return deleteByID(&systemModel.AISystemDept{}, id)
}

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

func deptColumns() map[string]string {
	return map[string]string{"parentId": "parent_id", "parent_id": "parent_id", "name": "name", "status": "status", "sort": "sort", "remark": "remark"}
}
