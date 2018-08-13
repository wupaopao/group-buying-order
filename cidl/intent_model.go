package cidl

import "encoding/json"

type GroupBuyingIndentTaskBriefType []*GroupBuyingIndentTasksBriefItem

func NewGroupBuyingIndentTaskBriefType() *GroupBuyingIndentTaskBriefType {
	return &GroupBuyingIndentTaskBriefType{}
}

func (m *GroupBuyingIndentTaskBriefType) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *GroupBuyingIndentTaskBriefType) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}

type GroupBuyingIndentStatisticsResultType map[string]*GroupBuyingIndentStatisticResultItem

func NewGroupBuyingIndentStatisticsResultType() *GroupBuyingIndentStatisticsResultType {
	return &GroupBuyingIndentStatisticsResultType{}
}

func (m *GroupBuyingIndentStatisticsResultType) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *GroupBuyingIndentStatisticsResultType) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}
