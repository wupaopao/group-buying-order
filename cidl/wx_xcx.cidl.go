package cidl

import "time"

// 团购任务列表结构参考admin.cidl和admin.md
// 今日团购
type AckWxXcxTaskListItem struct {
	TaskId                uint32                             `db:"TaskId"`
	OrganizationID        uint32                             `db:"OrganizationID"`
	Title                 string                             `db:"Title"`
	Specification         *GroupBuyingOrderTaskSpecification `db:"Specification"`
	CoverPicture          string                             `db:"CoverPicture"`
	MarketPriceRange      *SpecificationPriceRange           `db:"MarketPriceRange"`
	GroupBuyingPriceRange *SpecificationPriceRange           `db:"GroupBuyingPriceRange"`
	SettlementPriceRange  *SpecificationPriceRange           `db:"SettlementPriceRange"`
	CostPriceRange        *SpecificationPriceRange           `db:"CostPriceRange"`
	StartTime             time.Time                          `db:"StartTime"`
	EndTime               time.Time                          `db:"EndTime"`
	Sales                 uint32                             `db:"Sales"`
	Notes                 string                             `db:"Notes"`
	SellType              GroupBuyingOrderTaskSellType       `db:"SellType"`
	GroupState            GroupBuyingTaskGroupState          `db:"group_state"`
	OrderState            GroupBuyingTaskOrderState          `db:"order_state"`
}

func NewAckWxXcxTaskListItem() *AckWxXcxTaskListItem {
	return &AckWxXcxTaskListItem{
		Specification:         NewGroupBuyingOrderTaskSpecification(),
		MarketPriceRange:      NewSpecificationPriceRange(),
		GroupBuyingPriceRange: NewSpecificationPriceRange(),
		SettlementPriceRange:  NewSpecificationPriceRange(),
		CostPriceRange:        NewSpecificationPriceRange(),
	}
}

type AckWxXcxTaskTodayListByOrganizationID struct {
	Count uint32                  `db:"Count"`
	List  []*AckWxXcxTaskListItem `db:"List"`
}

func NewAckWxXcxTaskTodayListByOrganizationID() *AckWxXcxTaskTodayListByOrganizationID {
	return &AckWxXcxTaskTodayListByOrganizationID{
		List: make([]*AckWxXcxTaskListItem, 0),
	}
}

type MetaApiWxXcxTaskTodayListByOrganizationID struct {
}

var META_WX_XCX_TASK_TODAY_LIST_BY_ORGANIZATION_ID = &MetaApiWxXcxTaskTodayListByOrganizationID{}

func (m *MetaApiWxXcxTaskTodayListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxTaskTodayListByOrganizationID) GetURL() string {
	return "/group_buying_order/wx_xcx/task/today_list/:organization_id"
}
func (m *MetaApiWxXcxTaskTodayListByOrganizationID) GetName() string {
	return "WxXcxTaskTodayListByOrganizationID"
}
func (m *MetaApiWxXcxTaskTodayListByOrganizationID) GetType() string { return "json" }

// 今日团购
type ApiWxXcxTaskTodayListByOrganizationID struct {
	MetaApiWxXcxTaskTodayListByOrganizationID
	Ack    *AckWxXcxTaskTodayListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiWxXcxTaskTodayListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiWxXcxTaskTodayListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxTaskTodayListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxTaskTodayListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxTaskTodayListByOrganizationID() ApiWxXcxTaskTodayListByOrganizationID {
	return ApiWxXcxTaskTodayListByOrganizationID{
		Ack: NewAckWxXcxTaskTodayListByOrganizationID(),
	}
}

type AckWxXcxTaskFutureListByOrganizationID struct {
	Count uint32                  `db:"Count"`
	List  []*AckWxXcxTaskListItem `db:"List"`
}

func NewAckWxXcxTaskFutureListByOrganizationID() *AckWxXcxTaskFutureListByOrganizationID {
	return &AckWxXcxTaskFutureListByOrganizationID{
		List: make([]*AckWxXcxTaskListItem, 0),
	}
}

type MetaApiWxXcxTaskFutureListByOrganizationID struct {
}

var META_WX_XCX_TASK_FUTURE_LIST_BY_ORGANIZATION_ID = &MetaApiWxXcxTaskFutureListByOrganizationID{}

func (m *MetaApiWxXcxTaskFutureListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxTaskFutureListByOrganizationID) GetURL() string {
	return "/group_buying_order/wx_xcx/task/future_list/:organization_id"
}
func (m *MetaApiWxXcxTaskFutureListByOrganizationID) GetName() string {
	return "WxXcxTaskFutureListByOrganizationID"
}
func (m *MetaApiWxXcxTaskFutureListByOrganizationID) GetType() string { return "json" }

// 未来团购
type ApiWxXcxTaskFutureListByOrganizationID struct {
	MetaApiWxXcxTaskFutureListByOrganizationID
	Ack    *AckWxXcxTaskFutureListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiWxXcxTaskFutureListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiWxXcxTaskFutureListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxTaskFutureListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxTaskFutureListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxTaskFutureListByOrganizationID() ApiWxXcxTaskFutureListByOrganizationID {
	return ApiWxXcxTaskFutureListByOrganizationID{
		Ack: NewAckWxXcxTaskFutureListByOrganizationID(),
	}
}

type AckWxXcxTaskHistoryListByOrganizationID struct {
	Count uint32                  `db:"Count"`
	List  []*AckWxXcxTaskListItem `db:"List"`
}

func NewAckWxXcxTaskHistoryListByOrganizationID() *AckWxXcxTaskHistoryListByOrganizationID {
	return &AckWxXcxTaskHistoryListByOrganizationID{
		List: make([]*AckWxXcxTaskListItem, 0),
	}
}

type MetaApiWxXcxTaskHistoryListByOrganizationID struct {
}

var META_WX_XCX_TASK_HISTORY_LIST_BY_ORGANIZATION_ID = &MetaApiWxXcxTaskHistoryListByOrganizationID{}

func (m *MetaApiWxXcxTaskHistoryListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxTaskHistoryListByOrganizationID) GetURL() string {
	return "/group_buying_order/wx_xcx/task/history_list/:organization_id"
}
func (m *MetaApiWxXcxTaskHistoryListByOrganizationID) GetName() string {
	return "WxXcxTaskHistoryListByOrganizationID"
}
func (m *MetaApiWxXcxTaskHistoryListByOrganizationID) GetType() string { return "json" }

// 历史团购
type ApiWxXcxTaskHistoryListByOrganizationID struct {
	MetaApiWxXcxTaskHistoryListByOrganizationID
	Ack    *AckWxXcxTaskHistoryListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiWxXcxTaskHistoryListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiWxXcxTaskHistoryListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxTaskHistoryListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxTaskHistoryListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxTaskHistoryListByOrganizationID() ApiWxXcxTaskHistoryListByOrganizationID {
	return ApiWxXcxTaskHistoryListByOrganizationID{
		Ack: NewAckWxXcxTaskHistoryListByOrganizationID(),
	}
}

type MetaApiWxXcxTaskInfoByTaskID struct {
}

var META_WX_XCX_TASK_INFO_BY_TASK_ID = &MetaApiWxXcxTaskInfoByTaskID{}

func (m *MetaApiWxXcxTaskInfoByTaskID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxTaskInfoByTaskID) GetURL() string {
	return "/group_buying_order/wx_xcx/task/info/:task_id"
}
func (m *MetaApiWxXcxTaskInfoByTaskID) GetName() string { return "WxXcxTaskInfoByTaskID" }
func (m *MetaApiWxXcxTaskInfoByTaskID) GetType() string { return "json" }

// 获取团购任务
type ApiWxXcxTaskInfoByTaskID struct {
	MetaApiWxXcxTaskInfoByTaskID
	Ack    *GroupBuyingOrderTask
	Params struct {
		TaskID uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiWxXcxTaskInfoByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxTaskInfoByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxTaskInfoByTaskID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxTaskInfoByTaskID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxTaskInfoByTaskID() ApiWxXcxTaskInfoByTaskID {
	return ApiWxXcxTaskInfoByTaskID{
		Ack: NewGroupBuyingOrderTask(),
	}
}

type AskWxXcxTaskBatchWxSellTextByOrganizationID struct {
	TaskIds []uint32 `binding:"required,gt=0,lte=50" db:"TaskIds"`
}

func NewAskWxXcxTaskBatchWxSellTextByOrganizationID() *AskWxXcxTaskBatchWxSellTextByOrganizationID {
	return &AskWxXcxTaskBatchWxSellTextByOrganizationID{
		TaskIds: make([]uint32, 0),
	}
}

type AckWxXcxTaskBatchWxSellTextByOrganizationID struct {
	List []string `db:"List"`
}

func NewAckWxXcxTaskBatchWxSellTextByOrganizationID() *AckWxXcxTaskBatchWxSellTextByOrganizationID {
	return &AckWxXcxTaskBatchWxSellTextByOrganizationID{
		List: make([]string, 0),
	}
}

type MetaApiWxXcxTaskBatchWxSellTextByOrganizationID struct {
}

var META_WX_XCX_TASK_BATCH_WX_SELL_TEXT_BY_ORGANIZATION_ID = &MetaApiWxXcxTaskBatchWxSellTextByOrganizationID{}

func (m *MetaApiWxXcxTaskBatchWxSellTextByOrganizationID) GetMethod() string { return "POST" }
func (m *MetaApiWxXcxTaskBatchWxSellTextByOrganizationID) GetURL() string {
	return "/group_buying_order/wx_xcx/task/batch_wx_sell_text/:organization_id"
}
func (m *MetaApiWxXcxTaskBatchWxSellTextByOrganizationID) GetName() string {
	return "WxXcxTaskBatchWxSellTextByOrganizationID"
}
func (m *MetaApiWxXcxTaskBatchWxSellTextByOrganizationID) GetType() string { return "json" }

// 获取多个团购任务的微信销售文案
type ApiWxXcxTaskBatchWxSellTextByOrganizationID struct {
	MetaApiWxXcxTaskBatchWxSellTextByOrganizationID
	Ask    *AskWxXcxTaskBatchWxSellTextByOrganizationID
	Ack    *AckWxXcxTaskBatchWxSellTextByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
}

func (m *ApiWxXcxTaskBatchWxSellTextByOrganizationID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxTaskBatchWxSellTextByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxTaskBatchWxSellTextByOrganizationID) GetAsk() interface{}    { return m.Ask }
func (m *ApiWxXcxTaskBatchWxSellTextByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxTaskBatchWxSellTextByOrganizationID() ApiWxXcxTaskBatchWxSellTextByOrganizationID {
	return ApiWxXcxTaskBatchWxSellTextByOrganizationID{
		Ask: NewAskWxXcxTaskBatchWxSellTextByOrganizationID(),
		Ack: NewAckWxXcxTaskBatchWxSellTextByOrganizationID(),
	}
}

type AskWxXcxTaskAddCartByGroupIDByTaskID struct {
	IsCombination bool   `db:"IsCombination"`
	SkuId         string `binding:"required" db:"SkuId"`
	Count         uint32 `binding:"required,gt=0" db:"Count"`
}

func NewAskWxXcxTaskAddCartByGroupIDByTaskID() *AskWxXcxTaskAddCartByGroupIDByTaskID {
	return &AskWxXcxTaskAddCartByGroupIDByTaskID{}
}

type MetaApiWxXcxTaskAddCartByGroupIDByTaskID struct {
}

var META_WX_XCX_TASK_ADD_CART_BY_GROUP_ID_BY_TASK_ID = &MetaApiWxXcxTaskAddCartByGroupIDByTaskID{}

func (m *MetaApiWxXcxTaskAddCartByGroupIDByTaskID) GetMethod() string { return "POST" }
func (m *MetaApiWxXcxTaskAddCartByGroupIDByTaskID) GetURL() string {
	return "/group_buying_order/wx_xcx/cart/add_cart/:group_id/:task_id"
}
func (m *MetaApiWxXcxTaskAddCartByGroupIDByTaskID) GetName() string {
	return "WxXcxTaskAddCartByGroupIDByTaskID"
}
func (m *MetaApiWxXcxTaskAddCartByGroupIDByTaskID) GetType() string { return "json" }

// 添加到购物车
type ApiWxXcxTaskAddCartByGroupIDByTaskID struct {
	MetaApiWxXcxTaskAddCartByGroupIDByTaskID
	Ask    *AskWxXcxTaskAddCartByGroupIDByTaskID
	Params struct {
		GroupID uint32 `form:"group_id" binding:"required,gt=0" db:"GroupID"`
		TaskID  uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiWxXcxTaskAddCartByGroupIDByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxTaskAddCartByGroupIDByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxTaskAddCartByGroupIDByTaskID) GetAsk() interface{}    { return m.Ask }
func (m *ApiWxXcxTaskAddCartByGroupIDByTaskID) GetAck() interface{}    { return nil }
func MakeApiWxXcxTaskAddCartByGroupIDByTaskID() ApiWxXcxTaskAddCartByGroupIDByTaskID {
	return ApiWxXcxTaskAddCartByGroupIDByTaskID{
		Ask: NewAskWxXcxTaskAddCartByGroupIDByTaskID(),
	}
}

type AskWxXcxCartDeleteCartByGroupID struct {
	CartIds []string `db:"CartIds"`
}

func NewAskWxXcxCartDeleteCartByGroupID() *AskWxXcxCartDeleteCartByGroupID {
	return &AskWxXcxCartDeleteCartByGroupID{
		CartIds: make([]string, 0),
	}
}

type MetaApiWxXcxCartDeleteCartByGroupID struct {
}

var META_WX_XCX_CART_DELETE_CART_BY_GROUP_ID = &MetaApiWxXcxCartDeleteCartByGroupID{}

func (m *MetaApiWxXcxCartDeleteCartByGroupID) GetMethod() string { return "POST" }
func (m *MetaApiWxXcxCartDeleteCartByGroupID) GetURL() string {
	return "/group_buying_order/wx_xcx/cart/delete_cart/:group_id"
}
func (m *MetaApiWxXcxCartDeleteCartByGroupID) GetName() string { return "WxXcxCartDeleteCartByGroupID" }
func (m *MetaApiWxXcxCartDeleteCartByGroupID) GetType() string { return "json" }

// 从购物车中删除
type ApiWxXcxCartDeleteCartByGroupID struct {
	MetaApiWxXcxCartDeleteCartByGroupID
	Ask    *AskWxXcxCartDeleteCartByGroupID
	Params struct {
		GroupID uint32 `form:"group_id" binding:"required,gt=0" db:"GroupID"`
	}
}

func (m *ApiWxXcxCartDeleteCartByGroupID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxCartDeleteCartByGroupID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxCartDeleteCartByGroupID) GetAsk() interface{}    { return m.Ask }
func (m *ApiWxXcxCartDeleteCartByGroupID) GetAck() interface{}    { return nil }
func MakeApiWxXcxCartDeleteCartByGroupID() ApiWxXcxCartDeleteCartByGroupID {
	return ApiWxXcxCartDeleteCartByGroupID{
		Ask: NewAskWxXcxCartDeleteCartByGroupID(),
	}
}

type AskWxXcxCartChangeCountByGroupID struct {
	CartId string `db:"CartId"`
	Count  uint32 `db:"Count"`
}

func NewAskWxXcxCartChangeCountByGroupID() *AskWxXcxCartChangeCountByGroupID {
	return &AskWxXcxCartChangeCountByGroupID{}
}

type MetaApiWxXcxCartChangeCountByGroupID struct {
}

var META_WX_XCX_CART_CHANGE_COUNT_BY_GROUP_ID = &MetaApiWxXcxCartChangeCountByGroupID{}

func (m *MetaApiWxXcxCartChangeCountByGroupID) GetMethod() string { return "POST" }
func (m *MetaApiWxXcxCartChangeCountByGroupID) GetURL() string {
	return "/group_buying_order/wx_xcx/cart/change_count/:group_id"
}
func (m *MetaApiWxXcxCartChangeCountByGroupID) GetName() string {
	return "WxXcxCartChangeCountByGroupID"
}
func (m *MetaApiWxXcxCartChangeCountByGroupID) GetType() string { return "json" }

// 修改购买数目
type ApiWxXcxCartChangeCountByGroupID struct {
	MetaApiWxXcxCartChangeCountByGroupID
	Ask    *AskWxXcxCartChangeCountByGroupID
	Params struct {
		GroupID uint32 `form:"group_id" binding:"required,gt=0" db:"GroupID"`
	}
}

func (m *ApiWxXcxCartChangeCountByGroupID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxCartChangeCountByGroupID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxCartChangeCountByGroupID) GetAsk() interface{}    { return m.Ask }
func (m *ApiWxXcxCartChangeCountByGroupID) GetAck() interface{}    { return nil }
func MakeApiWxXcxCartChangeCountByGroupID() ApiWxXcxCartChangeCountByGroupID {
	return ApiWxXcxCartChangeCountByGroupID{
		Ask: NewAskWxXcxCartChangeCountByGroupID(),
	}
}

type AckWxXcxTaskCartCountByGroupID struct {
	Count uint32 `db:"Count"`
}

func NewAckWxXcxTaskCartCountByGroupID() *AckWxXcxTaskCartCountByGroupID {
	return &AckWxXcxTaskCartCountByGroupID{}
}

type MetaApiWxXcxTaskCartCountByGroupID struct {
}

var META_WX_XCX_TASK_CART_COUNT_BY_GROUP_ID = &MetaApiWxXcxTaskCartCountByGroupID{}

func (m *MetaApiWxXcxTaskCartCountByGroupID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxTaskCartCountByGroupID) GetURL() string {
	return "/group_buying_order/wx_xcx/cart/cart_count/:group_id"
}
func (m *MetaApiWxXcxTaskCartCountByGroupID) GetName() string { return "WxXcxTaskCartCountByGroupID" }
func (m *MetaApiWxXcxTaskCartCountByGroupID) GetType() string { return "json" }

// 购物车数目
type ApiWxXcxTaskCartCountByGroupID struct {
	MetaApiWxXcxTaskCartCountByGroupID
	Ack    *AckWxXcxTaskCartCountByGroupID
	Params struct {
		GroupID uint32 `form:"group_id" binding:"required,gt=0" db:"GroupID"`
	}
}

func (m *ApiWxXcxTaskCartCountByGroupID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxTaskCartCountByGroupID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxTaskCartCountByGroupID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxTaskCartCountByGroupID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxTaskCartCountByGroupID() ApiWxXcxTaskCartCountByGroupID {
	return ApiWxXcxTaskCartCountByGroupID{
		Ack: NewAckWxXcxTaskCartCountByGroupID(),
	}
}

type AckWxXcxCartCartListByGroupID struct {
	Count uint32                           `db:"Count"`
	List  []*GroupBuyingOrderCommunityCart `db:"List"`
}

func NewAckWxXcxCartCartListByGroupID() *AckWxXcxCartCartListByGroupID {
	return &AckWxXcxCartCartListByGroupID{
		List: make([]*GroupBuyingOrderCommunityCart, 0),
	}
}

type MetaApiWxXcxCartCartListByGroupID struct {
}

var META_WX_XCX_CART_CART_LIST_BY_GROUP_ID = &MetaApiWxXcxCartCartListByGroupID{}

func (m *MetaApiWxXcxCartCartListByGroupID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxCartCartListByGroupID) GetURL() string {
	return "/group_buying_order/wx_xcx/cart/cart_list/:group_id"
}
func (m *MetaApiWxXcxCartCartListByGroupID) GetName() string { return "WxXcxCartCartListByGroupID" }
func (m *MetaApiWxXcxCartCartListByGroupID) GetType() string { return "json" }

// 购物车列表
type ApiWxXcxCartCartListByGroupID struct {
	MetaApiWxXcxCartCartListByGroupID
	Ack    *AckWxXcxCartCartListByGroupID
	Params struct {
		GroupID uint32 `form:"group_id" db:"GroupID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiWxXcxCartCartListByGroupID) GetQuery() interface{}  { return &m.Query }
func (m *ApiWxXcxCartCartListByGroupID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxCartCartListByGroupID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxCartCartListByGroupID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxCartCartListByGroupID() ApiWxXcxCartCartListByGroupID {
	return ApiWxXcxCartCartListByGroupID{
		Ack: NewAckWxXcxCartCartListByGroupID(),
	}
}

type AskWxXcxOrderAddOrderByGroupID struct {
	CartIds []string `db:"CartIds"`
}

func NewAskWxXcxOrderAddOrderByGroupID() *AskWxXcxOrderAddOrderByGroupID {
	return &AskWxXcxOrderAddOrderByGroupID{
		CartIds: make([]string, 0),
	}
}

type AckWxXcxOrderAddOrderByGroupID struct {
	SuccessCartIds []string `db:"SuccessCartIds"`
}

func NewAckWxXcxOrderAddOrderByGroupID() *AckWxXcxOrderAddOrderByGroupID {
	return &AckWxXcxOrderAddOrderByGroupID{
		SuccessCartIds: make([]string, 0),
	}
}

type MetaApiWxXcxOrderAddOrderByGroupID struct {
}

var META_WX_XCX_ORDER_ADD_ORDER_BY_GROUP_ID = &MetaApiWxXcxOrderAddOrderByGroupID{}

func (m *MetaApiWxXcxOrderAddOrderByGroupID) GetMethod() string { return "POST" }
func (m *MetaApiWxXcxOrderAddOrderByGroupID) GetURL() string {
	return "/group_buying_order/wx_xcx/order/add_order/:group_id"
}
func (m *MetaApiWxXcxOrderAddOrderByGroupID) GetName() string { return "WxXcxOrderAddOrderByGroupID" }
func (m *MetaApiWxXcxOrderAddOrderByGroupID) GetType() string { return "json" }

// 从购物车提交订单
type ApiWxXcxOrderAddOrderByGroupID struct {
	MetaApiWxXcxOrderAddOrderByGroupID
	Ask    *AskWxXcxOrderAddOrderByGroupID
	Ack    *AckWxXcxOrderAddOrderByGroupID
	Params struct {
		GroupID uint32 `form:"group_id" binding:"required,gt=0" db:"GroupID"`
	}
}

func (m *ApiWxXcxOrderAddOrderByGroupID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxOrderAddOrderByGroupID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxOrderAddOrderByGroupID) GetAsk() interface{}    { return m.Ask }
func (m *ApiWxXcxOrderAddOrderByGroupID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxOrderAddOrderByGroupID() ApiWxXcxOrderAddOrderByGroupID {
	return ApiWxXcxOrderAddOrderByGroupID{
		Ask: NewAskWxXcxOrderAddOrderByGroupID(),
		Ack: NewAckWxXcxOrderAddOrderByGroupID(),
	}
}

// 直接提交订单
type AskDirectlyAddOrderItem struct {
	TaskId        uint32 `binding:"required,gt=0" db:"TaskId"`
	SkuId         string `binding:"required" db:"SkuId"`
	IsCombination bool   `db:"IsCombination"`
	Count         uint32 `binding:"required,gt=0" db:"Count"`
}

func NewAskDirectlyAddOrderItem() *AskDirectlyAddOrderItem {
	return &AskDirectlyAddOrderItem{}
}

type AckDirectlyAddOrderResultItem struct {
	TaskId  uint32 `db:"TaskId"`
	SkuId   string `db:"SkuId"`
	Message string `db:"Message"`
}

func NewAckDirectlyAddOrderResultItem() *AckDirectlyAddOrderResultItem {
	return &AckDirectlyAddOrderResultItem{}
}

type AskWxXcxOrderDirectlyAddByGroupID struct {
	Items []*AskDirectlyAddOrderItem `binding:"required,dive,required" db:"Items"`
}

func NewAskWxXcxOrderDirectlyAddByGroupID() *AskWxXcxOrderDirectlyAddByGroupID {
	return &AskWxXcxOrderDirectlyAddByGroupID{
		Items: make([]*AskDirectlyAddOrderItem, 0),
	}
}

type AckWxXcxOrderDirectlyAddByGroupID struct {
	OrderId   string                           `db:"OrderId"`
	ErrorList []*AckDirectlyAddOrderResultItem `db:"ErrorList"`
}

func NewAckWxXcxOrderDirectlyAddByGroupID() *AckWxXcxOrderDirectlyAddByGroupID {
	return &AckWxXcxOrderDirectlyAddByGroupID{
		ErrorList: make([]*AckDirectlyAddOrderResultItem, 0),
	}
}

type MetaApiWxXcxOrderDirectlyAddByGroupID struct {
}

var META_WX_XCX_ORDER_DIRECTLY_ADD_BY_GROUP_ID = &MetaApiWxXcxOrderDirectlyAddByGroupID{}

func (m *MetaApiWxXcxOrderDirectlyAddByGroupID) GetMethod() string { return "POST" }
func (m *MetaApiWxXcxOrderDirectlyAddByGroupID) GetURL() string {
	return "/group_buying_order/wx_xcx/order/directly_add/:group_id"
}
func (m *MetaApiWxXcxOrderDirectlyAddByGroupID) GetName() string {
	return "WxXcxOrderDirectlyAddByGroupID"
}
func (m *MetaApiWxXcxOrderDirectlyAddByGroupID) GetType() string { return "json" }

type ApiWxXcxOrderDirectlyAddByGroupID struct {
	MetaApiWxXcxOrderDirectlyAddByGroupID
	Ask    *AskWxXcxOrderDirectlyAddByGroupID
	Ack    *AckWxXcxOrderDirectlyAddByGroupID
	Params struct {
		GroupID uint32 `form:"group_id" binding:"required,gt=0" db:"GroupID"`
	}
}

func (m *ApiWxXcxOrderDirectlyAddByGroupID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxOrderDirectlyAddByGroupID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxOrderDirectlyAddByGroupID) GetAsk() interface{}    { return m.Ask }
func (m *ApiWxXcxOrderDirectlyAddByGroupID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxOrderDirectlyAddByGroupID() ApiWxXcxOrderDirectlyAddByGroupID {
	return ApiWxXcxOrderDirectlyAddByGroupID{
		Ask: NewAskWxXcxOrderDirectlyAddByGroupID(),
		Ack: NewAckWxXcxOrderDirectlyAddByGroupID(),
	}
}

type AckWxXcxOrderOrderListByGroupID struct {
	Count uint32                            `db:"Count"`
	List  []*GroupBuyingOrderCommunityOrder `db:"List"`
}

func NewAckWxXcxOrderOrderListByGroupID() *AckWxXcxOrderOrderListByGroupID {
	return &AckWxXcxOrderOrderListByGroupID{
		List: make([]*GroupBuyingOrderCommunityOrder, 0),
	}
}

type MetaApiWxXcxOrderOrderListByGroupID struct {
}

var META_WX_XCX_ORDER_ORDER_LIST_BY_GROUP_ID = &MetaApiWxXcxOrderOrderListByGroupID{}

func (m *MetaApiWxXcxOrderOrderListByGroupID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxOrderOrderListByGroupID) GetURL() string {
	return "/group_buying_order/wx_xcx/order/order_list/:group_id"
}
func (m *MetaApiWxXcxOrderOrderListByGroupID) GetName() string { return "WxXcxOrderOrderListByGroupID" }
func (m *MetaApiWxXcxOrderOrderListByGroupID) GetType() string { return "json" }

// 订单列表
type ApiWxXcxOrderOrderListByGroupID struct {
	MetaApiWxXcxOrderOrderListByGroupID
	Ack    *AckWxXcxOrderOrderListByGroupID
	Params struct {
		GroupID uint32 `form:"group_id" binding:"required,gt=0" db:"GroupID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiWxXcxOrderOrderListByGroupID) GetQuery() interface{}  { return &m.Query }
func (m *ApiWxXcxOrderOrderListByGroupID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxOrderOrderListByGroupID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxOrderOrderListByGroupID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxOrderOrderListByGroupID() ApiWxXcxOrderOrderListByGroupID {
	return ApiWxXcxOrderOrderListByGroupID{
		Ack: NewAckWxXcxOrderOrderListByGroupID(),
	}
}

type AckWxXcxBuyTaskListByGroupID struct {
	Count uint32                              `db:"Count"`
	List  []*GroupBuyingOrderCommunityBuyTask `db:"List"`
}

func NewAckWxXcxBuyTaskListByGroupID() *AckWxXcxBuyTaskListByGroupID {
	return &AckWxXcxBuyTaskListByGroupID{
		List: make([]*GroupBuyingOrderCommunityBuyTask, 0),
	}
}

type MetaApiWxXcxBuyTaskListByGroupID struct {
}

var META_WX_XCX_BUY_TASK_LIST_BY_GROUP_ID = &MetaApiWxXcxBuyTaskListByGroupID{}

func (m *MetaApiWxXcxBuyTaskListByGroupID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxBuyTaskListByGroupID) GetURL() string {
	return "/group_buying_order/wx_xcx/buy/task_list/:group_id"
}
func (m *MetaApiWxXcxBuyTaskListByGroupID) GetName() string { return "WxXcxBuyTaskListByGroupID" }
func (m *MetaApiWxXcxBuyTaskListByGroupID) GetType() string { return "json" }

// 商品购买列表
type ApiWxXcxBuyTaskListByGroupID struct {
	MetaApiWxXcxBuyTaskListByGroupID
	Ack    *AckWxXcxBuyTaskListByGroupID
	Params struct {
		GroupID uint32 `form:"group_id" binding:"required,gt=0" db:"GroupID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiWxXcxBuyTaskListByGroupID) GetQuery() interface{}  { return &m.Query }
func (m *ApiWxXcxBuyTaskListByGroupID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxBuyTaskListByGroupID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxBuyTaskListByGroupID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxBuyTaskListByGroupID() ApiWxXcxBuyTaskListByGroupID {
	return ApiWxXcxBuyTaskListByGroupID{
		Ack: NewAckWxXcxBuyTaskListByGroupID(),
	}
}

type AckWxXcxTaskStatusByTaskID struct {
	ShowState  GroupBuyingTaskShowState  `db:"show_state"`
	GroupState GroupBuyingTaskGroupState `db:"group_state"`
	OrderState GroupBuyingTaskOrderState `db:"order_state"`
	Sales      uint32                    `db:"sales"`
}

func NewAckWxXcxTaskStatusByTaskID() *AckWxXcxTaskStatusByTaskID {
	return &AckWxXcxTaskStatusByTaskID{}
}

type MetaApiWxXcxTaskStatusByTaskID struct {
}

var META_WX_XCX_TASK_STATUS_BY_TASK_ID = &MetaApiWxXcxTaskStatusByTaskID{}

func (m *MetaApiWxXcxTaskStatusByTaskID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxTaskStatusByTaskID) GetURL() string {
	return "/group_buying_order/wx_xcx/task/status/:task_id"
}
func (m *MetaApiWxXcxTaskStatusByTaskID) GetName() string { return "WxXcxTaskStatusByTaskID" }
func (m *MetaApiWxXcxTaskStatusByTaskID) GetType() string { return "json" }

// 商品状态
type ApiWxXcxTaskStatusByTaskID struct {
	MetaApiWxXcxTaskStatusByTaskID
	Ack    *AckWxXcxTaskStatusByTaskID
	Params struct {
		TaskID uint32 `form:"task_id" db:"TaskID"`
	}
}

func (m *ApiWxXcxTaskStatusByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxTaskStatusByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxTaskStatusByTaskID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxTaskStatusByTaskID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxTaskStatusByTaskID() ApiWxXcxTaskStatusByTaskID {
	return ApiWxXcxTaskStatusByTaskID{
		Ack: NewAckWxXcxTaskStatusByTaskID(),
	}
}

type AckWxXcxTaskInventoryByTaskID struct {
	Inventories map[string]*GroupBuyingOrderInventory `db:"Inventories"`
}

func NewAckWxXcxTaskInventoryByTaskID() *AckWxXcxTaskInventoryByTaskID {
	return &AckWxXcxTaskInventoryByTaskID{
		Inventories: make(map[string]*GroupBuyingOrderInventory),
	}
}

type MetaApiWxXcxTaskInventoryByTaskID struct {
}

var META_WX_XCX_TASK_INVENTORY_BY_TASK_ID = &MetaApiWxXcxTaskInventoryByTaskID{}

func (m *MetaApiWxXcxTaskInventoryByTaskID) GetMethod() string { return "GET" }
func (m *MetaApiWxXcxTaskInventoryByTaskID) GetURL() string {
	return "/group_buying_order/wx_xcx/task/inventory/:task_id"
}
func (m *MetaApiWxXcxTaskInventoryByTaskID) GetName() string { return "WxXcxTaskInventoryByTaskID" }
func (m *MetaApiWxXcxTaskInventoryByTaskID) GetType() string { return "json" }

// 团购任务库存
type ApiWxXcxTaskInventoryByTaskID struct {
	MetaApiWxXcxTaskInventoryByTaskID
	Ack    *AckWxXcxTaskInventoryByTaskID
	Params struct {
		TaskID uint32 `form:"task_id" db:"TaskID"`
	}
}

func (m *ApiWxXcxTaskInventoryByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiWxXcxTaskInventoryByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiWxXcxTaskInventoryByTaskID) GetAsk() interface{}    { return nil }
func (m *ApiWxXcxTaskInventoryByTaskID) GetAck() interface{}    { return m.Ack }
func MakeApiWxXcxTaskInventoryByTaskID() ApiWxXcxTaskInventoryByTaskID {
	return ApiWxXcxTaskInventoryByTaskID{
		Ack: NewAckWxXcxTaskInventoryByTaskID(),
	}
}
