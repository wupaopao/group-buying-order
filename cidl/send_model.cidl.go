package cidl

import "time"

// 送货单状态
type GroupBuyingSendState int

const (
	// 默认
	GroupBuyingSendStateDefault GroupBuyingSendState = 0
	// 送货单数据生成中
	GroupBuyingSendStateStatistic GroupBuyingSendState = 1
	// 送货单数据生成完成
	GroupBuyingSendStateFinishStatistic GroupBuyingSendState = 2
	// 送货单数据生成失败
	GroupBuyingSendStateFailStatistic GroupBuyingSendState = 3
)

func (m GroupBuyingSendState) String() string {
	switch m {

	case GroupBuyingSendStateDefault:
		return "GroupBuyingSendStateDefault<enum GroupBuyingSendState>"
	case GroupBuyingSendStateStatistic:
		return "GroupBuyingSendStateStatistic<enum GroupBuyingSendState>"
	case GroupBuyingSendStateFinishStatistic:
		return "GroupBuyingSendStateFinishStatistic<enum GroupBuyingSendState>"
	case GroupBuyingSendStateFailStatistic:
		return "GroupBuyingSendStateFailStatistic<enum GroupBuyingSendState>"
	default:
		return "UNKNOWN_Name_<GroupBuyingSendState>"
	}
}

// 配送单简介项
type GroupBuyingSendTaskBriefItem struct {
	TaskId    uint32    `db:"TaskId"`
	LineIds   []uint32  `db:"LineIds"`
	Title     string    `db:"Title"`
	StartTime time.Time `db:"StartTime"`
	EndTime   time.Time `db:"EndTime"`
}

func NewGroupBuyingSendTaskBriefItem() *GroupBuyingSendTaskBriefItem {
	return &GroupBuyingSendTaskBriefItem{
		LineIds: make([]uint32, 0),
	}
}

// 配送单
type GroupBuyingSend struct {
	SendId           string                         `db:"snd_id"`
	TasksBrief       *GroupBuyingSendTasksBriefType `db:"tasks_brief"`
	TasksDetail      string                         `db:"tasks_detail"`
	OrganizationId   uint32                         `db:"org_id"`
	OrganizationName string                         `db:"org_name"`
	State            GroupBuyingSendState           `db:"state"`
	ExcelUrl         string                         `db:"excel_url"`
	Version          uint32                         `db:"version"`
	CreateTime       time.Time                      `db:"create_time"`
}

func NewGroupBuyingSend() *GroupBuyingSend {
	return &GroupBuyingSend{
		TasksBrief: NewGroupBuyingSendTasksBriefType(),
	}
}

// 配送路线
type GroupBuyingSendLine struct {
	SendId           string                            `db:"send_id"`
	LineId           uint32                            `db:"lin_id"`
	LineName         string                            `db:"lin_name"`
	OrganizationId   uint32                            `db:"organization_id"`
	OrganizationName string                            `db:"organization_name"`
	CommunityCount   uint32                            `db:"community_count"`
	SettlementAmount float64                           `db:"settlement_amount"`
	Statistics       *GroupBuyingSendLineStatisticType `db:"statistics"`
	SendTime         time.Time                         `db:"send_time"`
	Version          uint32                            `db:"version"`
	CreateTime       time.Time                         `db:"create_time"`
}

func NewGroupBuyingSendLine() *GroupBuyingSendLine {
	return &GroupBuyingSendLine{
		Statistics: NewGroupBuyingSendLineStatisticType(),
	}
}

type GroupBuyingSendLineStatisticsItem struct {
	TaskId    uint32                                               `db:"TaskId"`
	TaskTitle string                                               `db:"TaskTitle"`
	Sku       map[string]*GroupBuyingSendLineStatisticsItemSkuItem `db:"Sku"`
}

func NewGroupBuyingSendLineStatisticsItem() *GroupBuyingSendLineStatisticsItem {
	return &GroupBuyingSendLineStatisticsItem{
		Sku: make(map[string]*GroupBuyingSendLineStatisticsItemSkuItem),
	}
}

type GroupBuyingSendLineStatisticsItemSkuItem struct {
	GroupBuyingSkuMapItem
	Sales           uint32  `db:"Sales"`
	CommunityCount  uint32  `db:"CommunityCount"`
	TotalCost       float64 `db:"TotalCost"`
	TotalSettlement float64 `db:"TotalSettlement"`
	TaskId          uint32  `db:"TaskId"`
	Title           string  `db:"Title"`
}

func NewGroupBuyingSendLineStatisticsItemSkuItem() *GroupBuyingSendLineStatisticsItemSkuItem {
	return &GroupBuyingSendLineStatisticsItemSkuItem{}
}

// 社群配送记录
type GroupBuyingSendCommunity struct {
	SendId                  string                                 `db:"send_id"`
	GroupId                 uint32                                 `db:"grp_id"`
	GroupName               string                                 `db:"grp_name"`
	GroupAddress            string                                 `db:"grp_address"`
	GroupManagerUid         string                                 `db:"grp_manager_uid"`
	GroupManagerName        string                                 `db:"grp_manager_name"`
	GroupManagerMobile      string                                 `db:"grp_manager_mobile"`
	OrganizationId          uint32                                 `db:"org_id"`
	OrganizationName        string                                 `db:"org_name"`
	OrganizationAddress     string                                 `db:"org_address"`
	OrganizationManagerUid  string                                 `db:"org_manager_uid"`
	OrganizationManagerName string                                 `db:"org_manager_name"`
	AuthorUid               string                                 `db:"author_uid"`
	AuthorName              string                                 `db:"author_name"`
	SettlementAmount        float64                                `db:"settlement_amount"`
	Statistics              *GroupBuyingSendCommunityStatisticType `db:"statistics"`
	SendTime                time.Time                              `db:"send_time"`
	Version                 uint32                                 `db:"version"`
	CreateTime              time.Time                              `db:"create_time"`
}

func NewGroupBuyingSendCommunity() *GroupBuyingSendCommunity {
	return &GroupBuyingSendCommunity{
		Statistics: NewGroupBuyingSendCommunityStatisticType(),
	}
}

type GroupBuyingSendCommunityStatisticsItem struct {
	TaskId    uint32                                                    `db:"TaskId"`
	TaskTitle string                                                    `db:"TaskTitle"`
	Sku       map[string]*GroupBuyingSendCommunityStatisticsItemSkuItem `db:"Sku"`
}

func NewGroupBuyingSendCommunityStatisticsItem() *GroupBuyingSendCommunityStatisticsItem {
	return &GroupBuyingSendCommunityStatisticsItem{
		Sku: make(map[string]*GroupBuyingSendCommunityStatisticsItemSkuItem),
	}
}

type GroupBuyingSendCommunityStatisticsItemSkuItem struct {
	GroupBuyingSkuMapItem
	Sales           uint32  `db:"Sales"`
	TotalCost       float64 `db:"TotalCost"`
	TotalSettlement float64 `db:"TotalSettlement"`
	TaskId          uint32  `db:"TaskId"`
	Title           string  `db:"Title"`
}

func NewGroupBuyingSendCommunityStatisticsItemSkuItem() *GroupBuyingSendCommunityStatisticsItemSkuItem {
	return &GroupBuyingSendCommunityStatisticsItemSkuItem{}
}
