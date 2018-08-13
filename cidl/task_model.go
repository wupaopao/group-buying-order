package cidl

import (
	"encoding/json"

	"github.com/mz-eco/mz/log"
)

// 团购任务配图类型
type TaskIllustrationPicturesType []string

func NewTaskIllustrationPicturesType() *TaskIllustrationPicturesType {
	return &TaskIllustrationPicturesType{}
}

func (m *TaskIllustrationPicturesType) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *TaskIllustrationPicturesType) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}

// 团购任务规格
type TaskInfoType []*GroupBuyingTaskInfo

func NewTaskInfoType() *TaskInfoType {
	return &TaskInfoType{}
}

func (m *TaskInfoType) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *TaskInfoType) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}

// 规格
func (m *GroupBuyingSpecification) ToString() (strSpecification string, err error) {
	bytesSpecification, err := json.Marshal(m)
	if err != nil {
		log.Warnf("marshal group buying specification failed. %s", err)
		return
	}
	strSpecification = string(bytesSpecification)
	return
}

func (m *GroupBuyingSpecification) FromString(strSpecification string) (err error) {
	err = json.Unmarshal([]byte(strSpecification), m)
	if err != nil {
		log.Warnf("unmarshal str authorization failed. %s", err)
		return
	}
	return
}

// 规格组合
func (m *GroupBuyingOrderTaskSpecification) ToString() (strCombination string, err error) {
	bytesCombination, err := json.Marshal(m)
	if err != nil {
		log.Warnf("marshal group buying combination failed. %s", err)
		return
	}
	strCombination = string(bytesCombination)
	return
}

func (m *GroupBuyingOrderTaskSpecification) FromString(strCombination string) (err error) {
	err = json.Unmarshal([]byte(strCombination), m)
	if err != nil {
		log.Warnf("unmarshal str authorization failed. %s", err)
		return
	}
	return
}

func (m *CommunityBuyDetailType) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *CommunityBuyDetailType) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}

type CommunityOrderGoodsDetail []*GroupBuyingOrderCommunityBuy

func NewCommunityOrderGoodsDetail() *CommunityOrderGoodsDetail {
	return &CommunityOrderGoodsDetail{}
}

func (m *CommunityOrderGoodsDetail) ToString() (strValue string, err error) {
	bytesValue, err := json.Marshal(m)
	if err != nil {
		return
	}
	strValue = string(bytesValue)
	return
}

func (m *CommunityOrderGoodsDetail) FromString(strValue string) (err error) {
	err = json.Unmarshal([]byte(strValue), m)
	if err != nil {
		return
	}
	return
}

func (m *SpecificationPriceRange) SetMin(price float64) (success bool) {
	if price < m.Min || (price > 0 && m.Min == 0) {
		m.Min = price
		success = true
	}
	return
}

func (m *SpecificationPriceRange) SetMax(price float64) (success bool) {
	if price > m.Max {
		m.Max = price
		success = true
	}
	return
}

func (m *SpecificationPriceRange) SetRange(value float64) {
	m.SetMin(value)
	m.SetMax(value)
}
