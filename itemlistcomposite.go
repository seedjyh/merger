// itemlistcomposite
package merger

import "fmt"

type itemListComposite struct {
	categoriesCount int
	itemLists     []itemList
}

func newItemListComposite(categoriesCount int) *itemListComposite {
	return &itemListComposite{
		categoriesCount: categoriesCount,
		itemLists:     make([]itemList, categoriesCount),
	}
}

// appendItem 添加一个Item。调用方需要保证不会传入重复项。
// 如果数据有问题，则返回error非空。
func (c *itemListComposite) appendItem(item Item) error {
	if item.CategoryID() < 0 || item.CategoryID() >= c.categoriesCount {
		return fmt.Errorf("no such category")
	}
	c.itemLists[item.CategoryID()] = append(c.itemLists[item.CategoryID()], item)
	return nil
}

func (c *itemListComposite) merge(maxDifference float64) []*ItemGroup {
	for _, il := range c.itemLists {
		il.sort()
	}
	var result []*ItemGroup
	for {
		// 每次选择各list的第一项
		var headItems itemList
		for _, il := range c.itemLists {
			if len(il) > 0 {
				headItems = append(headItems, il[0])
			}
		}
		if len(headItems) == 0 {
			break // 所有list都是空的了
		}
		// 按序取
		g := newItemGroup(c.categoriesCount)
		headItems.sort()
		var firstChosenItem Item // 也是“最小”的item
		for i, item := range headItems {
			if i == 0 {
				firstChosenItem = item
				g.Items[item.CategoryID()] = item
				c.itemLists[item.CategoryID()] = c.itemLists[item.CategoryID()][1:]
			} else {
				if item.DifferenceFrom(firstChosenItem) < maxDifference {
					g.Items[item.CategoryID()] = item
					c.itemLists[item.CategoryID()] = c.itemLists[item.CategoryID()][1:]
				} else {
					break
				}
			}
		}
		result = append(result, g)
	}
	return result
}
