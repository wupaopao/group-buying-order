package mq

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"business/agency/proxy/agency"
	"business/community/proxy/community"
	"business/group-buying-order/cidl"
	"business/group-buying-order/common/db"
	"business/group-buying-order/common/excel"
	"business/user/proxy/user"

	//"github.com/mz-eco/mz/errors"
	"github.com/mz-eco/mz/kafka"
	"github.com/mz-eco/mz/log"
	"github.com/mz-eco/mz/settings"
	"github.com/mz-eco/mz/utils"
)

var (
	topicGroupBuyingServiceSetting kafka.TopicGroupSetting
)

func init() {
	settings.RegisterWith(func(viper *settings.Viper) error {
		err := viper.Unmarshal(&topicGroupBuyingServiceSetting)
		if err != nil {
			panic(err)
			return err
		}
		return nil
	}, "kafka.subscribe.service_group_buying_order_service")
}

// 团购任务服务消息
type TopicGroupBuyingServiceHandler struct {
	kafka.TopicHandler
}

func NewTopicGroupBuyingServiceHandler() (handler *TopicGroupBuyingServiceHandler, err error) {
	handler = &TopicGroupBuyingServiceHandler{
		TopicHandler: kafka.TopicHandler{
			Topics:  []string{TOPIC_SERVICE_GROUP_BUYING_ORDER_SERVICE},
			Brokers: topicGroupBuyingServiceSetting.Address,
			Group:   topicGroupBuyingServiceSetting.Group,
		},
	}

	handler.MessageHandle = handler.handleMessage

	return
}

func (m *TopicGroupBuyingServiceHandler) handleMessage(identMessage *kafka.IdentMessage) (err error) {
	switch identMessage.Ident {
	case IDENT_SERVICE_GROUP_BUYING_ORDER_SERVICE_ADD_INDENT:
		addIndent := &AddIndentMessage{}
		err = json.Unmarshal(identMessage.Msg, addIndent)
		if err != nil {
			log.Warnf("unmarshal add indent message failed. %s", err)
			return
		}

		err = m.AddIndent(addIndent)
		if err != nil {
			log.Warnf("add indent failed. %s", err)
			return
		}

	case IDENT_SERVICE_GROUP_BUYING_ORDER_SERVICE_ADD_SEND:
		addSend := &AddSendMessage{}
		err = json.Unmarshal(identMessage.Msg, addSend)
		if err != nil {
			log.Warnf("unmarshal add send message failed. %s", err)
			return
		}

		err = m.AddSend(addSend)
		if err != nil {
			log.Warnf("add send failed. %s", err)
			return
		}

	case IDENT_SERVICE_GROUP_BUYING_ORDER_SERVICE_ADD_ORDER:
		addOrder := &AddOrderMessage{}
		err = json.Unmarshal(identMessage.Msg, addOrder)
		if err != nil {
			log.Warnf("unmarshal add order message failed. %s", err)
			return
		}

		err = m.AddOrder(addOrder)
		if err != nil {
			log.Warnf("add order failed. %s", err)
			return
		}

	}

	return
}

// 添加订货单
func (m *TopicGroupBuyingServiceHandler) AddIndent(msg *AddIndentMessage) (err error) {
	indentId := msg.IndentID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	defer func() {
		if err != nil {
			log.Warnf("add indent failed. %s", err)
			_, err = dbGroupBuying.UpdateIndentState(indentId, cidl.IndentStateFailStatistic)
			if err != nil {
				log.Warnf("update indent state failed. %s", err)
				return
			}
		}

	}()

	indent, err := dbGroupBuying.GetIndent(indentId)
	if err != nil {
		log.Warnf("get indent failed. %s", err)
		return
	}

	// 不在需要统计的状态
	if indent.State != cidl.IndentStateStatistic {
		log.Warnf("indent is not in need statistic state. %s", err)
		return
	}

	// excel文件
	indentExcel, err := excel.NewIndentExcel()
	if err != nil {
		log.Warnf("new indent excel failed. %s", err)
		return
	}
	indentExcel.Date = indent.CreateTime
	indentExcel.TickerNumber = fmt.Sprintf("%s%s",
		indentExcel.Date.Format("20060102150405"),
		indent.IndentId[:2])

	for _, taskBrief := range *indent.TasksBrief { // 逐个任务进行统计
		// 统计订货单下每个团购任务的情况
		indentStatistic, errGet := dbGroupBuying.GetIndentStatistics(indentId, taskBrief.TaskId)
		if errGet != nil {
			err = errGet
			log.Warnf("get indent statistics failed. %s", err)
			return
		}

		task := indentStatistic.TaskContent
		result := cidl.NewGroupBuyingIndentStatisticsResultType()
		communityCountMap := make(map[string]map[uint32]bool) // sku_id => grp_id => true
		for i := uint32(1); ; i++ {                           // 逐个购买记录进行统计

			// 获取团购任务购买
			list, errList := dbGroupBuying.CommunityBuyListOrderTaskIdSkuId(taskBrief.TaskId, i, 100)

			if errList != nil {
				err = errList
				log.Warnf("get community buy list failed. %s", err)
				return
			}

			if len(list) == 0 {
				break
			}

			for _, communityBuy := range list {
				buyCount := communityBuy.Count
				if buyCount <= 0 { // 购买数量不能为0
					continue
				}

				buyDetail := communityBuy.BuyDetail
				buySkuId := communityBuy.SkuId
				_, ok := communityCountMap[buySkuId]
				if !ok {
					communityCountMap[buySkuId] = make(map[uint32]bool)
				}
				_, ok = communityCountMap[buySkuId][communityBuy.GroupId]
				if !ok {
					communityCountMap[buySkuId][communityBuy.GroupId] = true
				}

				var buyItem *cidl.GroupBuyingSkuMapItem
				if !buyDetail.IsCombination { // 单品
					buyItem = buyDetail.SingleGoods

				} else { // 组合
					buyItem = cidl.NewGroupBuyingSkuMapItem()
					combinationItem := buyDetail.CombinationGoods
					label := cidl.NewGroupBuyingSpecificationItemLabel()
					label.LabelId = buySkuId
					label.Name = fmt.Sprintf("(%s)", combinationItem.Name)

					for _, subSkuItem := range combinationItem.SubItems {
						var labelKeys []string
						for key, _ := range subSkuItem.Labels {
							labelKeys = append(labelKeys, key)
						}

						sort.Strings(labelKeys)
						var labelName string
						for _, labelId := range labelKeys {
							subSkuItemLabel := subSkuItem.Labels[labelId]
							if label.Name == "" {
								labelName = subSkuItemLabel.Name
							} else {
								labelName = labelName + " " + subSkuItemLabel.Name
							}
						}
						labelName = fmt.Sprintf("%s %d份", labelName, subSkuItem.Count)
						label.Name = label.Name + "\n" + labelName
					}

					buyItem.Labels[label.LabelId] = label
					buyItem.SkuId = buySkuId
					buyItem.Name = combinationItem.Name
					buyItem.MarketPrice = combinationItem.MarketPrice
					buyItem.GroupBuyingPrice = combinationItem.GroupBuyingPrice
					buyItem.SettlementPrice = combinationItem.SettlementPrice
					buyItem.CostPrice = combinationItem.CostPrice
					buyItem.IllustrationPicture = combinationItem.IllustrationPicture
				}

				resultItem, ok := (*result)[buySkuId]
				if !ok { // 未统计
					resultItem = &cidl.GroupBuyingIndentStatisticResultItem{
						GroupBuyingSkuMapItem: cidl.GroupBuyingSkuMapItem{
							SkuId:               buyItem.SkuId,
							Name:                buyItem.Name,
							Labels:              buyItem.Labels,
							MarketPrice:         buyItem.MarketPrice,
							GroupBuyingPrice:    buyItem.GroupBuyingPrice,
							SettlementPrice:     buyItem.SettlementPrice,
							CostPrice:           buyItem.CostPrice,
							IllustrationPicture: buyItem.IllustrationPicture,
						},
						Sales:          buyCount,
						CommunityCount: 1,
						TaskId:         task.TaskId,
						StartTime:      task.StartTime,
						EndTime:        task.EndTime,
						Title:          task.Title,
					}

					(*result)[buySkuId] = resultItem

				} else {
					resultItem.Sales += buyCount

				}

			}

		}
		var resultKeys []string
		for key, _ := range *result {
			resultKeys = append(resultKeys, key)
		}
		sort.Strings(resultKeys)

		for _, skuId := range resultKeys {
			skuItem := (*result)[skuId]
			skuItem.CommunityCount = uint32(len(communityCountMap[skuId])) // 社区数目

			skuItem.TotalCost = skuItem.CostPrice * float64(skuItem.Sales)
			skuItem.TotalSettlement = skuItem.SettlementPrice * float64(skuItem.Sales)
			(*result)[skuId] = skuItem

			// excel 生成
			var labelKeys []string
			for labelKey, _ := range skuItem.Labels {
				labelKeys = append(labelKeys, labelKey)
			}
			sort.Strings(labelKeys)

			var strSpecification [3]string
			index := 0
			for _, labelKey := range labelKeys {
				label := skuItem.Labels[labelKey]
				if index > 2 {
					break
				}

				strSpecification[index] = label.Name
				index++
			}

			err = indentExcel.AddRow(
				skuItem.StartTime,
				skuItem.Title,
				strSpecification,
				skuItem.GroupBuyingPrice,
				skuItem.SettlementPrice,
				skuItem.CostPrice,
				skuItem.Sales,
				skuItem.CommunityCount,
				skuItem.TotalCost,
				skuItem.TotalSettlement)
			if err != nil {
				log.Warnf("add row failed. %s", err)
				return
			}
		}

		version := cidl.IndentStatisticsRecordVersion
		_, err = dbGroupBuying.UpdateIndentStatisticsResult(indentId, taskBrief.TaskId, result, version)
		if err != nil {
			log.Warnf("update indent statistics result failed. %s", err)
			return
		}

	}

	indentExcel.BeforeSave()

	// 生成excel
	today, err := utils.DayStartTime(time.Now())
	if err != nil {
		log.Warnf("get today start time failed. %s", err)
		return
	}

	excelUrl, err := indentExcel.SaveToQiniu(
		fmt.Sprintf("bilimall/byo/indent/%d/%d/", indent.OrganizationId, today.Unix()),
		fmt.Sprintf("销售概况单_%s_%s.xlsx", indentExcel.TickerNumber, indent.IndentId))

	if err != nil {
		log.Warnf("save indent excel file failed. %s", err)
		return
	}

	// 更新订货单状态
	_, err = dbGroupBuying.UpdateIndentStateAndExcelUrl(indentId, cidl.IndentStateFinishStatistic, excelUrl)
	if err != nil {
		log.Warnf("update indent state and excel url failed. %s", err)
		return
	}

	return
}

// 添加送货单
func (m *TopicGroupBuyingServiceHandler) AddSend(msg *AddSendMessage) (err error) {
	authorUserId := msg.AuthorUserId
	sendId := msg.SendID

	author, err := user.NewProxy("user-service").InnerUserInfoByUserID(authorUserId)
	if err != nil {
		log.Warnf("get author user failed. %s", err)
		return
	}

	dbGroupBuying := db.NewMallGroupBuyingOrder()
	defer func() {
		if err != nil {
			log.Warnf("add send failed. %s", err)
			// 更新送货单统计状态
			_, err = dbGroupBuying.UpdateSendState(sendId, cidl.GroupBuyingSendStateFailStatistic)
			if err != nil {
				log.Warnf("update send state failed. %s", err)
				return
			}
		}

	}()

	send, err := dbGroupBuying.GetSend(sendId)
	if err != nil {
		log.Warnf("get send failed. %s", err)
		return
	}

	if send.State != cidl.GroupBuyingSendStateStatistic {
		log.Warnf("send is not in need statistics state. %s", err)
		return
	}

	organization, err := agency.NewProxy("agency-service").InnerAgencyOrganizationInfoByOrganizationID(send.OrganizationId)
	if err != nil {
		log.Warnf("get organization failed. %s", err)
		return
	}

	var taskIds []uint32
	taskBriefMap := make(map[uint32]*cidl.GroupBuyingSendTaskBriefItem)
	for _, taskBrief := range *send.TasksBrief {
		taskIds = append(taskIds,taskBrief.TaskId)
		taskBriefMap[taskBrief.TaskId] = taskBrief
	}

	tasks, err := dbGroupBuying.GetFinishBuyingTasks(send.OrganizationId, taskIds)
	if err != nil {
		log.Warnf("get tasks failed. %s", err)
		return
	}

	// 社群配送
	// groupId -> taskId -> GroupBuyingSendCommunityStatisticsItem
	communityBuyMap := make(map[uint32]map[uint32]*cidl.GroupBuyingSendCommunityStatisticsItem)

	// 路线配送
	// lineId -> taskId -> GroupBuyingSendLineStatisticsItem
	lineBuyMap := make(map[uint32]map[uint32]*cidl.GroupBuyingSendLineStatisticsItem)
	// groupId -> GroupBuyingLineCommunity
	groupLineCommunityMap := make(map[uint32]*cidl.GroupBuyingLineCommunity)

	// 路线涉及社区数统计
	// lineId -> groupId -> bool
	lineCommunityCountMap := make(map[uint32]map[uint32]bool)

	for _,taskBrief := range taskBriefMap {
		//更新已导出线路	
		_, err = dbGroupBuying.UpdateTaskSelectedLines(taskBrief.TaskId, taskBrief.LineIds)
		if err != nil {
			log.Warnf("update task selected line failed. %s", err)
			return
		}

		list, errList := dbGroupBuying.CommunityBuyListByTaskIdLineIdsOrderSkuId(taskBrief.TaskId, taskBrief.LineIds)

		if errList != nil {
			err = errList
			log.Warnf("get community buy list by task ids failed. %s", err)
			return
		}

		if len(list) == 0 {
			// 没有记录
			continue	
		}

		for _, communityBuy := range list { // 逐个购买记录
			if communityBuy.Count <= 0 {
				continue
			}

			groupId := communityBuy.GroupId
			taskId := communityBuy.TaskId
			taskTitle := taskBriefMap[taskId].Title

			// 社群相关
			taskBuyMap, ok := communityBuyMap[groupId]
			if !ok {
				taskBuyMap = make(map[uint32]*cidl.GroupBuyingSendCommunityStatisticsItem)
				communityBuyMap[groupId] = taskBuyMap
			}

			skuBuyMap, ok := taskBuyMap[taskId]
			if !ok {
				skuBuyMap = cidl.NewGroupBuyingSendCommunityStatisticsItem()
				skuBuyMap.TaskId = taskId
				skuBuyMap.TaskTitle = taskTitle
				taskBuyMap[taskId] = skuBuyMap
			}

			// 路线相关
			lineCommunity, ok := groupLineCommunityMap[groupId]
			if !ok {
				lineCommunity, err = dbGroupBuying.GetLineCommunity(groupId)
				if err != nil {
					log.Warnf("get line community failed. %s", err)
					return
				}
				groupLineCommunityMap[groupId] = lineCommunity
			}

			lineId := lineCommunity.LineId
			lineTaskBuyMap, ok := lineBuyMap[lineId]
			if !ok {
				lineTaskBuyMap = make(map[uint32]*cidl.GroupBuyingSendLineStatisticsItem)
				lineBuyMap[lineId] = lineTaskBuyMap
			}

			lineSkuBuyMap, ok := lineTaskBuyMap[taskId]
			if !ok {
				lineSkuBuyMap = cidl.NewGroupBuyingSendLineStatisticsItem()
				lineSkuBuyMap.TaskId = taskId
				lineSkuBuyMap.TaskTitle = taskTitle
				lineTaskBuyMap[taskId] = lineSkuBuyMap
			}

			// 路线涉及社区数统计
			groupCountMap, ok := lineCommunityCountMap[lineId]
			if !ok {
				groupCountMap = make(map[uint32]bool)
				lineCommunityCountMap[lineId] = groupCountMap
			}
			groupCountMap[groupId] = true

			// 统计商品售卖
			buySkuId := communityBuy.SkuId
			buyCount := communityBuy.Count
			var buyItem *cidl.GroupBuyingSkuMapItem
			if !communityBuy.BuyDetail.IsCombination { // 单品
				buyItem = communityBuy.BuyDetail.SingleGoods

			} else { // 组合
				buyItem = cidl.NewGroupBuyingSkuMapItem()

				combinationItem := communityBuy.BuyDetail.CombinationGoods
				buyItem.SkuId = buySkuId
				label := cidl.NewGroupBuyingSpecificationItemLabel()
				label.LabelId = buySkuId
				label.Name = combinationItem.Name

				//label.Name = fmt.Sprintf("(%s)", combinationItem.Name)
				//
				//for _, subSkuItem := range combinationItem.SubItems {
				//	var labelKeys []string
				//	for key, _ := range subSkuItem.Labels {
				//		labelKeys = append(labelKeys, key)
				//	}
				//
				//	sort.Strings(labelKeys)
				//	var labelName string
				//	for _, labelId := range labelKeys {
				//		subSkuItemLabel := subSkuItem.Labels[labelId]
				//		if label.Name == "" {
				//			labelName = subSkuItemLabel.Name
				//		} else {
				//			labelName = labelName + " " + subSkuItemLabel.Name
				//		}
				//	}
				//	labelName = fmt.Sprintf("%s %d份", labelName, subSkuItem.Count)
				//	label.Name = label.Name + "\n" + labelName
				//}

				buyItem.Labels[label.LabelId] = label
				buyItem.Name = combinationItem.Name
				buyItem.MarketPrice = combinationItem.MarketPrice
				buyItem.GroupBuyingPrice = combinationItem.GroupBuyingPrice
				buyItem.SettlementPrice = combinationItem.SettlementPrice
				buyItem.CostPrice = combinationItem.CostPrice
				buyItem.IllustrationPicture = combinationItem.IllustrationPicture
			}

			// 社群相关
			buySkuItem, ok := skuBuyMap.Sku[buySkuId]
			if !ok {
				buySkuItem = &cidl.GroupBuyingSendCommunityStatisticsItemSkuItem{
					GroupBuyingSkuMapItem: cidl.GroupBuyingSkuMapItem{
						SkuId:               buyItem.SkuId,
						Name:                buyItem.Name,
						Labels:              buyItem.Labels,
						MarketPrice:         buyItem.MarketPrice,
						GroupBuyingPrice:    buyItem.GroupBuyingPrice,
						SettlementPrice:     buyItem.SettlementPrice,
						CostPrice:           buyItem.CostPrice,
						IllustrationPicture: buyItem.IllustrationPicture,
					},
					Sales:  buyCount,
					TaskId: taskId,
					Title:  taskTitle,
				}

				skuBuyMap.Sku[buySkuId] = buySkuItem

			} else {
				buySkuItem.Sales += buyCount

			}

			// 路线相关
			lineBuySkuItem, ok := lineSkuBuyMap.Sku[buySkuId]
			if !ok {
				lineBuySkuItem = &cidl.GroupBuyingSendLineStatisticsItemSkuItem{
					GroupBuyingSkuMapItem: cidl.GroupBuyingSkuMapItem{
						SkuId:               buyItem.SkuId,
						Name:                buyItem.Name,
						Labels:              buyItem.Labels,
						MarketPrice:         buyItem.MarketPrice,
						GroupBuyingPrice:    buyItem.GroupBuyingPrice,
						SettlementPrice:     buyItem.SettlementPrice,
						CostPrice:           buyItem.CostPrice,
						IllustrationPicture: buyItem.IllustrationPicture,
					},
					Sales:          buyCount,
					TaskId:         taskId,
					CommunityCount: 1,
					Title:          taskTitle,
				}

				lineSkuBuyMap.Sku[buySkuId] = lineBuySkuItem

			} else {
				lineBuySkuItem.Sales += buyCount

			}

		}

	}

	// 社群相关
	// 各个社群配送情况
	// 社群ID -> 购买情况
	sendCommunityMap := make(map[uint32]*cidl.GroupBuyingSendCommunity)

	// 社群map
	groupMap := make(map[uint32]*community.Group)
	for groupId, taskBuyMap := range communityBuyMap {
		group, ok := groupMap[groupId]
		if !ok {
			group, err = community.NewProxy("community-service").InnerCommunityGroupInfoByGroupID(groupId)
			if err != nil {
				log.Warnf("get community group info from proxy failed. %s", err)
				return
			}
			groupMap[groupId] = group
		}

		sendCommunity := &cidl.GroupBuyingSendCommunity{
			SendId:                  sendId,
			GroupId:                 groupId,
			GroupName:               group.Name,
			GroupAddress:            group.Address,
			GroupManagerUid:         group.ManagerUserId,
			GroupManagerName:        group.ManagerName,
			GroupManagerMobile:      group.ManagerMobile,
			OrganizationId:          organization.OrganizationId,
			OrganizationName:        organization.Name,
			OrganizationAddress:     organization.Address,
			OrganizationManagerUid:  organization.ManagerUserId,
			OrganizationManagerName: organization.ManagerName,
			AuthorUid:               author.UserID,
			AuthorName:              author.Name,
			Statistics:              cidl.NewGroupBuyingSendCommunityStatisticType(),
			SendTime:                time.Now(),
			Version:                 cidl.SendCommunityRecordVersion,
		}

		for _, statisticItem := range taskBuyMap {
			for skuId, skuItem := range statisticItem.Sku {
				skuItem.TotalCost = skuItem.CostPrice * float64(skuItem.Sales)
				skuItem.TotalSettlement = skuItem.SettlementPrice * float64(skuItem.Sales)
				statisticItem.Sku[skuId] = skuItem

				sendCommunity.SettlementAmount += skuItem.TotalSettlement
			}

			*sendCommunity.Statistics = append(*sendCommunity.Statistics, statisticItem)
		}

		_, err = dbGroupBuying.AddSendCommunity(sendCommunity)
		if err != nil {
			log.Warnf("add send community failed. %s", err)
			return
		}

		sendCommunityMap[groupId] = sendCommunity

	}

	// 产生表格
	sendExcel, err := excel.NewSendExcel()
	if err != nil {
		log.Warnf("new send excel failed. %s", err)
		return
	}

	sendExcel.Date = send.CreateTime
	sendExcel.OrganizationName = send.OrganizationName
	sendExcel.TicketNumber = fmt.Sprintf("%s%s",
		sendExcel.Date.Format("20060102150405"),
		send.SendId[:4])

	sendLineMap := make(map[uint32]*cidl.GroupBuyingSendLine)

	// 团购任务相关
	err = sendExcel.AddTaskSummarySheets("团购任务销售汇总", tasks, groupMap, sendCommunityMap)
	if err != nil {
		log.Warnf("add task summary sheets failed. %s", err)
		return
	}

	// 路线相关
	for lineId, taskBuyMap := range lineBuyMap {
		line, errGet := dbGroupBuying.GetLine(lineId)
		if errGet != nil {
			log.Warnf("get line from db failed. %s", err)
			return
		}

		groupCountMap := lineCommunityCountMap[lineId]
		communityCount := len(groupCountMap)
		sendLine := &cidl.GroupBuyingSendLine{
			SendId:           sendId,
			LineId:           lineId,
			LineName:         line.Name,
			OrganizationId:   line.OrganizationId,
			OrganizationName: line.OrganizationName,
			CommunityCount:   uint32(communityCount),
			SendTime:         time.Now(),
			Version:          cidl.SendCommunityRecordVersion,
		}

		sendLine.Statistics = cidl.NewGroupBuyingSendLineStatisticType()

		// 计算总结算金额
		settlementAmount := float64(0)
		for _, skuBuyMap := range taskBuyMap {
			for skuId, skuItem := range skuBuyMap.Sku {
				skuItem.TotalCost = skuItem.CostPrice * float64(skuItem.Sales)
				skuItem.TotalSettlement = skuItem.SettlementPrice * float64(skuItem.Sales)
				settlementAmount += skuItem.TotalSettlement
				skuBuyMap.Sku[skuId] = skuItem
			}

			(*sendLine.Statistics) = append((*sendLine.Statistics), skuBuyMap)
		}
		sendLine.SettlementAmount = settlementAmount

		_, err = dbGroupBuying.AddSendLine(sendLine)
		if err != nil {
			log.Warnf("add send line failed. %s", err)
			return
		}

		// excel 相关
		sendLineMap[lineId] = sendLine
	}

	// 团长销售额简报 
	groupSummarySheets := sendExcel.GetGroupSummarySheets()
	for _, groupCountMap := range lineCommunityCountMap {
		for groupId, _ := range groupCountMap {
			sendCommunity := sendCommunityMap[groupId]
		       	groupSummarySheets.AddLineRow(
				sendCommunity.GroupName,
				sendCommunity.GroupManagerName,
				sendCommunity.GroupManagerMobile,
				sendCommunity.SettlementAmount,
			)
		}
	}


	// 路线团购任务统计
	for lineId, groupCountMap := range lineCommunityCountMap {
		sendLine := sendLineMap[lineId]

		// 生成社群
		lineSendCommunityMap := make(map[uint32]*cidl.GroupBuyingSendCommunity)
		for groupId, _ := range groupCountMap {
			sendCommunity := sendCommunityMap[groupId]
			lineSendCommunityMap[sendCommunity.GroupId] = sendCommunity
		}

		sheetName := sendLine.LineName + "销售汇总"
		err = sendExcel.AddTaskSummarySheets(sheetName, tasks, groupMap, lineSendCommunityMap)
		if err != nil {
			log.Warnf("add line task summary sheets failed. %s", err)
			return
		}

	}

	// 路线总统计
	sendLineSummarySheets := sendExcel.GetSendLineSummarySheets()
	for lineId, _ := range lineCommunityCountMap {
		sendLine := sendLineMap[lineId]
		sendLineSummarySheets.AddLineRow(
			sendLine.LineName,
			sendLine.CommunityCount,
			sendLine.SettlementAmount,
		)
	}


	// 生成路线商品统计和社群商品统计
	for lineId, groupCountMap := range lineCommunityCountMap {
		sendLine := sendLineMap[lineId]
		err = sendExcel.AddSendLineSheets(sendLine)
		if err != nil {
			log.Warnf("add send line sheets failed. %s", err)
			return
		}

		// 生成社群
		var sendCommunities []*cidl.GroupBuyingSendCommunity
		for groupId, _ := range groupCountMap {
			sendCommunity := sendCommunityMap[groupId]
			err = sendExcel.AddSendCommunitySheets(sendCommunity)
			if err != nil {
				log.Warnf("add send community sheets failed. %s", err)
				return
			}

			sendCommunities = append(sendCommunities, sendCommunity)
		}

	}

	sendExcel.BeforeSave()

	today, err := utils.TodayStartTime()
	if err != nil {
		log.Warnf("get today start time failed. %s", err)
		return
	}

	// 生成excel到七牛
	excelUrl, err := sendExcel.SaveToQiniu(
		fmt.Sprintf("bilimall/byo/send/%d/%d/", send.OrganizationId, today.Unix()),
		fmt.Sprintf("送货单_%s_%s.xlsx", sendExcel.TicketNumber, send.SendId),
	)

	if err != nil {
		log.Warnf("save send excel file failed. %s", err)
		return
	}

	// 更新送货单状态
	_, err = dbGroupBuying.UpdateSendStatisticResult(sendId, cidl.GroupBuyingSendStateFinishStatistic, excelUrl, cidl.SendRecrodVersion)
	if err != nil {
		log.Warnf("update send state and excel url failed. %s", err)
		return
	}

	return
}
