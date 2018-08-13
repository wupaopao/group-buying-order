package cidl

type AckAdminTaskList struct {
	Count uint32                  `db:"Count"`
	List  []*GroupBuyingOrderTask `db:"List"`
}

func NewAckAdminTaskList() *AckAdminTaskList {
	return &AckAdminTaskList{
		List: make([]*GroupBuyingOrderTask, 0),
	}
}

type MetaApiAdminTaskList struct {
}

var META_ADMIN_TASK_LIST = &MetaApiAdminTaskList{}

func (m *MetaApiAdminTaskList) GetMethod() string { return "GET" }
func (m *MetaApiAdminTaskList) GetURL() string {
	return "/group_buying_order/admin/task/list/:organization_id"
}
func (m *MetaApiAdminTaskList) GetName() string { return "AdminTaskList" }
func (m *MetaApiAdminTaskList) GetType() string { return "json" }

// 团购任务默认列表
type ApiAdminTaskList struct {
	MetaApiAdminTaskList
	Ack    *AckAdminTaskList
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiAdminTaskList) GetQuery() interface{}  { return &m.Query }
func (m *ApiAdminTaskList) GetParams() interface{} { return &m.Params }
func (m *ApiAdminTaskList) GetAsk() interface{}    { return nil }
func (m *ApiAdminTaskList) GetAck() interface{}    { return m.Ack }
func MakeApiAdminTaskList() ApiAdminTaskList {
	return ApiAdminTaskList{
		Ack: NewAckAdminTaskList(),
	}
}

type AckAdminTaskMonthListByYearByMonth struct {
	Count uint32                  `db:"Count"`
	List  []*GroupBuyingOrderTask `db:"List"`
}

func NewAckAdminTaskMonthListByYearByMonth() *AckAdminTaskMonthListByYearByMonth {
	return &AckAdminTaskMonthListByYearByMonth{
		List: make([]*GroupBuyingOrderTask, 0),
	}
}

type MetaApiAdminTaskMonthListByYearByMonth struct {
}

var META_ADMIN_TASK_MONTH_LIST_BY_YEAR_BY_MONTH = &MetaApiAdminTaskMonthListByYearByMonth{}

func (m *MetaApiAdminTaskMonthListByYearByMonth) GetMethod() string { return "GET" }
func (m *MetaApiAdminTaskMonthListByYearByMonth) GetURL() string {
	return "/group_buying_order/admin/task/month_list/:organization_id"
}
func (m *MetaApiAdminTaskMonthListByYearByMonth) GetName() string {
	return "AdminTaskMonthListByYearByMonth"
}
func (m *MetaApiAdminTaskMonthListByYearByMonth) GetType() string { return "json" }

// 团购任务月份列表
type ApiAdminTaskMonthListByYearByMonth struct {
	MetaApiAdminTaskMonthListByYearByMonth
	Ack    *AckAdminTaskMonthListByYearByMonth
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

func (m *ApiAdminTaskMonthListByYearByMonth) GetQuery() interface{}  { return &m.Query }
func (m *ApiAdminTaskMonthListByYearByMonth) GetParams() interface{} { return &m.Params }
func (m *ApiAdminTaskMonthListByYearByMonth) GetAsk() interface{}    { return nil }
func (m *ApiAdminTaskMonthListByYearByMonth) GetAck() interface{}    { return m.Ack }
func MakeApiAdminTaskMonthListByYearByMonth() ApiAdminTaskMonthListByYearByMonth {
	return ApiAdminTaskMonthListByYearByMonth{
		Ack: NewAckAdminTaskMonthListByYearByMonth(),
	}
}

type AckAdminTaskFinishBuyingListByOrganizationID struct {
	Count uint32                  `db:"Count"`
	List  []*GroupBuyingOrderTask `db:"List"`
}

func NewAckAdminTaskFinishBuyingListByOrganizationID() *AckAdminTaskFinishBuyingListByOrganizationID {
	return &AckAdminTaskFinishBuyingListByOrganizationID{
		List: make([]*GroupBuyingOrderTask, 0),
	}
}

type MetaApiAdminTaskFinishBuyingListByOrganizationID struct {
}

var META_ADMIN_TASK_FINISH_BUYING_LIST_BY_ORGANIZATION_ID = &MetaApiAdminTaskFinishBuyingListByOrganizationID{}

func (m *MetaApiAdminTaskFinishBuyingListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiAdminTaskFinishBuyingListByOrganizationID) GetURL() string {
	return "/group_buying_order/admin/task/finish_buying_list/:organization_id"
}
func (m *MetaApiAdminTaskFinishBuyingListByOrganizationID) GetName() string {
	return "AdminTaskFinishBuyingListByOrganizationID"
}
func (m *MetaApiAdminTaskFinishBuyingListByOrganizationID) GetType() string { return "json" }

// 已结团团购任务列表
type ApiAdminTaskFinishBuyingListByOrganizationID struct {
	MetaApiAdminTaskFinishBuyingListByOrganizationID
	Ack    *AckAdminTaskFinishBuyingListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
		Search   string `form:"search" db:"Search"`
	}
}

func (m *ApiAdminTaskFinishBuyingListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiAdminTaskFinishBuyingListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiAdminTaskFinishBuyingListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiAdminTaskFinishBuyingListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiAdminTaskFinishBuyingListByOrganizationID() ApiAdminTaskFinishBuyingListByOrganizationID {
	return ApiAdminTaskFinishBuyingListByOrganizationID{
		Ack: NewAckAdminTaskFinishBuyingListByOrganizationID(),
	}
}

type AckAdminIndentListByOrganizationID struct {
	Count uint32               `db:"Count"`
	List  []*GroupBuyingIndent `db:"List"`
}

func NewAckAdminIndentListByOrganizationID() *AckAdminIndentListByOrganizationID {
	return &AckAdminIndentListByOrganizationID{
		List: make([]*GroupBuyingIndent, 0),
	}
}

type MetaApiAdminIndentListByOrganizationID struct {
}

var META_ADMIN_INDENT_LIST_BY_ORGANIZATION_ID = &MetaApiAdminIndentListByOrganizationID{}

func (m *MetaApiAdminIndentListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiAdminIndentListByOrganizationID) GetURL() string {
	return "/group_buying_order/admin/indent/list/:organization_id"
}
func (m *MetaApiAdminIndentListByOrganizationID) GetName() string {
	return "AdminIndentListByOrganizationID"
}
func (m *MetaApiAdminIndentListByOrganizationID) GetType() string { return "json" }

// 订货单列表
type ApiAdminIndentListByOrganizationID struct {
	MetaApiAdminIndentListByOrganizationID
	Ack    *AckAdminIndentListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiAdminIndentListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiAdminIndentListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiAdminIndentListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiAdminIndentListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiAdminIndentListByOrganizationID() ApiAdminIndentListByOrganizationID {
	return ApiAdminIndentListByOrganizationID{
		Ack: NewAckAdminIndentListByOrganizationID(),
	}
}

type AckAdminIndentSummaryListByIndentID struct {
	TaskStatistics []*GroupBuyingIndentStatisticResultItem `db:"TaskStatistics"`
}

func NewAckAdminIndentSummaryListByIndentID() *AckAdminIndentSummaryListByIndentID {
	return &AckAdminIndentSummaryListByIndentID{
		TaskStatistics: make([]*GroupBuyingIndentStatisticResultItem, 0),
	}
}

type MetaApiAdminIndentSummaryListByIndentID struct {
}

var META_ADMIN_INDENT_SUMMARY_LIST_BY_INDENT_ID = &MetaApiAdminIndentSummaryListByIndentID{}

func (m *MetaApiAdminIndentSummaryListByIndentID) GetMethod() string { return "GET" }
func (m *MetaApiAdminIndentSummaryListByIndentID) GetURL() string {
	return "/group_buying_order/admin/indent/summary/:indent_id"
}
func (m *MetaApiAdminIndentSummaryListByIndentID) GetName() string {
	return "AdminIndentSummaryListByIndentID"
}
func (m *MetaApiAdminIndentSummaryListByIndentID) GetType() string { return "json" }

// 订货单概要
type ApiAdminIndentSummaryListByIndentID struct {
	MetaApiAdminIndentSummaryListByIndentID
	Ack    *AckAdminIndentSummaryListByIndentID
	Params struct {
		IndentID string `form:"indent_id" binding:"required" db:"IndentID"`
	}
}

func (m *ApiAdminIndentSummaryListByIndentID) GetQuery() interface{}  { return nil }
func (m *ApiAdminIndentSummaryListByIndentID) GetParams() interface{} { return &m.Params }
func (m *ApiAdminIndentSummaryListByIndentID) GetAsk() interface{}    { return nil }
func (m *ApiAdminIndentSummaryListByIndentID) GetAck() interface{}    { return m.Ack }
func MakeApiAdminIndentSummaryListByIndentID() ApiAdminIndentSummaryListByIndentID {
	return ApiAdminIndentSummaryListByIndentID{
		Ack: NewAckAdminIndentSummaryListByIndentID(),
	}
}

type MetaApiAdminIndentInvoicesByIndentID struct {
}

var META_ADMIN_INDENT_INVOICES_BY_INDENT_ID = &MetaApiAdminIndentInvoicesByIndentID{}

func (m *MetaApiAdminIndentInvoicesByIndentID) GetMethod() string { return "GET" }
func (m *MetaApiAdminIndentInvoicesByIndentID) GetURL() string {
	return "/group_buying_order/admin/indent/invoices/:indent_id"
}
func (m *MetaApiAdminIndentInvoicesByIndentID) GetName() string {
	return "AdminIndentInvoicesByIndentID"
}
func (m *MetaApiAdminIndentInvoicesByIndentID) GetType() string { return "json" }

// 导出订货单
type ApiAdminIndentInvoicesByIndentID struct {
	MetaApiAdminIndentInvoicesByIndentID
	Params struct {
		IndentID string `form:"indent_id" binding:"required" db:"IndentID"`
	}
}

func (m *ApiAdminIndentInvoicesByIndentID) GetQuery() interface{}  { return nil }
func (m *ApiAdminIndentInvoicesByIndentID) GetParams() interface{} { return &m.Params }
func (m *ApiAdminIndentInvoicesByIndentID) GetAsk() interface{}    { return nil }
func (m *ApiAdminIndentInvoicesByIndentID) GetAck() interface{}    { return nil }
func MakeApiAdminIndentInvoicesByIndentID() ApiAdminIndentInvoicesByIndentID {
	return ApiAdminIndentInvoicesByIndentID{}
}

type AckAdminSendListByOrganizationID struct {
	Count uint32             `db:"Count"`
	List  []*GroupBuyingSend `db:"List"`
}

func NewAckAdminSendListByOrganizationID() *AckAdminSendListByOrganizationID {
	return &AckAdminSendListByOrganizationID{
		List: make([]*GroupBuyingSend, 0),
	}
}

type MetaApiAdminSendListByOrganizationID struct {
}

var META_ADMIN_SEND_LIST_BY_ORGANIZATION_ID = &MetaApiAdminSendListByOrganizationID{}

func (m *MetaApiAdminSendListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiAdminSendListByOrganizationID) GetURL() string {
	return "/group_buying_order/admin/send/list/:organization_id"
}
func (m *MetaApiAdminSendListByOrganizationID) GetName() string {
	return "AdminSendListByOrganizationID"
}
func (m *MetaApiAdminSendListByOrganizationID) GetType() string { return "json" }

// 配送单列表
type ApiAdminSendListByOrganizationID struct {
	MetaApiAdminSendListByOrganizationID
	Ack    *AckAdminSendListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiAdminSendListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiAdminSendListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiAdminSendListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiAdminSendListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiAdminSendListByOrganizationID() ApiAdminSendListByOrganizationID {
	return ApiAdminSendListByOrganizationID{
		Ack: NewAckAdminSendListByOrganizationID(),
	}
}

type MetaApiAdminSendInvoicesBySendID struct {
}

var META_ADMIN_SEND_INVOICES_BY_SEND_ID = &MetaApiAdminSendInvoicesBySendID{}

func (m *MetaApiAdminSendInvoicesBySendID) GetMethod() string { return "GET" }
func (m *MetaApiAdminSendInvoicesBySendID) GetURL() string {
	return "/group_buying_order/admin/send/invoices/:send_id"
}
func (m *MetaApiAdminSendInvoicesBySendID) GetName() string { return "AdminSendInvoicesBySendID" }
func (m *MetaApiAdminSendInvoicesBySendID) GetType() string { return "json" }

// 导出配送单
type ApiAdminSendInvoicesBySendID struct {
	MetaApiAdminSendInvoicesBySendID
	Params struct {
		SendID string `form:"send_id" binding:"required" db:"SendID"`
	}
}

func (m *ApiAdminSendInvoicesBySendID) GetQuery() interface{}  { return nil }
func (m *ApiAdminSendInvoicesBySendID) GetParams() interface{} { return &m.Params }
func (m *ApiAdminSendInvoicesBySendID) GetAsk() interface{}    { return nil }
func (m *ApiAdminSendInvoicesBySendID) GetAck() interface{}    { return nil }
func MakeApiAdminSendInvoicesBySendID() ApiAdminSendInvoicesBySendID {
	return ApiAdminSendInvoicesBySendID{}
}

type AckAdminLineListByOrganizationID struct {
	Count uint32             `db:"Count"`
	List  []*GroupBuyingLine `db:"List"`
}

func NewAckAdminLineListByOrganizationID() *AckAdminLineListByOrganizationID {
	return &AckAdminLineListByOrganizationID{
		List: make([]*GroupBuyingLine, 0),
	}
}

type MetaApiAdminLineListByOrganizationID struct {
}

var META_ADMIN_LINE_LIST_BY_ORGANIZATION_ID = &MetaApiAdminLineListByOrganizationID{}

func (m *MetaApiAdminLineListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiAdminLineListByOrganizationID) GetURL() string {
	return "/group_buying_order/admin/line/list/:organization_id"
}
func (m *MetaApiAdminLineListByOrganizationID) GetName() string {
	return "AdminLineListByOrganizationID"
}
func (m *MetaApiAdminLineListByOrganizationID) GetType() string { return "json" }

// 配送路线列表
type ApiAdminLineListByOrganizationID struct {
	MetaApiAdminLineListByOrganizationID
	Ack    *AckAdminLineListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiAdminLineListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiAdminLineListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiAdminLineListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiAdminLineListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiAdminLineListByOrganizationID() ApiAdminLineListByOrganizationID {
	return ApiAdminLineListByOrganizationID{
		Ack: NewAckAdminLineListByOrganizationID(),
	}
}

type AckAdminLineCommunityListByOrganizationID struct {
	Count uint32                      `db:"Count"`
	List  []*GroupBuyingLineCommunity `db:"List"`
}

func NewAckAdminLineCommunityListByOrganizationID() *AckAdminLineCommunityListByOrganizationID {
	return &AckAdminLineCommunityListByOrganizationID{
		List: make([]*GroupBuyingLineCommunity, 0),
	}
}

type MetaApiAdminLineCommunityListByOrganizationID struct {
}

var META_ADMIN_LINE_COMMUNITY_LIST_BY_ORGANIZATION_ID = &MetaApiAdminLineCommunityListByOrganizationID{}

func (m *MetaApiAdminLineCommunityListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiAdminLineCommunityListByOrganizationID) GetURL() string {
	return "/group_buying_order/admin/line_community/list/:organization_id/:line_id"
}
func (m *MetaApiAdminLineCommunityListByOrganizationID) GetName() string {
	return "AdminLineCommunityListByOrganizationID"
}
func (m *MetaApiAdminLineCommunityListByOrganizationID) GetType() string { return "json" }

// 配送路线社群列表
type ApiAdminLineCommunityListByOrganizationID struct {
	MetaApiAdminLineCommunityListByOrganizationID
	Ack    *AckAdminLineCommunityListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		LineID         uint32 `form:"line_id" binding:"required,gt=0" db:"LineID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiAdminLineCommunityListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiAdminLineCommunityListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiAdminLineCommunityListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiAdminLineCommunityListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiAdminLineCommunityListByOrganizationID() ApiAdminLineCommunityListByOrganizationID {
	return ApiAdminLineCommunityListByOrganizationID{
		Ack: NewAckAdminLineCommunityListByOrganizationID(),
	}
}
