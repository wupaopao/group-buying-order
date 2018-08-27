package cidl

// 产品信息
type GroupBuyingTaskInfo struct {
	Title   string `db:"Title"`
	Content string `db:"Content"`
}

func NewGroupBuyingTaskInfo() *GroupBuyingTaskInfo {
	return &GroupBuyingTaskInfo{}
}

// 规格项Label
type GroupBuyingSpecificationItemLabel struct {
	LabelId string `db:"LabelId"`
	Name    string `db:"Name"`
}

func NewGroupBuyingSpecificationItemLabel() *GroupBuyingSpecificationItemLabel {
	return &GroupBuyingSpecificationItemLabel{}
}

// 商品规格项
type GroupBuyingSpecificationItem struct {
	ItemId         string                                        `db:"ItemId"`
	Name           string                                        `db:"Name"`
	LabelIdCounter uint32                                        `db:"LabelIdCounter"`
	Labels         map[string]*GroupBuyingSpecificationItemLabel `db:"Labels"`
}

func NewGroupBuyingSpecificationItem() *GroupBuyingSpecificationItem {
	return &GroupBuyingSpecificationItem{
		Labels: make(map[string]*GroupBuyingSpecificationItemLabel),
	}
}

// SKU映射
type GroupBuyingSkuMapItem struct {
	SkuId               string                                        `db:"SkuId"`
	Name                string                                        `db:"Name"`
	Labels              map[string]*GroupBuyingSpecificationItemLabel `db:"Labels"`
	MarketPrice         float64                                       `db:"MarketPrice"`
	GroupBuyingPrice    float64                                       `db:"GroupBuyingPrice"`
	SettlementPrice     float64                                       `db:"SettlementPrice"`
	CostPrice           float64                                       `db:"CostPrice"`
	IllustrationPicture string                                        `db:"IllustrationPicture"`
	IsShow              bool                                          `db:"IsShow"`
	InventoryCount      uint32                                        `db:"InventoryCount"`
}

func NewGroupBuyingSkuMapItem() *GroupBuyingSkuMapItem {
	return &GroupBuyingSkuMapItem{
		Labels: make(map[string]*GroupBuyingSpecificationItemLabel),
	}
}

// 商品规格价格区间
type SpecificationPriceRange struct {
	Min float64 `db:"Min"`
	Max float64 `db:"Max"`
}

func NewSpecificationPriceRange() *SpecificationPriceRange {
	return &SpecificationPriceRange{}
}

// 商品规格
type GroupBuyingSpecification struct {
	ItemIdCounter         uint32                                   `db:"ItemIdCounter"`
	Items                 map[string]*GroupBuyingSpecificationItem `db:"Items"`
	SkuMap                map[string]*GroupBuyingSkuMapItem        `db:"SkuMap"`
	MarketPriceRange      *SpecificationPriceRange                 `db:"MarketPriceRange"`
	GroupBuyingPriceRange *SpecificationPriceRange                 `db:"GroupBuyingPriceRange"`
	SettlementPriceRange  *SpecificationPriceRange                 `db:"SettlementPriceRange"`
	CostPriceRange        *SpecificationPriceRange                 `db:"CostPriceRange"`
}

func NewGroupBuyingSpecification() *GroupBuyingSpecification {
	return &GroupBuyingSpecification{
		Items:                 make(map[string]*GroupBuyingSpecificationItem),
		SkuMap:                make(map[string]*GroupBuyingSkuMapItem),
		MarketPriceRange:      NewSpecificationPriceRange(),
		GroupBuyingPriceRange: NewSpecificationPriceRange(),
		SettlementPriceRange:  NewSpecificationPriceRange(),
		CostPriceRange:        NewSpecificationPriceRange(),
	}
}

// 显示状态
type GroupBuyingTaskShowState int

const (
	// 上架
	GroupBuyingTaskShowStateShow GroupBuyingTaskShowState = 1
	// 下架
	GroupBuyingTaskShowStateHidden GroupBuyingTaskShowState = 0
)

func (m GroupBuyingTaskShowState) String() string {
	switch m {

	case GroupBuyingTaskShowStateShow:
		return "GroupBuyingTaskShowStateShow<enum GroupBuyingTaskShowState>"
	case GroupBuyingTaskShowStateHidden:
		return "GroupBuyingTaskShowStateHidden<enum GroupBuyingTaskShowState>"
	default:
		return "UNKNOWN_Name_<GroupBuyingTaskShowState>"
	}
}

// 社团群组可见状态
type GroupBuyingTeamVisibleState int

const (
	// 全部可见
	GroupBuyingTeamVisibleStateAll GroupBuyingTeamVisibleState = 1
	// 部分可见
	GroupBuyingTeamVisibleStatePart GroupBuyingTeamVisibleState = 0
)

func (m GroupBuyingTeamVisibleState) String() string {
	switch m {

	case GroupBuyingTeamVisibleStateAll:
		return "GroupBuyingTeamVisibleStateAll<enum GroupBuyingTeamVisibleState>"
	case GroupBuyingTeamVisibleStatePart:
		return "GroupBuyingTeamVisibleStatePart<enum GroupBuyingTeamVisibleState>"
	default:
		return "UNKNOWN_Name_<GroupBuyingTeamVisibleState>"
	}
}

type GroupBuyingTaskAllowCancelState int

const (
	//允许取消订单
	GroupBuyingTaskAllowCancel GroupBuyingTaskAllowCancelState = 1
	//不允许取消订单
	GroupBuyingTaskNotAllowCancel GroupBuyingTaskAllowCancelState = 0
)

func (m GroupBuyingTaskAllowCancelState) String() string {
	switch m {

	case GroupBuyingTaskAllowCancel:
		return "GroupBuyingTaskAllowCancel<enum GroupBuyingTaskAllowCancelState>"
	case GroupBuyingTaskNotAllowCancel:
		return "GroupBuyingTaskNotAllowCancel<enum GroupBuyingTaskAllowCancelState>"
	default:
		return "UNKNOWN_Name_<GroupBuyingTaskAllowCancelState>"
	}
}

// 团购任务状态
type GroupBuyingTaskGroupState int

const (
	// 未开团
	GroupBuyingTaskGroupStateNotStart GroupBuyingTaskGroupState = 0
	// 进行中
	GroupBuyingTaskGroupStateInProgress GroupBuyingTaskGroupState = 1
	// 已截单
	GroupBuyingTaskGroupStateFinishOrdering GroupBuyingTaskGroupState = 2
	// 已结团
	GroupBuyingTaskGroupStateFinishBuying GroupBuyingTaskGroupState = 3
	// 已取消
	GroupBuyingTaskGroupStateCancel GroupBuyingTaskGroupState = 4
	// 已配送
	GroupBuyingTaskGroupStateDelivered GroupBuyingTaskGroupState = 5
)

func (m GroupBuyingTaskGroupState) String() string {
	switch m {

	case GroupBuyingTaskGroupStateNotStart:
		return "GroupBuyingTaskGroupStateNotStart<enum GroupBuyingTaskGroupState>"
	case GroupBuyingTaskGroupStateInProgress:
		return "GroupBuyingTaskGroupStateInProgress<enum GroupBuyingTaskGroupState>"
	case GroupBuyingTaskGroupStateFinishOrdering:
		return "GroupBuyingTaskGroupStateFinishOrdering<enum GroupBuyingTaskGroupState>"
	case GroupBuyingTaskGroupStateFinishBuying:
		return "GroupBuyingTaskGroupStateFinishBuying<enum GroupBuyingTaskGroupState>"
	case GroupBuyingTaskGroupStateCancel:
		return "GroupBuyingTaskGroupStateCancel<enum GroupBuyingTaskGroupState>"
	case GroupBuyingTaskGroupStateDelivered:
		return "GroupBuyingTaskGroupStateDelivered<enum GroupBuyingTaskGroupState>"
	default:
		return "UNKNOWN_Name_<GroupBuyingTaskGroupState>"
	}
}

// 下单送货状态
type GroupBuyingTaskOrderState int

const (
	// 未订货
	GroupBuyingTaskOrderStateNotOrder GroupBuyingTaskOrderState = 0
	// 已订货
	GroupBuyingTaskOrderStateOrdered GroupBuyingTaskOrderState = 1
	// 未配送
	GroupBuyingTaskOrderStateNotDeliver GroupBuyingTaskOrderState = 2
	// 已配送
	GroupBuyingTaskOrderStateDelivered GroupBuyingTaskOrderState = 3
	// 已完成
	GroupBuyingTaskOrderStateFinish GroupBuyingTaskOrderState = 4
)

func (m GroupBuyingTaskOrderState) String() string {
	switch m {

	case GroupBuyingTaskOrderStateNotOrder:
		return "GroupBuyingTaskOrderStateNotOrder<enum GroupBuyingTaskOrderState>"
	case GroupBuyingTaskOrderStateOrdered:
		return "GroupBuyingTaskOrderStateOrdered<enum GroupBuyingTaskOrderState>"
	case GroupBuyingTaskOrderStateNotDeliver:
		return "GroupBuyingTaskOrderStateNotDeliver<enum GroupBuyingTaskOrderState>"
	case GroupBuyingTaskOrderStateDelivered:
		return "GroupBuyingTaskOrderStateDelivered<enum GroupBuyingTaskOrderState>"
	case GroupBuyingTaskOrderStateFinish:
		return "GroupBuyingTaskOrderStateFinish<enum GroupBuyingTaskOrderState>"
	default:
		return "UNKNOWN_Name_<GroupBuyingTaskOrderState>"
	}
}
