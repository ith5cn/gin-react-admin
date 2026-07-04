package system

import (
	_ "embed"
	"encoding/json"
	"sync"
)

// region_data.json 是全国省市区三级行政区划数据（31 省 / 342 市 / 3056 区县），
// 数据源：国家统计局行政区划代码（经 github.com/modood/Administrative-divisions-of-China 整理）。
// 直接 embed 进二进制：数据一年才变一两次，没必要建表；更新时替换 JSON 文件重新编译即可。
//
//go:embed region_data.json
var regionRawData []byte

// RegionOption 是前端 Cascader 组件需要的选项结构：label 显示名，value 是行政区划代码。
type RegionOption struct {
	Label    string         `json:"label"`
	Value    string         `json:"value"`
	Children []RegionOption `json:"children,omitempty"`
}

// regionNode 对应原始 JSON 的 {code, name, children} 结构。
type regionNode struct {
	Code     string       `json:"code"`
	Name     string       `json:"name"`
	Children []regionNode `json:"children,omitempty"`
}

var (
	regionOnce    sync.Once
	regionOptions []RegionOption
	regionErr     error
)

// RegionOptions 返回省市区级联选项。
// 数据不可变，用 sync.Once 保证只解析一次，之后所有请求共享同一份内存结构。
func RegionOptions() ([]RegionOption, error) {
	regionOnce.Do(func() {
		var nodes []regionNode
		if err := json.Unmarshal(regionRawData, &nodes); err != nil {
			regionErr = err
			return
		}
		regionOptions = convertRegionNodes(nodes)
	})
	return regionOptions, regionErr
}

func convertRegionNodes(nodes []regionNode) []RegionOption {
	options := make([]RegionOption, 0, len(nodes))
	for _, node := range nodes {
		options = append(options, RegionOption{
			Label:    node.Name,
			Value:    node.Code,
			Children: convertRegionNodes(node.Children),
		})
	}
	return options
}
