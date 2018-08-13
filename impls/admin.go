package impls

import (
	"business/group-buying-order/cidl"
	"business/group-buying-order/common/db"
	"common/api"
	"common/file"

	"github.com/mz-eco/mz/http"
)

func init() {
	AddAdminTaskListHandler()
	AddAdminTaskMonthListByYearByMonthHandler()
	AddAdminTaskFinishBuyingListByOrganizationIDHandler()

	AddAdminIndentListByOrganizationIDHandler()
	AddAdminIndentSummaryListByIndentIDHandler()
	AddAdminIndentInvoicesByIndentIDHandler()

	AddAdminSendListByOrganizationIDHandler()
	AddAdminSendInvoicesBySendIDHandler()

	AddAdminLineListByOrganizationIDHandler()
	AddAdminLineCommunityListByOrganizationIDHandler()

}

// 团购任务默认列表
type AdminTaskListImpl struct {
	cidl.ApiAdminTaskList
}

func AddAdminTaskListHandler() {
	AddHandler(
		cidl.META_ADMIN_TASK_LIST,
		func() http.ApiHandler {
			return &AdminTaskListImpl{
				ApiAdminTaskList: cidl.MakeApiAdminTaskList(),
			}
		},
	)
}

func (m *AdminTaskListImpl) Handler(ctx *http.Context) {
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
type AdminTaskMonthListByYearByMonthImpl struct {
	cidl.ApiAdminTaskMonthListByYearByMonth
}

func AddAdminTaskMonthListByYearByMonthHandler() {
	AddHandler(
		cidl.META_ADMIN_TASK_MONTH_LIST_BY_YEAR_BY_MONTH,
		func() http.ApiHandler {
			return &AdminTaskMonthListByYearByMonthImpl{
				ApiAdminTaskMonthListByYearByMonth: cidl.MakeApiAdminTaskMonthListByYearByMonth(),
			}
		},
	)
}

func (m *AdminTaskMonthListByYearByMonthImpl) Handler(ctx *http.Context) {
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
type AdminTaskFinishBuyingListByOrganizationIDImpl struct {
	cidl.ApiAdminTaskFinishBuyingListByOrganizationID
}

func AddAdminTaskFinishBuyingListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ADMIN_TASK_FINISH_BUYING_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &AdminTaskFinishBuyingListByOrganizationIDImpl{
				ApiAdminTaskFinishBuyingListByOrganizationID: cidl.MakeApiAdminTaskFinishBuyingListByOrganizationID(),
			}
		},
	)
}

func (m *AdminTaskFinishBuyingListByOrganizationIDImpl) Handler(ctx *http.Context) {
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
type AdminIndentListByOrganizationIDImpl struct {
	cidl.ApiAdminIndentListByOrganizationID
}

func AddAdminIndentListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ADMIN_INDENT_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &AdminIndentListByOrganizationIDImpl{
				ApiAdminIndentListByOrganizationID: cidl.MakeApiAdminIndentListByOrganizationID(),
			}
		},
	)
}

func (m *AdminIndentListByOrganizationIDImpl) Handler(ctx *http.Context) {
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
type AdminSendListByOrganizationIDImpl struct {
	cidl.ApiAdminSendListByOrganizationID
}

func AddAdminSendListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ADMIN_SEND_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &AdminSendListByOrganizationIDImpl{
				ApiAdminSendListByOrganizationID: cidl.MakeApiAdminSendListByOrganizationID(),
			}
		},
	)
}

func (m *AdminSendListByOrganizationIDImpl) Handler(ctx *http.Context) {
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
type AdminSendInvoicesBySendIDImpl struct {
	cidl.ApiAdminSendInvoicesBySendID
}

func AddAdminSendInvoicesBySendIDHandler() {
	AddHandler(
		cidl.META_ADMIN_SEND_INVOICES_BY_SEND_ID,
		func() http.ApiHandler {
			return &AdminSendInvoicesBySendIDImpl{
				ApiAdminSendInvoicesBySendID: cidl.MakeApiAdminSendInvoicesBySendID(),
			}
		},
	)
}

func (m *AdminSendInvoicesBySendIDImpl) Handler(ctx *http.Context) {
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
type AdminIndentSummaryListByIndentIDImpl struct {
	cidl.ApiAdminIndentSummaryListByIndentID
}

func AddAdminIndentSummaryListByIndentIDHandler() {
	AddHandler(
		cidl.META_ADMIN_INDENT_SUMMARY_LIST_BY_INDENT_ID,
		func() http.ApiHandler {
			return &AdminIndentSummaryListByIndentIDImpl{
				ApiAdminIndentSummaryListByIndentID: cidl.MakeApiAdminIndentSummaryListByIndentID(),
			}
		},
	)
}

func (m *AdminIndentSummaryListByIndentIDImpl) Handler(ctx *http.Context) {
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
type AdminIndentInvoicesByIndentIDImpl struct {
	cidl.ApiAdminIndentInvoicesByIndentID
}

func AddAdminIndentInvoicesByIndentIDHandler() {
	AddHandler(
		cidl.META_ADMIN_INDENT_INVOICES_BY_INDENT_ID,
		func() http.ApiHandler {
			return &AdminIndentInvoicesByIndentIDImpl{
				ApiAdminIndentInvoicesByIndentID: cidl.MakeApiAdminIndentInvoicesByIndentID(),
			}
		},
	)
}

func (m *AdminIndentInvoicesByIndentIDImpl) Handler(ctx *http.Context) {
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
type AdminLineListByOrganizationIDImpl struct {
	cidl.ApiAdminLineListByOrganizationID
}

func AddAdminLineListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ADMIN_LINE_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &AdminLineListByOrganizationIDImpl{
				ApiAdminLineListByOrganizationID: cidl.MakeApiAdminLineListByOrganizationID(),
			}
		},
	)
}

func (m *AdminLineListByOrganizationIDImpl) Handler(ctx *http.Context) {
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

type AdminLineCommunityListByOrganizationIDImpl struct {
	cidl.ApiAdminLineCommunityListByOrganizationID
}

func AddAdminLineCommunityListByOrganizationIDHandler() {
	AddHandler(
		cidl.META_ADMIN_LINE_COMMUNITY_LIST_BY_ORGANIZATION_ID,
		func() http.ApiHandler {
			return &AdminLineCommunityListByOrganizationIDImpl{
				ApiAdminLineCommunityListByOrganizationID: cidl.MakeApiAdminLineCommunityListByOrganizationID(),
			}
		},
	)
}

func (m *AdminLineCommunityListByOrganizationIDImpl) Handler(ctx *http.Context) {
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
