package sched

import (
	"business/group-buying-order/cidl"
	"business/group-buying-order/common/db"
	"time"

	"github.com/mz-eco/mz/log"
)

func RunScheduler() {
	checkTaskGroupStateTicker := time.NewTicker(1 * time.Second) // 2s

	go func() {
		isDoingSchedule := false
		for {
			select {
			case <-checkTaskGroupStateTicker.C:
				if isDoingSchedule {
					return
				}

				isDoingSchedule = true
				err := CheckGroupBuyingOrderTaskGroupState()
				isDoingSchedule = false
				if err != nil {
					log.Warnf("ticker check task group state failed. %s", err)
				}

			}
		}

	}()

}

// 检查任务团购状态
func CheckGroupBuyingOrderTaskGroupState() (err error) {
	dbGroupBuying := db.NewMallGroupBuyingOrder()

	// 将自动开团的任务设置为团购状态进行中
	strSql := `
		UPDATE byo_task
		SET
			group_state=?
		WHERE
			group_state=?
			AND show_state=?
			AND start_time<=now()
			AND end_time>now()
	`
	_, err = dbGroupBuying.DB.Exec(strSql,
		cidl.GroupBuyingTaskGroupStateInProgress,
		cidl.GroupBuyingTaskGroupStateNotStart,
		cidl.GroupBuyingTaskShowStateShow)
	if err != nil {
		log.Warnf("update group state in progress failed. %s", err)
		return
	}

	// 团购状态已截单
	strSql = `
		UPDATE byo_task
		SET
			group_state=?
		WHERE
			group_state=?
			AND show_state=?
			AND end_time<=now()
	`
	_, err = dbGroupBuying.DB.Exec(strSql,
		cidl.GroupBuyingTaskGroupStateFinishOrdering,
		cidl.GroupBuyingTaskGroupStateInProgress,
		cidl.GroupBuyingTaskShowStateShow)
	if err != nil {
		log.Warnf("update group state finish ordering failed. %s", err)
		return
	}

	return
}
