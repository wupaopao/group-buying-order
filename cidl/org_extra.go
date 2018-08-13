package cidl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mz-eco/mz/errors"
)

func NewGroupBuyingTaskOrderSpecificationByAsk(
	askSpecificationItems []*AskTaskSpecificationItem,
	askSkuItems []*AskTaskSkuItem,
	askCombinationItems []*AskTaskCombinationItem,
) (
	specification *GroupBuyingOrderTaskSpecification,
	err error,
) {
	defer func() {
		if err != nil {
			specification = nil
		}

	}()

	// 旧ID映射新ID
	oldLabelIdMapNewLabelId := make(map[string]string)
	newLabelMap := make(map[string]*GroupBuyingSpecificationItemLabel)

	specification = NewGroupBuyingOrderTaskSpecificationCustom()
	for _, askSpecificationItem := range askSpecificationItems {
		if askSpecificationItem.Name == "" {
			err = errors.New("empty specification item name")
			return
		}

		specificationItem := specification.AddSpecificationItem(askSpecificationItem.Name)
		for _, askLabelItem := range askSpecificationItem.Labels {
			oldLabelId := askLabelItem.LabelId
			if askLabelItem.Name == "" || oldLabelId == "" {
				err = errors.New("empty label name or id")
				return
			}

			_, ok := oldLabelIdMapNewLabelId[oldLabelId]
			if ok { // ID需要唯一
				err = errors.New("wrong label id")
				return
			}

			label := specificationItem.AddLabel(askLabelItem.Name)
			oldLabelIdMapNewLabelId[oldLabelId] = label.LabelId
			newLabelMap[label.LabelId] = label
		}
	}

	if len(newLabelMap) == 0 {
		err = errors.New("empty legal specification")
		return
	}

	// 单品
	skuItemHasShow := false
	oldSkuIdMapNewSkuId := make(map[string]string) // 旧sku id 映射 新sku id
	for _, askSkuItem := range askSkuItems {

		// sku id 需要唯一
		_, ok := oldSkuIdMapNewSkuId[askSkuItem.SkuId]
		if ok {
			err = errors.New("sku id is conflicting")
			return
		}

		skuItem := NewGroupBuyingSkuMapItem()
		var skuItemLabelIds []string
		for _, askSkuItemLabelId := range askSkuItem.LabelIds {
			newLabelId, ok := oldLabelIdMapNewLabelId[askSkuItemLabelId]
			if !ok { // label id 不存在
				err = errors.New("label id not exist")
				return
			}

			skuItemLabelIds = append(skuItemLabelIds, newLabelId)
			skuItem.Labels[newLabelId] = newLabelMap[newLabelId]
		}

		if len(skuItemLabelIds) == 0 {
			err = errors.New("empty label")
			return
		}

		sort.Strings(skuItemLabelIds)

		var skuItemLabelNames []string
		for _, labelId := range skuItemLabelIds {
			skuItemLabelNames = append(skuItemLabelNames, skuItem.Labels[labelId].Name)
		}

		skuItem.SkuId = strings.Join(skuItemLabelIds, "-")
		skuItem.Name = strings.Join(skuItemLabelNames, ",")
		skuItem.MarketPrice = askSkuItem.MarketPrice
		skuItem.GroupBuyingPrice = askSkuItem.GroupBuyingPrice
		skuItem.SettlementPrice = askSkuItem.SettlementPrice
		skuItem.CostPrice = askSkuItem.CostPrice
		skuItem.IllustrationPicture = askSkuItem.IllustrationPicture
		skuItem.IsShow = askSkuItem.IsShow
		if skuItem.IsShow {
			skuItemHasShow = true
		}

		skuItem.InventoryCount = askSkuItem.InventoryCount
		if skuItem.MarketPrice < 0 || skuItem.GroupBuyingPrice < 0 || skuItem.SettlementPrice < 0 || skuItem.CostPrice < 0 /*|| skuItem.IllustrationPicture == ""*/ {
			err = errors.New("wrong sku item params")
			return
		}

		// 检查价格范围
		specification.MarketPriceRange.SetRange(skuItem.MarketPrice)
		specification.GroupBuyingPriceRange.SetRange(skuItem.GroupBuyingPrice)
		specification.SettlementPriceRange.SetRange(skuItem.SettlementPrice)
		specification.CostPriceRange.SetRange(skuItem.CostPrice)

		oldSkuIdMapNewSkuId[askSkuItem.SkuId] = skuItem.SkuId
		specification.SkuMap[skuItem.SkuId] = skuItem
	}

	// 组合
	for _, askCombinationItem := range askCombinationItems {
		combinationItem := NewGroupBuyingOrderTaskCombinationItem()
		var subSkuItemIds []string
		for _, askSubSkuItem := range askCombinationItem.SubSkuItems {

			// sku 单品项需要存在
			newSkuId, ok := oldSkuIdMapNewSkuId[askSubSkuItem.SkuId]
			if !ok {
				err = errors.New("sku id not exists")
				return
			}

			// 单品数目不能为0
			if askSubSkuItem.Count == 0 {
				err = errors.New("combination sub sku item count can not be 0")
				return
			}

			skuItem := specification.SkuMap[newSkuId]
			combinationSubSkuItem := NewGroupBuyingOrderCombinationSubItem()
			combinationSubSkuItem.SkuId = skuItem.SkuId
			combinationSubSkuItem.Labels = skuItem.Labels
			combinationSubSkuItem.Count = askSubSkuItem.Count
			combinationItem.SubItems[combinationSubSkuItem.SkuId] = combinationSubSkuItem
			subSkuItemIds = append(subSkuItemIds, fmt.Sprintf("%s-%d", combinationSubSkuItem.SkuId, combinationSubSkuItem.Count))
		}

		if len(subSkuItemIds) == 0 {
			err = errors.New("empty combination sub items")
			return
		}

		sort.Strings(subSkuItemIds)
		combinationItem.SkuId = fmt.Sprintf("C_%s", strings.Join(subSkuItemIds, "_"))

		combinationItem.Name = askCombinationItem.Name
		combinationItem.MarketPrice = askCombinationItem.MarketPrice
		combinationItem.GroupBuyingPrice = askCombinationItem.GroupBuyingPrice
		combinationItem.SettlementPrice = askCombinationItem.SettlementPrice
		combinationItem.CostPrice = askCombinationItem.CostPrice
		combinationItem.IllustrationPicture = askCombinationItem.IllustrationPicture
		combinationItem.IsShow = askCombinationItem.IsShow
		if combinationItem.IsShow {
			skuItemHasShow = true
		}

		if combinationItem.MarketPrice < 0 || combinationItem.GroupBuyingPrice < 0 || combinationItem.SettlementPrice < 0 || combinationItem.CostPrice < 0 /*|| skuItem.IllustrationPicture == ""*/ {
			err = errors.New("wrong combination item price")
			return
		}
		specification.MarketPriceRange.SetRange(combinationItem.MarketPrice)
		specification.GroupBuyingPriceRange.SetRange(combinationItem.GroupBuyingPrice)
		specification.SettlementPriceRange.SetRange(combinationItem.SettlementPrice)
		specification.CostPriceRange.SetRange(combinationItem.CostPrice)

		specification.AddCombinationItem(combinationItem)
	}

	if skuItemHasShow == false {
		err = errors.New("has not show items")
		return
	}

	return
}
