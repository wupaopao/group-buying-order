package cidl

import "encoding/json"

type GroupBuyingSendTasksBriefType []*GroupBuyingSendTaskBriefItem

func NewGroupBuyingSendTasksBriefType() *GroupBuyingSendTasksBriefType {
	return &GroupBuyingSendTasksBriefType{}
}

func (m *GroupBuyingSendTasksBriefType) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *GroupBuyingSendTasksBriefType) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}

// 未使用 GroupBuyingSend.TasksDetail
type GroupBuyingSendTasksDetailType []*GroupBuyingOrderTaskContent

func NewGroupBuyingSendTasksDetailType() *GroupBuyingSendTasksDetailType {
	return &GroupBuyingSendTasksDetailType{}
}

func (m *GroupBuyingSendTasksDetailType) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *GroupBuyingSendTasksDetailType) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}

type GroupBuyingSendLineStatisticType []*GroupBuyingSendLineStatisticsItem

func NewGroupBuyingSendLineStatisticType() *GroupBuyingSendLineStatisticType {
	return &GroupBuyingSendLineStatisticType{}
}

func (m *GroupBuyingSendLineStatisticType) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *GroupBuyingSendLineStatisticType) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}

type GroupBuyingSendCommunityStatisticType []*GroupBuyingSendCommunityStatisticsItem

func NewGroupBuyingSendCommunityStatisticType() *GroupBuyingSendCommunityStatisticType {
	return &GroupBuyingSendCommunityStatisticType{}
}

func (m *GroupBuyingSendCommunityStatisticType) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *GroupBuyingSendCommunityStatisticType) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}
