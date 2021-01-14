// item
package merger

import "sort"

// Item是要被分组的数据块。
type Item interface {
	// CategoryID 用于描述Item的来源。取值在0 <= CategoryID < Merger.category。
	CategoryID() int
	// Key 生成一个字符串。Key不同的Item必定被分到不同的ItemGroup。
	Key() string
	// Similar 比较两个Item的差异度。允许返回负数。差异度越接近0（绝对值越小），表示两个Item越应该属于同一个ItemGroup。
	DifferenceFrom(other Item) float64
}

type itemList []Item

func (l itemList) sort() {
	sort.SliceStable(l, func(i, j int) bool { return l[i].DifferenceFrom(l[j]) < 0.0 })
}

// ItemGroup 描述了一组Item。
type ItemGroup struct {
	// Items 的长度和Merger.category相同。但里面可能有一个nil。但至少一个非空的Item。
	Items []Item
}

func newItemGroup(category int) *ItemGroup {
	if category <= 0 {
		panic("require at least 1 category")
	}
	return &ItemGroup{
		Items: make([]Item, category),
	}
}