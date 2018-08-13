package mq

import (
	"business/group-buying-order/cidl"
	"business/group-buying-order/common/db"

	"github.com/mz-eco/mz/conn"
	"github.com/mz-eco/mz/log"
)

func (m *TopicGroupBuyingServiceHandler) AddOrder(msg *AddOrderMessage) (err error) {
	groupId := msg.GroupId
	taskId := msg.TaskId
	communityBuy := msg.CommunityBuy

	dbGroupBuyingOrder := db.NewMallGroupBuyingOrder()

	task, err := dbGroupBuyingOrder.GetTask(taskId)
	if err != nil {
		log.Warnf("db get task failed. %s", err)
		return
	}

	// 添加社群购买团购任务的订单统计
	buyTask, errGet := dbGroupBuyingOrder.GetCommunityBuyTask(groupId, taskId)
	if errGet != nil && errGet != conn.ErrNoRows {
		err = errGet
		log.Warnf("get community buy task failed. %s", errGet)
		return

	} else if errGet == conn.ErrNoRows {
		buyTask = cidl.NewGroupBuyingOrderCommunityBuyTask()
		buyTask.TaskId = taskId
		buyTask.GroupId = groupId
		buyTask.TaskDetail = &task.GroupBuyingOrderTaskContent
		buyTask.OrderCount = 1
		buyTask.GoodsCount = communityBuy.Count
		buyTask.TotalMarketPrice = communityBuy.TotalMarketPrice
		buyTask.TotalGroupBuyingPrice = communityBuy.TotalGroupBuyingPrice
		buyTask.TotalSettlementPrice = communityBuy.TotalSettlementPrice
		buyTask.TotalCostPrice = communityBuy.TotalCostPrice
		buyTask.Version = cidl.CommunityBuyTaskRecordVersion
		_, err = dbGroupBuyingOrder.AddCommunityBuyTask(buyTask)
		if err != nil {
			log.Warnf("add community buy task failed. %s", err)
		}

	} else {
		_, err = dbGroupBuyingOrder.IncrCommunityBuyTaskValues(
			taskId,
			groupId,
			1,
			communityBuy.Count,
			communityBuy.TotalMarketPrice,
			communityBuy.TotalGroupBuyingPrice,
			communityBuy.TotalSettlementPrice,
			communityBuy.TotalCostPrice,
		)
		if err != nil {
			log.Warnf("increase community buy task values failed. %s", err)
		}
	}

	// 增加任务的购买数量
	_, err = dbGroupBuyingOrder.IncrTaskSales(taskId, communityBuy.Count)
	if err != nil {
		log.Warnf("increase task sales failed. %s", err)
	}

	return
}
