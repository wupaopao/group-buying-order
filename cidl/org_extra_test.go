package cidl

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewGroupBuyingTaskOrderSpecificationByAsk(t *testing.T) {
	askSpecificationItems := []*AskTaskSpecificationItem{
		&AskTaskSpecificationItem{
			Name: "类型",
			Labels: []*AskTaskSpecificationLabelItem{
				&AskTaskSpecificationLabelItem{
					LabelId: "1:1",
					Name:    "上衣",
				},
				&AskTaskSpecificationLabelItem{
					LabelId: "1:2",
					Name:    "裤子",
				},
			},
		},
		&AskTaskSpecificationItem{
			Name: "尺寸",
			Labels: []*AskTaskSpecificationLabelItem{
				&AskTaskSpecificationLabelItem{
					LabelId: "2:1",
					Name:    "X",
				},
				&AskTaskSpecificationLabelItem{
					LabelId: "2:2",
					Name:    "L",
				},
			},
		},
	}

	askSkuItems := []*AskTaskSkuItem{
		&AskTaskSkuItem{
			SkuId: "1:1-2:1", // 上衣-X
			LabelIds: []string{
				"1:1",
				"2:1",
			},
			MarketPrice:         12,
			GroupBuyingPrice:    11,
			SettlementPrice:     10,
			CostPrice:           9,
			IllustrationPicture: "xxx.jpg",
			IsShow:              true,
			InventoryCount:      1000,
		},
		&AskTaskSkuItem{ // 上衣-L
			SkuId: "1:1-2:2",
			LabelIds: []string{
				"1:1",
				"2:2",
			},
			MarketPrice:         12,
			GroupBuyingPrice:    11,
			SettlementPrice:     10,
			CostPrice:           9,
			IllustrationPicture: "xxx.jpg",
			IsShow:              true,
			InventoryCount:      1000,
		},
		&AskTaskSkuItem{ // 裤子-X
			SkuId: "1:2-2:1",
			LabelIds: []string{
				"1:2",
				"2:1",
			},
			MarketPrice:         12,
			GroupBuyingPrice:    11,
			SettlementPrice:     10,
			CostPrice:           9,
			IllustrationPicture: "xxx.jpg",
			IsShow:              true,
			InventoryCount:      1000,
		},
		&AskTaskSkuItem{ // 裤子-L
			SkuId: "1:2-2:2",
			LabelIds: []string{
				"1:2",
				"2:2",
			},
			MarketPrice:         12,
			GroupBuyingPrice:    11,
			SettlementPrice:     10,
			CostPrice:           9,
			IllustrationPicture: "xxx.jpg",
			IsShow:              true,
			InventoryCount:      1000,
		},
	}

	askCombinationItems := []*AskTaskCombinationItem{
		&AskTaskCombinationItem{
			SubSkuItems: []*AskTaskCombinationItemSubSkuItem{
				&AskTaskCombinationItemSubSkuItem{
					SkuId: "1:1-2:1", // 上衣-X
					Count: 1,
				},
				&AskTaskCombinationItemSubSkuItem{
					SkuId: "1:2-2:1", // 裤子-x
					Count: 1,
				},
			},
			MarketPrice:         20,
			GroupBuyingPrice:    19,
			SettlementPrice:     18,
			CostPrice:           17,
			IllustrationPicture: "yyy.jpg",
			IsShow:              true,
		},
		&AskTaskCombinationItem{
			SubSkuItems: []*AskTaskCombinationItemSubSkuItem{
				&AskTaskCombinationItemSubSkuItem{
					SkuId: "1:1-2:2",
					Count: 1,
				},
				&AskTaskCombinationItemSubSkuItem{
					SkuId: "1:2-2:2",
					Count: 1,
				},
			},
			MarketPrice:         20,
			GroupBuyingPrice:    19,
			SettlementPrice:     18,
			CostPrice:           17,
			IllustrationPicture: "yyy.jpg",
			IsShow:              true,
		},
	}

	specification, err := NewGroupBuyingTaskOrderSpecificationByAsk(
		askSpecificationItems,
		askSkuItems,
		askCombinationItems,
	)
	if err != nil {
		t.Error(err)
		return
	}

	bytesSpecification, err := json.Marshal(specification)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(bytesSpecification))
}
