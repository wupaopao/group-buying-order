package cidl

import "time"

// 订货单状态
type IndentState int

const (
	// 默认
	IndentStateDefault IndentState = 0
	// 订货单数据生成中
	IndentStateStatistic IndentState = 1
	// 订货单数据生成完成
	IndentStateFinishStatistic IndentState = 2
	// 订货单数据生成失败
	IndentStateFailStatistic IndentState = 3
)

func (m IndentState) String() string {
	switch m {

	case IndentStateDefault:
		return "IndentStateDefault<enum IndentState>"
	case IndentStateStatistic:
		return "IndentStateStatistic<enum IndentState>"
	case IndentStateFinishStatistic:
		return "IndentStateFinishStatistic<enum IndentState>"
	case IndentStateFailStatistic:
		return "IndentStateFailStatistic<enum IndentState>"
	default:
		return "UNKNOWN_Name_<IndentState>"
	}
}

type GroupBuyingIndentTasksBriefItem struct {
	TaskId    uint32    `db:"TaskId"`
	Title     string    `db:"Title"`
	StartTime time.Time `db:"StartTime"`
	EndTime   time.Time `db:"EndTime"`
}

func NewGroupBuyingIndentTasksBriefItem() *GroupBuyingIndentTasksBriefItem {
	return &GroupBuyingIndentTasksBriefItem{}
}

// 订货单
type GroupBuyingIndent struct {
	IndentId       string                          `db:"indent_id"`
	OrganizationId uint32                          `db:"org_id"`
	TasksBrief     *GroupBuyingIndentTaskBriefType `db:"tasks_brief"`
	State          IndentState                     `db:"state"`
	ExcelUrl       string                          `db:"excel_url"`
	Version        uint32                          `db:"version"`
	CreateTime     time.Time                       `db:"create_time"`
}

func NewGroupBuyingIndent() *GroupBuyingIndent {
	return &GroupBuyingIndent{
		TasksBrief: NewGroupBuyingIndentTaskBriefType(),
	}
}

// 订货单统计
type GroupBuyingIndentStatistics struct {
	IndentId    string                                 `db:"indent_id"`
	TaskId      uint32                                 `db:"task_id"`
	TaskContent *GroupBuyingOrderTaskContent           `db:"task_content"`
	Result      *GroupBuyingIndentStatisticsResultType `db:"result"`
	Version     uint32                                 `db:"version"`
	CreateTime  time.Time                              `db:"create_time"`
}

func NewGroupBuyingIndentStatistics() *GroupBuyingIndentStatistics {
	return &GroupBuyingIndentStatistics{
		TaskContent: NewGroupBuyingOrderTaskContent(),
		Result:      NewGroupBuyingIndentStatisticsResultType(),
	}
}

// 订货单统计项
type GroupBuyingIndentStatisticResultItem struct {
	GroupBuyingSkuMapItem
	Sales           uint32    `db:"Sales"`
	CommunityCount  uint32    `db:"CommunityCount"`
	TotalCost       float64   `db:"TotalCost"`
	TotalSettlement float64   `db:"TotalSettlement"`
	TaskId          uint32    `db:"TaskId"`
	StartTime       time.Time `db:"StartTime"`
	EndTime         time.Time `db:"EndTime"`
	Title           string    `db:"Title"`
}

func NewGroupBuyingIndentStatisticResultItem() *GroupBuyingIndentStatisticResultItem {
	return &GroupBuyingIndentStatisticResultItem{}
}
