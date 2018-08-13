package cidl

import "fmt"

// 产生新标签id
func (m *GroupBuyingSpecificationItem) generateNewLabelId() (labelId string) {
	m.LabelIdCounter++
	labelId = fmt.Sprintf("%s:%d", m.ItemId, m.LabelIdCounter)
	return
}

// 添加规格项标签
func (m *GroupBuyingSpecificationItem) AddLabel(labelName string) (label *GroupBuyingSpecificationItemLabel) {
	label = NewGroupBuyingSpecificationItemLabel()
	label.Name = labelName
	label.LabelId = m.generateNewLabelId()
	m.Labels[label.LabelId] = label
	return
}

// 产生规则项ID
func (m *GroupBuyingSpecification) generateNewSpecificationItemId() (itemId string) {
	m.ItemIdCounter++
	itemId = fmt.Sprintf("%d", m.ItemIdCounter)
	return
}

// 添加规则项
func (m *GroupBuyingSpecification) AddSpecificationItem(itemName string) (item *GroupBuyingSpecificationItem) {
	item = NewGroupBuyingSpecificationItem()
	item.ItemId = m.generateNewSpecificationItemId()
	item.Name = itemName
	m.Items[item.ItemId] = item
	return
}

// 组合规格
// 添加组合
func (m *GroupBuyingOrderTaskSpecification) AddCombinationItem(combinationItem *GroupBuyingOrderTaskCombinationItem) (err error) {
	m.CombinationSkuMap[combinationItem.SkuId] = combinationItem
	return
}

// 组合规格项
