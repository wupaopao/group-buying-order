package impls

import (
	"fmt"

	"business/group-buying-order/cidl"
	"business/group-buying-order/common/db"
	"common/api"
	"common/file"

	"github.com/mz-eco/mz/http"
)

func init() {
	AddOrgTaskListHandler()
	AddOrgTaskMonthListByYearByMonthHandler()
	AddOrgTaskFinishBuyingListByOrganizationIDHandler()

	AddOrgIndentListByOrganizationIDHandler()
	AddOrgIndentSummaryListByIndentIDHandler()
	AddOrgIndentInvoicesByIndentIDHandler()

	AddOrgSendListByOrganizationIDHandler()
	AddOrgSendInvoicesBySendIDHandler()

	AddOrgLineListByOrganizationIDHandler()
	AddOrgLineCommunityListByOrganizationIDHandler()

}

// 团购任务默认列表
type OrgTaskListImpl struct {
	cidl.ApiOrgTaskList
}

func AddOrgTaskListHandler() {
	AddHandler(
		cidl.META_ORG_TASK_LIST,
		func() http.ApiHandler {
			return &OrgTaskListImpl{
				ApiOrgTaskList: cidl.MakeApiOrgTaskList(),
			}
		},
	)
}

func (m *OrgTaskListImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	ack.Count, err = dbGroupBuying.TaskNotReachEndTimeCount(organizationId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task count failed. %s", err)
		return
	}

	if ack.Count == 0 {
		ctx.Json(ack)
		return
	}

	ack.List, err = dbGroupBuying.TaskNotReachEndTimeList(organizationId, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task list failed. %s", err)
		return
	}

	ctx.Json(ack)
}

// 团购任务月份列表
type OrgTaskMonthListByYearByMonthImpl struct {
	cidl.ApiOrgTaskMonthListByYearByMonth
}

func AddOrgTaskMonthListByYearByMonthHandler() {
	AddHandler(
		cidl.META_ORG_TASK_MONTH_LIST_BY_YEAR_BY_MONTH,
		func() http.ApiHandler {
			return &OrgTaskMonthListByYearByMonthImpl{
				ApiOrgTaskMonthListByYearByMonth: cidl.MakeApiOrgTaskMonthListByYearByMonth(),
			}
		},
	)
}

func (m *OrgTaskMonthListByYearByMonthImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	ack.Count, err = dbGroupBuying.TaskMonthCount(organizationId, m.Query.Year, m.Query.Month)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task count failed. %s", err)
		return
	}

	fmt.Println(ack.Count)
	if ack.Count == 0 {
		ctx.Json(ack)
		return
	}

	ack.List, err = dbGroupBuying.TaskMonthList(organizationId, m.Query.Year, m.Query.Month, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task list failed. %s", err)
		return
	}

	ctx.Json(ack)
}

// 已结团团购任务列表
type OrgTaskFinishBuyingListByOrganizationIDImpl struct {
	cidl.ApiOrgTaskFinishBuyingListByOrganizationID
}

func AddOrgTaskFinishBuyingListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_TASK_FINISH_BUYING_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgTaskFinishBuyingListByOrganizationIDImpl{
				ApiOrgTaskFinishBuyingListByOrganizationID: cidl.MakeApiOrgTaskFinishBuyingListByOrganizationID(),
			}
		},
	)
}

func (m *OrgTaskFinishBuyingListByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	search := m.Query.Search
	if search == "" {
		ack.Count, err = dbGroupBuying.TaskFinishBuyingCount(organizationId)
	} else {
		ack.Count, err = dbGroupBuying.TaskFinishBuyingSearchCount(organizationId, search)
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
		ack.List, err = dbGroupBuying.TaskFinishBuyingList(organizationId, m.Query.Page, m.Query.PageSize, false)
	} else {
		ack.List, err = dbGroupBuying.TaskFinishBuyingSearchList(organizationId, search, m.Query.Page, m.Query.PageSize, false)
	}

	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get task list failed. %s", err)
		return
	}

	ctx.Json(ack)
}

// 订货单列表
type OrgIndentListByOrganizationIDImpl struct {
	cidl.ApiOrgIndentListByOrganizationID
}

func AddOrgIndentListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_INDENT_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgIndentListByOrganizationIDImpl{
				ApiOrgIndentListByOrganizationID: cidl.MakeApiOrgIndentListByOrganizationID(),
			}
		},
	)
}

func (m *OrgIndentListByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	organizationId := m.Params.OrganizationID
	ack.Count, err = dbGroupBuying.IndentCount(organizationId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get indent count failed. %s", err)
		return
	}

	if ack.Count == 0 {
		ctx.Json(ack)
		return
	}

	ack.List, err = dbGroupBuying.IndentList(organizationId, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get indent list failed. %s", err)
		return
	}

	ctx.Json(ack)
}

// 配送单列表
type OrgSendListByOrganizationIDImpl struct {
	cidl.ApiOrgSendListByOrganizationID
}

func AddOrgSendListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_SEND_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgSendListByOrganizationIDImpl{
				ApiOrgSendListByOrganizationID: cidl.MakeApiOrgSendListByOrganizationID(),
			}
		},
	)
}

func (m *OrgSendListByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	ack.Count, err = dbGroupBuying.SendCount(organizationId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get send count failed. %s", err)
		return
	}

	if ack.Count == 0 {
		ctx.Json(m.Ack)
		return
	}

	ack.List, err = dbGroupBuying.SendList(organizationId, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get send list failed. %s", err)
		return
	}

	ctx.Json(m.Ack)
}

// 导出配送单
type OrgSendInvoicesBySendIDImpl struct {
	cidl.ApiOrgSendInvoicesBySendID
}

func AddOrgSendInvoicesBySendIDHandler() {
	AddHandler(
		cidl.META_ORG_SEND_INVOICES_BY_SEND_ID,
		func() http.ApiHandler {
			return &OrgSendInvoicesBySendIDImpl{
				ApiOrgSendInvoicesBySendID: cidl.MakeApiOrgSendInvoicesBySendID(),
			}
		},
	)
}

func (m *OrgSendInvoicesBySendIDImpl) Handler(ctx *http.Context) {
	var (
		err error
		uri string
	)
	sendId := m.Params.SendID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	uri, err = dbGroupBuying.GetSendExcelUrl(sendId)
	if err != nil || uri == "" {
		ctx.Errorf(api.ErrDbQueryFailed, "get send excel url failed. %s", err)
		return
	}

	qiniu, err := file.GetQiniuPrivateBucket()
	if err != nil {
		ctx.Errorf(api.ErrServer, "get qiniu private bucket failed. %s", err)
		return
	}

	uri = qiniu.PrivateAccessUrl(uri, 3600)

	ctx.Redirect(uri)
}

// 订货单概要
type OrgIndentSummaryListByIndentIDImpl struct {
	cidl.ApiOrgIndentSummaryListByIndentID
}

func AddOrgIndentSummaryListByIndentIDHandler() {
	AddHandler(
		cidl.META_ORG_INDENT_SUMMARY_LIST_BY_INDENT_ID,
		func() http.ApiHandler {
			return &OrgIndentSummaryListByIndentIDImpl{
				ApiOrgIndentSummaryListByIndentID: cidl.MakeApiOrgIndentSummaryListByIndentID(),
			}
		},
	)
}

func (m *OrgIndentSummaryListByIndentIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	list, err := dbGroupBuying.IndentStatisticsAll(m.Params.IndentID, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "get indent statistics failed. %s", err)
		return
	}

	for _, item := range list {
		for _, skuItem := range *item.Result {
			m.Ack.TaskStatistics = append(m.Ack.TaskStatistics, skuItem)
		}
	}

	ctx.Json(m.Ack)
}

// 导出订货单
type OrgIndentInvoicesByIndentIDImpl struct {
	cidl.ApiOrgIndentInvoicesByIndentID
}

func AddOrgIndentInvoicesByIndentIDHandler() {
	AddHandler(
		cidl.META_ORG_INDENT_INVOICES_BY_INDENT_ID,
		func() http.ApiHandler {
			return &OrgIndentInvoicesByIndentIDImpl{
				ApiOrgIndentInvoicesByIndentID: cidl.MakeApiOrgIndentInvoicesByIndentID(),
			}
		},
	)
}

func (m *OrgIndentInvoicesByIndentIDImpl) Handler(ctx *http.Context) {
	var uri string
	var err error
	indentId := m.Params.IndentID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	uri, err = dbGroupBuying.GetIndentExcelUrl(indentId)
	if err != nil || uri == "" {
		ctx.Errorf(api.ErrDbQueryFailed, "get indent excel file failed. %s", err)
		return
	}

	qiniu, err := file.GetQiniuPrivateBucket()
	if err != nil {
		ctx.Errorf(api.ErrServer, "get qiniu private bucket failed. %s", err)
		return
	}

	uri = qiniu.PrivateAccessUrl(uri, 3600)

	ctx.Redirect(uri)
}

// 配送路线列表
type OrgLineListByOrganizationIDImpl struct {
	cidl.ApiOrgLineListByOrganizationID
}

func AddOrgLineListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_LINE_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgLineListByOrganizationIDImpl{
				ApiOrgLineListByOrganizationID: cidl.MakeApiOrgLineListByOrganizationID(),
			}
		},
	)
}

func (m *OrgLineListByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	ack.Count, err = dbGroupBuying.LineCount(organizationId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "query line count failed. %s", err)
		return
	}

	if ack.Count == 0 {
		ctx.Json(ack)
		return
	}

	ack.List, err = dbGroupBuying.LineList(organizationId, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "query line list failed. %s", err)
		return
	}

	ctx.Json(ack)
}

type OrgLineCommunityListByOrganizationIDImpl struct {
	cidl.ApiOrgLineCommunityListByOrganizationID
}

func AddOrgLineCommunityListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ORG_LINE_COMMUNITY_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &OrgLineCommunityListByOrganizationIDImpl{
				ApiOrgLineCommunityListByOrganizationID: cidl.MakeApiOrgLineCommunityListByOrganizationID(),
			}
		},
	)
}

func (m *OrgLineCommunityListByOrganizationIDImpl) Handler(ctx *http.Context) {
	var (
		err error
	)
	ack := m.Ack
	organizationId := m.Params.OrganizationID
	lineId := m.Params.LineID
	dbGroupBuying := db.NewMallGroupBuyingOrder()
	ack.Count, err = dbGroupBuying.LineCommunityCount(organizationId, lineId)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "query line community count failed. %s", err)
		return
	}

	if ack.Count == 0 {
		ctx.Json(ack)
		return
	}

	ack.List, err = dbGroupBuying.LineCommunityList(organizationId, lineId, m.Query.Page, m.Query.PageSize, false)
	if err != nil {
		ctx.Errorf(api.ErrDbQueryFailed, "query line community list failed. %s", err)
		return
	}

	ctx.Json(ack)
}
