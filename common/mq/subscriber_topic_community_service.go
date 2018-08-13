package mq

import (
	"encoding/json"
	"fmt"

	"business/community/common/mq"
	"business/community/proxy/community"
	"business/group-buying-order/cidl"
	"business/group-buying-order/common/db"
	"common/utils"

	"github.com/mz-eco/mz/conn"
	"github.com/mz-eco/mz/errors"
	"github.com/mz-eco/mz/kafka"
	"github.com/mz-eco/mz/log"
	"github.com/mz-eco/mz/settings"
)

var (
	topicCommunityServiceGroupSetting kafka.TopicGroupSetting
)

func init() {
	settings.RegisterWith(func(viper *settings.Viper) error {
		err := viper.Unmarshal(&topicCommunityServiceGroupSetting)
		if err != nil {
			panic(err)
			return err
		}
		return nil
	}, "kafka.subscribe.service_community_service")
}

// 社群服务消息
type TopicCommunityServiceHandler struct {
	kafka.TopicHandler
}

func NewTopicCommunityServiceHandler() (handler *TopicCommunityServiceHandler, err error) {
	handler = &TopicCommunityServiceHandler{
		TopicHandler: kafka.TopicHandler{
			Topics:  []string{mq.TOPIC_SERVICE_COMMUNITY_SERVICE},
			Brokers: topicCommunityServiceGroupSetting.Address,
			Group:   topicCommunityServiceGroupSetting.Group,
		},
	}

	handler.MessageHandle = handler.handleMessage

	return
}

func (m *TopicCommunityServiceHandler) handleMessage(identMessage *kafka.IdentMessage) (err error) {
	switch identMessage.Ident {
	case mq.IDENT_SERVICE_COMMUNITY_SERVICE_CHANGE_GROUP_AUDIT_STATE:
		changeInfo := &mq.ChangeGroupAuditStateMessage{}
		err = json.Unmarshal(identMessage.Msg, changeInfo)
		if err != nil {
			log.Warnf("unmarshal change audit state message failed. %s", err)
			return
		}

		err = m.ChangeGroupAuditState(changeInfo)
		if err != nil {
			log.Warnf("change group audit state failed. %s", err)
			return
		}

	case mq.IDENT_SERVICE_COMMUNITY_SERVICE_MODIFY_GROUP_INFO:
		modifyInfo := &mq.ModifyGroupInfoMessage{}
		err = json.Unmarshal(identMessage.Msg, modifyInfo)
		if err != nil {
			log.Warnf("unmarshal change audit state message failed. %s", err)
			return
		}

		err = m.ModifyGroupInfo(modifyInfo)
		if err != nil {
			log.Warnf("change group name and address failed. %s", err)
			return
		}

	}
	return
}

func (m *TopicCommunityServiceHandler) ChangeGroupAuditState(msg *mq.ChangeGroupAuditStateMessage) (err error) {
	groupId := msg.GroupId
	group, err := community.NewProxy("community-service").InnerCommunityGroupInfoByGroupID(groupId)
	if err != nil {
		log.Warnf("get group info from proxy failed. %s", err)
		return
	}

	// 通过
	if group.AuditState == community.GroupAuditStatePass {
		// 获取路线
		dbGroupBuying := db.NewMallGroupBuyingOrder()
		lineCommunity, errGet := dbGroupBuying.GetLineCommunity(groupId)
		if errGet != nil && errGet != conn.ErrNoRows {
			err = errGet
			log.Warnf("get line community failed. %s", err)
			return
		}

		if lineCommunity != nil {
			return
		}

		// 添加社群路线中的社群信息
		lineCommunity = &cidl.GroupBuyingLineCommunity{
			GroupId:        groupId,
			GroupName:      group.Name,
			ManagerUid:     group.ManagerUserId,
			ManagerName:    group.ManagerName,
			ManagerMobile:  group.ManagerMobile,
			OrganizationId: group.OrganizationId,
		}

		_, err = dbGroupBuying.AddLineCommunity(lineCommunity)
		if err != nil {
			log.Warnf("add line community failed. %s", err)
			return
		}
	}

	return
}

func (m *TopicCommunityServiceHandler) ModifyGroupInfo(msg *mq.ModifyGroupInfoMessage) (err error) {
	if len(msg.Values) == 0 {
		err = errors.New("empty update field.")
		return
	}

	strUpdateSql, updateValues := utils.DBUpdateFieldFilter(msg.Values, map[string]string{
		"name": "grp_name",
	})
	if len(updateValues) == 0 {
		return
	}
	updateValues = append(updateValues, msg.GroupId)

	strSql := `UPDATE gby_line_community SET %s WHERE grp_id=?`
	strSql = fmt.Sprintf(strSql, strUpdateSql)

	dbGroupBuying := db.NewMallGroupBuyingOrder()
	_, err = dbGroupBuying.DB.Exec(strSql, updateValues...)
	if err != nil {
		log.Warnf("update modified line community info failed. %s", err)
		return
	}

	return
}
