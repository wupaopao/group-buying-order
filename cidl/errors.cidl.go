package cidl

import "fmt"

type autoErrorserrorsTcidl int

const (
	ErrTaskNotInProgress                         autoErrorserrorsTcidl = 9000 //团购任务团购状态不在进行中
	ErrTaskGroupStateIsNotInNotStart             autoErrorserrorsTcidl = 9001 //团购任务团购状态不在未开始状态
	ErrTaskGroupStateIsNotInNotStartOrInProgress autoErrorserrorsTcidl = 9002 //团购任务不在进行中或者未开始状态
	ErrTaskAddEditIllegalTime                    autoErrorserrorsTcidl = 9003 //团购开始时间或结束时间不正确
	ErrTaskGoodsNotExists                        autoErrorserrorsTcidl = 9004 //所选团购任务商品不存在
	ErrInventoryShortage                         autoErrorserrorsTcidl = 9005 //商品库存不足
	ErrTaskGoodsUnavailableForSale               autoErrorserrorsTcidl = 9006 //商品不允许被购买
	ErrTaskAddOrderFailed                        autoErrorserrorsTcidl = 9007 //订单添加失败
	ErrSendCommunityNotBindLine                  autoErrorserrorsTcidl = 9010 //存在社群未绑定配送路线
)

func (m autoErrorserrorsTcidl) Number() int { return int(m) }
func (m autoErrorserrorsTcidl) Message() string {
	switch m {

	case ErrTaskNotInProgress:
		return "团购任务团购状态不在进行中"
	case ErrTaskGroupStateIsNotInNotStart:
		return "团购任务团购状态不在未开始状态"
	case ErrTaskGroupStateIsNotInNotStartOrInProgress:
		return "团购任务不在进行中或者未开始状态"
	case ErrTaskAddEditIllegalTime:
		return "团购开始时间或结束时间不正确"
	case ErrTaskGoodsNotExists:
		return "所选团购任务商品不存在"
	case ErrInventoryShortage:
		return "商品库存不足"
	case ErrTaskGoodsUnavailableForSale:
		return "商品不允许被购买"
	case ErrTaskAddOrderFailed:
		return "订单添加失败"
	case ErrSendCommunityNotBindLine:
		return "存在社群未绑定配送路线"
	default:
		return "UNKNOWN_MESSAGE_autoErrorserrorsTcidl"
	}
}
func (m autoErrorserrorsTcidl) Name() string {
	switch m {

	case ErrTaskNotInProgress:
		return "ErrTaskNotInProgress"
	case ErrTaskGroupStateIsNotInNotStart:
		return "ErrTaskGroupStateIsNotInNotStart"
	case ErrTaskGroupStateIsNotInNotStartOrInProgress:
		return "ErrTaskGroupStateIsNotInNotStartOrInProgress"
	case ErrTaskAddEditIllegalTime:
		return "ErrTaskAddEditIllegalTime"
	case ErrTaskGoodsNotExists:
		return "ErrTaskGoodsNotExists"
	case ErrInventoryShortage:
		return "ErrInventoryShortage"
	case ErrTaskGoodsUnavailableForSale:
		return "ErrTaskGoodsUnavailableForSale"
	case ErrTaskAddOrderFailed:
		return "ErrTaskAddOrderFailed"
	case ErrSendCommunityNotBindLine:
		return "ErrSendCommunityNotBindLine"
	default:
		return "UNKNOWN_Name_autoErrorserrorsTcidl"
	}
}
func (m autoErrorserrorsTcidl) String() string {
	return fmt.Sprintf("[%d:%s]%s", m, m.Name(), m.Message())

}
