package cidl

import "time"

// 送货路线
type GroupBuyingLine struct {
	LineId           uint32    `db:"lin_id"`
	OrganizationId   uint32    `db:"org_id"`
	OrganizationName string    `db:"org_name"`
	Name             string    `db:"name"`
	CommunityCount   uint32    `db:"community_count"`
	CreateTime       time.Time `db:"create_time"`
}

func NewGroupBuyingLine() *GroupBuyingLine {
	return &GroupBuyingLine{}
}

// 路线绑定社群
type GroupBuyingLineCommunity struct {
	GroupId        uint32    `db:"grp_id"`
	LineId         uint32    `db:"lin_id"`
	LineName       string    `db:"lin_name"`
	GroupName      string    `db:"grp_name"`
	ManagerUid     string    `db:"manager_uid"`
	ManagerName    string    `db:"manager_name"`
	ManagerMobile  string    `db:"manager_mobile"`
	OrganizationId uint32    `db:"org_id"`
	CreateTime     time.Time `db:"create_time"`
}

func NewGroupBuyingLineCommunity() *GroupBuyingLineCommunity {
	return &GroupBuyingLineCommunity{}
}

// 团购任务绑定路线
type GroupBuyingTaskLine struct {
	LineId     uint32    `db:"LineId"`
	LineName   string    `db:"LineName"`
	IsSelected bool      `db:"IsSelected"`
	UpdateTime time.Time `db:"UpdateTime"`
}

func NewGroupBuyingTaskLine() *GroupBuyingTaskLine {
	return &GroupBuyingTaskLine{}
}

type GroupBuyingTaskLineIDs struct {
	TaskId  uint32   `db:"TaskId"`
	LineIds []uint32 `db:"LineIds"`
}

func NewGroupBuyingTaskLineIDs() *GroupBuyingTaskLineIDs {
	return &GroupBuyingTaskLineIDs{
		LineIds: make([]uint32, 0),
	}
}
