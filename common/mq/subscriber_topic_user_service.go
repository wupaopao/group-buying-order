package mq

import (
	"encoding/json"
	"fmt"

	"business/group-buying-order/common/db"
	"business/user/common/mq"
	"common/utils"

	"github.com/mz-eco/mz/errors"
	"github.com/mz-eco/mz/kafka"
	"github.com/mz-eco/mz/log"
	"github.com/mz-eco/mz/settings"
)

var (
	topicUserServiceGroupSetting kafka.TopicGroupSetting
)

func init() {
	settings.RegisterWith(func(viper *settings.Viper) error {
		err := viper.Unmarshal(&topicUserServiceGroupSetting)
		if err != nil {
			panic(err)
			return err
		}
		return nil
	}, "kafka.subscribe.service_user_service")
}

type TopicUserServiceHandler struct {
	kafka.TopicHandler
}

func NewTopicUserServiceHandler() (handler *TopicUserServiceHandler, err error) {
	handler = &TopicUserServiceHandler{
		TopicHandler: kafka.TopicHandler{
			Topics:  []string{mq.TOPIC_SERVICE_USER_SERVICE},
			Brokers: topicUserServiceGroupSetting.Address,
			Group:   topicUserServiceGroupSetting.Group,
		},
	}

	handler.MessageHandle = handler.handleMessage

	return
}

func (m *TopicUserServiceHandler) handleMessage(identMessage *kafka.IdentMessage) (err error) {
	switch identMessage.Ident {
	case mq.IDENT_SERVICE_USER_SERVICE_MODIFY_USER_INFO:
		modifyInfo := &mq.ModifyUserInfoMessage{}
		err = json.Unmarshal(identMessage.Msg, modifyInfo)
		if err != nil {
			log.Warnf("unmarshal modify info message failed. %s", err)
			return
		}

		err = m.ModifyUserInfo(modifyInfo)
		if err != nil {
			log.Warnf("modify user info failed.")
			return
		}

	}

	return
}

func (m *TopicUserServiceHandler) ModifyUserInfo(msg *mq.ModifyUserInfoMessage) (err error) {
	if len(msg.Values) == 0 {
		err = errors.New("empty update field")
		return
	}

	// 更新路线社群
	strUpdateSql, updateValues := utils.DBUpdateFieldFilter(msg.Values, map[string]string{
		"mobile": "manager_mobile",
		"name":   "manager_name",
	})
	if len(updateValues) == 0 {
		return
	}

	dbGroupBuying := db.NewMallGroupBuyingOrder()
	strSql := `UPDATE gby_line_community SET %s WHERE manager_uid=?`
	strSql = fmt.Sprintf(strSql, strUpdateSql)
	updateValues = append(updateValues, msg.UserId)

	_, err = dbGroupBuying.DB.Exec(strSql, updateValues...)
	if err != nil {
		log.Warnf("update modified line community info failed. %s", err)
		return
	}

	return
}
