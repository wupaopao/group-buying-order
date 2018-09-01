package cidl

import "time"

type GroupBuyingOrderCombinationSubItem struct {
	SkuId  string                                        `db:"SkuId"`
	Labels map[string]*GroupBuyingSpecificationItemLabel `db:"Labels"`
	Count  uint32                                        `db:"Count"`
}

func NewGroupBuyingOrderCombinationSubItem() *GroupBuyingOrderCombinationSubItem {
	return &GroupBuyingOrderCombinationSubItem{
		Labels: make(map[string]*GroupBuyingSpecificationItemLabel),
	}
}

// 规格组合项
type GroupBuyingOrderTaskCombinationItem struct {
	SkuId               string                                         `db:"SkuId"`
	Name                string                                         `db:"Name"`
	SubItems            map[string]*GroupBuyingOrderCombinationSubItem `db:"SubItems"`
	MarketPrice         float64                                        `db:"MarketPrice"`
	GroupBuyingPrice    float64                                        `db:"GroupBuyingPrice"`
	SettlementPrice     float64                                        `db:"SettlementPrice"`
	CostPrice           float64                                        `db:"CostPrice"`
	IllustrationPicture string                                         `db:"IllustrationPicture"`
	IsShow              bool                                           `db:"IsShow"`
}

func NewGroupBuyingOrderTaskCombinationItem() *GroupBuyingOrderTaskCombinationItem {
	return &GroupBuyingOrderTaskCombinationItem{
		SubItems: make(map[string]*GroupBuyingOrderCombinationSubItem),
	}
}

// 商品规格组合
type GroupBuyingOrderTaskSpecification struct {
	GroupBuyingSpecification
	CombinationSkuMap map[string]*GroupBuyingOrderTaskCombinationItem `db:"CombinationSkuMap"`
}

func NewGroupBuyingOrderTaskSpecification() *GroupBuyingOrderTaskSpecification {
	return &GroupBuyingOrderTaskSpecification{
		CombinationSkuMap: make(map[string]*GroupBuyingOrderTaskCombinationItem),
	}
}

// 团购任务销售类型
type GroupBuyingOrderTaskSellType int

const (
	// 普通
	GroupBuyingOrderTaskSellTypeDefault GroupBuyingOrderTaskSellType = 1
	// 秒杀
	GroupBuyingOrderTaskSellTypeSeckill GroupBuyingOrderTaskSellType = 2
)

func (m GroupBuyingOrderTaskSellType) String() string {
	switch m {

	case GroupBuyingOrderTaskSellTypeDefault:
		return "GroupBuyingOrderTaskSellTypeDefault<enum GroupBuyingOrderTaskSellType>"
	case GroupBuyingOrderTaskSellTypeSeckill:
		return "GroupBuyingOrderTaskSellTypeSeckill<enum GroupBuyingOrderTaskSellType>"
	default:
		return "UNKNOWN_Name_<GroupBuyingOrderTaskSellType>"
	}
}

// 团购任务内容
type GroupBuyingOrderTaskContent struct {
	TaskId               uint32                             `db:"tsk_id"`
	OrganizationId       uint32                             `db:"org_id"`
	Title                string                             `db:"title"`
	Introduction         string                             `db:"introduction"`
	CoverPicture         string                             `db:"cover_picture"`
	IllustrationPictures *TaskIllustrationPicturesType      `db:"illustration_pictures"`
	Info                 *TaskInfoType                      `db:"info"`
	Specification        *GroupBuyingOrderTaskSpecification `db:"specification"`
	WxSellText           string                             `db:"wx_sell_text"`
	Notes                string                             `db:"notes"`
	ShowStartTime        time.Time                          `db:"show_start_time"`
	StartTime            time.Time                          `db:"start_time"`
	EndTime              time.Time                          `db:"end_time"`
	SellType             GroupBuyingOrderTaskSellType       `db:"sell_type"`
	Version              uint32                             `db:"version"`
	AllowCancel          GroupBuyingTaskAllowCancelState    `db:"allow_cancel"`
	TeamVisibleState     GroupBuyingTeamVisibleState        `db:"team_visible_state"`
	TeamIds              []uint32                           `db:"TeamIds"`
}

func NewGroupBuyingOrderTaskContent() *GroupBuyingOrderTaskContent {
	return &GroupBuyingOrderTaskContent{
		IllustrationPictures: NewTaskIllustrationPicturesType(),
		Info:                 NewTaskInfoType(),
		Specification:        NewGroupBuyingOrderTaskSpecification(),
		TeamIds:              make([]uint32, 0),
	}
}

// 团长下单模式团购任务
type GroupBuyingOrderTask struct {
	GroupBuyingOrderTaskContent
	ShowState          GroupBuyingTaskShowState  `db:"show_state"`
	GroupState         GroupBuyingTaskGroupState `db:"group_state"`
	OrderState         GroupBuyingTaskOrderState `db:"order_state"`
	Sales              uint32                    `db:"sales"`
	IsDelete           bool                      `db:"is_delete"`
	CreateTime         time.Time                 `db:"create_time"`
	LineList           []*GroupBuyingTaskLine    `db:"LineList"`
	IsSelectedAllLines bool                      `db:"IsSelectedAllLines"`
}

func NewGroupBuyingOrderTask() *GroupBuyingOrderTask {
	return &GroupBuyingOrderTask{
		LineList: make([]*GroupBuyingTaskLine, 0),
	}
}

// 团长下单模式库存
type GroupBuyingOrderInventory struct {
	TaskId  uint32 `db:"tsk_id"`
	SkuId   string `db:"sku_id"`
	Total   uint32 `db:"total"`
	Sales   uint32 `db:"sales"`
	Surplus uint32 `db:"surplus"`
}

func NewGroupBuyingOrderInventory() *GroupBuyingOrderInventory {
	return &GroupBuyingOrderInventory{}
}

// 社群订购记录表
type GroupBuyingOrderCommunityBuy struct {
	BuyId                 string                  `db:"cby_id"`
	OrderId               string                  `db:"order_id"`
	GroupId               uint32                  `db:"grp_id"`
	GroupOrderId          string                  `db:"grp_ord_id"`
	GroupName             string                  `db:"grp_name"`
	TaskId                uint32                  `db:"tsk_id"`
	TaskTitle             string                  `db:"tsk_title"`
	ManagerUserId         string                  `db:"manager_uid"`
	ManagerName           string                  `db:"manager_name"`
	ManagerMobile         string                  `db:"manager_mobile"`
	SkuId                 string                  `db:"sku_id"`
	BuyDetail             *CommunityBuyDetailType `db:"buy_detail"`
	Count                 uint32                  `db:"count"`
	TotalMarketPrice      float64                 `db:"total_market_price"`
	TotalGroupBuyingPrice float64                 `db:"total_group_buying_price"`
	TotalSettlementPrice  float64                 `db:"total_settlement_price"`
	TotalCostPrice        float64                 `db:"total_cost_price"`
	Version               uint32                  `db:"version"`
	CreateTime            time.Time               `db:"create_time"`
}

func NewGroupBuyingOrderCommunityBuy() *GroupBuyingOrderCommunityBuy {
	return &GroupBuyingOrderCommunityBuy{
		BuyDetail: NewCommunityBuyDetailType(),
	}
}

// 社群购买团购任务的订单统计
type GroupBuyingOrderCommunityBuyTask struct {
	TaskId                uint32                       `db:"tsk_id"`
	GroupId               uint32                       `db:"grp_id"`
	TaskDetail            *GroupBuyingOrderTaskContent `db:"task_detail"`
	OrderCount            uint32                       `db:"order_count"`
	GoodsCount            uint32                       `db:"goods_count"`
	TotalMarketPrice      float64                      `db:"total_market_price"`
	TotalGroupBuyingPrice float64                      `db:"total_group_buying_price"`
	TotalSettlementPrice  float64                      `db:"total_settlement_price"`
	TotalCostPrice        float64                      `db:"total_cost_price"`
	Version               uint32                       `db:"version"`
	CreateTime            time.Time                    `db:"create_time"`
}

func NewGroupBuyingOrderCommunityBuyTask() *GroupBuyingOrderCommunityBuyTask {
	return &GroupBuyingOrderCommunityBuyTask{
		TaskDetail: NewGroupBuyingOrderTaskContent(),
	}
}

// 订单
type GroupBuyingOrderCommunityOrder struct {
	OrderId               string                     `db:"ord_id"`
	GroupId               uint32                     `db:"grp_id"`
	GroupOrderId          string                     `db:"grp_ord_id"`
	GoodsDetail           *CommunityOrderGoodsDetail `db:"goods_detail"`
	Count                 uint32                     `db:"count"`
	TotalMarketPrice      float64                    `db:"total_market_price"`
	TotalGroupBuyingPrice float64                    `db:"total_group_buying_price"`
	TotalSettlementPrice  float64                    `db:"total_settlement_price"`
	TotalCostPrice        float64                    `db:"total_cost_price"`
	Version               uint32                     `db:"version"`
	CreateTime            time.Time                  `db:"create_time"`
	AllowCancel           bool                       `db:"allow_cancel"`
	Status                string                     `db:"Status"`
	IsCancel              bool                       `db:"is_cancel"`
}

func NewGroupBuyingOrderCommunityOrder() *GroupBuyingOrderCommunityOrder {
	return &GroupBuyingOrderCommunityOrder{
		GoodsDetail: NewCommunityOrderGoodsDetail(),
	}
}

// 社团购物车
type GroupBuyingOrderCommunityCart struct {
	CartId                string                  `db:"ccr_id"`
	GroupId               uint32                  `db:"grp_id"`
	TaskId                uint32                  `db:"tsk_id"`
	TaskTitle             string                  `db:"tsk_title"`
	SkuId                 string                  `db:"sku_id"`
	BuyDetail             *CommunityBuyDetailType `db:"buy_detail"`
	Count                 uint32                  `db:"count"`
	TotalMarketPrice      float64                 `db:"total_market_price"`
	TotalGroupBuyingPrice float64                 `db:"total_group_buying_price"`
	TotalSettlementPrice  float64                 `db:"total_settlement_price"`
	TotalCostPrice        float64                 `db:"total_cost_price"`
	Version               uint32                  `db:"version"`
	CreateTime            time.Time               `db:"create_time"`
}

func NewGroupBuyingOrderCommunityCart() *GroupBuyingOrderCommunityCart {
	return &GroupBuyingOrderCommunityCart{
		BuyDetail: NewCommunityBuyDetailType(),
	}
}

// 社群购买记录
type CommunityBuyDetailType struct {
	IsCombination    bool                                 `db:"IsCombination"`
	SingleGoods      *GroupBuyingSkuMapItem               `db:"SingleGoods"`
	CombinationGoods *GroupBuyingOrderTaskCombinationItem `db:"CombinationGoods"`
}

func NewCommunityBuyDetailType() *CommunityBuyDetailType {
	return &CommunityBuyDetailType{
		SingleGoods:      NewGroupBuyingSkuMapItem(),
		CombinationGoods: NewGroupBuyingOrderTaskCombinationItem(),
	}
}
