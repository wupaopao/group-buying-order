package impls

import (
	"fmt"
	"time"

	"business/community/proxy/community"
	"business/group-buying-order/cidl"
	"business/group-buying-order/common/db"
	"business/group-buying-order/common/mq"
	"common/api"

	"github.com/mz-eco/mz/http"
	"github.com/mz-eco/mz/log"
	"github.com/mz-eco/mz/utils"
)

func init() {
	AddWxXcxTaskTodayListByOrganizationIDHandler()
	AddWxXcxTaskFutureListByOrganizationIDHandler()
	AddWxXcxTaskHistoryListByOrganizationIDHandler()
	AddWxXcxTaskInfoByTaskIDHandler()
	AddWxXcxTaskBatchWxSellTextByOrganizationIDHandler()
	AddWxXcxTaskAddCartByGroupIDByTaskIDHandler()
	AddWxXcxCartDeleteCartByGroupIDHandler()
	AddWxXcxTaskCartCountByGroupIDHandler()
	AddWxXcxCartCartListByGroupIDHandler()
	AddWxXcxCartChangeCountByGroupIDHandler()
	AddWxXcxOrderAddOrderByGroupIDHandler()
	AddWxXcxOrderDirectlyAddByGroupIDHandler()
	AddWxXcxOrderOrderListByGroupIDHandler()
	AddWxXcxBuyTaskListByGroupIDHandler()
	AddWxXcxTaskStatusByTaskIDHandler()
	AddWxXcxTaskInventoryByTaskIDHandler()
	AddWxXcxOrderOrderCancelByGroupIDByOrderIDHandler()
}

// 今日团购
type WxXcxTaskTodayListByOrganizationIDImpl struct {
	cidl.ApiWxXcxTaskTodayListByOrganizationID
}

func AddWxXcxTaskTodayListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_TASK_TODAY_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &WxXcxTaskTodayListByOrganizationIDImpl{
				ApiWxXcxTaskTodayListByOrganizationID: cidl.MakeApiWxXcxTaskTodayListByOrganizationID(),
			}
		},
	)
}

func (m *WxXcxTaskTodayListByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	groupId := m.Params.GroupID
	

	
	team, err := community.NewProxy("community-service").InnerCommunityGroupTeamByGroupID(groupId)
	if err != nil {
		ctx.ProxyErrorf(err, "get community group team from proxy failed. %s", err)
		return
	}
	teamIds := team.TeamIDs 
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	ack.Count, err = dbGroupBuying.TaskTodayCount(organizationId, teamIds)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get today task count failed. %s", err)
		return
	}

	if ack.Count == 0 {
		ctx.Json(ack)
		return
	}

	list, err := dbGroupBuying.TaskTodayList(organizationId, teamIds, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get today task list failed. %s", err)
		return
	}

	for _, task := range list {
		ackListItem := cidl.NewAckWxXcxTaskListItem()
		ackListItem.TaskId = task.TaskId
		ackListItem.OrganizationID = task.OrganizationId
		ackListItem.Title = task.Title
		ackListItem.Specification = task.Specification

		newSkuMap := make(map[string]*cidl.GroupBuyingSkuMapItem, 0)
		for skuId, skuItem := range ackListItem.Specification.SkuMap {
			if skuItem.IsShow {
				newSkuMap[skuId] = skuItem
			}
		}

		ackListItem.Specification.SkuMap = newSkuMap

		newCombinationMap := make(map[string]*cidl.GroupBuyingOrderTaskCombinationItem, 0)
		for skuId, skuItem := range ackListItem.Specification.CombinationSkuMap {
			if skuItem.IsShow {
				newCombinationMap[skuId] = skuItem
			}
		}

		ackListItem.Specification.CombinationSkuMap = newCombinationMap

		ackListItem.CoverPicture = task.CoverPicture
		ackListItem.MarketPriceRange = task.Specification.MarketPriceRange
		ackListItem.GroupBuyingPriceRange = task.Specification.GroupBuyingPriceRange
		ackListItem.SettlementPriceRange = task.Specification.SettlementPriceRange
		ackListItem.CostPriceRange = task.Specification.CostPriceRange
		ackListItem.StartTime = task.StartTime
		ackListItem.EndTime = task.EndTime
		ackListItem.Sales = task.Sales
		ackListItem.Notes = task.Notes
		ackListItem.SellType = task.SellType
		ackListItem.GroupState = task.GroupState
		ackListItem.OrderState = task.OrderState
		ack.List = append(ack.List, ackListItem)
	}

	ctx.Json(ack)
}

// 未来团购
type WxXcxTaskFutureListByOrganizationIDImpl struct {
	cidl.ApiWxXcxTaskFutureListByOrganizationID
}

func AddWxXcxTaskFutureListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_TASK_FUTURE_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &WxXcxTaskFutureListByOrganizationIDImpl{
				ApiWxXcxTaskFutureListByOrganizationID: cidl.MakeApiWxXcxTaskFutureListByOrganizationID(),
			}
		},
	)
}

func (m *WxXcxTaskFutureListByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	groupId := m.Params.GroupID
	
	team, err := community.NewProxy("community-service").InnerCommunityGroupTeamByGroupID(groupId)
	if err != nil {
		ctx.ProxyErrorf(err, "get community group team from proxy failed. %s", err)
		return
	}
	teamIds := team.TeamIDs 


	ack.Count, err = dbGroupBuying.TaskFutureCount(organizationId,teamIds)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get today task count failed. %s", err)
		return
	}

	list, err := dbGroupBuying.TaskFutureList(organizationId, teamIds, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get today task list failed. %s", err)
		return
	}

	for _, task := range list {
		ackListItem := cidl.NewAckWxXcxTaskListItem()
		ackListItem.TaskId = task.TaskId
		ackListItem.OrganizationID = task.OrganizationId
		ackListItem.Title = task.Title
		ackListItem.CoverPicture = task.CoverPicture
		ackListItem.MarketPriceRange = task.Specification.MarketPriceRange
		ackListItem.GroupBuyingPriceRange = task.Specification.GroupBuyingPriceRange
		ackListItem.SettlementPriceRange = task.Specification.SettlementPriceRange
		ackListItem.CostPriceRange = task.Specification.CostPriceRange
		ackListItem.StartTime = task.StartTime
		ackListItem.EndTime = task.EndTime
		ackListItem.Sales = task.Sales
		ackListItem.Notes = task.Notes
		ackListItem.GroupState = task.GroupState
		ackListItem.OrderState = task.OrderState
		ack.List = append(ack.List, ackListItem)
	}

	ctx.Json(ack)
}

// 历史团购
type WxXcxTaskHistoryListByOrganizationIDImpl struct {
	cidl.ApiWxXcxTaskHistoryListByOrganizationID
}

func AddWxXcxTaskHistoryListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_TASK_HISTORY_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &WxXcxTaskHistoryListByOrganizationIDImpl{
				ApiWxXcxTaskHistoryListByOrganizationID: cidl.MakeApiWxXcxTaskHistoryListByOrganizationID(),
			}
		},
	)
}

func (m *WxXcxTaskHistoryListByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	groupId := m.Params.GroupID

	team, err := community.NewProxy("community-service").InnerCommunityGroupTeamByGroupID(groupId)
	if err != nil {
		ctx.ProxyErrorf(err, "get community group team from proxy failed. %s", err)
		return
	}
	teamIds := team.TeamIDs 

	ack.Count, err = dbGroupBuying.TaskHistoryCount(organizationId,teamIds)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get today task count failed. %s", err)
		return
	}

	list, err := dbGroupBuying.TaskHistoryList(organizationId, teamIds, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get today task list failed. %s", err)
		return
	}

	for _, task := range list {
		ackListItem := cidl.NewAckWxXcxTaskListItem()
		ackListItem.TaskId = task.TaskId
		ackListItem.OrganizationID = task.OrganizationId
		ackListItem.Title = task.Title
		ackListItem.CoverPicture = task.CoverPicture
		ackListItem.MarketPriceRange = task.Specification.MarketPriceRange
		ackListItem.GroupBuyingPriceRange = task.Specification.GroupBuyingPriceRange
		ackListItem.SettlementPriceRange = task.Specification.SettlementPriceRange
		ackListItem.CostPriceRange = task.Specification.CostPriceRange
		ackListItem.StartTime = task.StartTime
		ackListItem.EndTime = task.EndTime
		ackListItem.Sales = task.Sales
		ackListItem.Notes = task.Notes
		ackListItem.GroupState = task.GroupState
		ackListItem.OrderState = task.OrderState
		ack.List = append(ack.List, ackListItem)
	}

	ctx.Json(ack)
}

// 获取团购任务
type WxXcxTaskInfoByTaskIDImpl struct {
	cidl.ApiWxXcxTaskInfoByTaskID
}

func AddWxXcxTaskInfoByTaskIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_TASK_INFO_BY_TASK_ID,
		func() http.ApiHandler {
			return &WxXcxTaskInfoByTaskIDImpl{
				ApiWxXcxTaskInfoByTaskID: cidl.MakeApiWxXcxTaskInfoByTaskID(),
			}
		},
	)
}

func (m *WxXcxTaskInfoByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	m.Ack, err = dbGroupBuying.GetTask(m.Params.TaskID)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task failed. %s", err)
		return
	}

	ctx.Json(m.Ack)
}

// 获取多个团购任务的微信销售文案
type WxXcxTaskBatchWxSellTextByOrganizationIDImpl struct {
	cidl.ApiWxXcxTaskBatchWxSellTextByOrganizationID
}

func AddWxXcxTaskBatchWxSellTextByOrganizationIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_TASK_BATCH_WX_SELL_TEXT_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &WxXcxTaskBatchWxSellTextByOrganizationIDImpl{
				ApiWxXcxTaskBatchWxSellTextByOrganizationID: cidl.MakeApiWxXcxTaskBatchWxSellTextByOrganizationID(),
			}
		},
	)
}

func (m *WxXcxTaskBatchWxSellTextByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	organizationId := m.Params.OrganizationID
	taskIds := m.Ask.TaskIds
	if len(taskIds) > 50 {
		ctx.Errorf(api.ErrWrongParams, "task ids less than or equal 50.")
		return
	}

	dbGroupBuying := db.NewMallGroupBuyingOrder()
	m.Ack.List, err = dbGroupBuying.WxXcxTaskBatchWxSellTextList(organizationId, taskIds)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get wx xcx task batch list failed. %s", err)
		return
	}

	ctx.Json(m.Ack)
}

// 添加到购物车
type WxXcxTaskAddCartByGroupIDByTaskIDImpl struct {
	cidl.ApiWxXcxTaskAddCartByGroupIDByTaskID
}

func AddWxXcxTaskAddCartByGroupIDByTaskIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_TASK_ADD_CART_BY_GROUP_ID_BY_TASK_ID,
		func() http.ApiHandler {
			return &WxXcxTaskAddCartByGroupIDByTaskIDImpl{
				ApiWxXcxTaskAddCartByGroupIDByTaskID: cidl.MakeApiWxXcxTaskAddCartByGroupIDByTaskID(),
			}
		},
	)
}

func (m *WxXcxTaskAddCartByGroupIDByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	taskId := m.Params.TaskID
	ask := m.Ask

	group, err := community.NewProxy("community-service").InnerCommunityGroupInfoByGroupID(m.Params.GroupID)
	if err != nil {
		ctx.ProxyErrorf(err, "get community group info from proxy failed. %s", err)
		return
	}
	groupId := group.GroupId

	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	task, err := dbGroupBuyingOrder.GetTask(taskId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task failed. %s", err)
		return
	}

	// 不在进行中不能上报销量
	now := time.Now()
	if task.GroupState != cidl.GroupBuyingTaskGroupStateInProgress || !((now.Equal(task.StartTime) || now.After(task.StartTime)) && (now.Before(task.EndTime))) {
		ctx.Errorf(cidl.ErrTaskNotInProgress, "task group state is not in progress")
		return
	}

	communityCart := cidl.NewGroupBuyingOrderCommunityCart()
	communityCart.CartId = utils.UniqueID()
	communityCart.GroupId = groupId
	communityCart.TaskId = taskId
	communityCart.TaskTitle = task.Title
	communityCart.SkuId = ask.SkuId
	communityCart.Version = cidl.CommunityCartRecordVersion
	communityCart.BuyDetail.IsCombination = ask.IsCombination
	communityCart.Count = ask.Count
	count := float64(communityCart.Count)
	if !ask.IsCombination { // 单品
		communityCart.BuyDetail.CombinationGoods = nil

		skuItem, ok := task.Specification.SkuMap[ask.SkuId]
		if !ok {
			ctx.Errorf(cidl.ErrTaskGoodsNotExists, "sku id not exists. %s", err)
			return
		}

		if !skuItem.IsShow {
			ctx.Errorf(cidl.ErrTaskGoodsUnavailableForSale, "sku not available for sale. %s", err)
			return
		}

		communityCart.BuyDetail.SingleGoods = skuItem

		inventory, errGet := dbGroupBuyingOrder.GetInventory(taskId, skuItem.SkuId)
		err = errGet
		if err != nil {
			ctx.Errorf(api.ErrDbQueryFailed, "get inventory failed. %s", err)
			return
		}

		if inventory.Surplus < ask.Count {
			ctx.Errorf(cidl.ErrInventoryShortage, "inventory shortage")
			return
		}

		communityCart.TotalMarketPrice = skuItem.MarketPrice * count
		communityCart.TotalGroupBuyingPrice = skuItem.GroupBuyingPrice * count
		communityCart.TotalSettlementPrice = skuItem.SettlementPrice * count
		communityCart.TotalCostPrice = skuItem.CostPrice * count

	} else { // 组合
		communityCart.BuyDetail.SingleGoods = nil

		skuItem, ok := task.Specification.CombinationSkuMap[ask.SkuId]
		if !ok {
			ctx.Errorf(cidl.ErrTaskGoodsNotExists, "sku id not exists. %s", err)
			return
		}

		if !skuItem.IsShow {
			ctx.Errorf(cidl.ErrTaskGoodsUnavailableForSale, "sku not available for sale. %s", err)
			return
		}

		communityCart.BuyDetail.CombinationGoods = skuItem

		var subSkuIds []string
		for _, subSku := range skuItem.SubItems {
			subSkuIds = append(subSkuIds, subSku.SkuId)
		}

		inventories, errGet := dbGroupBuyingOrder.GetInventoriesBySkuIds(taskId, subSkuIds)
		err = errGet
		if err != nil {
			ctx.Errorf(api.ErrDbQueryFailed, "get inventories failed. %s", err)
			return
		}

		for _, inventory := range inventories {
			if inventory.Surplus < skuItem.SubItems[inventory.SkuId].Count*m.Ask.Count {
				ctx.Errorf(cidl.ErrInventoryShortage, "inventory shortage")
				return
			}
		}

		communityCart.TotalMarketPrice = skuItem.MarketPrice * count
		communityCart.TotalGroupBuyingPrice = skuItem.GroupBuyingPrice * count
		communityCart.TotalSettlementPrice = skuItem.SettlementPrice * count
		communityCart.TotalCostPrice = skuItem.CostPrice * count

	}

	_, err = dbGroupBuyingOrder.AddCommunityCart(communityCart)
	if err != nil {
		ctx.Errorf(api.ErrDBInsertFailed, "add community cart failed. %s", err)
		return
	}

	ctx.Succeed()
}

// 从购物车中删除
type WxXcxCartDeleteCartByGroupIDImpl struct {
	cidl.ApiWxXcxCartDeleteCartByGroupID
}

func AddWxXcxCartDeleteCartByGroupIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_CART_DELETE_CART_BY_GROUP_ID,
		func() http.ApiHandler {
			return &WxXcxCartDeleteCartByGroupIDImpl{
				ApiWxXcxCartDeleteCartByGroupID: cidl.MakeApiWxXcxCartDeleteCartByGroupID(),
			}
		},
	)
}

func (m *WxXcxCartDeleteCartByGroupIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	_, err = dbGroupBuyingOrder.DeleteCommunityCarts(m.Params.GroupID, m.Ask.CartIds)
	if err != nil {
		ctx.Errorf(api.ErrDbDeleteRecordFailed, "delete community carts failed. %s", err)
		return
	}

	ctx.Succeed()
}

// 修改购买数目
type WxXcxCartChangeCountByGroupIDImpl struct {
	cidl.ApiWxXcxCartChangeCountByGroupID
}

func AddWxXcxCartChangeCountByGroupIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_CART_CHANGE_COUNT_BY_GROUP_ID,
		func() http.ApiHandler {
			return &WxXcxCartChangeCountByGroupIDImpl{
				ApiWxXcxCartChangeCountByGroupID: cidl.MakeApiWxXcxCartChangeCountByGroupID(),
			}
		},
	)
}

func (m *WxXcxCartChangeCountByGroupIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	groupId := m.Params.GroupID
	cartId := m.Ask.CartId
	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	cart, err := dbGroupBuyingOrder.GetCommunityCart(groupId, cartId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community cart failed. %s", err)
		return
	}

	taskId := cart.TaskId
	task, err := dbGroupBuyingOrder.GetTask(taskId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task failed. %s", err)
		return
	}

	// 不在进行中不能上报销量
	now := time.Now()
	if task.GroupState != cidl.GroupBuyingTaskGroupStateInProgress || !((now.Equal(task.StartTime) || now.After(task.StartTime)) && (now.Before(task.EndTime))) {
		ctx.Errorf(cidl.ErrTaskNotInProgress, "task group state is not in progress")
		return
	}

	// 单品
	if !cart.BuyDetail.IsCombination {
		skuId := cart.SkuId
		inventory, errGet := dbGroupBuyingOrder.GetInventory(taskId, skuId)
		err = errGet
		if err != nil {
			ctx.Errorf(api.ErrDbQueryFailed, "get inventory failed. %s", err)
			return
		}

		if inventory.Surplus < m.Ask.Count {
			ctx.Errorf(cidl.ErrInventoryShortage, "inventory shortage.")
			return
		}

	} else { // 组合
		var subSkuIds []string
		subSkuItems := cart.BuyDetail.CombinationGoods.SubItems
		for _, subSku := range subSkuItems {
			subSkuIds = append(subSkuIds, subSku.SkuId)
		}

		inventories, errGet := dbGroupBuyingOrder.GetInventoriesBySkuIds(taskId, subSkuIds)
		err = errGet
		if err != nil {
			ctx.Errorf(api.ErrDbQueryFailed, "get inventories failed. %s", err)
			return
		}

		for _, inventory := range inventories {
			if inventory.Surplus < subSkuItems[inventory.SkuId].Count*m.Ask.Count {
				ctx.Errorf(cidl.ErrInventoryShortage, "inventory shortage")
				return
			}
		}
	}

	_, err = dbGroupBuyingOrder.ChangeCommunityCartBuyCount(groupId, cartId, m.Ask.Count)
	if err != nil {
		ctx.Errorf(api.ErrDBUpdateFailed, "change community cart buy count failed. %s", err)
		return
	}

	ctx.Succeed()
}

// 购物车数目
type WxXcxTaskCartCountByGroupIDImpl struct {
	cidl.ApiWxXcxTaskCartCountByGroupID
}

func AddWxXcxTaskCartCountByGroupIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_TASK_CART_COUNT_BY_GROUP_ID,
		func() http.ApiHandler {
			return &WxXcxTaskCartCountByGroupIDImpl{
				ApiWxXcxTaskCartCountByGroupID: cidl.MakeApiWxXcxTaskCartCountByGroupID(),
			}
		},
	)
}

func (m *WxXcxTaskCartCountByGroupIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	m.Ack.Count, err = dbGroupBuyingOrder.GetCommunityCartCount(m.Params.GroupID)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community cart count failed. %s", err)
		return
	}

	ctx.Json(m.Ack)
}

// 购物车列表
type WxXcxCartCartListByGroupIDImpl struct {
	cidl.ApiWxXcxCartCartListByGroupID
}

func AddWxXcxCartCartListByGroupIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_CART_CART_LIST_BY_GROUP_ID,
		func() http.ApiHandler {
			return &WxXcxCartCartListByGroupIDImpl{
				ApiWxXcxCartCartListByGroupID: cidl.MakeApiWxXcxCartCartListByGroupID(),
			}
		},
	)
}

func (m *WxXcxCartCartListByGroupIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	groupId := m.Params.GroupID
	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	m.Ack.Count, err = dbGroupBuyingOrder.GetCommunityCartCount(groupId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community cart count failed. %s", err)
		return
	}

	if m.Ack.Count == 0 {
		ctx.Json(m.Ack)
		return
	}

	m.Ack.List, err = dbGroupBuyingOrder.CommunityCartList(groupId, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community cart list failed. %s", err)
		return
	}

	ctx.Json(m.Ack)
}

// 提交订单
type WxXcxOrderAddOrderByGroupIDImpl struct {
	cidl.ApiWxXcxOrderAddOrderByGroupID
}

func AddWxXcxOrderAddOrderByGroupIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_ORDER_ADD_ORDER_BY_GROUP_ID,
		func() http.ApiHandler {
			return &WxXcxOrderAddOrderByGroupIDImpl{
				ApiWxXcxOrderAddOrderByGroupID: cidl.MakeApiWxXcxOrderAddOrderByGroupID(),
			}
		},
	)
}

func (m *WxXcxOrderAddOrderByGroupIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	groupId := m.Params.GroupID
	cartIds := m.Ask.CartIds

	group, err := community.NewProxy("community-service").InnerCommunityGroupInfoByGroupID(m.Params.GroupID)
	if err != nil {
		ctx.ProxyErrorf(err, "get community group info from proxy failed. %s", err)
		return
	}

	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	carts, err := dbGroupBuyingOrder.GetCommunityCarts(groupId, cartIds)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community carts failed. %s", err)
		return
	}

	orderId := utils.UniqueID()
	groupOrderId := time.Now().Format("20060102150405")
	communityOrder := cidl.NewGroupBuyingOrderCommunityOrder()
	communityOrder.OrderId = orderId
	communityOrder.GroupId = groupId
	communityOrder.GroupOrderId = groupOrderId
	communityOrder.Version = cidl.CommunityOrderRecordVersion

	for _, cart := range carts {
		taskId := cart.TaskId
		task, errGet := dbGroupBuyingOrder.GetTask(taskId)
		if errGet != nil {
			err = errGet
			log.Warnf("db get task failed. %s", err)
			continue
		}

		// 不在进行中不能提交订单
		now := time.Now()
		if task.GroupState != cidl.GroupBuyingTaskGroupStateInProgress || !((now.Equal(task.StartTime) || now.After(task.StartTime)) && (now.Before(task.EndTime))) {
			log.Warnf("task is not in progress. task_id:%d", taskId)
			continue
		}

		tx, err := dbGroupBuyingOrder.DB.Begin()
		if err != nil {
			log.Warnf("begin transaction failed. %s", err)
			continue
		}

		communityBuy := cidl.NewGroupBuyingOrderCommunityBuy()
		communityBuy.BuyId = utils.UniqueID()
		communityBuy.OrderId = orderId
		communityBuy.GroupId = groupId
		communityBuy.GroupOrderId = groupOrderId
		communityBuy.GroupName = group.Name
		communityBuy.TaskId = cart.TaskId
		communityBuy.TaskTitle = cart.TaskTitle
		communityBuy.ManagerUserId = group.ManagerUserId
		communityBuy.ManagerName = group.ManagerName
		communityBuy.ManagerMobile = group.ManagerMobile
		communityBuy.SkuId = cart.SkuId
		communityBuy.Count = cart.Count
		communityBuy.TotalMarketPrice = cart.TotalMarketPrice
		communityBuy.TotalGroupBuyingPrice = cart.TotalGroupBuyingPrice
		communityBuy.TotalSettlementPrice = cart.TotalSettlementPrice
		communityBuy.TotalCostPrice = cart.TotalCostPrice
		communityBuy.Version = cidl.CommunityBuyRecordVersion
		communityBuy.BuyDetail.IsCombination = cart.BuyDetail.IsCombination
		communityBuy.CreateTime = time.Now()

		// 单品
		if !cart.BuyDetail.IsCombination {
			skuItem, ok := task.Specification.SkuMap[cart.SkuId]
			if !ok {
				log.Warnf("sku id not exist. sku_id: %s", cart.SkuId)
				continue
			}

			if !skuItem.IsShow {
				log.Warnf("sku not available for sale. %s", err)
				continue
			}

			// 减少库存
			success, errSub := dbGroupBuyingOrder.TxSubtractInventorySurplus(tx, taskId, cart.SkuId, cart.Count)
			if errSub != nil || !success {
				log.Warnf("subtract inventory failed. %s", errSub)
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}

			communityBuy.BuyDetail.SingleGoods = skuItem
			communityBuy.BuyDetail.CombinationGoods = nil

		} else { // 组合
			skuItem, ok := task.Specification.CombinationSkuMap[cart.SkuId]
			if !ok {
				log.Warnf("sku id not exist. sku_id: %s", cart.SkuId)
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}

			if !skuItem.IsShow {
				log.Warnf("sku not available for sale. %s", err)
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}

			// 减少库存
			var subSkuIds []string
			for subSkuId, _ := range skuItem.SubItems {
				subSkuIds = append(subSkuIds, subSkuId)
			}

			// 锁定库存
			_, err = dbGroupBuyingOrder.TxLockInventorySurplus(tx, taskId, subSkuIds)
			if err != nil {
				log.Warnf("lock inventory surplus failed. %s", err)
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}

			// 事务减库存
			buySuccess := true
			for _, subSkuItem := range skuItem.SubItems {
				success, errSub := dbGroupBuyingOrder.TxSubtractInventorySurplus(tx, taskId, subSkuItem.SkuId, subSkuItem.Count*cart.Count)
				if errSub != nil || success == false {
					log.Warnf("subtract inventory failed. %s", err)
					buySuccess = false
					break
				}
			}

			if !buySuccess {
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}

			communityBuy.BuyDetail.CombinationGoods = skuItem
			communityBuy.BuyDetail.SingleGoods = nil

		}

		_, err = dbGroupBuyingOrder.TxAddCommunityBuy(tx, communityBuy)
		if err != nil {
			log.Warnf("tx add community buy failed. %s", err)
			err = tx.Rollback()
			if err != nil {
				log.Warnf("rollback community buy failed. %s", err)
			}
			continue
		}

		// 提交事务
		err = tx.Commit()
		if err != nil {
			log.Warnf("commit community buy failed. %s", err)
			continue
		}

		// 成功购买返回ID
		m.Ack.SuccessCartIds = append(m.Ack.SuccessCartIds, cart.CartId)

		// 添加成功的订单
		*communityOrder.GoodsDetail = append(*communityOrder.GoodsDetail, communityBuy)
		communityOrder.Count += communityBuy.Count
		communityOrder.TotalMarketPrice += communityBuy.TotalMarketPrice
		communityOrder.TotalGroupBuyingPrice += communityBuy.TotalGroupBuyingPrice
		communityOrder.TotalSettlementPrice += communityBuy.TotalSettlementPrice
		communityOrder.TotalCostPrice += communityBuy.TotalCostPrice

		// kafka消息广播添加订单消息
		topic, err := mq.GetTopicServiceGroupBuyingOrderService()
		if err != nil {
			log.Warnf("get topic service-group-buying-order-service failed. %s", err)
			continue
		}

		err = topic.AddOrder(&mq.AddOrderMessage{
			GroupId:      groupId,
			TaskId:       taskId,
			CommunityBuy: communityBuy,
		})
		if err != nil {
			log.Warnf("publish topic service-group-buying-order-service add order message failed. %s", err)
			continue
		}

	}

	if communityOrder.Count > 0 {
		_, err = dbGroupBuyingOrder.AddCommunityOrder(communityOrder)
		if err != nil {
			ctx.Errorf(api.ErrDBInsertFailed, "add community order failed. %s", err)
			return
		}

		// 删除购物车
		_, err = dbGroupBuyingOrder.DeleteCommunityCarts(groupId, m.Ack.SuccessCartIds)
		if err != nil {
			ctx.Errorf(api.ErrDbDeleteRecordFailed, "delete community carts failed. %s", err)
			return
		}

	}

	ctx.Json(m.Ack)
}

// 直接提交订单
type WxXcxOrderDirectlyAddByGroupIDImpl struct {
	cidl.ApiWxXcxOrderDirectlyAddByGroupID
}

func AddWxXcxOrderDirectlyAddByGroupIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_ORDER_DIRECTLY_ADD_BY_GROUP_ID,
		func() http.ApiHandler {
			return &WxXcxOrderDirectlyAddByGroupIDImpl{
				ApiWxXcxOrderDirectlyAddByGroupID: cidl.MakeApiWxXcxOrderDirectlyAddByGroupID(),
			}
		},
	)
}

func (m *WxXcxOrderDirectlyAddByGroupIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	groupId := m.Params.GroupID
	group, err := community.NewProxy("community-service").InnerCommunityGroupInfoByGroupID(m.Params.GroupID)
	if err != nil {
		ctx.ProxyErrorf(err, "get community group info from proxy failed. %s", err)
		return
	}

	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()

	orderId := utils.UniqueID()
	groupOrderId := time.Now().Format("20060102150405")
	m.Ack.OrderId = orderId
	communityOrder := cidl.NewGroupBuyingOrderCommunityOrder()
	communityOrder.OrderId = orderId
	communityOrder.GroupId = groupId
	communityOrder.GroupOrderId = groupOrderId
	communityOrder.Version = cidl.CommunityOrderRecordVersion

	buyItems := m.Ask.Items
	for index, buyItem := range buyItems {

		taskId := buyItem.TaskId
		skuId := buyItem.SkuId
		count := float64(buyItem.Count)

		m.Ack.ErrorList = append(m.Ack.ErrorList, &cidl.AckDirectlyAddOrderResultItem{
			TaskId:  taskId,
			SkuId:   skuId,
			Message: "购买失败",
		})

		task, err := dbGroupBuyingOrder.GetTask(taskId)
		if err != nil {
			log.Warnf("get task failed. %s", err)
			continue
		}

		// 不在进行中不能上报销量
		now := time.Now()
		if task.GroupState != cidl.GroupBuyingTaskGroupStateInProgress || !((now.Equal(task.StartTime) || now.After(task.StartTime)) && (now.Before(task.EndTime))) {
			ctx.Errorf(cidl.ErrTaskNotInProgress, "task group state is not in progress")
			return
		}

		tx, err := dbGroupBuyingOrder.DB.Begin()
		if err != nil {
			log.Warnf("begin transaction failed. %s", err)
			continue
		}

		communityBuy := cidl.NewGroupBuyingOrderCommunityBuy()
		communityBuy.BuyId = utils.UniqueID()
		communityBuy.OrderId = orderId
		communityBuy.GroupId = groupId
		communityBuy.GroupOrderId = groupOrderId
		communityBuy.GroupName = group.Name
		communityBuy.TaskId = taskId
		communityBuy.TaskTitle = task.Title
		communityBuy.ManagerUserId = group.ManagerUserId
		communityBuy.ManagerName = group.ManagerName
		communityBuy.ManagerMobile = group.ManagerMobile
		communityBuy.SkuId = skuId
		communityBuy.Count = buyItem.Count
		communityBuy.BuyDetail.IsCombination = buyItem.IsCombination
		communityBuy.CreateTime = time.Now()
		communityBuy.Version = cidl.CommunityBuyRecordVersion

		if !buyItem.IsCombination { // 单品
			skuItem, ok := task.Specification.SkuMap[skuId]
			if !ok {
				log.Warnf("sku id not exist. sku_id: %s", skuId)
				continue
			}

			if !skuItem.IsShow {
				log.Warnf("sku not available for sale. %s", err)
				continue
			}

			// 减少库存
			success, errSub := dbGroupBuyingOrder.TxSubtractInventorySurplus(tx, taskId, skuId, uint32(count))
			if errSub != nil || !success {
				inventory, errGet := dbGroupBuyingOrder.GetInventory(taskId, skuId)
				err = errGet
				if err != nil {
					log.Warnf("get inventory failed.")
					m.Ack.ErrorList[index].Message = "库存不足"
				} else {
					m.Ack.ErrorList[index].Message = fmt.Sprintf("库存不足，剩余库存%d", inventory.Surplus)
				}

				log.Warnf("subtract inventory failed. %s", errSub)
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}

			communityBuy.BuyDetail.SingleGoods = skuItem
			communityBuy.BuyDetail.CombinationGoods = nil
			communityBuy.TotalMarketPrice = skuItem.MarketPrice * count
			communityBuy.TotalGroupBuyingPrice = skuItem.GroupBuyingPrice * count
			communityBuy.TotalSettlementPrice = skuItem.SettlementPrice * count
			communityBuy.TotalCostPrice = skuItem.CostPrice * count

		} else { // 组合
			skuItem, ok := task.Specification.CombinationSkuMap[skuId]
			if !ok {
				log.Warnf("sku id not exist. sku_id: %s", skuId)
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}

			if !skuItem.IsShow {
				log.Warnf("sku not available for sale. %s", err)
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}
			// 减少库存
			var subSkuIds []string
			for subSkuId, _ := range skuItem.SubItems {
				subSkuIds = append(subSkuIds, subSkuId)
			}

			// 锁定库存
			_, err = dbGroupBuyingOrder.TxLockInventorySurplus(tx, taskId, subSkuIds)
			if err != nil {
				log.Warnf("lock inventory surplus failed. %s", err)
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}

			// 事务减库存
			buySuccess := true
			for _, subSkuItem := range skuItem.SubItems {
				success, errSub := dbGroupBuyingOrder.TxSubtractInventorySurplus(tx, taskId, subSkuItem.SkuId, subSkuItem.Count*uint32(count))
				if errSub != nil || success == false {
					m.Ack.ErrorList[index].Message = "库存不足"
					log.Warnf("subtract inventory failed. %s", err)
					buySuccess = false
					break
				}
			}

			if !buySuccess {
				err = tx.Rollback()
				if err != nil {
					log.Warnf("rollback community buy failed. %s", err)
				}
				continue
			}

			communityBuy.BuyDetail.CombinationGoods = skuItem
			communityBuy.BuyDetail.SingleGoods = nil

			communityBuy.TotalMarketPrice = skuItem.MarketPrice * count
			communityBuy.TotalGroupBuyingPrice = skuItem.GroupBuyingPrice * count
			communityBuy.TotalSettlementPrice = skuItem.SettlementPrice * count
			communityBuy.TotalCostPrice = skuItem.CostPrice * count

		}

		_, err = dbGroupBuyingOrder.TxAddCommunityBuy(tx, communityBuy)
		if err != nil {
			log.Warnf("tx add community buy failed. %s", err)
			err = tx.Rollback()
			if err != nil {
				log.Warnf("rollback community buy failed. %s", err)
			}
			continue
		}

		// 提交事务
		err = tx.Commit()
		if err != nil {
			log.Warnf("commit community buy failed. %s", err)
			continue
		}

		m.Ack.ErrorList[index] = nil

		// 添加成功的订单
		*communityOrder.GoodsDetail = append(*communityOrder.GoodsDetail, communityBuy)
		communityOrder.Count += communityBuy.Count
		communityOrder.TotalMarketPrice += communityBuy.TotalMarketPrice
		communityOrder.TotalGroupBuyingPrice += communityBuy.TotalGroupBuyingPrice
		communityOrder.TotalSettlementPrice += communityBuy.TotalSettlementPrice
		communityOrder.TotalCostPrice += communityBuy.TotalCostPrice

		// kafka消息广播添加订单消息
		topic, err := mq.GetTopicServiceGroupBuyingOrderService()
		if err != nil {
			log.Warnf("get topic service-group-buying-order-service failed. %s", err)
			continue
		}

		err = topic.AddOrder(&mq.AddOrderMessage{
			GroupId:      groupId,
			TaskId:       taskId,
			CommunityBuy: communityBuy,
		})
		if err != nil {
			log.Warnf("publish topic service-group-buying-order-service add order message failed. %s", err)
			continue
		}

	}

	if communityOrder.Count > 0 {
		// 添加订单数目
		_, err = dbGroupBuyingOrder.AddCommunityOrder(communityOrder)
		if err != nil {
			ctx.Errorf(api.ErrDBInsertFailed, "add community order failed. %s", err)
			return
		}

	}

	newErrList := make([]*cidl.AckDirectlyAddOrderResultItem, 0)
	for _, errInfo := range m.Ack.ErrorList {
		if errInfo != nil {
			newErrList = append(newErrList, errInfo)
		}
	}

	m.Ack.ErrorList = newErrList
	ctx.Json(m.Ack)
}

// 订单列表
type WxXcxOrderOrderListByGroupIDImpl struct {
	cidl.ApiWxXcxOrderOrderListByGroupID
}

func AddWxXcxOrderOrderListByGroupIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_ORDER_ORDER_LIST_BY_GROUP_ID,
		func() http.ApiHandler {
			return &WxXcxOrderOrderListByGroupIDImpl{
				ApiWxXcxOrderOrderListByGroupID: cidl.MakeApiWxXcxOrderOrderListByGroupID(),
			}
		},
	)
}

func (m *WxXcxOrderOrderListByGroupIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	groupId := m.Params.GroupID
	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	m.Ack.Count, err = dbGroupBuyingOrder.CommunityOrderCount(groupId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community order count failed. %s", err)
		return
	}

	if m.Ack.Count == 0 {
		ctx.Json(m.Ack)
		return
	}
	mNotAllowCancelList, err := dbGroupBuyingOrder.CommunityNotAllowCancelOrderList(groupId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get allow cancel order failed. %s", err)
		return
	}

	list, err := dbGroupBuyingOrder.CommunityOrderList(groupId, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community order failed. %s", err)
		return
	}

	for _, communityOrder := range list {
		orderId := communityOrder.OrderId
		if status,ok := mNotAllowCancelList[orderId];ok{
			communityOrder.AllowCancel = false
			communityOrder.Status = status
		}else {
			communityOrder.AllowCancel = true 
		}
		if communityOrder.IsCancel {
			communityOrder.Status = "订单已取消"
		}
		for _, communityBuy := range *communityOrder.GoodsDetail {
			buyDetail := communityBuy.BuyDetail
			if buyDetail.IsCombination { // 组合
				goods := buyDetail.CombinationGoods
				communityBuy.BuyDetail.SingleGoods = &cidl.GroupBuyingSkuMapItem{
					SkuId:               goods.SkuId,
					Name:                goods.Name,
					Labels:              nil,
					MarketPrice:         goods.MarketPrice,
					GroupBuyingPrice:    goods.GroupBuyingPrice,
					SettlementPrice:     goods.SettlementPrice,
					CostPrice:           goods.CostPrice,
					IllustrationPicture: goods.IllustrationPicture,
					IsShow:              goods.IsShow,
				}
				buyDetail.CombinationGoods = nil
			}
		}
	}

	m.Ack.List = list
	ctx.Json(m.Ack)
}

// 商品购买列表
type WxXcxBuyTaskListByGroupIDImpl struct {
	cidl.ApiWxXcxBuyTaskListByGroupID
}

func AddWxXcxBuyTaskListByGroupIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_BUY_TASK_LIST_BY_GROUP_ID,
		func() http.ApiHandler {
			return &WxXcxBuyTaskListByGroupIDImpl{
				ApiWxXcxBuyTaskListByGroupID: cidl.MakeApiWxXcxBuyTaskListByGroupID(),
			}
		},
	)
}

func (m *WxXcxBuyTaskListByGroupIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	groupId := m.Params.GroupID
	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	m.Ack.Count, err = dbGroupBuyingOrder.CommunityBuyTaskCount(groupId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community or")
		return
	}

	if m.Ack.Count == 0 {
		ctx.Json(m.Ack)
		return
	}

	m.Ack.List, err = dbGroupBuyingOrder.CommunityBuyTaskList(groupId, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community buy task list failed. %s", err)
		return
	}

	ctx.Json(m.Ack)
}

// 商品状态
type WxXcxTaskStatusByTaskIDImpl struct {
	cidl.ApiWxXcxTaskStatusByTaskID
}

func AddWxXcxTaskStatusByTaskIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_TASK_STATUS_BY_TASK_ID,
		func() http.ApiHandler {
			return &WxXcxTaskStatusByTaskIDImpl{
				ApiWxXcxTaskStatusByTaskID: cidl.MakeApiWxXcxTaskStatusByTaskID(),
			}
		},
	)
}

func (m *WxXcxTaskStatusByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	taskId := m.Params.TaskID
	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	task, err := dbGroupBuyingOrder.GetTask(taskId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task failed. %s", err)
		return
	}

	m.Ack.ShowState = task.ShowState
	m.Ack.GroupState = task.GroupState
	m.Ack.OrderState = task.OrderState
	m.Ack.Sales = task.Sales

	ctx.Json(m.Ack)
}

// 团购任务库存
type WxXcxTaskInventoryByTaskIDImpl struct {
	cidl.ApiWxXcxTaskInventoryByTaskID
}

func AddWxXcxTaskInventoryByTaskIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_TASK_INVENTORY_BY_TASK_ID,
		func() http.ApiHandler {
			return &WxXcxTaskInventoryByTaskIDImpl{
				ApiWxXcxTaskInventoryByTaskID: cidl.MakeApiWxXcxTaskInventoryByTaskID(),
			}
		},
	)
}

func (m *WxXcxTaskInventoryByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	taskId := m.Params.TaskID
	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	inventories, err := dbGroupBuyingOrder.GetInventories(taskId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get inventories failed. %s", err)
		return
	}

	for _, inventory := range inventories {
		m.Ack.Inventories[inventory.SkuId] = inventory
	}

	ctx.Json(m.Ack)
}

// 取消订单 
type WxXcxOrderOrderCancelByGroupIDByOrderIDImpl struct {
	cidl.ApiWxXcxOrderOrderCancelByGroupIDByOrderID
}

func AddWxXcxOrderOrderCancelByGroupIDByOrderIDHandler() {
	AddHandler(
		cidl.META_WX_XCX_ORDER_ORDER_CANCEL_BY_GROUP_ID_BY_ORDER_ID,
		func() http.ApiHandler {
			return &WxXcxOrderOrderCancelByGroupIDByOrderIDImpl{
				ApiWxXcxOrderOrderCancelByGroupIDByOrderID: cidl.MakeApiWxXcxOrderOrderCancelByGroupIDByOrderID(),
			}
		},
	)
}

func (m *WxXcxOrderOrderCancelByGroupIDByOrderIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	groupId := m.Params.GroupID
	orderId := m.Params.OrderID
	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()
	
	//订单是否支持取消
	/*
	isAllowCancel, _ := dbGroupBuyingOrder.IsOrderAllowCancel(groupId,orderId)	
	if !isAllowCancel {
		ctx.Errorf(cidl.ErrOrderNotAllowCancel, "order is not allow cancel.")
		return
	}*/

	communityOrder, err := dbGroupBuyingOrder.CommunityOrderInfo(orderId)
	
	//for test
	fmt.Println(communityOrder.GoodsDetail.ToString())

	goodsDetail := *communityOrder.GoodsDetail
	for _, communityBuy := range goodsDetail{
		taskId := communityBuy.TaskId
		groupId := communityBuy.GroupId
		buyDetail := communityBuy.BuyDetail
		count := communityBuy.Count
		skuId := communityBuy.SkuId

		//开始事务
		tx, err := dbGroupBuyingOrder.DB.Begin()
		if err != nil {
			log.Warnf("begin transaction failed. %s", err)
			continue
		}
		if !buyDetail.IsCombination {
			//singleGoods := buyDetail.SingleGoods
			success, errAdd := dbGroupBuyingOrder.TxAddtractInventorySurplus(tx, taskId, skuId, count)
                        if errAdd != nil || !success {
                                log.Warnf("addtract inventory failed. %s", errAdd)
                                err = tx.Rollback()
                                if err != nil {
                                        log.Warnf("rollback community buy failed. %s", err)
                                }
                                continue
                        }
		} else {
			combinationGoods := buyDetail.CombinationGoods
			subItems := combinationGoods.SubItems
			var skuIds []string
			for skuId,_ := range subItems {
				skuIds = append(skuIds,skuId)	
			}
			// 锁定库存
                        _, err = dbGroupBuyingOrder.TxLockInventorySurplus(tx, taskId, skuIds)
                        if err != nil {
                                log.Warnf("lock inventory surplus failed. %s", err)
                                err = tx.Rollback()
                                if err != nil {
                                        log.Warnf("rollback community buy failed. %s", err)
                                }
                                continue
                        }	

                        // 事务加库存
                        buySuccess := true
                        for _, subItem := range subItems {
                                success, errAdd := dbGroupBuyingOrder.TxAddtractInventorySurplus(tx, taskId, subItem.SkuId, subItem.Count*count)
                                if errAdd != nil || success == false {
                                        log.Warnf("addtract inventory failed. %s", err)
                                        buySuccess = false
                                        break
                                }
                        }

                        if !buySuccess {
                                err = tx.Rollback()
                                if err != nil {
                                        log.Warnf("rollback community buy failed. %s", err)
                                }
                                continue
                        }	
               }

               _, err = dbGroupBuyingOrder.TxDeleteCommunityBuy(tx, orderId, groupId, taskId)
                if err != nil {
                        log.Warnf("tx add community buy failed. %s", err)
                        err = tx.Rollback()
                        if err != nil {
                                log.Warnf("rollback community buy failed. %s", err)
                        }
                        continue
                }

                // 提交事务
                err = tx.Commit()
                if err != nil {
                        log.Warnf("commit community buy failed. %s", err)
                        continue
                }


                // 减少社群购买团购任务的订单统计
		goodsCount := communityBuy.Count
		totalMarketPrice := communityBuy.TotalMarketPrice
		totalGroupBuyingPrice := communityBuy.TotalGroupBuyingPrice
		totalSettlementPrice := communityBuy.TotalSettlementPrice
		totalCostPrice := communityBuy.TotalCostPrice

                _, err = dbGroupBuyingOrder.DecrCommunityBuyTaskValues(taskId,groupId,1,goodsCount,totalMarketPrice,totalGroupBuyingPrice,totalSettlementPrice,totalCostPrice)
                if err != nil {
                        log.Warnf("decrease community buy task values failed. %s", err)
                }

                // 减少任务的购买数量
                _, err = dbGroupBuyingOrder.DecrTaskSales(taskId, goodsCount)
                if err != nil {
                        log.Warnf("decrease task sales failed. %s", err)
                }
	}
	_, err = dbGroupBuyingOrder.DeleteCommunityOrder(groupId, orderId)
	if err != nil {
		ctx.Errorf(api.ErrDBInsertFailed, "delete community order failed. %s", err)
		return
	}
	
	ctx.Succeed()
}
