package mq

import (
	"github.com/mz-eco/mz/kafka"
	"github.com/mz-eco/mz/log"
)

func NewSubscriber() (subscriber *kafka.Subscriber, err error) {
	subscriber = &kafka.Subscriber{}

	// 社群
	topicHandler, err := NewTopicCommunityServiceHandler()
	if err != nil {
		log.Warnf("new topic community service handler failed. %s", err)
		return
	}

	subscriber.TopicHandlers = append(subscriber.TopicHandlers, topicHandler)

	// 团购任务
	topicGroupBuyingHandler, err := NewTopicGroupBuyingServiceHandler()
	if err != nil {
		log.Warnf("new topic group buying service handler failed. %s", err)
		return
	}

	subscriber.TopicHandlers = append(subscriber.TopicHandlers, topicGroupBuyingHandler)

	// 用户
	topicUserHandler, err := NewTopicUserServiceHandler()
	if err != nil {
		log.Warnf("new topic user service handler failed. %s", err)
		return
	}

	subscriber.TopicHandlers = append(subscriber.TopicHandlers, topicUserHandler)

	return
}
