package impls

import (
	"fmt"
	"time"

	"business/agency/proxy/agency"
	"business/community/proxy/community"
	"business/group-buying-order/cidl"
	"business/group-buying-order/common/db"
	"business/group-buying-order/common/mq"
	"common/api"
	"common/file"

	"github.com/mz-eco/mz/conn"
	"github.com/mz-eco/mz/http"
	"github.com/mz-eco/mz/utils"
)

func init() {
	AddOrgTaskPicTokenByOrganizationIDHandler()
	AddOrgTaskAddByOrganizationIDHandler()
	AddOrgTaskEditByOrganizationIDByTaskIDHandler()
	AddOrgTaskInfoByTaskIDHandler()

	AddOrgTaskShowByOrganizationIDByTaskIDHandler()
	AddOrgTaskHideByOrganizationIDByTaskIDHandler()
	AddOrgTaskDeleteByOrganizationIDByTaskIDHandler()

	AddOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskIDHandler()
	AddOrgTaskSoldListByOrganizationIDHandler()
	AddOrgTaskFinishBuyingGroupListByOrganizationIDByTaskIDHandler()
	AddOrgTaskFinishBuyingByOrganizationIDByTaskIDHandler()

	AddOrgSendAddByOrganizationIDHandler()

	AddOrgIndentAddByOrganizationIDHandler()

}

// 团购任务图片上传TOKEN
type OrgTaskPicTokenByOrganizationIDImpl struct {
	cidl.ApiOrgTaskPicTokenByOrganizationID
}

func AddOrgTaskPicTokenByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_PIC_TOKEN_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgTaskPicTokenByOrganizationIDImpl{
				ApiOrgTaskPicTokenByOrganizationID: cidl.MakeApiOrgTaskPicTokenByOrganizationID(),
			}
		},
	)
}

func (m *OrgTaskPicTokenByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	organizationId := m.Params.OrganizationID
	today, err := utils.DayStartTime(time.Now())
	if err != nil {
		ctx.Errorf(api.ErrServer, "get day start time failed. %s", err)
		return
	}

	qiniu, err := file.GetQiniuPublicBucket()
	if err != nil {
		ctx.Errorf(api.ErrServer, "get qiniu public bucket failed. %s", err)
		return
	}

	prefix := fmt.Sprintf("billimall/byo/task/%d/%d/", organizationId, today.Unix())
	for _, fileName := range m.Ask.FileNames {
		if fileName == "" {
			ctx.Errorf(api.ErrWrongParams, "empty pic file name. %s", err)
			return
		}

		token, key, err := qiniu.GenerateUploadToken(fileName, prefix)
		if err != nil {
			return
		}

		m.Ack.Tokens = append(m.Ack.Tokens, &cidl.AckPicToken{
			OriginalFileName: fileName,
			Token:            token,
			Key:              key,
			StoreUrl:         qiniu.StoreUrl(key),
			AccessUrl:        qiniu.StoreUrl(key),
		})

	}

	ctx.Json(m.Ack)
}

// 添加团购任务
type OrgTaskAddByOrganizationIDImpl struct {
	cidl.ApiOrgTaskAddByOrganizationID
}

func AddOrgTaskAddByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_ADD_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgTaskAddByOrganizationIDImpl{
				ApiOrgTaskAddByOrganizationID: cidl.MakeApiOrgTaskAddByOrganizationID(),
			}
		},
	)
}

func (m *OrgTaskAddByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	task := cidl.NewGroupBuyingOrderTaskCustom()

	ask := m.Ask
	task.OrganizationId = m.Params.OrganizationID
	task.Title = ask.Title
	task.Introduction = ask.Introduction
	switch ask.SellType {
	case cidl.GroupBuyingOrderTaskSellTypeDefault:
		task.SellType = cidl.GroupBuyingOrderTaskSellTypeDefault
		task.ShowStartTime = ask.ShowStartTime
		task.StartTime = ask.ShowStartTime

	case cidl.GroupBuyingOrderTaskSellTypeSeckill:
		if ask.ShowStartTime.After(ask.StartTime) {
			ctx.Errorf(api.ErrWrongParams, "show start time can not be after start time")
			return
		}

		task.SellType = cidl.GroupBuyingOrderTaskSellTypeSeckill
		task.ShowStartTime = ask.ShowStartTime
		task.StartTime = ask.StartTime
	default:
		ctx.Errorf(api.ErrWrongParams, "unknown sell type.")
		return
	}
	task.EndTime = ask.EndTime
	task.CoverPicture = ask.CoverPicture

	now := time.Now()
	if !task.StartTime.Before(task.EndTime) || !now.Before(task.EndTime) {
		ctx.Errorf(api.ErrWrongParams, "wrong start time or end time")
		return
	}

	// 检查配图
	for _, ackIllustrationPicture := range *ask.IllustrationPictures {
		if ackIllustrationPicture == "" {
			ctx.Errorf(api.ErrWrongParams, "empty illustration picture")
			return
		}
	}

	task.IllustrationPictures = ask.IllustrationPictures
	for _, infoItem := range *ask.Info {
		if infoItem.Title == "" || infoItem.Content == "" {
			ctx.Errorf(api.ErrWrongParams, "empty info title or content")
			return
		}
	}

	task.Info = ask.Info
	task.WxSellText = ask.WxSellText
	task.Notes = ask.Notes
	task.ShowState = ask.ShowState
	task.Version = cidl.GroupBuyingTaskRecordVersion
	task.Specification, err = cidl.NewGroupBuyingTaskOrderSpecificationByAsk(ask.Specification, ask.Sku, ask.Combination)
	if err != nil {
		ctx.Errorf(api.ErrWrongParams, "illegal specification. %s", err)
		return
	}

	dbGroupBuying := db.NewMallGroupBuyingOrder()
	addResult, err := dbGroupBuying.AddTask(task)
	if err != nil {
		ctx.Errorf(api.ErrDBInsertFailed, "add task failed. %s", err)
		return
	}
	taskId, err := addResult.LastInsertId()
	if err != nil {
		ctx.Errorf(api.ErrServer, "get add task last insert id failed. %s", err)
		return
	}

	// 添加库存
	var inventories []*cidl.GroupBuyingOrderInventory
	for _, skuItem := range task.Specification.SkuMap {
		inventory := cidl.NewGroupBuyingOrderInventory()
		inventory.TaskId = uint32(taskId)
		inventory.SkuId = skuItem.SkuId
		inventory.Total = skuItem.InventoryCount
		inventory.Surplus = inventory.Total
		inventories = append(inventories, inventory)
	}

	_, err = dbGroupBuying.AddInventories(inventories)
	if err != nil {
		ctx.Errorf(api.ErrDBInsertFailed, "add task inventories failed. %s", err)
		return
	}

	ctx.Succeed()
}

// 编辑团购任务
type OrgTaskEditByOrganizationIDByTaskIDImpl struct {
	cidl.ApiOrgTaskEditByOrganizationIDByTaskID
}

func AddOrgTaskEditByOrganizationIDByTaskIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_EDIT_BY_ORGANIZATION_ID_BY_TASK_ID,
		func() http.ApiHandler {
			return &OrgTaskEditByOrganizationIDByTaskIDImpl{
				ApiOrgTaskEditByOrganizationIDByTaskID: cidl.MakeApiOrgTaskEditByOrganizationIDByTaskID(),
			}
		},
	)
}

func (m *OrgTaskEditByOrganizationIDByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	taskId := m.Params.TaskID
	ask := m.Ask

	// 获取原团购任务
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	oldTask, err := dbGroupBuying.GetTask(taskId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task failed. %s", err)
		return
	}

	task := cidl.NewGroupBuyingOrderTaskCustom()
	task.TaskId = taskId
	task.OrganizationId = m.Params.OrganizationID
	task.Title = ask.Title
	task.Introduction = ask.Introduction
	switch ask.SellType {
	case cidl.GroupBuyingOrderTaskSellTypeDefault:
		task.SellType = cidl.GroupBuyingOrderTaskSellTypeDefault
		task.ShowStartTime = ask.ShowStartTime
		task.StartTime = ask.ShowStartTime

	case cidl.GroupBuyingOrderTaskSellTypeSeckill:
		if ask.ShowStartTime.After(ask.StartTime) {
			ctx.Errorf(api.ErrWrongParams, "show start time can not be after start time")
			return
		}

		task.SellType = cidl.GroupBuyingOrderTaskSellTypeSeckill
		task.ShowStartTime = ask.ShowStartTime
		task.StartTime = ask.StartTime
	default:
		ctx.Errorf(api.ErrWrongParams, "unknown sell type.")
		return
	}
	task.EndTime = ask.EndTime
	task.CoverPicture = ask.CoverPicture

	now := time.Now()
	if !task.StartTime.Before(task.EndTime) || !now.Before(task.EndTime) {
		ctx.Errorf(cidl.ErrTaskAddEditIllegalTime, "wrong start time or end time")
		return
	}

	// 检查配图
	for _, ackIllustrationPicture := range *ask.IllustrationPictures {
		if ackIllustrationPicture == "" {
			ctx.Errorf(api.ErrWrongParams, "empty illustration picture")
			return
		}
	}

	task.IllustrationPictures = ask.IllustrationPictures
	for _, infoItem := range *ask.Info {
		if infoItem.Title == "" || infoItem.Content == "" {
			ctx.Errorf(api.ErrWrongParams, "empty info title or content")
			return
		}
	}

	task.Info = ask.Info
	task.WxSellText = ask.WxSellText
	task.Notes = ask.Notes
	task.ShowState = ask.ShowState
	task.SellType = ask.SellType
	task.Version = cidl.GroupBuyingTaskRecordVersion
	task.Specification, err = cidl.NewGroupBuyingTaskOrderSpecificationByAsk(ask.Specification, ask.Sku, ask.Combination)
	if err != nil {
		ctx.Errorf(api.ErrWrongParams, "illegal specification. %s", err)
		return
	}

	// 团购状态为进行中时，只允许编辑部分信息
	if oldTask.ShowState == cidl.GroupBuyingTaskShowStateShow && oldTask.GroupStateShowIsInProgress(now) {
		_, err = dbGroupBuying.UpdateInProgressTask(task)
		if err != nil {
			ctx.Errorf(api.ErrDBUpdateFailed, "add task failed. %s", err)
			return
		}

		ctx.Succeed()
		return

	} else if oldTask.GroupState == cidl.GroupBuyingTaskGroupStateNotStart && (oldTask.ShowState == cidl.GroupBuyingTaskShowStateHidden || now.Before(oldTask.StartTime)) { // 团购状态为未开始
		_, err = dbGroupBuying.UpdateNotStartTask(task)
		if err != nil {
			ctx.Errorf(api.ErrDBUpdateFailed, "update task failed. %s", err)
			return
		}

		// 删除原来库存
		_, err = dbGroupBuying.DeleteInventory(taskId)
		if err != nil {
			ctx.Errorf(api.ErrDbDeleteRecordFailed, "delete task inventory failed. %s", err)
			return
		}

		// 添加库存
		var inventories []*cidl.GroupBuyingOrderInventory
		for _, skuItem := range task.Specification.SkuMap {
			inventory := cidl.NewGroupBuyingOrderInventory()
			inventory.TaskId = uint32(taskId)
			inventory.SkuId = skuItem.SkuId
			inventory.Total = skuItem.InventoryCount
			inventory.Surplus = inventory.Total
			inventories = append(inventories, inventory)
		}

		_, err = dbGroupBuying.AddInventories(inventories)
		if err != nil {
			ctx.Errorf(api.ErrDBInsertFailed, "add task inventories failed. %s", err)
			return
		}

		ctx.Succeed()
		return

	}

	ctx.Errorf(cidl.ErrTaskGroupStateIsNotInNotStartOrInProgress, "task is not in editable state")
}

// 获取团购任务
type OrgTaskInfoByTaskIDImpl struct {
	cidl.ApiOrgTaskInfoByTaskID
}

func AddOrgTaskInfoByTaskIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_INFO_BY_TASK_ID,
		func() http.ApiHandler {
			return &OrgTaskInfoByTaskIDImpl{
				ApiOrgTaskInfoByTaskID: cidl.MakeApiOrgTaskInfoByTaskID(),
			}
		},
	)
}

func (m *OrgTaskInfoByTaskIDImpl) Handler(ctx *http.Context) {
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

// 上架
type OrgTaskShowByOrganizationIDByTaskIDImpl struct {
	cidl.ApiOrgTaskShowByOrganizationIDByTaskID
}

func AddOrgTaskShowByOrganizationIDByTaskIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_SHOW_BY_ORGANIZATION_ID_BY_TASK_ID,
		func() http.ApiHandler {
			return &OrgTaskShowByOrganizationIDByTaskIDImpl{
				ApiOrgTaskShowByOrganizationIDByTaskID: cidl.MakeApiOrgTaskShowByOrganizationIDByTaskID(),
			}
		},
	)
}

func (m *OrgTaskShowByOrganizationIDByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	organizationId := m.Params.OrganizationID
	taskId := m.Params.TaskID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	task, err := dbGroupBuying.GetTaskByOrganizationIdAndTaskId(organizationId, taskId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task failed. %s", err)
		return
	}

	if task.ShowState == cidl.GroupBuyingTaskShowStateShow {
		ctx.Succeed()
		return
	}

	// 只对未开始但是未过结束时间的任务进行上架
	now := time.Now()
	if !(task.GroupState == cidl.GroupBuyingTaskGroupStateNotStart && task.EndTime.After(now)) {
		ctx.Errorf(cidl.ErrTaskGroupStateIsNotInNotStart, "task group state is not in not start state.")
		return
	}

	// 只能对未进行且未超过结束时间的任务进行上架
	if task.StartTime.After(now) { // 还未到开团时间
		strSql := `
			UPDATE byo_task
			SET
				show_state=?
			WHERE
				tsk_id=?
				AND org_id=?
				AND show_state=?
				AND group_state=?
				AND now()<end_time
		`
		_, err = dbGroupBuying.DB.Exec(strSql,
			cidl.GroupBuyingTaskShowStateShow,
			taskId,
			organizationId,
			cidl.GroupBuyingTaskShowStateHidden,
			cidl.GroupBuyingTaskGroupStateNotStart)
		if err != nil {
			ctx.Errorf(api.ErrDBUpdateFailed, "update task show state failed. %s", err)
			return
		}

	} else { // 在开团时间范围内
		strSql := `
			UPDATE byo_task
			SET
				show_state=?,
				group_state=?
			WHERE
				tsk_id=?
				AND org_id=?
				AND show_state=?
				AND group_state=?
				AND now()>=start_time
				AND now()<end_time
		`
		_, err = dbGroupBuying.DB.Exec(strSql,
			cidl.GroupBuyingTaskShowStateShow,
			cidl.GroupBuyingTaskGroupStateInProgress,
			taskId,
			organizationId,
			cidl.GroupBuyingTaskShowStateHidden,
			cidl.GroupBuyingTaskGroupStateNotStart)
		if err != nil {
			ctx.Errorf(api.ErrDBUpdateFailed, "update task show state failed. %s", err)
			return
		}

	}

	ctx.Succeed()
}

// 下架
type OrgTaskHideByOrganizationIDByTaskIDImpl struct {
	cidl.ApiOrgTaskHideByOrganizationIDByTaskID
}

func AddOrgTaskHideByOrganizationIDByTaskIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_HIDE_BY_ORGANIZATION_ID_BY_TASK_ID,
		func() http.ApiHandler {
			return &OrgTaskHideByOrganizationIDByTaskIDImpl{
				ApiOrgTaskHideByOrganizationIDByTaskID: cidl.MakeApiOrgTaskHideByOrganizationIDByTaskID(),
			}
		},
	)
}

func (m *OrgTaskHideByOrganizationIDByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	organizationId := m.Params.OrganizationID
	taskId := m.Params.TaskID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	task, err := dbGroupBuying.GetTaskByOrganizationIdAndTaskId(organizationId, taskId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task failed. %s", err)
		return
	}

	if task.ShowState == cidl.GroupBuyingTaskShowStateHidden {
		ctx.Succeed()
		return
	}

	now := time.Now()
	if task.GroupStateShowIsNotStart(now) { // 如果任务不在进行中，则直接下架
		strSql := `
			UPDATE byo_task
			SET	show_state=?
			WHERE
				tsk_id=?
				AND org_id=?
				AND show_state=?
				AND ( group_state=? AND now()<start_time )
		`
		_, err = dbGroupBuying.DB.Exec(strSql,
			cidl.GroupBuyingTaskShowStateHidden,
			taskId,
			organizationId,
			cidl.GroupBuyingTaskShowStateShow,
			cidl.GroupBuyingTaskGroupStateNotStart)

		if err != nil {
			ctx.Errorf(api.ErrDBUpdateFailed, "update not start task hidden failed. %s", err)
			return
		}

		ctx.Succeed()
		return

	} else if task.GroupStateShowIsInProgress(now) { // 如果任务正在进行中，则下架并截单
		strSql := `
			UPDATE byo_task
			SET
				show_state=?,
				group_state=?
			WHERE
				tsk_id=?
				AND org_id=?
				AND show_state=?
				AND ( group_state=? AND now()>=start_time AND now()<end_time )
		`
		_, err = dbGroupBuying.DB.Exec(strSql,
			cidl.GroupBuyingTaskShowStateHidden,
			cidl.GroupBuyingTaskGroupStateFinishOrdering,
			taskId,
			organizationId,
			cidl.GroupBuyingTaskShowStateShow,
			cidl.GroupBuyingTaskGroupStateInProgress)

		if err != nil {
			ctx.Errorf(api.ErrDBUpdateFailed, "update in progress task hidden failed. %s", err)
			return
		}

		ctx.Succeed()
		return

	} else {
		ctx.Errorf(cidl.ErrTaskGroupStateIsNotInNotStartOrInProgress, "group state is not in not start or in progress . %s", err)
	}

}

// 删除团购任务
type OrgTaskDeleteByOrganizationIDByTaskIDImpl struct {
	cidl.ApiOrgTaskDeleteByOrganizationIDByTaskID
}

func AddOrgTaskDeleteByOrganizationIDByTaskIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_DELETE_BY_ORGANIZATION_ID_BY_TASK_ID,
		func() http.ApiHandler {
			return &OrgTaskDeleteByOrganizationIDByTaskIDImpl{
				ApiOrgTaskDeleteByOrganizationIDByTaskID: cidl.MakeApiOrgTaskDeleteByOrganizationIDByTaskID(),
			}
		},
	)
}

func (m *OrgTaskDeleteByOrganizationIDByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	organizationId := m.Params.OrganizationID
	taskId := m.Params.TaskID

	dbGroupBuying := db.NewMallGroupBuyingOrder()
	_, err = dbGroupBuying.UpdateTaskIsDelete(organizationId, taskId, true)
	if err != nil {
		ctx.Errorf(api.ErrDBUpdateFailed, "update task is delete failed. %s", err)
		return
	}

	ctx.Succeed()
}

// 获取结团信息
type OrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskIDImpl struct {
	cidl.ApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID
}

func AddOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_FINISH_BUYING_CONFIRM_COUNT_BY_ORGANIZATION_ID_BY_TASK_ID,
		func() http.ApiHandler {
			return &OrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskIDImpl{
				ApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID: cidl.MakeApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID(),
			}
		},
	)
}

func (m *OrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	organizationId := m.Params.OrganizationID
	taskId := m.Params.TaskID
	ackCount, err := community.NewProxy("community-service").InnerCommunityGroupCountByOrganizationID(organizationId)
	if err != nil {
		ctx.Errorf(api.ErrProxyFailed, "get community group count from proxy failed. %s", err)
		return
	}
	m.Ack.TotalGroupCount = ackCount.Count

	strSql := `SELECT COUNT(DISTINCT grp_id) FROM byo_community_buy WHERE tsk_id=?`
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	var buyGroupCount uint32
	err = dbGroupBuying.DB.Get(&buyGroupCount, strSql, taskId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community group buy count failed. %s", err)
		return
	}

	m.Ack.BuyGroupCount = buyGroupCount
	m.Ack.NotBuyGroupCount = m.Ack.TotalGroupCount - m.Ack.BuyGroupCount

	ctx.Json(m.Ack)
}

// 已销售出去并且已经结团的团购任务列表
type OrgTaskSoldListByOrganizationIDImpl struct {
	cidl.ApiOrgTaskSoldListByOrganizationID
}

func AddOrgTaskSoldListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_SOLD_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgTaskSoldListByOrganizationIDImpl{
				ApiOrgTaskSoldListByOrganizationID: cidl.MakeApiOrgTaskSoldListByOrganizationID(),
			}
		},
	)
}

func (m *OrgTaskSoldListByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	search := m.Query.Search
	if search == "" {
		ack.Count, err = dbGroupBuying.TaskSoldFinishBuyingCount(organizationId)
	} else {
		ack.Count, err = dbGroupBuying.TaskSoldFinishBuyingSearchCount(organizationId, search)
	}

	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task count failed. %s", err)
		return
	}

	if ack.Count == 0 {
		ctx.Json(ack)
		return
	}

	if search == "" {
		ack.List, err = dbGroupBuying.TaskSoldFinishBuyingList(organizationId, m.Query.Page, m.Query.PageSize, false)
	} else {
		ack.List, err = dbGroupBuying.TaskSoldFinishBuyingSearchList(organizationId, search, m.Query.Page, m.Query.PageSize, false)
	}

	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task list failed. %s", err)
		return
	}

	ctx.Json(ack)
}

// 社群团购信息列表
type OrgTaskFinishBuyingGroupListByOrganizationIDByTaskIDImpl struct {
	cidl.ApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID
}

func AddOrgTaskFinishBuyingGroupListByOrganizationIDByTaskIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_FINISH_BUYING_GROUP_LIST_BY_ORGANIZATION_ID_BY_TASK_ID,
		func() http.ApiHandler {
			return &OrgTaskFinishBuyingGroupListByOrganizationIDByTaskIDImpl{
				ApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID: cidl.MakeApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID(),
			}
		},
	)
}

func (m *OrgTaskFinishBuyingGroupListByOrganizationIDByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	taskId := m.Params.TaskID

	dbGroupBuying := db.NewMallGroupBuyingOrder()
	m.Ack.Count, err = dbGroupBuying.CommunityBuyCount(taskId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community buy count failed. %s", err)
		return
	}

	if m.Ack.Count == 0 {
		ctx.Json(m.Ack)
		return
	}

	m.Ack.List, err = dbGroupBuying.CommunityBuyList(taskId, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get community buy list failed. %s", err)
		return
	}

	ctx.Json(m.Ack)
}

// 确定结团
type OrgTaskFinishBuyingByOrganizationIDByTaskIDImpl struct {
	cidl.ApiOrgTaskFinishBuyingByOrganizationIDByTaskID
}

func AddOrgTaskFinishBuyingByOrganizationIDByTaskIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_FINISH_BUYING_BY_ORGANIZATION_ID_BY_TASK_ID,
		func() http.ApiHandler {
			return &OrgTaskFinishBuyingByOrganizationIDByTaskIDImpl{
				ApiOrgTaskFinishBuyingByOrganizationIDByTaskID: cidl.MakeApiOrgTaskFinishBuyingByOrganizationIDByTaskID(),
			}
		},
	)
}

func (m *OrgTaskFinishBuyingByOrganizationIDByTaskIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	organizationId := m.Params.OrganizationID
	taskId := m.Params.TaskID

	dbGroupBuying := db.NewMallGroupBuyingOrder()
	result, err := dbGroupBuying.UpdateTaskGroupState(organizationId, taskId, cidl.GroupBuyingTaskGroupStateFinishBuying)
	if err != nil {
		ctx.Errorf(api.ErrDBUpdateFailed, "set task finish group buying failed. %s", err)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		ctx.Errorf(api.ErrDBUpdateFailed, "get rows affected failed.")
		return
	}

	ctx.Succeed()
}

// 绑定订货单团购任务
type OrgIndentAddByOrganizationIDImpl struct {
	cidl.ApiOrgIndentAddByOrganizationID
}

func AddOrgIndentAddByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_INDENT_ADD_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgIndentAddByOrganizationIDImpl{
				ApiOrgIndentAddByOrganizationID: cidl.MakeApiOrgIndentAddByOrganizationID(),
			}
		},
	)
}

func (m *OrgIndentAddByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	indent := cidl.NewGroupBuyingIndent()
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	tasks, err := dbGroupBuying.GetFinishBuyingTasks(organizationId, m.Ask.TaskIds)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get tasks failed. %s", err)
		return
	}

	if len(tasks) == 0 {
		ctx.Errorf(api.ErrDBQueryNoRecords, "no finish buying task.")
		return
	}

	indentId := utils.UniqueID()
	if err != nil {
		ctx.Errorf(api.ErrServer, "generate indent id failed. %s", err)
		return
	}

	var tasksBrief cidl.GroupBuyingIndentTaskBriefType
	for _, task := range tasks {
		briefItem := cidl.NewGroupBuyingIndentTasksBriefItem()
		briefItem.TaskId = task.TaskId
		briefItem.Title = task.Title
		briefItem.StartTime = task.StartTime
		briefItem.EndTime = task.EndTime
		tasksBrief = append(tasksBrief, briefItem)

		indentStatistics := cidl.NewGroupBuyingIndentStatistics()
		indentStatistics.IndentId = indentId
		indentStatistics.TaskId = task.TaskId
		indentStatistics.TaskContent = cidl.NewGroupBuyingOrderTaskContentByTask(task)
		indentStatistics.Version = cidl.IndentStatisticsRecordVersion

		_, err = dbGroupBuying.AddIndentStatistics(indentStatistics)
		if err != nil {
			ctx.Errorf(api.ErrDBInsertFailed, "add indent statistics failed. %s", err)
			return
		}

	}

	indent.IndentId = indentId
	indent.OrganizationId = organizationId
	indent.TasksBrief = &tasksBrief
	indent.State = cidl.IndentStateStatistic
	indent.Version = cidl.IndentRecordVersion

	_, err = dbGroupBuying.AddIndent(indent)
	if err != nil {
		ctx.Errorf(api.ErrDBInsertFailed, "add indent failed. %s", err)
		return
	}

	// 广播
	topic, err := mq.GetTopicServiceGroupBuyingOrderService()
	if err != nil {
		ctx.Errorf(api.ErrMqConnectFailed, "get topic service-group-buying-service failed. %s", err)
		return
	}

	err = topic.AddIndent(&mq.AddIndentMessage{
		IndentID: indentId,
	})

	if err != nil {
		ctx.Errorf(api.ErrMqPublishFailed, "publish topic service-group-buying-service failed. %s", err)
		return
	}

	m.Ack.IndentId = indentId
	ctx.Json(m.Ack)
}

// 绑定送货单团购任务
type OrgSendAddByOrganizationIDImpl struct {
	cidl.ApiOrgSendAddByOrganizationID
}

func AddOrgSendAddByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_SEND_ADD_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgSendAddByOrganizationIDImpl{
				ApiOrgSendAddByOrganizationID: cidl.MakeApiOrgSendAddByOrganizationID(),
			}
		},
	)
}

func (m *OrgSendAddByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)

	userId := ctx.Session.Uid
	send := cidl.NewGroupBuyingSend()
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	tasks, err := dbGroupBuying.GetFinishBuyingTasks(organizationId, m.Ask.TaskIds)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get finish ordering tasks failed. %s", err)
		return
	}

	if len(tasks) == 0 {
		ctx.Errorf(api.ErrDBQueryNoRecords, "no finish buying task.")
		return
	}

	// 判断是否有未绑定路线但是上报了销量的社团
	strSql := `
		SELECT COUNT(*)
		FROM gby_line_community
		WHERE
				lin_id = 0
				AND grp_id IN (
					SELECT DISTINCT grp_id
					FROM byo_community_buy
					WHERE tsk_id IN (?)
		)
	`
	strSql, args, err := conn.In(strSql, m.Ask.TaskIds)
	if err != nil {
		ctx.Errorf(api.ErrServer, "transform in array sql failed. %s", err)
		return
	}

	notBindLineCommunityCount := uint32(0)
	err = dbGroupBuying.DB.Get(&notBindLineCommunityCount, strSql, args...)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get unbind line community count failed. %s", err)
		return
	}

	if notBindLineCommunityCount > 0 {
		ctx.Errorf(cidl.ErrSendCommunityNotBindLine, "community was not bind with line")
		return
	}

	send.SendId = utils.UniqueID()
	if err != nil {
		ctx.Errorf(api.ErrServer, "generate send id failed. %s", err)
		return
	}

	var tasksBrief cidl.GroupBuyingSendTasksBriefType
	tasksContents := cidl.NewGroupBuyingSendTasksDetailType()

	for _, task := range tasks {
		briefItem := cidl.NewGroupBuyingSendTaskBriefItem()
		briefItem.TaskId = task.TaskId
		briefItem.Title = task.Title
		briefItem.StartTime = task.StartTime
		briefItem.EndTime = task.EndTime
		tasksBrief = append(tasksBrief, briefItem)
		taskContent := cidl.NewGroupBuyingOrderTaskContentByTask(task)
		*tasksContents = append(*tasksContents, taskContent)
	}

	send.TasksBrief = &tasksBrief
	send.TasksDetail, err = tasksContents.ToString()
	if err != nil {
		ctx.Errorf(api.ErrServer, "get string tasks contents failed. %s", err)
		return
	}

	organization, err := agency.NewProxy("agency-service").InnerAgencyOrganizationInfoByOrganizationID(organizationId)
	if err != nil {
		ctx.Errorf(api.ErrProxyFailed, "get organization from proxy failed. %s", err)
		return
	}

	send.OrganizationId = organizationId
	send.OrganizationName = organization.Name
	send.State = cidl.GroupBuyingSendStateStatistic

	_, err = dbGroupBuying.AddSend(send)
	if err != nil {
		ctx.Errorf(api.ErrDBInsertFailed, "add send failed. %s", err)
		return
	}

	// 广播
	topic, err := mq.GetTopicServiceGroupBuyingOrderService()
	if err != nil {
		ctx.Errorf(api.ErrMqConnectFailed, "get topic service-group-buying-service failed. %s", err)
		return
	}

	err = topic.AddSend(&mq.AddSendMessage{
		SendID:       send.SendId,
		AuthorUserId: userId,
	})
	if err != nil {
		ctx.Errorf(api.ErrMqPublishFailed, "publish topic service-group-buying-service failed. %s", err)
		return
	}

	m.Ack.SendId = send.SendId
	ctx.Json(m.Ack)
}
