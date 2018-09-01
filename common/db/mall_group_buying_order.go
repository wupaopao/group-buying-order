package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"strconv"

	"business/group-buying-order/cidl"

	"github.com/mz-eco/mz/conn"
	"github.com/mz-eco/mz/log"
)

type MallGroupBuyingOrder struct {
	DB *conn.DB
}

func NewMallGroupBuyingOrder() *MallGroupBuyingOrder {
	return &MallGroupBuyingOrder{
		DB: conn.NewDB("mal_group_buying"),
	}
}

func (m *MallGroupBuyingOrder) AddTask(task *cidl.GroupBuyingOrderTask) (result sql.Result, err error) {
	illustrationPictures, err := task.IllustrationPictures.ToString()
	if err != nil {
		log.Warnf("get string illustration pictures failed. %s", err)
		return
	}

	info, err := task.Info.ToString()
	if err != nil {
		log.Warnf("get string info failed. %s", err)
		return
	}

	specification, err := task.Specification.ToString()
	if err != nil {
		log.Warnf("get string specification failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO byo_task
			(
				org_id,
				title,
				introduction,
				cover_picture,
				illustration_pictures,
				info,
				specification,
				wx_sell_text,
				notes,
				show_start_time,
				start_time,
				end_time,
				sell_type,
				show_state,
				group_state,
				order_state,
				sales,
				is_delete,
				version,
				allow_cancel,
				team_visible	
			)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?)
	`
	result, err = m.DB.Exec(strSql,
		task.OrganizationId,
		task.Title,
		task.Introduction,
		task.CoverPicture,
		illustrationPictures,
		info,
		specification,
		task.WxSellText,
		task.Notes,
		task.ShowStartTime,
		task.StartTime,
		task.EndTime,
		task.SellType,
		task.ShowState,
		task.GroupState,
		task.OrderState,
		task.Sales,
		task.IsDelete,
		task.Version,
		task.AllowCancel,
		task.TeamVisibleState)

	return
}

func (m *MallGroupBuyingOrder) AddTaskVisibleTeam(taskId uint32, teamIds []uint32) (result sql.Result, err error) {
	if len(teamIds) <= 0 {
		return nil,nil
	}
	strSql := `
		INSERT INTO byo_task_team
		(
		 	tsk_id,
		 	team_id
		)
		VALUES
			%s	
		`
	var args []interface{}
	var sliceStrValue []string
	for _, teamId := range teamIds {
		sliceStrValue = append(sliceStrValue, "(?, ?)")
		args = append(args, taskId)
		args = append(args, teamId)
	}

	strValues := strings.Join(sliceStrValue, ",")
	strSql = fmt.Sprintf(strSql, strValues)

	result, err = m.DB.Exec(strSql, args...)

	return
}

func (m *MallGroupBuyingOrder) UpdateNotStartTask(task *cidl.GroupBuyingOrderTask) (result sql.Result, err error) {
	illustrationPictures, err := task.IllustrationPictures.ToString()
	if err != nil {
		log.Warnf("get string illustration pictures failed. %s", err)
		return
	}

	info, err := task.Info.ToString()
	if err != nil {
		log.Warnf("get string info failed. %s", err)
		return
	}

	specification, err := task.Specification.ToString()
	if err != nil {
		log.Warnf("get string specification failed. %s", err)
		return
	}

	strSql := `
		UPDATE byo_task
		SET
			title=?,
			introduction=?,
			cover_picture=?,
			illustration_pictures=?,
			info=?,
			specification=?,
			wx_sell_text=?,
			notes=?,
			show_start_time=?,
			start_time=?,
			end_time=?,
			sell_type=?,
			show_state=?,
			version=?,
			allow_cancel=?
		WHERE
			tsk_id=?
			AND org_id=?
			AND group_state=?
			AND (show_state=? OR now()<start_time)
	`
	// 如果show_state=0时，group_state是绝对准确的
	// 如果show_state=1时，因为group_state是由定时器设置的，所以需要增加now()<start_time来确保状态为未开始
	result, err = m.DB.Exec(strSql,
		task.Title,
		task.Introduction,
		task.CoverPicture,
		illustrationPictures,
		info,
		specification,
		task.WxSellText,
		task.Notes,
		task.ShowStartTime,
		task.StartTime,
		task.EndTime,
		task.SellType,
		task.ShowState,
		task.Version,
		task.AllowCancel,
		task.TaskId,
		task.OrganizationId,
		cidl.GroupBuyingTaskGroupStateNotStart,
		cidl.GroupBuyingTaskShowStateHidden)

	return
}

// 团购状态为进行中的任务不准修改开团日期、sku信息、团购状态、上下架状态、version
func (m *MallGroupBuyingOrder) UpdateInProgressTask(task *cidl.GroupBuyingOrderTask) (result sql.Result, err error) {
	illustrationPictures, err := task.IllustrationPictures.ToString()
	if err != nil {
		log.Warnf("get string illustration pictures failed. %s", err)
		return
	}

	info, err := task.Info.ToString()
	if err != nil {
		log.Warnf("get string info failed. %s", err)
		return
	}

	strSql := `
		UPDATE byo_task
		SET
			title=?,
			introduction=?,
			cover_picture=?,
			illustration_pictures=?,
			info=?,
			wx_sell_text=?,
			notes=?,
			end_time=?
		WHERE
			tsk_id=?
			AND org_id=?
			AND group_state=?
			AND now()>=start_time
			AND now()<end_time
	`
	// group_state为进行中时，show_state肯定为1
	result, err = m.DB.Exec(strSql,
		task.Title,
		task.Introduction,
		task.CoverPicture,
		illustrationPictures,
		info,
		task.WxSellText,
		task.Notes,
		task.EndTime,
		task.TaskId,
		task.OrganizationId,
		cidl.GroupBuyingTaskGroupStateInProgress)

	return
}

func (m *MallGroupBuyingOrder) UpdateTaskGroupState(organizationId uint32, taskId uint32, state cidl.GroupBuyingTaskGroupState) (result sql.Result, err error) {
	strSql := `UPDATE byo_task SET group_state=? WHERE tsk_id=? AND org_id=?`
	result, err = m.DB.Exec(strSql, state, taskId, organizationId)
	return
}

func (m *MallGroupBuyingOrder) UpdateTaskShowState(organizationId uint32, taskId uint32, showState cidl.GroupBuyingTaskShowState) (result sql.Result, err error) {
	strSql := `
		UPDATE
			byo_task
		SET
			show_state=?
		WHERE
			tsk_id=? AND org_id=?
	`
	result, err = m.DB.Exec(strSql, showState, taskId, organizationId)
	return
}

func (m *MallGroupBuyingOrder) UpdateTaskIsDelete(organizationId uint32, taskId uint32, isDelete bool) (result sql.Result, err error) {
	strSql := `
		UPDATE
			byo_task
		SET
			is_delete=?
		WHERE
			org_id=?
			AND tsk_id=?
			AND group_state=?
			AND now()<end_time`
	result, err = m.DB.Exec(strSql, isDelete, organizationId, taskId, cidl.GroupBuyingTaskGroupStateNotStart)
	return
}

func (m *MallGroupBuyingOrder) IncrTaskSales(taskId uint32, incrSales uint32) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_task
		SET
			sales=sales+?
		WHERE tsk_id=?
	`
	result, err = m.DB.Exec(strSql, incrSales, taskId)
	return
}

func (m *MallGroupBuyingOrder) DecrTaskSales(taskId uint32, decrSales uint32) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_task
		SET
			sales=sales-?
		WHERE tsk_id=?
	`
	result, err = m.DB.Exec(strSql, decrSales, taskId)
	return
}

func (m *MallGroupBuyingOrder) GetTask(taskId uint32) (task *cidl.GroupBuyingOrderTask, err error) {
	defer func() {
		if err != nil {
			task = nil
		}
	}()

	task = cidl.NewGroupBuyingOrderTaskCustom()
	var (
		illustrationPictures string
		info                 string
		specification        string
	)
	strSql := `
		SELECT
			tsk_id,
			org_id,
			title,
			introduction,
			cover_picture,
			illustration_pictures,
			info,
			specification,
			wx_sell_text,
			notes,
			show_start_time,
			start_time,
			end_time,
			sell_type,
			show_state,
			group_state,
			order_state,
			sales,
			is_delete,
			version,
			create_time,
			allow_cancel,
			team_visible
		FROM byo_task
		WHERE tsk_id=?
	`
	queryRow, err := m.DB.QueryRow(strSql, taskId)
	if err != nil {
		log.Warnf("get query row failed. %s", err)
		return
	}

	err = queryRow.Scan(
		&task.TaskId,
		&task.OrganizationId,
		&task.Title,
		&task.Introduction,
		&task.CoverPicture,
		&illustrationPictures,
		&info,
		&specification,
		&task.WxSellText,
		&task.Notes,
		&task.ShowStartTime,
		&task.StartTime,
		&task.EndTime,
		&task.SellType,
		&task.ShowState,
		&task.GroupState,
		&task.OrderState,
		&task.Sales,
		&task.IsDelete,
		&task.Version,
		&task.CreateTime,
		&task.AllowCancel,
		&task.TeamVisibleState,
	)
	if err != nil {
		if err != conn.ErrNoRows {
			log.Warnf("query task failed. %s", err)
		}
		return
	}

	err = task.IllustrationPictures.FromString(illustrationPictures)
	if err != nil {
		log.Warnf("init illustrationPicture from string failed. %s", err)
		return
	}

	err = task.Info.FromString(info)
	if err != nil {
		log.Warnf("init info from string failed. %s", err)
		return
	}

	err = task.Specification.FromString(specification)
	if err != nil {
		log.Warnf("init specification from string failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) GetTaskByOrganizationIdAndTaskId(organizationId uint32, taskId uint32) (task *cidl.GroupBuyingOrderTask, err error) {
	defer func() {
		if err != nil {
			task = nil
		}
	}()

	task = cidl.NewGroupBuyingOrderTaskCustom()
	var (
		illustrationPictures string
		info                 string
		specification        string
	)
	strSql := `
		SELECT
			tsk_id,
			org_id,
			title,
			introduction,
			cover_picture,
			illustration_pictures,
			info,
			specification,
			wx_sell_text,
			notes,
			show_start_time,
			start_time,
			end_time,
			sell_type,
			show_state,
			group_state,
			order_state,
			sales,
			is_delete,
			version,
			create_time
		FROM byo_task
		WHERE tsk_id=? AND org_id=?
	`
	queryRow, err := m.DB.QueryRow(strSql, taskId, organizationId)
	if err != nil {
		log.Warnf("get query row failed. %s", err)
		return
	}

	err = queryRow.Scan(
		&task.TaskId,
		&task.OrganizationId,
		&task.Title,
		&task.Introduction,
		&task.CoverPicture,
		&illustrationPictures,
		&info,
		&specification,
		&task.WxSellText,
		&task.Notes,
		&task.ShowStartTime,
		&task.StartTime,
		&task.EndTime,
		&task.SellType,
		&task.ShowState,
		&task.GroupState,
		&task.OrderState,
		&task.Sales,
		&task.IsDelete,
		&task.Version,
		&task.CreateTime,
	)
	if err != nil {
		if err != conn.ErrNoRows {
			log.Warnf("query task failed. %s", err)
		}
		return
	}

	err = task.IllustrationPictures.FromString(illustrationPictures)
	if err != nil {
		log.Warnf("init illustrationPicture from string failed. %s", err)
		return
	}

	err = task.Info.FromString(info)
	if err != nil {
		log.Warnf("init info from string failed. %s", err)
		return
	}

	err = task.Specification.FromString(specification)
	if err != nil {
		log.Warnf("init specification from string failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) GetFinishBuyingTasks(organizationId uint32, taskIds []uint32) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	strSql := `
		SELECT
			tsk_id,
			org_id,
			title,
			introduction,
			cover_picture,
			illustration_pictures,
			info,
			specification,
			wx_sell_text,
			notes,
			show_start_time,
			start_time,
			end_time,
			sell_type,
			show_state,
			group_state,
			order_state,
			sales,
			is_delete,
			version,
			create_time
		FROM byo_task
		WHERE
			org_id=?
			AND tsk_id IN (?)
			AND group_state=?
	`
	strSql, args, err := conn.In(strSql, organizationId, taskIds, cidl.GroupBuyingTaskGroupStateFinishBuying)
	if err != nil {
		log.Warnf("transform sql in array failed. %s", err)
		return
	}

	rows, err := m.DB.Query(strSql, args...)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

// 今天及未来的团购任务数目
func (m *MallGroupBuyingOrder) TaskNotReachEndTimeCount(organizationId uint32) (count uint32, err error) {
	strSql := `
		SELECT
			COUNT(*)
		FROM
			byo_task
		WHERE
			org_id=?
			AND now()<=end_time
			AND is_delete=0
	`
	err = m.DB.Get(&count, strSql, organizationId)
	return
}

// 今天及未来的团购任务列表
func (m *MallGroupBuyingOrder) TaskNotReachEndTimeList(organizationId uint32, page uint32, pageSize uint32, idAsc bool) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			tsk_id,
			org_id,
			title,
			introduction,
			cover_picture,
			illustration_pictures,
			info,
			specification,
			wx_sell_text,
			notes,
			show_start_time,
			start_time,
			end_time,
			sell_type,
			show_state,
			group_state,
			order_state,
			sales,
			is_delete,
			version,
			create_time,
			allow_cancel
		FROM byo_task
		WHERE
			org_id=?
			AND now()<=end_time
			AND is_delete=0
		ORDER BY tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, pageSize, offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
			&task.AllowCancel,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

func (m *MallGroupBuyingOrder) TaskMonthCount(organizationId uint32, year uint32, month uint32) (count uint32, err error) {
	strSql := `
		SELECT
			COUNT(*)
		FROM byo_task
		WHERE
			org_id=?
			AND YEAR(start_time)=?
			AND MONTH(start_time)=?
			AND is_delete=0
	`
	err = m.DB.Get(&count, strSql, organizationId, year, month)
	return
}

func (m *MallGroupBuyingOrder) TaskMonthList(organizationId uint32, year uint32, month uint32, page uint32, pageSize uint32, idAsc bool) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			tsk_id,
			org_id,
			title,
			introduction,
			cover_picture,
			illustration_pictures,
			info,
			specification,
			wx_sell_text,
			notes,
			show_start_time,
			start_time,
			end_time,
			sell_type,
			show_state,
			group_state,
			order_state,
			sales,
			is_delete,
			version,
			create_time,
			allow_cancel
		FROM byo_task
		WHERE
			org_id=?
			AND YEAR(start_time)=?
			AND MONTH(start_time)=?
			AND is_delete=0
		ORDER BY tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, year, month, pageSize, offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
			&task.AllowCancel,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

func (m *MallGroupBuyingOrder) TaskFinishBuyingCount(organizationId uint32) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM byo_task WHERE org_id=? AND group_state=?`
	err = m.DB.Get(&count, strSql, organizationId, cidl.GroupBuyingTaskGroupStateFinishBuying)
	return
}

func (m *MallGroupBuyingOrder) TaskFinishBuyingList(organizationId uint32, page uint32, pageSize uint32, idAsc bool) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			tsk_id,
			org_id,
			title,
			introduction,
			cover_picture,
			illustration_pictures,
			info,
			specification,
			wx_sell_text,
			notes,
			show_start_time,
			start_time,
			end_time,
			sell_type,
			show_state,
			group_state,
			order_state,
			sales,
			is_delete,
			version,
			create_time,
			allow_cancel
		FROM byo_task
		WHERE org_id=? AND group_state=?
		ORDER BY tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, cidl.GroupBuyingTaskGroupStateFinishBuying, pageSize, offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
			&task.AllowCancel,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

func (m *MallGroupBuyingOrder) TaskFinishBuyingSearchCount(organizationId uint32, search string) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM byo_task WHERE org_id=? AND group_state=? AND title LIKE ?`
	search = "%" + search + "%"
	err = m.DB.Get(&count, strSql, organizationId, cidl.GroupBuyingTaskGroupStateFinishBuying, search)
	return
}

func (m *MallGroupBuyingOrder) TaskFinishBuyingSearchList(organizationId uint32, search string, page uint32, pageSize uint32, idAsc bool) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0.")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			tsk_id,
			org_id,
			title,
			introduction,
			cover_picture,
			illustration_pictures,
			info,
			specification,
			wx_sell_text,
			notes,
			show_start_time,
			start_time,
			end_time,
			sell_type,
			show_state,
			group_state,
			order_state,
			sales,
			is_delete,
			version,
			create_time,
			allow_cancel
		FROM byo_task
		WHERE org_id=? AND group_state=? AND title LIKE ?
		ORDER BY tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	search = "%" + search + "%"
	rows, err := m.DB.Query(strSql, organizationId, cidl.GroupBuyingTaskGroupStateFinishBuying, search, pageSize, offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
			&task.AllowCancel,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

// 已经销售出去，并且已经结团的团购任务
func (m *MallGroupBuyingOrder) TaskSoldFinishBuyingCount(organizationId uint32) (count uint32, err error) {
	strSql := `
		SELECT
			COUNT(*)
		FROM
			byo_task
		WHERE
			org_id=?
			AND group_state=?
			AND sales>0
	`
	err = m.DB.Get(&count, strSql, organizationId, cidl.GroupBuyingTaskGroupStateFinishBuying)
	return
}

func (m *MallGroupBuyingOrder) TaskSoldFinishBuyingList(organizationId uint32, page uint32, pageSize uint32, idAsc bool) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			tsk_id,
			org_id,
			title,
			introduction,
			cover_picture,
			illustration_pictures,
			info,
			specification,
			wx_sell_text,
			notes,
			show_start_time,
			start_time,
			end_time,
			sell_type,
			show_state,
			group_state,
			order_state,
			sales,
			is_delete,
			version,
			create_time,
			allow_cancel
		FROM byo_task
		WHERE
			org_id=?
			AND group_state=?
			AND sales>0
		ORDER BY tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, cidl.GroupBuyingTaskGroupStateFinishBuying, pageSize, offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
			&task.AllowCancel,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

func (m *MallGroupBuyingOrder) TaskSoldFinishBuyingSearchCount(organizationId uint32, search string) (count uint32, err error) {
	strSql := `
		SELECT
			COUNT(*)
		FROM
			byo_task
		WHERE
			org_id=?
			AND group_state=?
			AND sales>0
			AND title LIKE ?`
	search = "%" + search + "%"
	err = m.DB.Get(&count, strSql, organizationId, cidl.GroupBuyingTaskGroupStateFinishBuying, search)
	return
}

func (m *MallGroupBuyingOrder) TaskSoldFinishBuyingSearchList(organizationId uint32, search string, page uint32, pageSize uint32, idAsc bool) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0.")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			tsk_id,
			org_id,
			title,
			introduction,
			cover_picture,
			illustration_pictures,
			info,
			specification,
			wx_sell_text,
			notes,
			show_start_time,
			start_time,
			end_time,
			sell_type,
			show_state,
			group_state,
			order_state,
			sales,
			is_delete,
			version,
			create_time,
			allow_cancel
		FROM byo_task
		WHERE
			org_id=?
			AND group_state=?
			AND sales>0
			AND title LIKE ?
		ORDER BY tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	search = "%" + search + "%"
	rows, err := m.DB.Query(strSql, organizationId, cidl.GroupBuyingTaskGroupStateFinishBuying, search, pageSize, offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
			&task.AllowCancel,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

// 微信小程序今日团购
func (m *MallGroupBuyingOrder) TaskTodayCount(organizationId uint32, teamIds []uint32) (count uint32, err error) {

	strSql := `
		SELECT
			COUNT(*)
		FROM
			byo_task a
		LEFT JOIN
			byo_task_team b
		ON
			a.tsk_id=b.tsk_id
		WHERE
			a.org_id=?
			AND a.show_state=?
			AND (TO_DAYS(a.show_start_time)-TO_DAYS(now())<=0)
			AND (TO_DAYS(a.end_time)-TO_DAYS(now())>=0)
			AND a.is_delete=0
			AND (
				a.team_visible=1 
				OR
				b.team_id in (?) 

			)
		`
	var sliceStrValues []string
	for _, teamId := range teamIds {
		sliceStrValues = append(sliceStrValues,strconv.FormatUint(uint64(teamId), 10))
	}
	strValues := strings.Join(sliceStrValues, ",")

	err = m.DB.Get(&count, strSql, organizationId, cidl.GroupBuyingTaskShowStateShow, strValues)
	return
}

func (m *MallGroupBuyingOrder) TaskTodayList(organizationId uint32, teamIds []uint32,  page uint32, pageSize uint32, idAsc bool) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			a.tsk_id,
			a.org_id,
			a.title,
			a.introduction,
			a.cover_picture,
			a.illustration_pictures,
			a.info,
			a.specification,
			a.wx_sell_text,
			a.notes,
			a.show_start_time,
			a.start_time,
			a.end_time,
			a.sell_type,
			a.show_state,
			a.group_state,
			a.order_state,
			a.sales,
			a.is_delete,
			a.version,
			a.create_time
		FROM
			byo_task a
		LEFT JOIN
			byo_task_team b
		ON
			a.tsk_id=b.tsk_id
		WHERE
			a.org_id=?
			AND a.show_state=?
			AND (TO_DAYS(a.show_start_time)-TO_DAYS(now())<=0)
			AND (TO_DAYS(a.end_time)-TO_DAYS(now())>=0)
			AND a.is_delete=0
			AND (
				a.team_visible=1 
				OR
				b.team_id in (?) 

			)
		ORDER BY a.tsk_id %s
		LIMIT ? OFFSET ?
	`

	var sliceStrValues []string
	for _, teamId := range teamIds {
		sliceStrValues = append(sliceStrValues,strconv.FormatUint(uint64(teamId), 10))
	}
	strValues := strings.Join(sliceStrValues, ",")

	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, cidl.GroupBuyingTaskShowStateShow, strValues, pageSize, offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

// 未来团购
func (m *MallGroupBuyingOrder) TaskFutureCount(organizationId uint32, teamIds []uint32) (count uint32, err error) {
	strSql := `
		SELECT
			COUNT(*)
		FROM
			byo_task a
		LEFT JOIN
			byo_task_team b
		ON
			a.tsk_id=b.tsk_id
		WHERE
			a.org_id=?
			AND a.show_state=?
			AND (TO_DAYS(a.show_start_time)-TO_DAYS(now())>0)
			AND a.is_delete=0
			AND (
				a.team_visible=1 
				OR
				b.team_id in (?) 
			)
	`
	var sliceStrValues []string
	for _, teamId := range teamIds {
		sliceStrValues = append(sliceStrValues,strconv.FormatUint(uint64(teamId), 10))
	}
	strValues := strings.Join(sliceStrValues, ",")

	err = m.DB.Get(&count, strSql, organizationId, cidl.GroupBuyingTaskShowStateShow, strValues)
	return
}

func (m *MallGroupBuyingOrder) TaskFutureList(organizationId uint32, teamIds []uint32, page uint32, pageSize uint32, idAsc bool) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			a.tsk_id,
			a.org_id,
			a.title,
			a.introduction,
			a.cover_picture,
			a.illustration_pictures,
			a.info,
			a.specification,
			a.wx_sell_text,
			a.notes,
			a.show_start_time,
			a.start_time,
			a.end_time,
			a.sell_type,
			a.show_state,
			a.group_state,
			a.order_state,
			a.sales,
			a.is_delete,
			a.version,
			a.create_time
		FROM 
			byo_task a
		LEFT JOIN
			byo_task_team b
		ON
			a.tsk_id=b.tsk_id
		WHERE
			a.org_id=?
			AND a.show_state=?
			AND (TO_DAYS(a.show_start_time)-TO_DAYS(now())>0)
			AND a.is_delete=0
			AND (
				a.team_visible=1 
				OR
				b.team_id in (?) 
			)
		ORDER BY a.tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	var sliceStrValues []string
	for _, teamId := range teamIds {
		sliceStrValues = append(sliceStrValues,strconv.FormatUint(uint64(teamId), 10))
	}
	strValues := strings.Join(sliceStrValues, ",")
	rows, err := m.DB.Query(strSql, organizationId, cidl.GroupBuyingTaskShowStateShow, strValues, pageSize, offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

// 历史团购
func (m *MallGroupBuyingOrder) TaskHistoryCount(organizationId uint32, teamIds []uint32) (count uint32, err error) {
	strSql := `
		SELECT
			COUNT(*)
		FROM
			byo_task a
		LEFT JOIN
			byo_task_team b
		ON
			a.tsk_id=b.tsk_id
		WHERE
			a.org_id=?
			AND a.show_state=?
			AND a.group_state=?
			AND a.is_delete=0
			AND (
				a.team_visible=1 
				OR
				b.team_id in (?) 
			)
	`
	var sliceStrValues []string
	for _, teamId := range teamIds {
		sliceStrValues = append(sliceStrValues,strconv.FormatUint(uint64(teamId), 10))
	}
	strValues := strings.Join(sliceStrValues, ",")
	err = m.DB.Get(&count,
		strSql,
		organizationId,
		cidl.GroupBuyingTaskShowStateShow,
		cidl.GroupBuyingTaskGroupStateFinishBuying,
		strValues)

	return
}

func (m *MallGroupBuyingOrder) TaskHistoryList(organizationId uint32,teamIds []uint32, page uint32, pageSize uint32, idAsc bool) (tasks []*cidl.GroupBuyingOrderTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			a.tsk_id,
			a.org_id,
			a.title,
			a.introduction,
			a.cover_picture,
			a.illustration_pictures,
			a.info,
			a.specification,
			a.wx_sell_text,
			a.notes,
			a.show_start_time,
			a.start_time,
			a.end_time,
			a.sell_type,
			a.show_state,
			a.group_state,
			a.order_state,
			a.sales,
			a.is_delete,
			a.version,
			a.create_time
		FROM 
			byo_task a
		LEFT JOIN
			byo_task_team b
		ON
			a.tsk_id=b.tsk_id
		WHERE
			a.org_id=?
			AND a.show_state=?
			AND a.group_state=?
			AND a.is_delete=0
			AND (
				a.team_visible=1 
				OR
				b.team_id in (?) 
			)
		ORDER BY a.tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)

	var sliceStrValues []string
	for _, teamId := range teamIds {
		sliceStrValues = append(sliceStrValues,strconv.FormatUint(uint64(teamId), 10))
	}
	strValues := strings.Join(sliceStrValues, ",")
	rows, err := m.DB.Query(strSql,
		organizationId,
		cidl.GroupBuyingTaskShowStateShow,
		cidl.GroupBuyingTaskGroupStateFinishBuying,
		strValues,
		pageSize,
		offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		task := cidl.NewGroupBuyingOrderTaskCustom()
		var (
			illustrationPictures string
			info                 string
			specification        string
		)
		err = rows.Scan(
			&task.TaskId,
			&task.OrganizationId,
			&task.Title,
			&task.Introduction,
			&task.CoverPicture,
			&illustrationPictures,
			&info,
			&specification,
			&task.WxSellText,
			&task.Notes,
			&task.ShowStartTime,
			&task.StartTime,
			&task.EndTime,
			&task.SellType,
			&task.ShowState,
			&task.GroupState,
			&task.OrderState,
			&task.Sales,
			&task.IsDelete,
			&task.Version,
			&task.CreateTime,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query task failed. %s", err)
			}
			return
		}

		err = task.IllustrationPictures.FromString(illustrationPictures)
		if err != nil {
			log.Warnf("init illustrationPicture from string failed. %s", err)
			return
		}

		err = task.Info.FromString(info)
		if err != nil {
			log.Warnf("init info from string failed. %s", err)
			return
		}

		err = task.Specification.FromString(specification)
		if err != nil {
			log.Warnf("init specification from string failed. %s", err)
			return
		}

		tasks = append(tasks, task)
	}

	return
}

// 批量获取任务信息
func (m *MallGroupBuyingOrder) WxXcxTaskBatchWxSellTextList(organizationId uint32, taskIds []uint32) (texts []string, err error) {
	strSql := `
		SELECT
			wx_sell_text
		FROM byo_task
		WHERE
			org_id=?
			AND tsk_id IN (?)
			AND show_state=?
	`
	strSql, args, err := conn.In(strSql, organizationId, taskIds, cidl.GroupBuyingTaskShowStateShow)
	if err != nil {
		log.Warnf("transform sql in array failed. %s", err)
		return
	}

	err = m.DB.Select(&texts, strSql, args...)
	if err != nil {
		log.Warnf("query task wx_sell_text list failed. %s", err)
		return
	}

	return
}

// 社群购买商品
func (m *MallGroupBuyingOrder) TxAddCommunityBuy(tx *conn.Tx, communityBuy *cidl.GroupBuyingOrderCommunityBuy) (result sql.Result, err error) {
	buyDetail, err := communityBuy.BuyDetail.ToString()
	if err != nil {
		log.Warnf("get string buy_detail failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO	byo_community_buy
			(
				cby_id,
				order_id,
				grp_id,
				grp_ord_id,
				grp_name,
				tsk_id,
				tsk_title,
				manager_uid,
				manager_name,
				manager_mobile,
				sku_id,
				buy_detail,
				count,
				total_market_price,
				total_group_buying_price,
				total_settlement_price,
				total_cost_price,
				version
			)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err = tx.Exec(strSql,
		communityBuy.BuyId,
		communityBuy.OrderId,
		communityBuy.GroupId,
		communityBuy.GroupOrderId,
		communityBuy.GroupName,
		communityBuy.TaskId,
		communityBuy.TaskTitle,
		communityBuy.ManagerUserId,
		communityBuy.ManagerName,
		communityBuy.ManagerMobile,
		communityBuy.SkuId,
		buyDetail,
		communityBuy.Count,
		communityBuy.TotalMarketPrice,
		communityBuy.TotalGroupBuyingPrice,
		communityBuy.TotalSettlementPrice,
		communityBuy.TotalCostPrice,
		communityBuy.Version,
	)

	return
}

func (m *MallGroupBuyingOrder) GetCommunityBuy(buyId string) (communityBuy *cidl.GroupBuyingOrderCommunityBuy, err error) {
	defer func() {
		if err != nil {
			communityBuy = nil
		}
	}()

	communityBuy = cidl.NewGroupBuyingOrderCommunityBuy()
	var buyDetail string
	strSql := `
		SELECT
			cby_id,
			order_id,
			grp_id,
			grp_ord_id,
			grp_name,
			tsk_id,
			tsk_title,
			manager_uid,
			manager_name,
			manager_mobile,
			sku_id,
			buy_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version
		FROM
			byo_community_buy
		WHERE
			cby_id=?
	`
	queryRow, err := m.DB.QueryRow(strSql, buyId)
	if err != nil {
		log.Warnf("get query row failed. %s", err)
		return
	}

	err = queryRow.Scan(
		&communityBuy.BuyId,
		&communityBuy.OrderId,
		&communityBuy.GroupId,
		&communityBuy.GroupOrderId,
		&communityBuy.GroupName,
		&communityBuy.TaskId,
		&communityBuy.TaskTitle,
		&communityBuy.ManagerUserId,
		&communityBuy.ManagerName,
		&communityBuy.ManagerMobile,
		&communityBuy.SkuId,
		&buyDetail,
		&communityBuy.Count,
		&communityBuy.TotalMarketPrice,
		&communityBuy.TotalGroupBuyingPrice,
		&communityBuy.TotalSettlementPrice,
		&communityBuy.TotalCostPrice,
		&communityBuy.Version,
	)
	if err != nil {
		if err != conn.ErrNoRows {
			log.Warnf("query community buy failed. %s", err)
		}
		return
	}

	err = communityBuy.BuyDetail.FromString(buyDetail)
	if err != nil {
		log.Warnf("init community buy buy detail failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) CommunityBuyListByGroupIdTaskId(groupId uint32, taskId uint32) (buys []*cidl.GroupBuyingOrderCommunityBuy, err error) {
	strSql := `
		SELECT
			cby_id,
			order_id,
			grp_id,
			grp_ord_id,
			grp_name,
			tsk_id,
			tsk_title,
			manager_uid,
			manager_name,
			manager_mobile,
			sku_id,
			buy_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version
		FROM
			byo_community_buy
		WHERE
			grp_id=?
			AND tsk_id=?
	`
	rows, err := m.DB.Query(strSql, groupId, taskId)
	if err != nil {
		log.Warnf("query community buy failed. %s", err)
		return
	}

	for rows.Next() {
		var buyDetail string
		communityBuy := cidl.NewGroupBuyingOrderCommunityBuy()
		err = rows.Scan(
			&communityBuy.BuyId,
			&communityBuy.OrderId,
			&communityBuy.GroupId,
			&communityBuy.GroupOrderId,
			&communityBuy.GroupName,
			&communityBuy.TaskId,
			&communityBuy.TaskTitle,
			&communityBuy.ManagerUserId,
			&communityBuy.ManagerName,
			&communityBuy.ManagerMobile,
			&communityBuy.SkuId,
			&buyDetail,
			&communityBuy.Count,
			&communityBuy.TotalMarketPrice,
			&communityBuy.TotalGroupBuyingPrice,
			&communityBuy.TotalSettlementPrice,
			&communityBuy.TotalCostPrice,
			&communityBuy.Version,
		)

		if err != nil {
			log.Warnf("scan community buy failed. %s", err)
			return
		}

		err = communityBuy.BuyDetail.FromString(buyDetail)
		if err != nil {
			log.Warnf("init buy detail from string failed. %s", err)
			return
		}

		buys = append(buys, communityBuy)
	}

	return
}

func (m *MallGroupBuyingOrder) CommunityBuyCount(taskId uint32) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM byo_community_buy WHERE tsk_id=?`
	err = m.DB.Get(&count, strSql, taskId)
	return
}

func (m *MallGroupBuyingOrder) CommunityBuyList(taskId uint32, page uint32, pageSize uint32, idAsc bool) (communityBuys []*cidl.GroupBuyingOrderCommunityBuy, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}

	strSql := `
		SELECT
			cby_id,
			order_id,
			grp_id,
			grp_ord_id,
			grp_name,
			tsk_id,
			tsk_title,
			manager_uid,
			manager_name,
			manager_mobile,
			sku_id,
			buy_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version
		FROM
			byo_community_buy
		WHERE
			tsk_id=?
		ORDER BY
			tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, taskId, pageSize, offset)
	if err != nil {
		log.Warnf("query community buy failed.")
		return
	}

	for rows.Next() {
		communityBuy := cidl.NewGroupBuyingOrderCommunityBuy()
		var buyDetail string
		err = rows.Scan(
			&communityBuy.BuyId,
			&communityBuy.OrderId,
			&communityBuy.GroupId,
			&communityBuy.GroupOrderId,
			&communityBuy.GroupName,
			&communityBuy.TaskId,
			&communityBuy.TaskTitle,
			&communityBuy.ManagerUserId,
			&communityBuy.ManagerName,
			&communityBuy.ManagerMobile,
			&communityBuy.SkuId,
			&buyDetail,
			&communityBuy.Count,
			&communityBuy.TotalMarketPrice,
			&communityBuy.TotalGroupBuyingPrice,
			&communityBuy.TotalSettlementPrice,
			&communityBuy.TotalCostPrice,
			&communityBuy.Version,
		)
		if err != nil {
			log.Warnf("scan community buy failed. %s", err)
			return
		}

		err = communityBuy.BuyDetail.FromString(buyDetail)
		if err != nil {
			log.Warnf("init buy detail from string failed. %s", err)
			return
		}

		communityBuys = append(communityBuys, communityBuy)
	}

	return
}

func (m *MallGroupBuyingOrder) CommunityBuyListOrderTaskIdSkuId(taskId uint32, page uint32, pageSize uint32) (communityBuys []*cidl.GroupBuyingOrderCommunityBuy, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strSql := `
		SELECT
			cby_id,
			order_id,
			grp_id,
			grp_ord_id,
			grp_name,
			tsk_id,
			tsk_title,
			manager_uid,
			manager_name,
			manager_mobile,
			sku_id,
			buy_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version
		FROM
			byo_community_buy
		WHERE
			tsk_id=?
		ORDER BY
			tsk_id ASC, sku_id ASC
		LIMIT ? OFFSET ?
	`
	rows, err := m.DB.Query(strSql, taskId, pageSize, offset)
	if err != nil {
		log.Warnf("query community buy failed.")
		return
	}

	for rows.Next() {
		communityBuy := cidl.NewGroupBuyingOrderCommunityBuy()
		var buyDetail string
		err = rows.Scan(
			&communityBuy.BuyId,
			&communityBuy.OrderId,
			&communityBuy.GroupId,
			&communityBuy.GroupOrderId,
			&communityBuy.GroupName,
			&communityBuy.TaskId,
			&communityBuy.TaskTitle,
			&communityBuy.ManagerUserId,
			&communityBuy.ManagerName,
			&communityBuy.ManagerMobile,
			&communityBuy.SkuId,
			&buyDetail,
			&communityBuy.Count,
			&communityBuy.TotalMarketPrice,
			&communityBuy.TotalGroupBuyingPrice,
			&communityBuy.TotalSettlementPrice,
			&communityBuy.TotalCostPrice,
			&communityBuy.Version,
		)
		if err != nil {
			log.Warnf("scan community buy failed. %s", err)
			return
		}

		err = communityBuy.BuyDetail.FromString(buyDetail)
		if err != nil {
			log.Warnf("init buy detail from string failed. %s", err)
			return
		}

		communityBuys = append(communityBuys, communityBuy)
	}

	return
}

func (m *MallGroupBuyingOrder) CommunityBuyListByTaskIds(taskIds []uint32, page uint32, pageSize uint32, idAsc bool) (communityBuys []*cidl.GroupBuyingOrderCommunityBuy, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}

	strSql := `
		SELECT
			cby_id,
			order_id,
			grp_id,
			grp_ord_id,
			grp_name,
			tsk_id,
			tsk_title,
			manager_uid,
			manager_name,
			manager_mobile,
			sku_id,
			buy_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version
		FROM byo_community_buy
		WHERE tsk_id IN (?)
		ORDER BY tsk_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	strSql, args, err := conn.In(strSql, taskIds, pageSize, offset)
	if err != nil {
		log.Warnf("transform sql in array failed. %s", err)
		return
	}

	rows, err := m.DB.Query(strSql, args...)
	if err != nil {
		log.Warnf("query community buy failed.")
		return
	}

	for rows.Next() {
		communityBuy := cidl.NewGroupBuyingOrderCommunityBuy()
		var buyDetail string
		err = rows.Scan(
			&communityBuy.BuyId,
			&communityBuy.OrderId,
			&communityBuy.GroupId,
			&communityBuy.GroupOrderId,
			&communityBuy.GroupName,
			&communityBuy.TaskId,
			&communityBuy.TaskTitle,
			&communityBuy.ManagerUserId,
			&communityBuy.ManagerName,
			&communityBuy.ManagerMobile,
			&communityBuy.SkuId,
			&buyDetail,
			&communityBuy.Count,
			&communityBuy.TotalMarketPrice,
			&communityBuy.TotalGroupBuyingPrice,
			&communityBuy.TotalSettlementPrice,
			&communityBuy.TotalCostPrice,
			&communityBuy.Version,
		)
		if err != nil {
			log.Warnf("scan community buy failed. %s", err)
			return
		}

		err = communityBuy.BuyDetail.FromString(buyDetail)
		if err != nil {
			log.Warnf("init buy detail from string failed. %s", err)
			return
		}

		communityBuys = append(communityBuys, communityBuy)
	}

	return
}

func (m *MallGroupBuyingOrder) CommunityBuyListByTaskIdLineIdsOrderSkuId(taskId uint32, lineIds []uint32) (communityBuys []*cidl.GroupBuyingOrderCommunityBuy, err error) {
	strSql := `
		SELECT
			a.cby_id,
			a.order_id,
			a.grp_id,
			a.grp_ord_id,
			a.grp_name,
			a.tsk_id,
			a.tsk_title,
			a.manager_uid,
			a.manager_name,
			a.manager_mobile,
			a.sku_id,
			a.buy_detail,
			a.count,
			a.total_market_price,
			a.total_group_buying_price,
			a.total_settlement_price,
			a.total_cost_price,
			a.version
		FROM byo_community_buy a, gby_line_community b
		WHERE a.tsk_id = ? 
		AND b.lin_id in (?)
		AND a.grp_id = b.grp_id 
		ORDER BY sku_id ASC
	`
	strSql, args, err := conn.In(strSql, taskId, lineIds)
	if err != nil {
		log.Warnf("transform sql in array failed. %s", err)
		return
	}

	rows, err := m.DB.Query(strSql, args...)
	if err != nil {
		log.Warnf("query community buy failed.")
		return
	}

	for rows.Next() {
		communityBuy := cidl.NewGroupBuyingOrderCommunityBuy()
		var buyDetail string
		err = rows.Scan(
			&communityBuy.BuyId,
			&communityBuy.OrderId,
			&communityBuy.GroupId,
			&communityBuy.GroupOrderId,
			&communityBuy.GroupName,
			&communityBuy.TaskId,
			&communityBuy.TaskTitle,
			&communityBuy.ManagerUserId,
			&communityBuy.ManagerName,
			&communityBuy.ManagerMobile,
			&communityBuy.SkuId,
			&buyDetail,
			&communityBuy.Count,
			&communityBuy.TotalMarketPrice,
			&communityBuy.TotalGroupBuyingPrice,
			&communityBuy.TotalSettlementPrice,
			&communityBuy.TotalCostPrice,
			&communityBuy.Version,
		)
		if err != nil {
			log.Warnf("scan community buy failed. %s", err)
			return
		}

		err = communityBuy.BuyDetail.FromString(buyDetail)
		if err != nil {
			log.Warnf("init buy detail from string failed. %s", err)
			return
		}

		communityBuys = append(communityBuys, communityBuy)
	}

	return
}

// 社群购买团购任务统计
func (m *MallGroupBuyingOrder) GetCommunityBuyTask(groupId uint32, taskId uint32) (buyTask *cidl.GroupBuyingOrderCommunityBuyTask, err error) {
	strSql := `
		SELECT
			tsk_id,
			grp_id,
			task_detail,
			order_count,
			goods_count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version,
			create_time
		FROM
			byo_community_buy_task
		WHERE
			tsk_id=?
			AND grp_id=?
	`
	buyTask = cidl.NewGroupBuyingOrderCommunityBuyTask()
	queryRow, err := m.DB.QueryRow(strSql, taskId, groupId)
	if err != nil {
		log.Warnf("get community buy task failed. %s", err)
		return
	}

	var taskDetail string
	err = queryRow.Scan(
		&buyTask.TaskId,
		&buyTask.GroupId,
		&taskDetail,
		&buyTask.OrderCount,
		&buyTask.GoodsCount,
		&buyTask.TotalMarketPrice,
		&buyTask.TotalGroupBuyingPrice,
		&buyTask.TotalSettlementPrice,
		&buyTask.TotalCostPrice,
		&buyTask.Version,
		&buyTask.CreateTime,
	)
	if err != nil {
		if err != conn.ErrNoRows {
			log.Warnf("query community buy failed. %s", err)
		}
		return
	}

	err = buyTask.TaskDetail.FromString(taskDetail)
	if err != nil {
		log.Warnf("init community buy task task_detail failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) AddCommunityBuyTask(buyTask *cidl.GroupBuyingOrderCommunityBuyTask) (result sql.Result, err error) {
	taskDetail, err := buyTask.TaskDetail.ToString()
	if err != nil {
		log.Warnf("get string task_detail failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO	byo_community_buy_task
			(
				tsk_id,
				grp_id,
				task_detail,
				order_count,
				goods_count,
				total_market_price,
				total_group_buying_price,
				total_settlement_price,
				total_cost_price,
				version
			)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err = m.DB.Exec(
		strSql,
		buyTask.TaskId,
		buyTask.GroupId,
		taskDetail,
		buyTask.OrderCount,
		buyTask.GoodsCount,
		buyTask.TotalMarketPrice,
		buyTask.TotalGroupBuyingPrice,
		buyTask.TotalSettlementPrice,
		buyTask.TotalCostPrice,
		buyTask.Version,
	)
	return
}

func (m *MallGroupBuyingOrder) IncrCommunityBuyTaskValues(
	taskId uint32,
	groupId uint32,
	incrOrderCount uint32,
	incrGoodsCount uint32,
	incrTotalMarketPrice float64,
	incrTotalGroupBuyingPrice float64,
	incrTotalSettlementPrice float64,
	incrTotalCostPrice float64,
) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_community_buy_task
		SET
			order_count=order_count+?,
			goods_count=goods_count+?,
			total_market_price=total_market_price+?,
			total_group_buying_price=total_group_buying_price+?,
			total_settlement_price=total_settlement_price+?,
			total_cost_price=total_cost_price+?
		WHERE tsk_id=? AND grp_id=?
	`
	result, err = m.DB.Exec(
		strSql,
		incrOrderCount,
		incrGoodsCount,
		incrTotalMarketPrice,
		incrTotalGroupBuyingPrice,
		incrTotalSettlementPrice,
		incrTotalCostPrice,
		taskId,
		groupId,
	)
	return
}

func (m *MallGroupBuyingOrder) DecrCommunityBuyTaskValues(
	taskId uint32,
	groupId uint32,
	decrOrderCount uint32,
	decrGoodsCount uint32,
	decrTotalMarketPrice float64,
	decrTotalGroupBuyingPrice float64,
	decrTotalSettlementPrice float64,
	decrTotalCostPrice float64,
) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_community_buy_task
		SET
			order_count=order_count-?,
			goods_count=goods_count-?,
			total_market_price=total_market_price-?,
			total_group_buying_price=total_group_buying_price-?,
			total_settlement_price=total_settlement_price-?,
			total_cost_price=total_cost_price-?
		WHERE tsk_id=? AND grp_id=?
	`
	result, err = m.DB.Exec(
		strSql,
		decrOrderCount,
		decrGoodsCount,
		decrTotalMarketPrice,
		decrTotalGroupBuyingPrice,
		decrTotalSettlementPrice,
		decrTotalCostPrice,
		taskId,
		groupId,
	)
	return
}

func (m *MallGroupBuyingOrder) CommunityBuyTaskCount(groupId uint32) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM byo_community_buy_task WHERE grp_id=?`
	err = m.DB.Get(&count, strSql, groupId)
	return
}

func (m *MallGroupBuyingOrder) CommunityBuyTaskList(groupId uint32, page uint32, pageSize uint32, createTimeAsc bool) (buyTasks []*cidl.GroupBuyingOrderCommunityBuyTask, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == createTimeAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			tsk_id,
			grp_id,
			task_detail,
			order_count,
			goods_count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version,
			create_time
		FROM byo_community_buy_task
		WHERE
			grp_id=?
		ORDER BY create_time %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, groupId, pageSize, offset)
	if err != nil {
		log.Warnf("query community buy task list failed. %s", err)
		return
	}

	for rows.Next() {
		communityBuyTask := cidl.NewGroupBuyingOrderCommunityBuyTask()
		var taskDetail string
		err = rows.Scan(
			&communityBuyTask.TaskId,
			&communityBuyTask.GroupId,
			&taskDetail,
			&communityBuyTask.OrderCount,
			&communityBuyTask.GoodsCount,
			&communityBuyTask.TotalMarketPrice,
			&communityBuyTask.TotalGroupBuyingPrice,
			&communityBuyTask.TotalSettlementPrice,
			&communityBuyTask.TotalCostPrice,
			&communityBuyTask.Version,
			&communityBuyTask.CreateTime,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query community buy task failed. %s", err)
			}
			return
		}

		err = communityBuyTask.TaskDetail.FromString(taskDetail)
		if err != nil {
			log.Warnf("init community buy task task_detail failed. %s", err)
			return
		}

		buyTasks = append(buyTasks, communityBuyTask)
	}

	return
}

func (m *MallGroupBuyingOrder) AddIndent(indent *cidl.GroupBuyingIndent) (result sql.Result, err error) {
	taskBrief, err := indent.TasksBrief.ToString()
	if err != nil {
		log.Warnf("get string tasks_brief failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO byo_indent
			(
				idt_id,
				org_id,
				tasks_brief,
				state,
				excel_url,
				version
			)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err = m.DB.Exec(strSql,
		indent.IndentId,
		indent.OrganizationId,
		taskBrief,
		indent.State,
		indent.ExcelUrl,
		indent.Version)
	return
}

func (m *MallGroupBuyingOrder) UpdateIndentState(indentId string, state cidl.IndentState) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_indent
		SET state=?
		WHERE idt_id=?
	`
	result, err = m.DB.Exec(strSql, state, indentId)
	return
}

func (m *MallGroupBuyingOrder) UpdateIndentStateAndExcelUrl(indentId string, state cidl.IndentState, excelUrl string) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_indent
		SET
			state=?,
			excel_url=?
		WHERE idt_id=?
	`
	result, err = m.DB.Exec(strSql, state, excelUrl, indentId)
	return
}

func (m *MallGroupBuyingOrder) GetIndent(indentId string) (indent *cidl.GroupBuyingIndent, err error) {
	defer func() {
		if err != nil {
			indent = nil
		}
	}()

	indent = cidl.NewGroupBuyingIndent()
	var tasksBrief string
	strSql := `
		SELECT
			idt_id,
			org_id,
			tasks_brief,
			state,
			excel_url,
			version,
			create_time
		FROM byo_indent
		WHERE idt_id=?
	`
	queryRow, err := m.DB.QueryRow(strSql, indentId)
	if err != nil {
		log.Warnf("get query row failed. %s", err)
		return
	}

	err = queryRow.Scan(
		&indent.IndentId,
		&indent.OrganizationId,
		&tasksBrief,
		&indent.State,
		&indent.ExcelUrl,
		&indent.Version,
		&indent.CreateTime,
	)
	if err != nil {
		if err != conn.ErrNoRows {
			log.Warnf("query indent failed. %s", err)
		}
		return
	}

	err = indent.TasksBrief.FromString(tasksBrief)
	if err != nil {
		log.Warnf("init indent tasks_brief from string failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) GetIndentExcelUrl(indentId string) (excelUrl string, err error) {
	strSql := `SELECT excel_url FROM byo_indent WHERE idt_id=?`
	err = m.DB.Get(&excelUrl, strSql, indentId)
	return
}

func (m *MallGroupBuyingOrder) IndentCount(organizationId uint32) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM byo_indent WHERE org_id=?`
	err = m.DB.Get(&count, strSql, organizationId)
	return
}

func (m *MallGroupBuyingOrder) IndentList(organizationId uint32, page uint32, pageSize uint32, createTimeAsc bool) (indents []*cidl.GroupBuyingIndent, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == createTimeAsc {
		strOrderBy = "DESC"
	}

	strSql := `
		SELECT
			idt_id,
			org_id,
			tasks_brief,
			state,
			excel_url,
			version,
			create_time
		FROM byo_indent
		WHERE org_id=?
		ORDER BY create_time %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, pageSize, offset)
	if err != nil {
		log.Warnf("query indent list failed. %s", err)
		return
	}

	for rows.Next() {
		indent := cidl.NewGroupBuyingIndent()
		var tasksBrief string
		err = rows.Scan(
			&indent.IndentId,
			&indent.OrganizationId,
			&tasksBrief,
			&indent.State,
			&indent.ExcelUrl,
			&indent.Version,
			&indent.CreateTime,
		)

		if err != nil {
			log.Warnf("scan indent failed. %s", err)
			return
		}

		err = indent.TasksBrief.FromString(tasksBrief)
		if err != nil {
			log.Warnf("init indent tasks_brief from string failed. %s", err)
			return
		}

		indents = append(indents, indent)
	}

	return
}

func (m *MallGroupBuyingOrder) AddIndentStatistics(indentStatistics *cidl.GroupBuyingIndentStatistics) (result sql.Result, err error) {
	taskContent, err := indentStatistics.TaskContent.ToString()
	if err != nil {
		log.Warnf("get string task content failed. %s", err)
		return
	}

	statisticsResult, err := indentStatistics.Result.ToString()
	if err != nil {
		log.Warnf("get string indent statistics result failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO byo_indent_statistics
			(
				idt_id,
				tsk_id,
				task_content,
				result,
				version
			)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err = m.DB.Exec(strSql,
		indentStatistics.IndentId,
		indentStatistics.TaskId,
		taskContent,
		statisticsResult,
		indentStatistics.Version)

	return
}

func (m *MallGroupBuyingOrder) UpdateIndentStatisticsResult(indentId string, taskId uint32, indentStatisticsResult *cidl.GroupBuyingIndentStatisticsResultType, version int) (result sql.Result, err error) {
	statisticsResult, err := indentStatisticsResult.ToString()
	if err != nil {
		log.Warnf("get string indent statistics result failed. %s", err)
		return
	}

	strSql := `
		UPDATE byo_indent_statistics
		SET
			result=?,
			version=?
		WHERE idt_id=? AND tsk_id=?
	`
	result, err = m.DB.Exec(strSql, statisticsResult, version, indentId, taskId)
	return
}

func (m *MallGroupBuyingOrder) GetIndentStatistics(indentId string, taskId uint32) (indentStatistics *cidl.GroupBuyingIndentStatistics, err error) {
	defer func() {
		if err != nil {
			indentStatistics = nil
		}
	}()

	indentStatistics = cidl.NewGroupBuyingIndentStatistics()
	var (
		taskContent string
		result      string
	)
	strSql := `
		SELECT
			idt_id,
			tsk_id,
			task_content,
			result,
			version,
			create_time
		FROM byo_indent_statistics
		WHERE idt_id=? AND tsk_id=?
	`
	queryRow, err := m.DB.QueryRow(strSql, indentId, taskId)
	if err != nil {
		log.Warnf("get query row failed. %s", err)
		return
	}

	err = queryRow.Scan(
		&indentStatistics.IndentId,
		&indentStatistics.TaskId,
		&taskContent,
		&result,
		&indentStatistics.Version,
		&indentStatistics.CreateTime,
	)

	if err != nil {
		if err != conn.ErrNoRows {
			log.Warnf("query indent statistics failed. %s", err)
		}
		return
	}

	err = indentStatistics.TaskContent.FromString(taskContent)
	if err != nil {
		log.Warnf("init indent_statistics task_content from string failed. %s", err)
		return
	}

	err = indentStatistics.Result.FromString(result)
	if err != nil {
		log.Warnf("init indent_statistics result from string failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) IndentStatisticsCount(indentId string) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM agc_indent_statistics WHERE idt_id=?`
	err = m.DB.Get(&count, strSql, indentId)
	return
}

func (m *MallGroupBuyingOrder) IndentStatisticsAll(indentId string, idAsc bool) (list []*cidl.GroupBuyingIndentStatistics, err error) {
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}

	strSql := `
		SELECT
			idt_id,
			tsk_id,
			task_content,
			result,
			version,
			create_time
		FROM byo_indent_statistics
		WHERE idt_id=?
		ORDER BY tsk_id %s
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, indentId)
	if err != nil {
		log.Warnf("query indent statistics list failed. %s", err)
		return
	}

	for rows.Next() {
		indentStatistics := cidl.NewGroupBuyingIndentStatistics()
		var (
			taskContent string
			result      string
		)
		err = rows.Scan(
			&indentStatistics.IndentId,
			&indentStatistics.TaskId,
			&taskContent,
			&result,
			&indentStatistics.Version,
			&indentStatistics.CreateTime,
		)

		if err != nil {
			log.Warnf("scan indent statistics failed. %s", err)
			return
		}

		err = indentStatistics.TaskContent.FromString(taskContent)
		if err != nil {
			log.Warnf("init indent_statistics task_content from string failed. %s", err)
			return
		}

		err = indentStatistics.Result.FromString(result)
		if err != nil {
			log.Warnf("init indent_statistics result from string failed. %s", err)
			return
		}

		list = append(list, indentStatistics)
	}

	return
}

func (m *MallGroupBuyingOrder) AddLine(line *cidl.GroupBuyingLine) (result sql.Result, err error) {
	strSql := `
		INSERT INTO byo_line
			(
				lin_id,
				org_id,
				org_name,
				name,
				community_count
			)
		VALUES
			(
				:lin_id,
				:org_id,
				:org_name,
				:name,
				:community_count
			)
	`
	result, err = m.DB.NamedExec(strSql, line)
	return
}

func (m *MallGroupBuyingOrder) GetLine(lineId uint32) (line *cidl.GroupBuyingLine, err error) {
	defer func() {
		if err != nil {
			line = nil
		}
	}()

	line = cidl.NewGroupBuyingLine()
	strSql := `
		SELECT
			lin_id,
			org_id,
			org_name,
			name,
			community_count,
			create_time
		FROM gby_line
		WHERE lin_id=?
	`

	err = m.DB.Get(line, strSql, lineId)

	return
}

func (m *MallGroupBuyingOrder) LineCount(organizationId uint32) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM gby_line WHERE org_id=?`
	err = m.DB.Get(&count, strSql, organizationId)
	return
}

func (m *MallGroupBuyingOrder) LineList(organizationId uint32, page uint32, pageSize uint32, idAsc bool) (lines []*cidl.GroupBuyingLine, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}

	strSql := `
		SELECT
			lin_id,
			org_id,
			org_name,
			name,
			community_count,
			create_time
		FROM gby_line
		WHERE org_id=?
		ORDER BY lin_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, pageSize, offset)
	if err != nil {
		log.Warnf("query line list failed. %s", err)
		return
	}

	for rows.Next() {
		line := cidl.NewGroupBuyingLine()
		err = rows.StructScan(line)
		if err != nil {
			log.Warnf("scan line failed. %s", err)
			return
		}

		lines = append(lines, line)
	}

	return
}

func (m *MallGroupBuyingOrder) AddLineCommunity(lineCommunity *cidl.GroupBuyingLineCommunity) (result sql.Result, err error) {
	strSql := `
		INSERT INTO gby_line_community
			(
				grp_id,
				lin_id,
				lin_name,
				grp_name,
				manager_uid,
				manager_name,
				manager_mobile,
				org_id
			)
		VALUES
			(
				:grp_id,
				:lin_id,
				:lin_name,
				:grp_name,
				:manager_uid,
				:manager_name,
				:manager_mobile,
				:org_id
			)
	`
	result, err = m.DB.NamedExec(strSql, lineCommunity)
	return
}

func (m *MallGroupBuyingOrder) GetLineCommunity(groupId uint32) (lineCommunity *cidl.GroupBuyingLineCommunity, err error) {
	defer func() {
		if err != nil {
			lineCommunity = nil
		}

	}()

	lineCommunity = cidl.NewGroupBuyingLineCommunity()

	strSql := `
		SELECT
			grp_id,
			lin_id,
			lin_name,
			grp_name,
			manager_uid,
			manager_name,
			manager_mobile,
			org_id,
			create_time
		FROM gby_line_community
		WHERE grp_id=?
    `

	err = m.DB.Get(lineCommunity, strSql, groupId)
	return
}

func (m *MallGroupBuyingOrder) LineCommunityCount(organizationId uint32, lineId uint32) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM gby_line_community WHERE org_id=? AND lin_id=?`
	err = m.DB.Get(&count, strSql, organizationId, lineId)
	return
}

func (m *MallGroupBuyingOrder) LineCommunityList(organizationId uint32, lineId uint32, page uint32, pageSize uint32, idAsc bool) (communities []*cidl.GroupBuyingLineCommunity, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}

	strSql := `
		SELECT
			grp_id,
			lin_id,
			lin_name,
			grp_name,
			manager_uid,
			manager_name,
			manager_mobile,
			org_id,
			create_time
		FROM gby_line_community
		WHERE org_id=? AND lin_id=?
		ORDER BY grp_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, lineId, pageSize, offset)
	if err != nil {
		log.Warnf("query line community list failed. %s", err)
		return
	}

	for rows.Next() {
		lineCommunity := cidl.NewGroupBuyingLineCommunity()
		err = rows.StructScan(lineCommunity)
		if err != nil {
			log.Warnf("struct scan line community failed. %s", err)
			return
		}
		communities = append(communities, lineCommunity)
	}

	return
}

func (m *MallGroupBuyingOrder) LineCommunityUnbindCount(organizationId uint32) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM gby_line_community WHERE org_id=? AND lin_id=0`
	err = m.DB.Get(&count, strSql, organizationId)
	return
}

func (m *MallGroupBuyingOrder) LineCommunityUnbindList(organizationId uint32, page uint32, pageSize uint32, idAsc bool) (communities []*cidl.GroupBuyingLineCommunity, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == idAsc {
		strOrderBy = "DESC"
	}

	strSql := `
		SELECT
			grp_id,
			lin_id,
			lin_name,
			grp_name,
			manager_uid,
			manager_name,
			manager_mobile,
			org_id,
			create_time
		FROM gby_line_community
		WHERE org_id=? AND lin_id=0
		ORDER BY grp_id %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, pageSize, offset)
	if err != nil {
		log.Warnf("query line community list failed. %s", err)
		return
	}

	for rows.Next() {
		lineCommunity := cidl.NewGroupBuyingLineCommunity()
		err = rows.StructScan(lineCommunity)
		if err != nil {
			log.Warnf("struct scan line community failed. %s", err)
			return
		}
		communities = append(communities, lineCommunity)
	}

	return
}

func (m *MallGroupBuyingOrder) AddSend(send *cidl.GroupBuyingSend) (result sql.Result, err error) {
	tasksBrief, err := send.TasksBrief.ToString()
	if err != nil {
		log.Warnf("get string tasks brief failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO byo_send
			(
				snd_id,
				tasks_brief,
				tasks_detail,
				org_id,
				org_name,
				state,
				excel_url,
				version
			)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err = m.DB.Exec(strSql,
		send.SendId,
		tasksBrief,
		send.TasksDetail,
		send.OrganizationId,
		send.OrganizationName,
		send.State,
		send.ExcelUrl,
		send.Version)

	return
}

func (m *MallGroupBuyingOrder) UpdateSendState(sendId string, state cidl.GroupBuyingSendState) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_send
		SET
			state=?
		WHERE snd_id=?
	`
	result, err = m.DB.Exec(strSql, state, sendId)
	return
}

func (m *MallGroupBuyingOrder) UpdateSendStatisticResult(sendId string, state cidl.GroupBuyingSendState, excelUrl string, version uint32) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_send
		SET
			state=?,
			excel_url=?,
			version=?
		WHERE
			snd_id=?
	`
	result, err = m.DB.Exec(strSql, state, excelUrl, version, sendId)
	return
}

func (m *MallGroupBuyingOrder) GetSend(sendId string) (send *cidl.GroupBuyingSend, err error) {
	defer func() {
		if err != nil {
			send = nil
		}
	}()

	send = cidl.NewGroupBuyingSend()
	var tasksBrief string
	strSql := `
		SELECT
			snd_id,
			tasks_brief,
			org_id,
			org_name,
			state,
			excel_url,
			version,
			create_time
		FROM byo_send
		WHERE snd_id=?
	`
	queryRow, err := m.DB.QueryRow(strSql, sendId)
	if err != nil {
		log.Warnf("get query row failed. %s", err)
		return
	}

	err = queryRow.Scan(
		&send.SendId,
		&tasksBrief,
		&send.OrganizationId,
		&send.OrganizationName,
		&send.State,
		&send.ExcelUrl,
		&send.Version,
		&send.CreateTime,
	)

	if err != nil {
		log.Warnf("query send failed. %s", err)
		return
	}

	err = send.TasksBrief.FromString(tasksBrief)
	if err != nil {
		log.Warnf("init tasks_brief failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) GetSendExcelUrl(sendId string) (excelUrl string, err error) {
	strSql := `SELECT excel_url FROM byo_send WHERE snd_id=?`
	err = m.DB.Get(&excelUrl, strSql, sendId)
	return
}

func (m *MallGroupBuyingOrder) SendCount(organizationId uint32) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM byo_send WHERE org_id=?`
	err = m.DB.Get(&count, strSql, organizationId)
	return
}

func (m *MallGroupBuyingOrder) SendList(organizationId uint32, page uint32, pageSize uint32, createTimeAsc bool) (sends []*cidl.GroupBuyingSend, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == createTimeAsc {
		strOrderBy = "DESC"
	}

	strSql := `
		SELECT
			snd_id,
			tasks_brief,
			org_id,
			org_name,
			state,
			excel_url,
			version,
			create_time
		FROM byo_send
		WHERE org_id=?
		ORDER BY create_time %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, organizationId, pageSize, offset)
	if err != nil {
		log.Warnf("query send list failed. %s", err)
		return
	}

	for rows.Next() {
		send := cidl.NewGroupBuyingSend()
		var tasksBrief string
		err = rows.Scan(
			&send.SendId,
			&tasksBrief,
			&send.OrganizationId,
			&send.OrganizationName,
			&send.State,
			&send.ExcelUrl,
			&send.Version,
			&send.CreateTime,
		)

		if err != nil {
			log.Warnf("scan send failed. %s", err)
			return
		}

		err = send.TasksBrief.FromString(tasksBrief)
		if err != nil {
			log.Warnf("init tasks_brief failed. %s", err)
			return
		}

		sends = append(sends, send)
	}

	return
}

func (m *MallGroupBuyingOrder) AddSendLine(sendLine *cidl.GroupBuyingSendLine) (result sql.Result, err error) {
	statistics, err := sendLine.Statistics.ToString()
	if err != nil {
		log.Warnf("get string send line statistics failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO byo_send_line
			(
				snd_id,
				lin_id,
				lin_name,
				org_id,
				org_name,
				community_count,
				settlement_amount,
				statistics,
				send_time,
				version
			)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err = m.DB.Exec(strSql,
		sendLine.SendId,
		sendLine.LineId,
		sendLine.LineName,
		sendLine.OrganizationId,
		sendLine.OrganizationName,
		sendLine.CommunityCount,
		sendLine.SettlementAmount,
		statistics,
		sendLine.SendTime,
		sendLine.Version,
	)
	return
}

func (m *MallGroupBuyingOrder) GetSendLine(sendId string, lineId uint32) (sendLine *cidl.GroupBuyingSendLine, err error) {
	defer func() {
		if err != nil {
			sendLine = nil
		}
	}()

	sendLine = cidl.NewGroupBuyingSendLine()
	var statistics string
	strSql := `
		SELECT
			snd_id,
			lin_id,
			lin_name,
			org_id,
			org_name,
			community_count,
			settlement_amount,
			statistics,
			send_time,
			version,
			create_time
		FROM byo_send_line
		WHERE snd_id=? AND lin_id=?
	`

	queryRow, err := m.DB.QueryRow(strSql, sendId, lineId)
	if err != nil {
		log.Warnf("get query row failed. %s", err)
		return
	}

	err = queryRow.Scan(
		&sendLine.SendId,
		&sendLine.LineId,
		&sendLine.LineName,
		&sendLine.OrganizationId,
		&sendLine.OrganizationName,
		&sendLine.CommunityCount,
		&sendLine.SettlementAmount,
		&statistics,
		&sendLine.SendTime,
		&sendLine.Version,
		&sendLine.CreateTime,
	)

	if err != nil {
		if err != conn.ErrNoRows {
			log.Warnf("query send line failed. %s", err)
		}
		return
	}

	err = sendLine.Statistics.FromString(statistics)
	if err != nil {
		log.Warnf("init send line statistics from string failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) AddSendCommunity(sendCommunity *cidl.GroupBuyingSendCommunity) (result sql.Result, err error) {
	statistics, err := sendCommunity.Statistics.ToString()
	if err != nil {
		log.Warnf("get string send community statistics failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO byo_send_community
			(
				snd_id,
				grp_id,
				grp_name,
				grp_address,
				grp_manager_uid,
				grp_manager_name,
				grp_manager_mobile,
				org_id,
				org_name,
				org_address,
				org_manager_uid,
				org_manager_name,
				author_uid,
				author_name,
				settlement_amount,
				statistics,
				send_time,
				version
			)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err = m.DB.Exec(strSql,
		sendCommunity.SendId,
		sendCommunity.GroupId,
		sendCommunity.GroupName,
		sendCommunity.GroupAddress,
		sendCommunity.GroupManagerUid,
		sendCommunity.GroupManagerName,
		sendCommunity.GroupManagerMobile,
		sendCommunity.OrganizationId,
		sendCommunity.OrganizationName,
		sendCommunity.OrganizationAddress,
		sendCommunity.OrganizationManagerUid,
		sendCommunity.OrganizationManagerName,
		sendCommunity.AuthorUid,
		sendCommunity.AuthorName,
		sendCommunity.SettlementAmount,
		statistics,
		sendCommunity.SendTime,
		sendCommunity.Version,
	)

	return
}

func (m *MallGroupBuyingOrder) GetSendCommunity(sendId string, groupId uint32) (sendCommunity *cidl.GroupBuyingSendCommunity, err error) {
	defer func() {
		if err != nil {
			sendCommunity = nil
		}
	}()

	sendCommunity = cidl.NewGroupBuyingSendCommunity()
	var statistics string
	strSql := `
		SELECT
			snd_id,
			grp_id,
			grp_name,
			grp_address,
			grp_manager_uid,
			grp_manager_name,
			grp_manager_mobile,
			org_id,
			org_name,
			org_address,
			org_manager_uid,
			org_manager_name,
			author_uid,
			author_name,
			settlement_amount,
			statistics,
			send_time,
			version,
			create_time
		FROM byo_send_community
		WHERE snd_id=? AND grp_id=?
	`

	queryRow, err := m.DB.QueryRow(strSql, sendId, groupId)
	if err != nil {
		log.Warnf("get query row failed. %s", err)
		return
	}

	err = queryRow.Scan(
		&sendCommunity.SendId,
		&sendCommunity.GroupId,
		&sendCommunity.GroupName,
		&sendCommunity.GroupAddress,
		&sendCommunity.GroupManagerUid,
		&sendCommunity.GroupManagerName,
		&sendCommunity.GroupManagerMobile,
		&sendCommunity.OrganizationId,
		&sendCommunity.OrganizationName,
		&sendCommunity.OrganizationAddress,
		&sendCommunity.OrganizationManagerUid,
		&sendCommunity.OrganizationManagerName,
		&sendCommunity.AuthorUid,
		&sendCommunity.AuthorName,
		&sendCommunity.SettlementAmount,
		&statistics,
		&sendCommunity.SendTime,
		&sendCommunity.Version,
		&sendCommunity.CreateTime,
	)

	if err != nil {
		log.Warnf("query send community failed. %s", err)
		return
	}

	err = sendCommunity.Statistics.FromString(statistics)
	if err != nil {
		log.Warnf("init send community statistics from string failed. %s", err)
		return
	}

	return
}

// 库存
func (m *MallGroupBuyingOrder) AddInventories(inventories []*cidl.GroupBuyingOrderInventory) (result sql.Result, err error) {
	strSql := `
		INSERT INTO byo_inventory
			(
				tsk_id,
				sku_id,
				total,
				surplus
			)
		VALUES
			%s
	`
	var args []interface{}
	var sliceStrValue []string
	for _, inventory := range inventories {
		sliceStrValue = append(sliceStrValue, "(?, ?, ?, ?)")
		args = append(args, inventory.TaskId)
		args = append(args, inventory.SkuId)
		args = append(args, inventory.Total)
		args = append(args, inventory.Surplus)
	}

	strValues := strings.Join(sliceStrValue, ",")
	strSql = fmt.Sprintf(strSql, strValues)
	result, err = m.DB.Exec(strSql, args...)
	return
}

func (m *MallGroupBuyingOrder) GetInventory(taskId uint32, skuId string) (inventory *cidl.GroupBuyingOrderInventory, err error) {
	defer func() {
		if err != nil {
			inventory = nil
			return
		}
	}()

	inventory = cidl.NewGroupBuyingOrderInventory()
	strSql := `
		SELECT
			tsk_id,
			sku_id,
			total,
			sales,
			surplus
		FROM
			byo_inventory
		WHERE
			tsk_id=? AND sku_id=?
	`
	err = m.DB.Get(inventory, strSql, taskId, skuId)
	return
}

func (m *MallGroupBuyingOrder) GetInventories(taskId uint32) (inventories []*cidl.GroupBuyingOrderInventory, err error) {
	strSql := `
		SELECT
			tsk_id,
			sku_id,
			total,
			sales,
			surplus
		FROM
			byo_inventory
		WHERE
			tsk_id=?
	`
	err = m.DB.Select(&inventories, strSql, taskId)
	if err != nil {
		log.Warnf("query inventories failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) GetInventoriesBySkuIds(taskId uint32, skuIds []string) (inventories []*cidl.GroupBuyingOrderInventory, err error) {
	strSql := `
		SELECT
			tsk_id,
			sku_id,
			total,
			sales,
			surplus
		FROM
			byo_inventory
		WHERE
			tsk_id=? AND sku_id IN (?)
	`
	strSql, args, err := conn.In(strSql, taskId, skuIds)
	if err != nil {
		log.Warnf("transform sql in array failed. %s", err)
		return
	}

	err = m.DB.Select(&inventories, strSql, args...)
	if err != nil {
		log.Warnf("query inventories failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) DeleteInventory(taskId uint32) (result sql.Result, err error) {
	strSql := `
		DELETE FROM byo_inventory WHERE tsk_id=?
	`
	result, err = m.DB.Exec(strSql, taskId)
	return
}

// 事务锁库存
func (m *MallGroupBuyingOrder) TxLockInventorySurplus(tx *conn.Tx, taskId uint32, skuIds []string) (result sql.Result, err error) {
	strSql := `
		SELECT surplus FROM byo_inventory WHERE tsk_id=? AND sku_id IN (?) FOR UPDATE
	`
	strSql, args, err := conn.In(strSql, taskId, skuIds)
	result, err = tx.Exec(strSql, args...)
	return
}

// 事务锁库存(taskIds and skuIds)
func (m *MallGroupBuyingOrder) TxLockInventorySurplusByTaskIdsSkuIds(tx *conn.Tx, taskId uint32, skuIds []string) (result sql.Result, err error) {
	strSql := `
		SELECT surplus FROM byo_inventory WHERE tsk_id=? AND sku_id IN (?) FOR UPDATE
	`
	strSql, args, err := conn.In(strSql, taskId, skuIds)
	result, err = tx.Exec(strSql, args...)
	return
}

// 减库存
func (m *MallGroupBuyingOrder) SubtractInventorySurplus(taskId uint32, skuId string, count uint32) (success bool, err error) {
	strSql := `
		UPDATE byo_inventory
		SET
			sales=sales+?,
			surplus=surplus-?
		WHERE
			tsk_id=? AND sku_id=? AND surplus>=?
	`
	result, err := m.DB.Exec(strSql, count, count, taskId, skuId, count)
	if err != nil {
		log.Warnf("update byo_inventory failed. %s", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("get rows affected failed. %s", err)
		return
	}

	if rowsAffected == 0 {
		log.Warnf("no rows affected. %s", err)
		return
	}

	success = true

	return
}

func (m *MallGroupBuyingOrder) TxSubtractInventorySurplus(tx *conn.Tx, taskId uint32, skuId string, count uint32) (success bool, err error) {
	strSql := `
		UPDATE byo_inventory
		SET
			sales=sales+?,
			surplus=surplus-?
		WHERE
			tsk_id=?
			AND sku_id=?
			AND surplus>=?
	`
	result, err := tx.Exec(strSql, count, count, taskId, skuId, count)
	if err != nil {
		log.Warnf("update byo_inventory failed. %s", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("get rows affected failed. %s", err)
		return
	}

	if rowsAffected == 0 {
		log.Warnf("no rows affected. %s", err)
		return
	}

	success = true

	return
}

// 增加库存
func (m *MallGroupBuyingOrder) TxAddtractInventorySurplus(tx *conn.Tx, taskId uint32, skuId string, count uint32) (success bool, err error) {
	strSql := `
		UPDATE byo_inventory
		SET
			sales=sales-?,
			surplus=surplus+?
		WHERE
			tsk_id=?
			AND sku_id=?
			AND surplus <= total
	`
	result, err := tx.Exec(strSql, count, count, taskId, skuId)
	if err != nil {
		log.Warnf("update byo_inventory failed. %s", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("get rows affected failed. %s", err)
		return
	}

	if rowsAffected == 0 {
		log.Warnf("no rows affected. %s", err)
		return
	}

	success = true

	return
}


// 添加购物车
func (m *MallGroupBuyingOrder) AddCommunityCart(cart *cidl.GroupBuyingOrderCommunityCart) (result sql.Result, err error) {
	buyDetail, err := cart.BuyDetail.ToString()
	if err != nil {
		log.Warnf("get string cart buy detail failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO byo_community_cart
			(
				ccr_id,
				grp_id,
				tsk_id,
				tsk_title,
				sku_id,
				buy_detail,
				count,
				total_market_price,
				total_group_buying_price,
				total_settlement_price,
				total_cost_price,
				version
			)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err = m.DB.Exec(strSql,
		cart.CartId,
		cart.GroupId,
		cart.TaskId,
		cart.TaskTitle,
		cart.SkuId,
		buyDetail,
		cart.Count,
		cart.TotalMarketPrice,
		cart.TotalGroupBuyingPrice,
		cart.TotalSettlementPrice,
		cart.TotalCostPrice,
		cart.Version,
	)

	return
}

func (m *MallGroupBuyingOrder) ChangeCommunityCartBuyCount(groupId uint32, cartId string, count uint32) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_community_cart
		SET
			count=?
		WHERE grp_id=? AND ccr_id=?
	`
	result, err = m.DB.Exec(strSql, count, groupId, cartId)
	return
}

func (m *MallGroupBuyingOrder) GetCommunityCartCount(groupId uint32) (count uint32, err error) {
	strSql := `SELECT COUNT(*) FROM byo_community_cart WHERE grp_id=?`
	err = m.DB.Get(&count, strSql, groupId)
	return
}

func (m *MallGroupBuyingOrder) CommunityCartList(groupId uint32, page uint32, pageSize uint32, createTimeAsc bool) (carts []*cidl.GroupBuyingOrderCommunityCart, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == createTimeAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			ccr_id,
			grp_id,
			tsk_id,
			tsk_title,
			sku_id,
			buy_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version,
			create_time
		FROM
			byo_community_cart
		WHERE
			grp_id=?
		ORDER BY create_time %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, groupId, pageSize, offset)
	if err != nil {
		log.Warnf("query task list failed. %s", err)
		return
	}

	for rows.Next() {
		cart := cidl.NewGroupBuyingOrderCommunityCart()
		var buyDetail string
		err = rows.Scan(
			&cart.CartId,
			&cart.GroupId,
			&cart.TaskId,
			&cart.TaskTitle,
			&cart.SkuId,
			&buyDetail,
			&cart.Count,
			&cart.TotalMarketPrice,
			&cart.TotalGroupBuyingPrice,
			&cart.TotalSettlementPrice,
			&cart.TotalCostPrice,
			&cart.Version,
			&cart.CreateTime,
		)
		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query community cart failed. %s", err)
			}
			return
		}

		err = cart.BuyDetail.FromString(buyDetail)
		if err != nil {
			log.Warnf("init buy detail from string failed. %s", err)
			return
		}

		carts = append(carts, cart)
	}

	return
}

func (m *MallGroupBuyingOrder) GetCommunityCart(groupId uint32, cartId string) (cart *cidl.GroupBuyingOrderCommunityCart, err error) {
	defer func() {
		if err != nil {
			cart = nil
			return
		}
	}()

	strSql := `
		SELECT
			ccr_id,
			grp_id,
			tsk_id,
			tsk_title,
			sku_id,
			buy_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version,
			create_time
		FROM
			byo_community_cart
		WHERE
			grp_id=? AND ccr_id=?
	`
	cart = cidl.NewGroupBuyingOrderCommunityCart()
	queryRow, err := m.DB.QueryRow(strSql, groupId, cartId)
	if err != nil {
		log.Warnf("get community cart failed. %s", err)
		return
	}
	var buyDetail string
	err = queryRow.Scan(
		&cart.CartId,
		&cart.GroupId,
		&cart.TaskId,
		&cart.TaskTitle,
		&cart.SkuId,
		&buyDetail,
		&cart.Count,
		&cart.TotalMarketPrice,
		&cart.TotalGroupBuyingPrice,
		&cart.TotalSettlementPrice,
		&cart.TotalCostPrice,
		&cart.Version,
		&cart.CreateTime,
	)
	if err != nil {
		if err != conn.ErrNoRows {
			log.Warnf("query community cart failed. %s", err)
		}
		return
	}

	err = cart.BuyDetail.FromString(buyDetail)
	if err != nil {
		log.Warnf("init community cart failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) GetCommunityCarts(groupId uint32, cartsId []string) (carts []*cidl.GroupBuyingOrderCommunityCart, err error) {
	strSql := `
		SELECT
			ccr_id,
			grp_id,
			tsk_id,
			tsk_title,
			sku_id,
			buy_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version,
			create_time
		FROM
			byo_community_cart
		WHERE
			grp_id=? AND ccr_id IN (?)
	`
	strSql, args, err := conn.In(strSql, groupId, cartsId)
	if err != nil {
		log.Warnf("transform sql in array failed. %s", err)
		return
	}

	rows, err := m.DB.Query(strSql, args...)
	if err != nil {
		log.Warnf("query community carts failed. %s", err)
		return
	}

	for rows.Next() {
		cart := cidl.NewGroupBuyingOrderCommunityCart()
		var buyDetail string
		err = rows.Scan(
			&cart.CartId,
			&cart.GroupId,
			&cart.TaskId,
			&cart.TaskTitle,
			&cart.SkuId,
			&buyDetail,
			&cart.Count,
			&cart.TotalMarketPrice,
			&cart.TotalGroupBuyingPrice,
			&cart.TotalSettlementPrice,
			&cart.TotalCostPrice,
			&cart.Version,
			&cart.CreateTime,
		)
		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query community cart failed. %s", err)
			}
			return
		}

		err = cart.BuyDetail.FromString(buyDetail)
		if err != nil {
			log.Warnf("init buy detail from string failed. %s", err)
			return
		}

		carts = append(carts, cart)
	}

	return
}

func (m *MallGroupBuyingOrder) DeleteCommunityCarts(groupId uint32, cartIds []string) (result sql.Result, err error) {
	strSql := `DELETE FROM byo_community_cart WHERE grp_id=? AND ccr_id IN (?)`
	strSql, args, err := conn.In(strSql, groupId, cartIds)
	if err != nil {
		log.Warnf("transform in array sql failed. %s", err)
		return
	}

	result, err = m.DB.Exec(strSql, args...)
	return
}

// 添加订单
func (m *MallGroupBuyingOrder) AddCommunityOrder(communityOrder *cidl.GroupBuyingOrderCommunityOrder) (result sql.Result, err error) {
	goodsDetail, err := communityOrder.GoodsDetail.ToString()
	if err != nil {
		log.Warnf("get community order goods detail failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO byo_community_order
			(
				ord_id,
				grp_id,
				grp_ord_id,
				goods_detail,
				count,
				total_market_price,
				total_group_buying_price,
				total_settlement_price,
				total_cost_price,
				version
			)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err = m.DB.Exec(strSql,
		communityOrder.OrderId,
		communityOrder.GroupId,
		communityOrder.GroupOrderId,
		goodsDetail,
		communityOrder.Count,
		communityOrder.TotalMarketPrice,
		communityOrder.TotalGroupBuyingPrice,
		communityOrder.TotalSettlementPrice,
		communityOrder.TotalCostPrice,
		communityOrder.Version,
	)

	return
}

func (m *MallGroupBuyingOrder) TxAddCommunityOrder(tx *conn.Tx, communityOrder *cidl.GroupBuyingOrderCommunityOrder) (result sql.Result, err error) {
	goodsDetail, err := communityOrder.GoodsDetail.ToString()
	if err != nil {
		log.Warnf("get community order goods detail failed. %s", err)
		return
	}

	strSql := `
		INSERT INTO byo_community_order
			(
				ord_id,
				grp_id,
				grp_ord_id,
				goods_detail,
				count,
				total_market_price,
				total_group_buying_price,
				total_settlement_price,
				total_cost_price,
				version
			)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err = m.DB.Exec(strSql,
		communityOrder.OrderId,
		communityOrder.GroupId,
		communityOrder.GroupOrderId,
		goodsDetail,
		communityOrder.Count,
		communityOrder.TotalMarketPrice,
		communityOrder.TotalGroupBuyingPrice,
		communityOrder.TotalSettlementPrice,
		communityOrder.TotalCostPrice,
		communityOrder.Version,
	)

	return
}

// 订单数目
func (m *MallGroupBuyingOrder) CommunityOrderCount(groupId uint32) (count uint32, err error) {
	strSql := `
		SELECT COUNT(*) FROM byo_community_order WHERE grp_id=?
	`
	err = m.DB.Get(&count, strSql, groupId)
	return
}

// 订单列表
func (m *MallGroupBuyingOrder) CommunityOrderList(groupId uint32, page uint32, pageSize uint32, createTimeAsc bool) (communityOrders []*cidl.GroupBuyingOrderCommunityOrder, err error) {
	if page <= 0 || pageSize <= 0 {
		err = errors.New("page or pageSize should be greater than 0")
		return
	}

	offset := (page - 1) * pageSize
	strOrderBy := "ASC"
	if false == createTimeAsc {
		strOrderBy = "DESC"
	}
	strSql := `
		SELECT
			ord_id,
			grp_id,
			grp_ord_id,
			goods_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version,
			create_time,
			is_cancel
		FROM
			byo_community_order
		WHERE
			grp_id=?
		ORDER BY create_time %s
		LIMIT ? OFFSET ?
	`
	strSql = fmt.Sprintf(strSql, strOrderBy)
	rows, err := m.DB.Query(strSql, groupId, pageSize, offset)
	if err != nil {
		log.Warnf("query community order list failed. %s", err)
		return
	}

	for rows.Next() {
		communityOrder := cidl.NewGroupBuyingOrderCommunityOrder()
		var goodsDetail string
		err = rows.Scan(
			&communityOrder.OrderId,
			&communityOrder.GroupId,
			&communityOrder.GroupOrderId,
			&goodsDetail,
			&communityOrder.Count,
			&communityOrder.TotalMarketPrice,
			&communityOrder.TotalGroupBuyingPrice,
			&communityOrder.TotalSettlementPrice,
			&communityOrder.TotalCostPrice,
			&communityOrder.Version,
			&communityOrder.CreateTime,
			&communityOrder.IsCancel,
		)
		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query community order failed. %s", err)
			}
			return
		}

		err = communityOrder.GoodsDetail.FromString(goodsDetail)
		if err != nil {
			log.Warnf("init goods detail from string failed. %s", err)
			return
		}

		communityOrders = append(communityOrders, communityOrder)
	}

	return
}

// 不允许取消订单列表
func (m *MallGroupBuyingOrder) CommunityNotAllowCancelOrderList(groupId uint32) (mOrderIds map[string]string, err error) {
	strSql := `
		SELECT
			a.order_id,
			b.allow_cancel,
			b.group_state,
			b.show_state
		FROM
			byo_community_buy a
		LEFT JOIN
			byo_task b
		ON 
			a.tsk_id = b.tsk_id
		WHERE
			a.grp_id = ?
		AND  (
			b.show_state != ? 
			OR
			b.group_state != ?
			OR
			b.allow_cancel = ?
		)
	`
	rows, err := m.DB.Query(strSql, groupId, cidl.GroupBuyingTaskShowStateShow, cidl.GroupBuyingTaskGroupStateInProgress,cidl.GroupBuyingTaskNotAllowCancel)
	if err != nil {
		log.Warnf("query order state failed. %s", err)
		return
	}

	mOrderIds = make(map[string]string)
	
	for rows.Next() {
		var (
			orderId  string
			allowCancel bool
			showState cidl.GroupBuyingTaskShowState
			groupState cidl.GroupBuyingTaskGroupState
		)
		err = rows.Scan(&orderId,&allowCancel,&groupState,&showState)
		if err != nil {
			log.Warnf("query allow cancel order failed. %s", err)
			return
		}

		var status string
		if !allowCancel {
			status = "不支持取消订单"
		} else if groupState != cidl.GroupBuyingTaskGroupStateInProgress {
			status = "订单已截团"
		} else if showState != cidl.GroupBuyingTaskShowStateShow {
			status = "已下架"
		}
	 	mOrderIds[orderId] = status;	
	}
	return 

}

func (m *MallGroupBuyingOrder) GetTaskLineList(taskId uint32) (lines []*cidl.GroupBuyingTaskLine, err error) {
	strSql := `
		SELECT
			distinct(a.lin_id),
			a.lin_name 
		FROM gby_line_community a,byo_community_buy_task b
		WHERE
			a.grp_id=b.grp_id
			AND b.tsk_id = ?
	`

	rows, err := m.DB.Query(strSql, taskId)
	if err != nil {
		log.Warnf("query line list failed. %s", err)
		return
	}

	for rows.Next() {
		line := cidl.NewGroupBuyingTaskLine()
		err = rows.Scan(
			&line.LineId,
			&line.LineName,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query line failed. %s", err)
			}
			return
		}
		lines = append(lines, line)
	}

	var lineIds []uint32
	for _,line := range lines {
		lineIds = append(lineIds,line.LineId)
	}

	strSql = `
		SELECT
			lin_id,
			update_time
		FROM byo_task_selected_line 
		WHERE
			tsk_id = ?
			AND lin_id in (?)
	`

	strSql, args, err := conn.In(strSql, taskId, lineIds)
	if err != nil {
		log.Warnf("transform sql in array failed. %s", err)
		return
	}

	rows, err = m.DB.Query(strSql, args...)
	if err != nil {
		log.Warnf("query line list failed. %s", err)
		return
	}

	m_line := make(map[uint32]time.Time)
	for rows.Next() {
		var lin_id uint32
		var time time.Time
		err = rows.Scan(
			&lin_id,
			&time,
		)

		if err != nil {
			if err != conn.ErrNoRows {
				log.Warnf("query line failed. %s", err)
			}
			continue
		}
		m_line[lin_id] = time
	}

	for i, line := range lines {
		if _,ok := m_line[line.LineId]; ok{
			lines[i].IsSelected = false  //曾经被选中过，则下次默认不选中 
			lines[i].UpdateTime = m_line[line.LineId]		
		} else {
			lines[i].IsSelected = true
		} 
	}

	return
}

func (m *MallGroupBuyingOrder) UpdateTaskSelectedLines(taskId uint32, lineIds []uint32) (result sql.Result, err error) {

	strSql := `
		INSERT INTO byo_task_selected_line 
			(
				tsk_id, 
				lin_id
			)
	    	values	
			%s
		ON DUPLICATE KEY UPDATE update_time = CURRENT_TIMESTAMP

		`

	var args []interface{}
	var sliceStrValue []string
	for _, lineId := range lineIds {
		sliceStrValue = append(sliceStrValue, "(?, ?)")
		args = append(args, taskId)
		args = append(args, lineId)
	}

	strValues := strings.Join(sliceStrValue, ",")
	strSql = fmt.Sprintf(strSql, strValues)
	result, err = m.DB.Exec(strSql, args...)
	return
}


func (m *MallGroupBuyingOrder) CommunityOrderTaskCount (groupId uint32, orderId string) (mTaskCount map[uint32]map[string]uint32, err error) {

	strSql := `
		SELECT tsk_id,sku_id,count
		FROM byo_community_buy 
		WHERE groupId = ? and orderId = ?
		`
	rows, err := m.DB.Query(strSql, groupId, orderId)
	if err != nil {
		log.Warnf("query byo_community_buy failed. %s", err)
		return
	}

	mTaskCount = make(map[uint32]map[string]uint32)
	for rows.Next() {
		var (
			tsk_id,count uint32
			sku_id string
		)
		err = rows.Scan(&tsk_id,&sku_id, &count)
		if err != nil {
			log.Warnf("scan taskid count failed. %s", err)
			return
		}
		if _, ok := mTaskCount[tsk_id]; !ok {
			mSkuCount := make(map[string]uint32)
			mSkuCount[sku_id] = count
			mTaskCount[tsk_id] = mSkuCount
			continue
		}
		mSkuCount := mTaskCount[tsk_id]
		mSkuCount[sku_id] = count
		mTaskCount[tsk_id] = mSkuCount 
	}
	return
}

func (m *MallGroupBuyingOrder) IsOrderAllowCancel(groupId uint32, orderId string) (isAllowCancel bool, err error) {

	strSql := `
		SELECT 1 
		FROM byo_community_buy a, byo_task b
		WHERE a.order_id = ?
		AND a.grp_id = ?
		AND a.tsk_id = b.tsk_id
		AND b.show_state = ?
		AND b.group_state = ?
		AND b.allow_cancel = ?
		`
	queryRow, err := m.DB.QueryRow(strSql, groupId, orderId, cidl.GroupBuyingTaskShowStateShow, cidl.GroupBuyingTaskGroupStateInProgress,cidl.GroupBuyingTaskAllowCancel)
	if err != nil {
		log.Warnf("query is order allow cancel failed. %s", err)
		return
	}
	
	var tmp uint32
	err = queryRow.Scan(
		&tmp,
	)

	isAllowCancel = false
	if err != nil {
		//for test
		log.Warnf("IsOrderAllowCancel failed")
		if err == conn.ErrNoRows {
			return false, nil
			log.Warnf("err no rows")
		}
		return
	}
	isAllowCancel = true
	return
}


// 取消社群购买商品
func (m *MallGroupBuyingOrder) TxDeleteCommunityBuy(tx *conn.Tx, orderId string, groupId uint32, taskId uint32) (result sql.Result, err error) {

	strSql := `
		DELETE FROM byo_community_buy
		WHERE order_id = ? AND grp_id = ? AND tsk_id = ?		
	`
	result, err = m.DB.Exec(strSql, orderId, groupId, taskId)
	return
}


func (m *MallGroupBuyingOrder) DeleteCommunityOrder(groupId uint32, orderId string) (result sql.Result, err error) {
	strSql := `
		UPDATE byo_community_order
		SET is_cancel = 1
		WHERE ord_id = ? 
	`
	result, err = m.DB.Exec(strSql, orderId)
	return
}

// 订单详情
func (m *MallGroupBuyingOrder) CommunityOrderInfo(orderId string) (communityOrder *cidl.GroupBuyingOrderCommunityOrder, err error) {

	strSql := `
		SELECT
			ord_id,
			grp_id,
			grp_ord_id,
			goods_detail,
			count,
			total_market_price,
			total_group_buying_price,
			total_settlement_price,
			total_cost_price,
			version,
			create_time,
			is_cancel
		FROM
			byo_community_order
		WHERE
			ord_id=?
	`
	row, err := m.DB.QueryRow(strSql, orderId)
	if err != nil {
		log.Warnf("query community order list failed. %s", err)
		return
	}

	communityOrder = cidl.NewGroupBuyingOrderCommunityOrder()
	var goodsDetail string
	err = row.Scan(
		&communityOrder.OrderId,
		&communityOrder.GroupId,
		&communityOrder.GroupOrderId,
		&goodsDetail,
		&communityOrder.Count,
		&communityOrder.TotalMarketPrice,
		&communityOrder.TotalGroupBuyingPrice,
		&communityOrder.TotalSettlementPrice,
		&communityOrder.TotalCostPrice,
		&communityOrder.Version,
		&communityOrder.CreateTime,
		&communityOrder.IsCancel,
	)
	if err != nil {
		if err != conn.ErrNoRows {
			log.Warnf("query community order failed. %s", err)
		}
		return
	}

	err = communityOrder.GoodsDetail.FromString(goodsDetail)
	if err != nil {
		log.Warnf("init goods detail from string failed. %s", err)
		return
	}

	return
}

func (m *MallGroupBuyingOrder) GetTaskVisibleTeamIDs(taskId uint32) (teamIds []uint32, err error) {
	defer func() {
		if err != nil {
			teamIds = nil
		}
	}()

	strSql := `
		SELECT
			team_id
		FROM byo_task_team
		WHERE tsk_id=?
	`
	rows, err := m.DB.Query(strSql, taskId)
	if err != nil {
		log.Warnf("query task teamId failed. %s", err)
		return
	}

	for rows.Next() {
		var teamId uint32
		err = rows.Scan(&teamId)
		if err != nil {
			log.Warnf("scan task teamId failed. %s", err)
			return
		}

		teamIds = append(teamIds,teamId)
	}
	return
}


