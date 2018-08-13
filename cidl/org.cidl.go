package cidl

type AckOrgTaskList struct {
	Count uint32                  `db:"Count"`
	List  []*GroupBuyingOrderTask `db:"List"`
}

func NewAckOrgTaskList() *AckOrgTaskList {
	return &AckOrgTaskList{
		List: make([]*GroupBuyingOrderTask, 0),
	}
}

type MetaApiOrgTaskList struct {
}

var META_ORG_TASK_LIST = &MetaApiOrgTaskList{}

func (m *MetaApiOrgTaskList) GetMethod() string { return "GET" }
func (m *MetaApiOrgTaskList) GetURL() string {
	return "/group_buying_order/org/task/list/:organization_id"
}
func (m *MetaApiOrgTaskList) GetName() string { return "OrgTaskList" }
func (m *MetaApiOrgTaskList) GetType() string { return "json" }

// 团购任务默认列表
type ApiOrgTaskList struct {
	MetaApiOrgTaskList
	Ack    *AckOrgTaskList
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiOrgTaskList) GetQuery() interface{}  { return &m.Query }
func (m *ApiOrgTaskList) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskList) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskList) GetAck() interface{}    { return m.Ack }
func MakeApiOrgTaskList() ApiOrgTaskList {
	return ApiOrgTaskList{
		Ack: NewAckOrgTaskList(),
	}
}

type AckOrgTaskMonthListByYearByMonth struct {
	Count uint32                  `db:"Count"`
	List  []*GroupBuyingOrderTask `db:"List"`
}

func NewAckOrgTaskMonthListByYearByMonth() *AckOrgTaskMonthListByYearByMonth {
	return &AckOrgTaskMonthListByYearByMonth{
		List: make([]*GroupBuyingOrderTask, 0),
	}
}

type MetaApiOrgTaskMonthListByYearByMonth struct {
}

var META_ORG_TASK_MONTH_LIST_BY_YEAR_BY_MONTH = &MetaApiOrgTaskMonthListByYearByMonth{}

func (m *MetaApiOrgTaskMonthListByYearByMonth) GetMethod() string { return "GET" }
func (m *MetaApiOrgTaskMonthListByYearByMonth) GetURL() string {
	return "/group_buying_order/org/task/month_list/:organization_id"
}
func (m *MetaApiOrgTaskMonthListByYearByMonth) GetName() string {
	return "OrgTaskMonthListByYearByMonth"
}
func (m *MetaApiOrgTaskMonthListByYearByMonth) GetType() string { return "json" }

// 团购任务月份列表
type ApiOrgTaskMonthListByYearByMonth struct {
	MetaApiOrgTaskMonthListByYearByMonth
	Ack    *AckOrgTaskMonthListByYearByMonth
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Year     uint32 `form:"year" binding:"required,gt=0" db:"Year"`
		Month    uint32 `form:"month" binding:"required,gt=0,lte=12" db:"Month"`
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiOrgTaskMonthListByYearByMonth) GetQuery() interface{}  { return &m.Query }
func (m *ApiOrgTaskMonthListByYearByMonth) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskMonthListByYearByMonth) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskMonthListByYearByMonth) GetAck() interface{}    { return m.Ack }
func MakeApiOrgTaskMonthListByYearByMonth() ApiOrgTaskMonthListByYearByMonth {
	return ApiOrgTaskMonthListByYearByMonth{
		Ack: NewAckOrgTaskMonthListByYearByMonth(),
	}
}

type AckOrgTaskFinishBuyingListByOrganizationID struct {
	Count uint32                  `db:"Count"`
	List  []*GroupBuyingOrderTask `db:"List"`
}

func NewAckOrgTaskFinishBuyingListByOrganizationID() *AckOrgTaskFinishBuyingListByOrganizationID {
	return &AckOrgTaskFinishBuyingListByOrganizationID{
		List: make([]*GroupBuyingOrderTask, 0),
	}
}

type MetaApiOrgTaskFinishBuyingListByOrganizationID struct {
}

var META_ORG_TASK_FINISH_BUYING_LIST_BY_ORGANIZATION_ID = &MetaApiOrgTaskFinishBuyingListByOrganizationID{}

func (m *MetaApiOrgTaskFinishBuyingListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiOrgTaskFinishBuyingListByOrganizationID) GetURL() string {
	return "/group_buying_order/org/task/finish_buying_list/:organization_id"
}
func (m *MetaApiOrgTaskFinishBuyingListByOrganizationID) GetName() string {
	return "OrgTaskFinishBuyingListByOrganizationID"
}
func (m *MetaApiOrgTaskFinishBuyingListByOrganizationID) GetType() string { return "json" }

// 已结团团购任务列表
type ApiOrgTaskFinishBuyingListByOrganizationID struct {
	MetaApiOrgTaskFinishBuyingListByOrganizationID
	Ack    *AckOrgTaskFinishBuyingListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
		Search   string `form:"search" db:"Search"`
	}
}

func (m *ApiOrgTaskFinishBuyingListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiOrgTaskFinishBuyingListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskFinishBuyingListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskFinishBuyingListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgTaskFinishBuyingListByOrganizationID() ApiOrgTaskFinishBuyingListByOrganizationID {
	return ApiOrgTaskFinishBuyingListByOrganizationID{
		Ack: NewAckOrgTaskFinishBuyingListByOrganizationID(),
	}
}

type AckOrgIndentListByOrganizationID struct {
	Count uint32               `db:"Count"`
	List  []*GroupBuyingIndent `db:"List"`
}

func NewAckOrgIndentListByOrganizationID() *AckOrgIndentListByOrganizationID {
	return &AckOrgIndentListByOrganizationID{
		List: make([]*GroupBuyingIndent, 0),
	}
}

type MetaApiOrgIndentListByOrganizationID struct {
}

var META_ORG_INDENT_LIST_BY_ORGANIZATION_ID = &MetaApiOrgIndentListByOrganizationID{}

func (m *MetaApiOrgIndentListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiOrgIndentListByOrganizationID) GetURL() string {
	return "/group_buying_order/org/indent/list/:organization_id"
}
func (m *MetaApiOrgIndentListByOrganizationID) GetName() string {
	return "OrgIndentListByOrganizationID"
}
func (m *MetaApiOrgIndentListByOrganizationID) GetType() string { return "json" }

// 订货单列表
type ApiOrgIndentListByOrganizationID struct {
	MetaApiOrgIndentListByOrganizationID
	Ack    *AckOrgIndentListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiOrgIndentListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiOrgIndentListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgIndentListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiOrgIndentListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgIndentListByOrganizationID() ApiOrgIndentListByOrganizationID {
	return ApiOrgIndentListByOrganizationID{
		Ack: NewAckOrgIndentListByOrganizationID(),
	}
}

type AckOrgIndentSummaryListByIndentID struct {
	TaskStatistics []*GroupBuyingIndentStatisticResultItem `db:"TaskStatistics"`
}

func NewAckOrgIndentSummaryListByIndentID() *AckOrgIndentSummaryListByIndentID {
	return &AckOrgIndentSummaryListByIndentID{
		TaskStatistics: make([]*GroupBuyingIndentStatisticResultItem, 0),
	}
}

type MetaApiOrgIndentSummaryListByIndentID struct {
}

var META_ORG_INDENT_SUMMARY_LIST_BY_INDENT_ID = &MetaApiOrgIndentSummaryListByIndentID{}

func (m *MetaApiOrgIndentSummaryListByIndentID) GetMethod() string { return "GET" }
func (m *MetaApiOrgIndentSummaryListByIndentID) GetURL() string {
	return "/group_buying_order/org/indent/summary/:indent_id"
}
func (m *MetaApiOrgIndentSummaryListByIndentID) GetName() string {
	return "OrgIndentSummaryListByIndentID"
}
func (m *MetaApiOrgIndentSummaryListByIndentID) GetType() string { return "json" }

// 订货单概要
type ApiOrgIndentSummaryListByIndentID struct {
	MetaApiOrgIndentSummaryListByIndentID
	Ack    *AckOrgIndentSummaryListByIndentID
	Params struct {
		IndentID string `form:"indent_id" binding:"required" db:"IndentID"`
	}
}

func (m *ApiOrgIndentSummaryListByIndentID) GetQuery() interface{}  { return nil }
func (m *ApiOrgIndentSummaryListByIndentID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgIndentSummaryListByIndentID) GetAsk() interface{}    { return nil }
func (m *ApiOrgIndentSummaryListByIndentID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgIndentSummaryListByIndentID() ApiOrgIndentSummaryListByIndentID {
	return ApiOrgIndentSummaryListByIndentID{
		Ack: NewAckOrgIndentSummaryListByIndentID(),
	}
}

type MetaApiOrgIndentInvoicesByIndentID struct {
}

var META_ORG_INDENT_INVOICES_BY_INDENT_ID = &MetaApiOrgIndentInvoicesByIndentID{}

func (m *MetaApiOrgIndentInvoicesByIndentID) GetMethod() string { return "GET" }
func (m *MetaApiOrgIndentInvoicesByIndentID) GetURL() string {
	return "/group_buying_order/org/indent/invoices/:indent_id"
}
func (m *MetaApiOrgIndentInvoicesByIndentID) GetName() string { return "OrgIndentInvoicesByIndentID" }
func (m *MetaApiOrgIndentInvoicesByIndentID) GetType() string { return "json" }

// 导出订货单
type ApiOrgIndentInvoicesByIndentID struct {
	MetaApiOrgIndentInvoicesByIndentID
	Params struct {
		IndentID string `form:"indent_id" binding:"required" db:"IndentID"`
	}
}

func (m *ApiOrgIndentInvoicesByIndentID) GetQuery() interface{}  { return nil }
func (m *ApiOrgIndentInvoicesByIndentID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgIndentInvoicesByIndentID) GetAsk() interface{}    { return nil }
func (m *ApiOrgIndentInvoicesByIndentID) GetAck() interface{}    { return nil }
func MakeApiOrgIndentInvoicesByIndentID() ApiOrgIndentInvoicesByIndentID {
	return ApiOrgIndentInvoicesByIndentID{}
}

type AckOrgSendListByOrganizationID struct {
	Count uint32             `db:"Count"`
	List  []*GroupBuyingSend `db:"List"`
}

func NewAckOrgSendListByOrganizationID() *AckOrgSendListByOrganizationID {
	return &AckOrgSendListByOrganizationID{
		List: make([]*GroupBuyingSend, 0),
	}
}

type MetaApiOrgSendListByOrganizationID struct {
}

var META_ORG_SEND_LIST_BY_ORGANIZATION_ID = &MetaApiOrgSendListByOrganizationID{}

func (m *MetaApiOrgSendListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiOrgSendListByOrganizationID) GetURL() string {
	return "/group_buying_order/org/send/list/:organization_id"
}
func (m *MetaApiOrgSendListByOrganizationID) GetName() string { return "OrgSendListByOrganizationID" }
func (m *MetaApiOrgSendListByOrganizationID) GetType() string { return "json" }

// 配送单列表
type ApiOrgSendListByOrganizationID struct {
	MetaApiOrgSendListByOrganizationID
	Ack    *AckOrgSendListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiOrgSendListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiOrgSendListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgSendListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiOrgSendListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgSendListByOrganizationID() ApiOrgSendListByOrganizationID {
	return ApiOrgSendListByOrganizationID{
		Ack: NewAckOrgSendListByOrganizationID(),
	}
}

type MetaApiOrgSendInvoicesBySendID struct {
}

var META_ORG_SEND_INVOICES_BY_SEND_ID = &MetaApiOrgSendInvoicesBySendID{}

func (m *MetaApiOrgSendInvoicesBySendID) GetMethod() string { return "GET" }
func (m *MetaApiOrgSendInvoicesBySendID) GetURL() string {
	return "/group_buying_order/org/send/invoices/:send_id"
}
func (m *MetaApiOrgSendInvoicesBySendID) GetName() string { return "OrgSendInvoicesBySendID" }
func (m *MetaApiOrgSendInvoicesBySendID) GetType() string { return "json" }

// 导出配送单
type ApiOrgSendInvoicesBySendID struct {
	MetaApiOrgSendInvoicesBySendID
	Params struct {
		SendID string `form:"send_id" binding:"required" db:"SendID"`
	}
}

func (m *ApiOrgSendInvoicesBySendID) GetQuery() interface{}  { return nil }
func (m *ApiOrgSendInvoicesBySendID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgSendInvoicesBySendID) GetAsk() interface{}    { return nil }
func (m *ApiOrgSendInvoicesBySendID) GetAck() interface{}    { return nil }
func MakeApiOrgSendInvoicesBySendID() ApiOrgSendInvoicesBySendID {
	return ApiOrgSendInvoicesBySendID{}
}

type AckOrgLineListByOrganizationID struct {
	Count uint32             `db:"Count"`
	List  []*GroupBuyingLine `db:"List"`
}

func NewAckOrgLineListByOrganizationID() *AckOrgLineListByOrganizationID {
	return &AckOrgLineListByOrganizationID{
		List: make([]*GroupBuyingLine, 0),
	}
}

type MetaApiOrgLineListByOrganizationID struct {
}

var META_ORG_LINE_LIST_BY_ORGANIZATION_ID = &MetaApiOrgLineListByOrganizationID{}

func (m *MetaApiOrgLineListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiOrgLineListByOrganizationID) GetURL() string {
	return "/group_buying_order/org/line/list/:organization_id"
}
func (m *MetaApiOrgLineListByOrganizationID) GetName() string { return "OrgLineListByOrganizationID" }
func (m *MetaApiOrgLineListByOrganizationID) GetType() string { return "json" }

// 配送路线列表
type ApiOrgLineListByOrganizationID struct {
	MetaApiOrgLineListByOrganizationID
	Ack    *AckOrgLineListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiOrgLineListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiOrgLineListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgLineListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiOrgLineListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgLineListByOrganizationID() ApiOrgLineListByOrganizationID {
	return ApiOrgLineListByOrganizationID{
		Ack: NewAckOrgLineListByOrganizationID(),
	}
}

type AckOrgLineCommunityListByOrganizationID struct {
	Count uint32                      `db:"Count"`
	List  []*GroupBuyingLineCommunity `db:"List"`
}

func NewAckOrgLineCommunityListByOrganizationID() *AckOrgLineCommunityListByOrganizationID {
	return &AckOrgLineCommunityListByOrganizationID{
		List: make([]*GroupBuyingLineCommunity, 0),
	}
}

type MetaApiOrgLineCommunityListByOrganizationID struct {
}

var META_ORG_LINE_COMMUNITY_LIST_BY_ORGANIZATION_ID = &MetaApiOrgLineCommunityListByOrganizationID{}

func (m *MetaApiOrgLineCommunityListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiOrgLineCommunityListByOrganizationID) GetURL() string {
	return "/group_buying_order/org/line_community/list/:organization_id/:line_id"
}
func (m *MetaApiOrgLineCommunityListByOrganizationID) GetName() string {
	return "OrgLineCommunityListByOrganizationID"
}
func (m *MetaApiOrgLineCommunityListByOrganizationID) GetType() string { return "json" }

// 配送路线社群列表
type ApiOrgLineCommunityListByOrganizationID struct {
	MetaApiOrgLineCommunityListByOrganizationID
	Ack    *AckOrgLineCommunityListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		LineID         uint32 `form:"line_id" binding:"required,gt=0" db:"LineID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiOrgLineCommunityListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiOrgLineCommunityListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgLineCommunityListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiOrgLineCommunityListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgLineCommunityListByOrganizationID() ApiOrgLineCommunityListByOrganizationID {
	return ApiOrgLineCommunityListByOrganizationID{
		Ack: NewAckOrgLineCommunityListByOrganizationID(),
	}
}
