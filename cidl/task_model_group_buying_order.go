package cidl

import (
	"encoding/json"
	"time"
)

func NewGroupBuyingOrderTaskCustom() *GroupBuyingOrderTask {
	return &GroupBuyingOrderTask{
		GroupBuyingOrderTaskContent: *NewGroupBuyingOrderTaskContent(),
	}
}

func NewGroupBuyingOrderTaskSpecificationCustom() *GroupBuyingOrderTaskSpecification {
	return &GroupBuyingOrderTaskSpecification{
		GroupBuyingSpecification: *NewGroupBuyingSpecification(),
		CombinationSkuMap:        make(map[string]*GroupBuyingOrderTaskCombinationItem),
	}
}

// 团购任务内容
func (m *GroupBuyingOrderTaskContent) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *GroupBuyingOrderTaskContent) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}

func NewGroupBuyingOrderTaskContentByTask(task *GroupBuyingOrderTask) *GroupBuyingOrderTaskContent {
	return &task.GroupBuyingOrderTaskContent
}

// 已上架，团购任务未开始
func (m *GroupBuyingOrderTask) GroupStateShowIsNotStart(checkTime time.Time) bool {
	if m.GroupState == GroupBuyingTaskGroupStateNotStart && m.StartTime.After(checkTime) {
		return true
	}

	return false
}

// 已上架，团购任务正在进行中
func (m *GroupBuyingOrderTask) GroupStateShowIsInProgress(checkTime time.Time) bool {
	if !m.StartTime.After(checkTime) && m.EndTime.After(checkTime) {
		if m.GroupState == GroupBuyingTaskGroupStateInProgress || m.GroupState == GroupBuyingTaskGroupStateNotStart {
			return true
		}
	}

	return false
}
