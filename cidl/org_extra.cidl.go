package cidl

import "time"

// ####################团购任务
type AskTaskSpecificationLabelItem struct {
	LabelId string `binding:"required" db:"LabelId"`
	Name    string `binding:"required" db:"Name"`
}

func NewAskTaskSpecificationLabelItem() *AskTaskSpecificationLabelItem {
	return &AskTaskSpecificationLabelItem{}
}

type AskTaskSpecificationItem struct {
	Name   string                           `binding:"required" db:"Name"`
	Labels []*AskTaskSpecificationLabelItem `binding:"required,gt=0,dive,required" db:"Labels"`
}

func NewAskTaskSpecificationItem() *AskTaskSpecificationItem {
	return &AskTaskSpecificationItem{
		Labels: make([]*AskTaskSpecificationLabelItem, 0),
	}
}

type AskTaskSkuItem struct {
	SkuId               string   `binding:"required" db:"SkuId"`
	LabelIds            []string `binding:"required,gt=0" db:"LabelIds"`
	MarketPrice         float64  `binding:"required,gte=0" db:"MarketPrice"`
	GroupBuyingPrice    float64  `binding:"required,gte=0" db:"GroupBuyingPrice"`
	SettlementPrice     float64  `binding:"required,gte=0" db:"SettlementPrice"`
	CostPrice           float64  `binding:"required,gte=0" db:"CostPrice"`
	IllustrationPicture string   `db:"IllustrationPicture"`
	InventoryCount      uint32   `binding:"required,gt=0" db:"InventoryCount"`
	IsShow              bool     `db:"IsShow"`
}

func NewAskTaskSkuItem() *AskTaskSkuItem {
	return &AskTaskSkuItem{
		LabelIds: make([]string, 0),
	}
}

// 组合Item Sku子项
type AskTaskCombinationItemSubSkuItem struct {
	SkuId string `binding:"required" db:"SkuId"`
	Count uint32 `binding:"required" db:"Count"`
}

func NewAskTaskCombinationItemSubSkuItem() *AskTaskCombinationItemSubSkuItem {
	return &AskTaskCombinationItemSubSkuItem{}
}

// 组合Item
type AskTaskCombinationItem struct {
	Name                string                              `binding:"required" db:"Name"`
	SubSkuItems         []*AskTaskCombinationItemSubSkuItem `binding:"required,gt=0,dive,required" db:"SubSkuItems"`
	MarketPrice         float64                             `binding:"required,gte=0" db:"MarketPrice"`
	GroupBuyingPrice    float64                             `binding:"required,gte=0" db:"GroupBuyingPrice"`
	SettlementPrice     float64                             `binding:"required,gte=0" db:"SettlementPrice"`
	CostPrice           float64                             `binding:"required,gte=0" db:"CostPrice"`
	IllustrationPicture string                              `db:"IllustrationPicture"`
	IsShow              bool                                `db:"IsShow"`
}

func NewAskTaskCombinationItem() *AskTaskCombinationItem {
	return &AskTaskCombinationItem{
		SubSkuItems: make([]*AskTaskCombinationItemSubSkuItem, 0),
	}
}

// 团购任务图片上传TOKEN获取
type AckPicToken struct {
	OriginalFileName string `db:"OriginalFileName"`
	Token            string `db:"Token"`
	Key              string `db:"Key"`
	StoreUrl         string `db:"StoreUrl"`
	AccessUrl        string `db:"AccessUrl"`
}

func NewAckPicToken() *AckPicToken {
	return &AckPicToken{}
}

type AskOrgTaskPicTokenByOrganizationID struct {
	FileNames []string `db:"FileNames"`
}

func NewAskOrgTaskPicTokenByOrganizationID() *AskOrgTaskPicTokenByOrganizationID {
	return &AskOrgTaskPicTokenByOrganizationID{
		FileNames: make([]string, 0),
	}
}

type AckOrgTaskPicTokenByOrganizationID struct {
	Tokens []*AckPicToken `db:"Tokens"`
}

func NewAckOrgTaskPicTokenByOrganizationID() *AckOrgTaskPicTokenByOrganizationID {
	return &AckOrgTaskPicTokenByOrganizationID{
		Tokens: make([]*AckPicToken, 0),
	}
}

type MetaApiOrgTaskPicTokenByOrganizationID struct {
}

var META_ORG_TASK_PIC_TOKEN_BY_ORGANIZATION_ID = &MetaApiOrgTaskPicTokenByOrganizationID{}

func (m *MetaApiOrgTaskPicTokenByOrganizationID) GetMethod() string { return "POST" }
func (m *MetaApiOrgTaskPicTokenByOrganizationID) GetURL() string {
	return "/group_buying_order/org/task/pic_token/:organization_id"
}
func (m *MetaApiOrgTaskPicTokenByOrganizationID) GetName() string {
	return "OrgTaskPicTokenByOrganizationID"
}
func (m *MetaApiOrgTaskPicTokenByOrganizationID) GetType() string { return "json" }

// 图片token
type ApiOrgTaskPicTokenByOrganizationID struct {
	MetaApiOrgTaskPicTokenByOrganizationID
	Ask    *AskOrgTaskPicTokenByOrganizationID
	Ack    *AckOrgTaskPicTokenByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
}

func (m *ApiOrgTaskPicTokenByOrganizationID) GetQuery() interface{}  { return nil }
func (m *ApiOrgTaskPicTokenByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskPicTokenByOrganizationID) GetAsk() interface{}    { return m.Ask }
func (m *ApiOrgTaskPicTokenByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgTaskPicTokenByOrganizationID() ApiOrgTaskPicTokenByOrganizationID {
	return ApiOrgTaskPicTokenByOrganizationID{
		Ask: NewAskOrgTaskPicTokenByOrganizationID(),
		Ack: NewAckOrgTaskPicTokenByOrganizationID(),
	}
}

type AskOrgTaskAddByOrganizationID struct {
	ShowStartTime        time.Time                       `binding:"required" db:"ShowStartTime"`
	StartTime            time.Time                       `db:"StartTime"`
	EndTime              time.Time                       `binding:"required" db:"EndTime"`
	SellType             GroupBuyingOrderTaskSellType    `binding:"required" db:"SellType"`
	Notes                string                          `binding:"required,lte=1000" db:"Notes"`
	ShowState            GroupBuyingTaskShowState        `db:"ShowState"`
	AllowCancel          GroupBuyingTaskAllowCancelState `db:"AllowCancel"`
	TeamVisibleState     GroupBuyingTeamVisibleState     `db:"TeamVisibleState"`
	TeamIds              []uint32                        `db:"TeamIds"`
	Title                string                          `binding:"required,lte=64" db:"Title"`
	Introduction         string                          `binding:"required,lte=255" db:"Introduction"`
	CoverPicture         string                          `binding:"required,lte=255" db:"CoverPicture"`
	IllustrationPictures *TaskIllustrationPicturesType   `binding:"required" db:"IllustrationPictures"`
	Info                 *TaskInfoType                   `binding:"required" db:"Info"`
	Specification        []*AskTaskSpecificationItem     `binding:"required,gt=0,dive,required" db:"Specification"`
	Sku                  []*AskTaskSkuItem               `binding:"required,gt=0,dive,required" db:"Sku"`
	Combination          []*AskTaskCombinationItem       `binding:"required,dive,required" db:"Combination"`
	WxSellText           string                          `binding:"required" db:"WxSellText"`
}

func NewAskOrgTaskAddByOrganizationID() *AskOrgTaskAddByOrganizationID {
	return &AskOrgTaskAddByOrganizationID{
		TeamIds:              make([]uint32, 0),
		IllustrationPictures: NewTaskIllustrationPicturesType(),
		Info:                 NewTaskInfoType(),
		Specification:        make([]*AskTaskSpecificationItem, 0),
		Sku:                  make([]*AskTaskSkuItem, 0),
		Combination:          make([]*AskTaskCombinationItem, 0),
	}
}

type MetaApiOrgTaskAddByOrganizationID struct {
}

var META_ORG_TASK_ADD_BY_ORGANIZATION_ID = &MetaApiOrgTaskAddByOrganizationID{}

func (m *MetaApiOrgTaskAddByOrganizationID) GetMethod() string { return "POST" }
func (m *MetaApiOrgTaskAddByOrganizationID) GetURL() string {
	return "/group_buying_order/org/task/add/:organization_id"
}
func (m *MetaApiOrgTaskAddByOrganizationID) GetName() string { return "OrgTaskAddByOrganizationID" }
func (m *MetaApiOrgTaskAddByOrganizationID) GetType() string { return "json" }

// 添加团购任务
type ApiOrgTaskAddByOrganizationID struct {
	MetaApiOrgTaskAddByOrganizationID
	Ask    *AskOrgTaskAddByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
}

func (m *ApiOrgTaskAddByOrganizationID) GetQuery() interface{}  { return nil }
func (m *ApiOrgTaskAddByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskAddByOrganizationID) GetAsk() interface{}    { return m.Ask }
func (m *ApiOrgTaskAddByOrganizationID) GetAck() interface{}    { return nil }
func MakeApiOrgTaskAddByOrganizationID() ApiOrgTaskAddByOrganizationID {
	return ApiOrgTaskAddByOrganizationID{
		Ask: NewAskOrgTaskAddByOrganizationID(),
	}
}

type AskOrgTaskEditByOrganizationIDByTaskID struct {
	ShowStartTime        time.Time                     `binding:"required" db:"ShowStartTime"`
	StartTime            time.Time                     `db:"StartTime"`
	EndTime              time.Time                     `binding:"required" db:"EndTime"`
	SellType             GroupBuyingOrderTaskSellType  `binding:"required" db:"SellType"`
	Notes                string                        `binding:"required,lte=1000" db:"Notes"`
	ShowState            GroupBuyingTaskShowState      `db:"ShowState"`
	Title                string                        `binding:"required,lte=64" db:"Title"`
	Introduction         string                        `binding:"required,lte=255" db:"Introduction"`
	CoverPicture         string                        `binding:"required,lte=255" db:"CoverPicture"`
	IllustrationPictures *TaskIllustrationPicturesType `binding:"required" db:"IllustrationPictures"`
	Info                 *TaskInfoType                 `binding:"required" db:"Info"`
	Specification        []*AskTaskSpecificationItem   `binding:"required,gt=0,dive,required" db:"Specification"`
	Sku                  []*AskTaskSkuItem             `binding:"required,gt=0,dive,required" db:"Sku"`
	Combination          []*AskTaskCombinationItem     `binding:"required,dive,required" db:"Combination"`
	WxSellText           string                        `binding:"required" db:"WxSellText"`
}

func NewAskOrgTaskEditByOrganizationIDByTaskID() *AskOrgTaskEditByOrganizationIDByTaskID {
	return &AskOrgTaskEditByOrganizationIDByTaskID{
		IllustrationPictures: NewTaskIllustrationPicturesType(),
		Info:                 NewTaskInfoType(),
		Specification:        make([]*AskTaskSpecificationItem, 0),
		Sku:                  make([]*AskTaskSkuItem, 0),
		Combination:          make([]*AskTaskCombinationItem, 0),
	}
}

type MetaApiOrgTaskEditByOrganizationIDByTaskID struct {
}

var META_ORG_TASK_EDIT_BY_ORGANIZATION_ID_BY_TASK_ID = &MetaApiOrgTaskEditByOrganizationIDByTaskID{}

func (m *MetaApiOrgTaskEditByOrganizationIDByTaskID) GetMethod() string { return "POST" }
func (m *MetaApiOrgTaskEditByOrganizationIDByTaskID) GetURL() string {
	return "/group_buying_order/org/task/edit/:organization_id/:task_id"
}
func (m *MetaApiOrgTaskEditByOrganizationIDByTaskID) GetName() string {
	return "OrgTaskEditByOrganizationIDByTaskID"
}
func (m *MetaApiOrgTaskEditByOrganizationIDByTaskID) GetType() string { return "json" }

// 编辑团购任务
type ApiOrgTaskEditByOrganizationIDByTaskID struct {
	MetaApiOrgTaskEditByOrganizationIDByTaskID
	Ask    *AskOrgTaskEditByOrganizationIDByTaskID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		TaskID         uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiOrgTaskEditByOrganizationIDByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiOrgTaskEditByOrganizationIDByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskEditByOrganizationIDByTaskID) GetAsk() interface{}    { return m.Ask }
func (m *ApiOrgTaskEditByOrganizationIDByTaskID) GetAck() interface{}    { return nil }
func MakeApiOrgTaskEditByOrganizationIDByTaskID() ApiOrgTaskEditByOrganizationIDByTaskID {
	return ApiOrgTaskEditByOrganizationIDByTaskID{
		Ask: NewAskOrgTaskEditByOrganizationIDByTaskID(),
	}
}

type MetaApiOrgTaskInfoByTaskID struct {
}

var META_ORG_TASK_INFO_BY_TASK_ID = &MetaApiOrgTaskInfoByTaskID{}

func (m *MetaApiOrgTaskInfoByTaskID) GetMethod() string { return "GET" }
func (m *MetaApiOrgTaskInfoByTaskID) GetURL() string {
	return "/group_buying_order/org/task/info/:task_id"
}
func (m *MetaApiOrgTaskInfoByTaskID) GetName() string { return "OrgTaskInfoByTaskID" }
func (m *MetaApiOrgTaskInfoByTaskID) GetType() string { return "json" }

// 获取团购任务
type ApiOrgTaskInfoByTaskID struct {
	MetaApiOrgTaskInfoByTaskID
	Ack    *GroupBuyingOrderTask
	Params struct {
		TaskID uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiOrgTaskInfoByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiOrgTaskInfoByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskInfoByTaskID) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskInfoByTaskID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgTaskInfoByTaskID() ApiOrgTaskInfoByTaskID {
	return ApiOrgTaskInfoByTaskID{
		Ack: NewGroupBuyingOrderTask(),
	}
}

type MetaApiOrgTaskShowByOrganizationIDByTaskID struct {
}

var META_ORG_TASK_SHOW_BY_ORGANIZATION_ID_BY_TASK_ID = &MetaApiOrgTaskShowByOrganizationIDByTaskID{}

func (m *MetaApiOrgTaskShowByOrganizationIDByTaskID) GetMethod() string { return "POST" }
func (m *MetaApiOrgTaskShowByOrganizationIDByTaskID) GetURL() string {
	return "/group_buying_order/org/task/show/:organization_id/:task_id"
}
func (m *MetaApiOrgTaskShowByOrganizationIDByTaskID) GetName() string {
	return "OrgTaskShowByOrganizationIDByTaskID"
}
func (m *MetaApiOrgTaskShowByOrganizationIDByTaskID) GetType() string { return "json" }

// 上架
type ApiOrgTaskShowByOrganizationIDByTaskID struct {
	MetaApiOrgTaskShowByOrganizationIDByTaskID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		TaskID         uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiOrgTaskShowByOrganizationIDByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiOrgTaskShowByOrganizationIDByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskShowByOrganizationIDByTaskID) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskShowByOrganizationIDByTaskID) GetAck() interface{}    { return nil }
func MakeApiOrgTaskShowByOrganizationIDByTaskID() ApiOrgTaskShowByOrganizationIDByTaskID {
	return ApiOrgTaskShowByOrganizationIDByTaskID{}
}

type MetaApiOrgTaskHideByOrganizationIDByTaskID struct {
}

var META_ORG_TASK_HIDE_BY_ORGANIZATION_ID_BY_TASK_ID = &MetaApiOrgTaskHideByOrganizationIDByTaskID{}

func (m *MetaApiOrgTaskHideByOrganizationIDByTaskID) GetMethod() string { return "POST" }
func (m *MetaApiOrgTaskHideByOrganizationIDByTaskID) GetURL() string {
	return "/group_buying_order/org/task/hide/:organization_id/:task_id"
}
func (m *MetaApiOrgTaskHideByOrganizationIDByTaskID) GetName() string {
	return "OrgTaskHideByOrganizationIDByTaskID"
}
func (m *MetaApiOrgTaskHideByOrganizationIDByTaskID) GetType() string { return "json" }

// 下架
type ApiOrgTaskHideByOrganizationIDByTaskID struct {
	MetaApiOrgTaskHideByOrganizationIDByTaskID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		TaskID         uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiOrgTaskHideByOrganizationIDByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiOrgTaskHideByOrganizationIDByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskHideByOrganizationIDByTaskID) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskHideByOrganizationIDByTaskID) GetAck() interface{}    { return nil }
func MakeApiOrgTaskHideByOrganizationIDByTaskID() ApiOrgTaskHideByOrganizationIDByTaskID {
	return ApiOrgTaskHideByOrganizationIDByTaskID{}
}

type MetaApiOrgTaskDeleteByOrganizationIDByTaskID struct {
}

var META_ORG_TASK_DELETE_BY_ORGANIZATION_ID_BY_TASK_ID = &MetaApiOrgTaskDeleteByOrganizationIDByTaskID{}

func (m *MetaApiOrgTaskDeleteByOrganizationIDByTaskID) GetMethod() string { return "POST" }
func (m *MetaApiOrgTaskDeleteByOrganizationIDByTaskID) GetURL() string {
	return "/group_buying_order/org/task/delete/:organization_id/:task_id"
}
func (m *MetaApiOrgTaskDeleteByOrganizationIDByTaskID) GetName() string {
	return "OrgTaskDeleteByOrganizationIDByTaskID"
}
func (m *MetaApiOrgTaskDeleteByOrganizationIDByTaskID) GetType() string { return "json" }

// 删除团购任务
type ApiOrgTaskDeleteByOrganizationIDByTaskID struct {
	MetaApiOrgTaskDeleteByOrganizationIDByTaskID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		TaskID         uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiOrgTaskDeleteByOrganizationIDByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiOrgTaskDeleteByOrganizationIDByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskDeleteByOrganizationIDByTaskID) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskDeleteByOrganizationIDByTaskID) GetAck() interface{}    { return nil }
func MakeApiOrgTaskDeleteByOrganizationIDByTaskID() ApiOrgTaskDeleteByOrganizationIDByTaskID {
	return ApiOrgTaskDeleteByOrganizationIDByTaskID{}
}

type AckOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID struct {
	TotalGroupCount  uint32 `db:"TotalGroupCount"`
	BuyGroupCount    uint32 `db:"BuyGroupCount"`
	NotBuyGroupCount uint32 `db:"NotBuyGroupCount"`
}

func NewAckOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID() *AckOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID {
	return &AckOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID{}
}

type MetaApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID struct {
}

var META_ORG_TASK_FINISH_BUYING_CONFIRM_COUNT_BY_ORGANIZATION_ID_BY_TASK_ID = &MetaApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID{}

func (m *MetaApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID) GetMethod() string {
	return "GET"
}
func (m *MetaApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID) GetURL() string {
	return "/group_buying_order/org/task/finish_buying_confirm_count/:organization_id/:task_id"
}
func (m *MetaApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID) GetName() string {
	return "OrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID"
}
func (m *MetaApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID) GetType() string {
	return "json"
}

// 获取结团确认信息
type ApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID struct {
	MetaApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID
	Ack    *AckOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		TaskID         uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID) GetQuery() interface{} {
	return nil
}
func (m *ApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID) GetParams() interface{} {
	return &m.Params
}
func (m *ApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID) GetAsk() interface{} { return nil }
func (m *ApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID) GetAck() interface{} {
	return m.Ack
}
func MakeApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID() ApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID {
	return ApiOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID{
		Ack: NewAckOrgTaskFinishBuyingConfirmCountByOrganizationIDByTaskID(),
	}
}

type AckOrgTaskSoldListByOrganizationID struct {
	Count uint32                  `db:"Count"`
	List  []*GroupBuyingOrderTask `db:"List"`
}

func NewAckOrgTaskSoldListByOrganizationID() *AckOrgTaskSoldListByOrganizationID {
	return &AckOrgTaskSoldListByOrganizationID{
		List: make([]*GroupBuyingOrderTask, 0),
	}
}

type MetaApiOrgTaskSoldListByOrganizationID struct {
}

var META_ORG_TASK_SOLD_LIST_BY_ORGANIZATION_ID = &MetaApiOrgTaskSoldListByOrganizationID{}

func (m *MetaApiOrgTaskSoldListByOrganizationID) GetMethod() string { return "GET" }
func (m *MetaApiOrgTaskSoldListByOrganizationID) GetURL() string {
	return "/group_buying_order/org/task/sold_list/:organization_id"
}
func (m *MetaApiOrgTaskSoldListByOrganizationID) GetName() string {
	return "OrgTaskSoldListByOrganizationID"
}
func (m *MetaApiOrgTaskSoldListByOrganizationID) GetType() string { return "json" }

// 已结团并且已经销售出去的团购任务列表
type ApiOrgTaskSoldListByOrganizationID struct {
	MetaApiOrgTaskSoldListByOrganizationID
	Ack    *AckOrgTaskSoldListByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
		Search   string `form:"search" db:"Search"`
	}
}

func (m *ApiOrgTaskSoldListByOrganizationID) GetQuery() interface{}  { return &m.Query }
func (m *ApiOrgTaskSoldListByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskSoldListByOrganizationID) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskSoldListByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgTaskSoldListByOrganizationID() ApiOrgTaskSoldListByOrganizationID {
	return ApiOrgTaskSoldListByOrganizationID{
		Ack: NewAckOrgTaskSoldListByOrganizationID(),
	}
}

type AckOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID struct {
	Count uint32                          `db:"Count"`
	List  []*GroupBuyingOrderCommunityBuy `db:"List"`
}

func NewAckOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID() *AckOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID {
	return &AckOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID{
		List: make([]*GroupBuyingOrderCommunityBuy, 0),
	}
}

type MetaApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID struct {
}

var META_ORG_TASK_FINISH_BUYING_GROUP_LIST_BY_ORGANIZATION_ID_BY_TASK_ID = &MetaApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID{}

func (m *MetaApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID) GetMethod() string { return "GET" }
func (m *MetaApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID) GetURL() string {
	return "/group_buying_order/org/task/finish_buying_group_list/:organization_id/:task_id"
}
func (m *MetaApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID) GetName() string {
	return "OrgTaskFinishBuyingGroupListByOrganizationIDByTaskID"
}
func (m *MetaApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID) GetType() string { return "json" }

// 社群团购信息列表
type ApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID struct {
	MetaApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID
	Ack    *AckOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		TaskID         uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
	Query struct {
		Page     uint32 `form:"page" binding:"required,gt=0" db:"Page"`
		PageSize uint32 `form:"page_size" binding:"required,gt=0,lt=50" db:"PageSize"`
	}
}

func (m *ApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID) GetQuery() interface{} {
	return &m.Query
}
func (m *ApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID) GetParams() interface{} {
	return &m.Params
}
func (m *ApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID) GetAsk() interface{} { return nil }
func (m *ApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID) GetAck() interface{} { return m.Ack }
func MakeApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID() ApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID {
	return ApiOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID{
		Ack: NewAckOrgTaskFinishBuyingGroupListByOrganizationIDByTaskID(),
	}
}

type MetaApiOrgTaskFinishBuyingByOrganizationIDByTaskID struct {
}

var META_ORG_TASK_FINISH_BUYING_BY_ORGANIZATION_ID_BY_TASK_ID = &MetaApiOrgTaskFinishBuyingByOrganizationIDByTaskID{}

func (m *MetaApiOrgTaskFinishBuyingByOrganizationIDByTaskID) GetMethod() string { return "POST" }
func (m *MetaApiOrgTaskFinishBuyingByOrganizationIDByTaskID) GetURL() string {
	return "/group_buying_order/org/task/finish_buying/:organization_id/:task_id"
}
func (m *MetaApiOrgTaskFinishBuyingByOrganizationIDByTaskID) GetName() string {
	return "OrgTaskFinishBuyingByOrganizationIDByTaskID"
}
func (m *MetaApiOrgTaskFinishBuyingByOrganizationIDByTaskID) GetType() string { return "json" }

// 确定结团
type ApiOrgTaskFinishBuyingByOrganizationIDByTaskID struct {
	MetaApiOrgTaskFinishBuyingByOrganizationIDByTaskID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		TaskID         uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiOrgTaskFinishBuyingByOrganizationIDByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiOrgTaskFinishBuyingByOrganizationIDByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskFinishBuyingByOrganizationIDByTaskID) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskFinishBuyingByOrganizationIDByTaskID) GetAck() interface{}    { return nil }
func MakeApiOrgTaskFinishBuyingByOrganizationIDByTaskID() ApiOrgTaskFinishBuyingByOrganizationIDByTaskID {
	return ApiOrgTaskFinishBuyingByOrganizationIDByTaskID{}
}

type AskOrgIndentAddByOrganizationID struct {
	TaskIds []uint32 `binding:"required,gt=0" db:"TaskIds"`
}

func NewAskOrgIndentAddByOrganizationID() *AskOrgIndentAddByOrganizationID {
	return &AskOrgIndentAddByOrganizationID{
		TaskIds: make([]uint32, 0),
	}
}

type AckOrgIndentAddByOrganizationID struct {
	IndentId string `db:"IndentId"`
}

func NewAckOrgIndentAddByOrganizationID() *AckOrgIndentAddByOrganizationID {
	return &AckOrgIndentAddByOrganizationID{}
}

type MetaApiOrgIndentAddByOrganizationID struct {
}

var META_ORG_INDENT_ADD_BY_ORGANIZATION_ID = &MetaApiOrgIndentAddByOrganizationID{}

func (m *MetaApiOrgIndentAddByOrganizationID) GetMethod() string { return "POST" }
func (m *MetaApiOrgIndentAddByOrganizationID) GetURL() string {
	return "/group_buying_order/org/indent/add/:organization_id"
}
func (m *MetaApiOrgIndentAddByOrganizationID) GetName() string { return "OrgIndentAddByOrganizationID" }
func (m *MetaApiOrgIndentAddByOrganizationID) GetType() string { return "json" }

// ####################订货单
// 添加订货单团购任务
type ApiOrgIndentAddByOrganizationID struct {
	MetaApiOrgIndentAddByOrganizationID
	Ask    *AskOrgIndentAddByOrganizationID
	Ack    *AckOrgIndentAddByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
}

func (m *ApiOrgIndentAddByOrganizationID) GetQuery() interface{}  { return nil }
func (m *ApiOrgIndentAddByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgIndentAddByOrganizationID) GetAsk() interface{}    { return m.Ask }
func (m *ApiOrgIndentAddByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgIndentAddByOrganizationID() ApiOrgIndentAddByOrganizationID {
	return ApiOrgIndentAddByOrganizationID{
		Ask: NewAskOrgIndentAddByOrganizationID(),
		Ack: NewAckOrgIndentAddByOrganizationID(),
	}
}

type AskOrgSendAddByOrganizationID struct {
	TaskLineIds []*GroupBuyingTaskLineIDs `binding:"required,gt=0" db:"TaskLineIds"`
}

func NewAskOrgSendAddByOrganizationID() *AskOrgSendAddByOrganizationID {
	return &AskOrgSendAddByOrganizationID{
		TaskLineIds: make([]*GroupBuyingTaskLineIDs, 0),
	}
}

type AckOrgSendAddByOrganizationID struct {
	SendId string `db:"SendId"`
}

func NewAckOrgSendAddByOrganizationID() *AckOrgSendAddByOrganizationID {
	return &AckOrgSendAddByOrganizationID{}
}

type MetaApiOrgSendAddByOrganizationID struct {
}

var META_ORG_SEND_ADD_BY_ORGANIZATION_ID = &MetaApiOrgSendAddByOrganizationID{}

func (m *MetaApiOrgSendAddByOrganizationID) GetMethod() string { return "POST" }
func (m *MetaApiOrgSendAddByOrganizationID) GetURL() string {
	return "/group_buying_order/org/send/add/:organization_id"
}
func (m *MetaApiOrgSendAddByOrganizationID) GetName() string { return "OrgSendAddByOrganizationID" }
func (m *MetaApiOrgSendAddByOrganizationID) GetType() string { return "json" }

// ####################配送单
// 绑定配送单团购任务
type ApiOrgSendAddByOrganizationID struct {
	MetaApiOrgSendAddByOrganizationID
	Ask    *AskOrgSendAddByOrganizationID
	Ack    *AckOrgSendAddByOrganizationID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
	}
}

func (m *ApiOrgSendAddByOrganizationID) GetQuery() interface{}  { return nil }
func (m *ApiOrgSendAddByOrganizationID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgSendAddByOrganizationID) GetAsk() interface{}    { return m.Ask }
func (m *ApiOrgSendAddByOrganizationID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgSendAddByOrganizationID() ApiOrgSendAddByOrganizationID {
	return ApiOrgSendAddByOrganizationID{
		Ask: NewAskOrgSendAddByOrganizationID(),
		Ack: NewAckOrgSendAddByOrganizationID(),
	}
}

type AckOrgTaskLineListByOrganizationIDByTaskID struct {
	List []*GroupBuyingTaskLine `db:"List"`
}

func NewAckOrgTaskLineListByOrganizationIDByTaskID() *AckOrgTaskLineListByOrganizationIDByTaskID {
	return &AckOrgTaskLineListByOrganizationIDByTaskID{
		List: make([]*GroupBuyingTaskLine, 0),
	}
}

type MetaApiOrgTaskLineListByOrganizationIDByTaskID struct {
}

var META_ORG_TASK_LINE_LIST_BY_ORGANIZATION_ID_BY_TASK_ID = &MetaApiOrgTaskLineListByOrganizationIDByTaskID{}

func (m *MetaApiOrgTaskLineListByOrganizationIDByTaskID) GetMethod() string { return "GET" }
func (m *MetaApiOrgTaskLineListByOrganizationIDByTaskID) GetURL() string {
	return "/group_buying_order/org/task/line_list/:organization_id/:task_id"
}
func (m *MetaApiOrgTaskLineListByOrganizationIDByTaskID) GetName() string {
	return "OrgTaskLineListByOrganizationIDByTaskID"
}
func (m *MetaApiOrgTaskLineListByOrganizationIDByTaskID) GetType() string { return "json" }

// ################## 路线
// # 团购任务路线
type ApiOrgTaskLineListByOrganizationIDByTaskID struct {
	MetaApiOrgTaskLineListByOrganizationIDByTaskID
	Ack    *AckOrgTaskLineListByOrganizationIDByTaskID
	Params struct {
		OrganizationID uint32 `form:"organization_id" binding:"required,gt=0" db:"OrganizationID"`
		TaskID         uint32 `form:"task_id" binding:"required,gt=0" db:"TaskID"`
	}
}

func (m *ApiOrgTaskLineListByOrganizationIDByTaskID) GetQuery() interface{}  { return nil }
func (m *ApiOrgTaskLineListByOrganizationIDByTaskID) GetParams() interface{} { return &m.Params }
func (m *ApiOrgTaskLineListByOrganizationIDByTaskID) GetAsk() interface{}    { return nil }
func (m *ApiOrgTaskLineListByOrganizationIDByTaskID) GetAck() interface{}    { return m.Ack }
func MakeApiOrgTaskLineListByOrganizationIDByTaskID() ApiOrgTaskLineListByOrganizationIDByTaskID {
	return ApiOrgTaskLineListByOrganizationIDByTaskID{
		Ack: NewAckOrgTaskLineListByOrganizationIDByTaskID(),
	}
}
