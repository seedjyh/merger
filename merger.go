package merger

import (
	"fmt"
)

// Merger 允许输入不同类型的若干Item列表并分组返回。通常用于核对多个大规模数据列表（如话单、账单），找出差异。
type Merger struct {
	category int
	k2c      map[string]*itemListComposite
}

// category表示，即将输入的Item一共来自category个来源。
func NewMerger(category int) *Merger {
	return &Merger{
		category: category,
		k2c:      make(map[string]*itemListComposite),
	}
}

// AppendItem 添加一个Item。调用方需要保证不会传入重复项。
// 如果数据有问题，则返回error非空。
func (m *Merger) AppendItem(item Item) error {
	if item.CategoryID() < 0 || item.CategoryID() >= m.category {
		return fmt.Errorf("no such category")
	}
	key := item.Key()
	if _, ok := m.k2c[key]; !ok {
		m.k2c[key] = newItemListComposite(m.category)
	}
	return m.k2c[key].appendItem(item)
}

// Merge 返回一个ItemGroup切片。
// maxDifference是在同一ItemGroup内允许的最大Item差异度。
func (m *Merger) Merge(maxDifference float64) []*ItemGroup {
	var result []*ItemGroup
	for _, v := range m.k2c {
		result = append(result, v.merge(maxDifference)...)
	}
	return result
}
