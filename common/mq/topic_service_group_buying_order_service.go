package mq

import (
	"business/group-buying-order/cidl"

	"github.com/mz-eco/mz/kafka"
	"github.com/mz-eco/mz/log"
)

const (
	TOPIC_SERVICE_GROUP_BUYING_ORDER_SERVICE = "service-group-buying-order-service"
)

const (
	IDENT_SERVICE_GROUP_BUYING_ORDER_SERVICE_ADD_INDENT = "add_indent"
	IDENT_SERVICE_GROUP_BUYING_ORDER_SERVICE_ADD_SEND   = "add_send"
	IDENT_SERVICE_GROUP_BUYING_ORDER_SERVICE_ADD_ORDER  = "add_order"
)

var (
	topicServiceGroupBuyingOrderService *TopicServiceGroupBuyingOrderService = nil
)

func GetTopicServiceGroupBuyingOrderService() (topic *TopicServiceGroupBuyingOrderService, err error) {
	if topicServiceGroupBuyingOrderService != nil {
		topic = topicServiceGroupBuyingOrderService
		return
	}

	producer, err := kafka.NewAsyncProducer()
	if err != nil {
		log.Warnf("new async producer failed. %s", err)
		return
	}

	topicServiceGroupBuyingOrderService = &TopicServiceGroupBuyingOrderService{
		Producer: producer,
	}

	topic = topicServiceGroupBuyingOrderService

	return
}

type TopicServiceGroupBuyingOrderService struct {
	Producer *kafka.AsyncProducer
}

func (m *TopicServiceGroupBuyingOrderService) send(ident string, msg interface{}) (err error) {
	err = m.Producer.SendMessage(TOPIC_SERVICE_GROUP_BUYING_ORDER_SERVICE, ident, msg)
	if err != nil {
		log.Warnf("send topic message failed. %s", err)
		return
	}
	return
}

// 添加订货单
type AddIndentMessage struct {
	IndentID string
}

func (m *TopicServiceGroupBuyingOrderService) AddIndent(msg *AddIndentMessage) (err error) {
	return m.send(IDENT_SERVICE_GROUP_BUYING_ORDER_SERVICE_ADD_INDENT, msg)
}

// 添加配送单
type AddSendMessage struct {
	SendID       string
	AuthorUserId string
}

func (m *TopicServiceGroupBuyingOrderService) AddSend(msg *AddSendMessage) (err error) {
	return m.send(IDENT_SERVICE_GROUP_BUYING_ORDER_SERVICE_ADD_SEND, msg)
}

// 提交订单
type AddOrderMessage struct {
	GroupId      uint32
	TaskId       uint32
	CommunityBuy *cidl.GroupBuyingOrderCommunityBuy
}

func (m *TopicServiceGroupBuyingOrderService) AddOrder(msg *AddOrderMessage) (err error) {
	return m.send(IDENT_SERVICE_GROUP_BUYING_ORDER_SERVICE_ADD_ORDER, msg)
}
