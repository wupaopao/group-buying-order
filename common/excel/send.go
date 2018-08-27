package excel

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"business/community/proxy/community"
	"business/group-buying-order/cidl"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/mz-eco/mz/errors"
	"github.com/mz-eco/mz/log"
)

// 模板
var (
	sendTemplateBytes []byte
)

func init() {
	template, err := excelize.OpenFile("./data/excel/send_template.xlsx")
	if err != nil {
		panic(err)
		return
	}

	buffer := bytes.NewBuffer(nil)
	err = template.Write(buffer)
	if err != nil {
		panic(err)
		return
	}

	sendTemplateBytes = buffer.Bytes()
}

/**
生成配送单xlsx文件
*/
type SendExcel struct {
	GroupBuyingExcelFile

	// 单据号
	TicketNumber string

	// 时间
	Date time.Time

	// 组织名称
	OrganizationName string
}

func NewSendExcel() (sendExcel *SendExcel, err error) {
	excel, err := excelize.OpenReader(bytes.NewReader(sendTemplateBytes))
	if err != nil {
		log.Warnf("open excel reader failed. %s", err)
		return
	}

	sendExcel = &SendExcel{
		GroupBuyingExcelFile: GroupBuyingExcelFile{
			ExcelFile: excel,
		},
	}
	return
}

func (m *SendExcel) sendLineSummarySheetTemplateName() string {
	return "模板-汇总对账单"
}

func (m *SendExcel) sendLineSheetTemplateName() string {
	return "模板-配送路线"
}

func (m *SendExcel) sendCommunitySheetTemplateName() string {
	return "模板-社群"
}

func (m *SendExcel) groupSummarySheetTemplateName() string {
	return "模板-团长销售额简报"
}

func (m *SendExcel) AddTaskSummarySheets(sheetName string, tasks []*cidl.GroupBuyingOrderTask, groupMap map[uint32]*community.Group, sendCommunityMap map[uint32]*cidl.GroupBuyingSendCommunity) (err error) {
	taskSummarySheets := &TaskSummarySheets{
		SendExcel:        m,
		ExcelFile:        m.ExcelFile,
		Tasks:            tasks,
		GroupMap:         groupMap,
		SendCommunityMap: sendCommunityMap,
		SheetName:        sheetName,
	}
	err = taskSummarySheets.InitSheets()
	if err != nil {
		log.Warnf("add task summary sheets failed. %s", err)
		return
	}

	return
}

func (m *SendExcel) GetSendLineSummarySheets() (sendLineSummarySheets *SendLineSummarySheets) {
	sendLineSummarySheets = &SendLineSummarySheets{
		ExcelFile: m.ExcelFile,
		SendExcel: m,
	}
	return
}

func (m *SendExcel) GetGroupSummarySheets() (groupSummarySheets *GroupSummarySheets) {
	groupSummarySheets = &GroupSummarySheets{
		ExcelFile: m.ExcelFile,
		SendExcel: m,
	}
	return
}

func (m *SendExcel) AddSendLineSheets(sendLine *cidl.GroupBuyingSendLine) (err error) {
	sendLineSheets := &SendLineSheets{
		SendExcel: m,
		ExcelFile: m.ExcelFile,
		SendLine:  sendLine,
	}

	err = sendLineSheets.InitSheets()
	if err != nil {
		log.Warnf("send line sheets init sheets failed. %s", err)
		return
	}

	return
}

func (m *SendExcel) AddSendCommunitySheets(sendCommunity *cidl.GroupBuyingSendCommunity) (err error) {
	sendCommunitySheets := &SendCommunitySheets{
		SendExcel:     m,
		ExcelFile:     m.ExcelFile,
		SendCommunity: sendCommunity,
	}

	err = sendCommunitySheets.InitSheets()
	if err != nil {
		log.Warnf("send community sheets init sheets failed. %s", err)
		return
	}

	return
}

func (m *SendExcel) BeforeSave() {
	xlsx := m.ExcelFile
	xlsx.DeleteSheet(m.sendLineSummarySheetTemplateName())
	xlsx.DeleteSheet(m.sendLineSheetTemplateName())
	xlsx.DeleteSheet(m.sendCommunitySheetTemplateName())
}

/**
团购任务销售统计
*/

type TaskSummarySheets struct {
	SendExcel *SendExcel

	// excel文件
	ExcelFile *excelize.File

	// 任务
	Tasks []*cidl.GroupBuyingOrderTask

	// 社群
	GroupMap map[uint32]*community.Group

	// 社群配送map
	SendCommunityMap map[uint32]*cidl.GroupBuyingSendCommunity

	// 概况表单名称
	SheetName string
}

func (m *TaskSummarySheets) InitSheets() (err error) {
	xlsx := m.ExcelFile
	sheet := m.SheetName
	xlsx.NewSheet(sheet)

	// taskId->skuId->columnIndex
	taskSkuColumnIndexMap := make(map[uint32]map[string]int)

	// 设置头
	var taskColumnIndex int = 66
	for _, task := range m.Tasks {
		xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", string(taskColumnIndex), 1), task.Title)

		skuColumnIndexMap, ok := taskSkuColumnIndexMap[task.TaskId]
		if !ok {
			skuColumnIndexMap = make(map[string]int)
			taskSkuColumnIndexMap[task.TaskId] = skuColumnIndexMap
		}

		for _, skuItem := range task.Specification.SkuMap {
			if !skuItem.IsShow {
				continue
			}

			axis := fmt.Sprintf("%s%d", string(taskColumnIndex), 2)
			xlsx.SetCellValue(sheet, axis, skuItem.Name)

			skuColumnIndexMap[skuItem.SkuId] = taskColumnIndex

			taskColumnIndex++
		}

		for _, skuItem := range task.Specification.CombinationSkuMap {
			if !skuItem.IsShow {
				continue
			}

			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", string(taskColumnIndex), 2), skuItem.Name)
			skuColumnIndexMap[skuItem.SkuId] = taskColumnIndex
			taskColumnIndex++
		}
	}

	// 设置内容
	rowIndex := 3
	for groupId, sendCommunity := range m.SendCommunityMap {
		group, ok := m.GroupMap[groupId]
		if !ok {
			err = errors.New("no needed group")
			return
		}

		xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", string(65), rowIndex), group.Name)
		for _, statisticItem := range *sendCommunity.Statistics {
			for skuId, skuItem := range statisticItem.Sku {
				colIndex, ok := taskSkuColumnIndexMap[skuItem.TaskId][skuId]
				if !ok {
					err = errors.New("no needed task sku")
					return
				}

				xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", string(colIndex), rowIndex), skuItem.Sales)
			}
		}

		rowIndex++
	}

	return
}

/**
团长销售额简报
*/
type GroupSummarySheets struct {
	SendExcel *SendExcel

	// excel文件
	ExcelFile *excelize.File

	// 当前有效的sheet
	activeSheet string

	// 路线项行编号
	activeLineRowsCount uint32

	// sheet数目
	//sheetCount uint32

	// 配送社区
	SendCommunity *cidl.GroupBuyingSendCommunity
}

func (m *GroupSummarySheets) newSheet() (sheet string, err error) {
	xlsx := m.ExcelFile
	//m.sheetCount++
	sheet = fmt.Sprintf("团长销售额简报")

	index := xlsx.NewSheet(sheet)
	templateSheetIndex := xlsx.GetSheetIndex(m.SendExcel.groupSummarySheetTemplateName())
	err = xlsx.CopySheet(templateSheetIndex, index)
	if err != nil {
		log.Warnf("copy group summary sheet failed. %s", err)
		return
	}

	m.activeSheet = sheet

	return
}

func (m *GroupSummarySheets) AddLineRow(groupName string, groupManagerName string, groupManagerMobile string, settlementAmount float64) (err error) {
	if m.activeSheet == "" {
		_, err = m.newSheet()
		if err != nil {
			log.Warnf("new sheet failed. %s", err)
			return
		}

		m.activeLineRowsCount = 0
	}

	xlsx := m.ExcelFile
	sheet := m.activeSheet
	rowIndex := m.activeLineRowsCount + 2

	xlsx.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), m.SendExcel.TicketNumber)
	xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), groupName)
	xlsx.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), groupManagerName)
	xlsx.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), groupManagerMobile)
	xlsx.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), settlementAmount)
	xlsx.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), m.SendExcel.Date)

	m.activeLineRowsCount++
	return
}


/**
配送路线汇总对账单
*/
type SendLineSummarySheets struct {
	SendExcel *SendExcel

	// excel文件
	ExcelFile *excelize.File

	// 当前有效的sheet
	activeSheet string

	// 路线项行编号
	activeLineRowsCount uint32

	// 当前有效的小写金额
	activeTotalAmount float64

	// sheet数目
	sheetCount uint32
}

func (m *SendLineSummarySheets) newSheet() (sheet string, err error) {
	xlsx := m.ExcelFile
	m.sheetCount++
	sheet = fmt.Sprintf("汇总对账单-%d", m.sheetCount)

	index := xlsx.NewSheet(sheet)
	templateSheetIndex := xlsx.GetSheetIndex(m.SendExcel.sendLineSummarySheetTemplateName())
	err = xlsx.CopySheet(templateSheetIndex, index)
	if err != nil {
		log.Warnf("copy send line summary sheet failed. %s", err)
		return
	}

	m.activeSheet = sheet

	// 画基本的文字
	xlsx.SetCellValue(sheet, "G2", fmt.Sprintf("NO:%s", m.SendExcel.TicketNumber))
	xlsx.SetCellValue(sheet, "A3", m.SendExcel.Date)
	xlsx.SetCellValue(sheet, "B23", m.SendExcel.OrganizationName)

	return
}

func (m *SendLineSummarySheets) AddLineRow(lineName string, communityCount uint32, totalSettlement float64) (err error) {
	if m.activeLineRowsCount >= 16 || m.activeSheet == "" {
		_, err = m.newSheet()
		if err != nil {
			log.Warnf("new sheet failed. %s", err)
			return
		}

		m.activeTotalAmount = 0
		m.activeLineRowsCount = 0
	}

	xlsx := m.ExcelFile
	sheet := m.activeSheet
	rowIndex := m.activeLineRowsCount + 6

	xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), lineName)
	xlsx.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), communityCount)
	xlsx.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), totalSettlement)
	m.activeTotalAmount += totalSettlement

	m.activeLineRowsCount++

	xlsx.SetCellValue(sheet, "A4", fmt.Sprintf("包含%d条配送路线", m.activeLineRowsCount))
	xlsx.SetCellValue(sheet, "F22", fmt.Sprintf("小写金额:￥%.2f", m.activeTotalAmount))

	return
}

/**
配送路线单
*/
type SendLineSheets struct {
	SendExcel *SendExcel

	// excel文件
	ExcelFile *excelize.File

	// 当前有效的sheet
	activeSheet string

	// 路线项行编号
	activeLineRowsCount uint32

	// sheet数目
	sheetCount uint32

	// 配送路线
	SendLine *cidl.GroupBuyingSendLine
}

func (m *SendLineSheets) newSheet() (sheet string, err error) {
	xlsx := m.ExcelFile
	m.sheetCount++
	sheet = fmt.Sprintf("路线-%d-%s-%d", m.SendLine.LineId, m.SendLine.LineName, m.sheetCount)
	index := xlsx.NewSheet(sheet)
	templateSheetIndex := xlsx.GetSheetIndex(m.SendExcel.sendLineSheetTemplateName())
	err = xlsx.CopySheet(templateSheetIndex, index)
	if err != nil {
		log.Warnf("copy send line sheet failed. %s", err)
		return
	}

	m.activeSheet = sheet

	// 画基本文字
	xlsx.SetCellValue(sheet, "A1", m.SendLine.LineName)
	xlsx.SetCellValue(sheet, "G2", fmt.Sprintf("NO:%s", m.SendExcel.TicketNumber))
	xlsx.SetCellValue(sheet, "A3", m.SendExcel.Date)
	xlsx.SetCellValue(sheet, "A4", fmt.Sprintf("覆盖%d个配送点", m.SendLine.CommunityCount))
	xlsx.SetCellValue(sheet, "F24", fmt.Sprintf("小写金额:￥%.2f", m.SendLine.SettlementAmount))
	xlsx.SetCellValue(sheet, "B25", m.SendExcel.OrganizationName)

	return
}

func (m *SendLineSheets) InitSheets() (err error) {
	xlsx := m.ExcelFile
	sendLine := m.SendLine
	for _, taskBuy := range *sendLine.Statistics { // 逐个任务

		var skuIds []string
		for skuId, _ := range taskBuy.Sku {
			skuIds = append(skuIds, skuId)
		}
		sort.Strings(skuIds)

		for _, skuId := range skuIds {
			skuBuy := taskBuy.Sku[skuId]
			if m.activeLineRowsCount >= 18 || m.activeSheet == "" {
				_, err = m.newSheet()
				if err != nil {
					log.Warnf("new send line sheet failed. %s", err)
					return
				}

				m.activeLineRowsCount = 0
			}

			sheet := m.activeSheet
			rowIndex := m.activeLineRowsCount + 6

			xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), taskBuy.TaskTitle)

			var labelKeys []string
			for labelkey, _ := range skuBuy.Labels {
				labelKeys = append(labelKeys, labelkey)
			}
			sort.Strings(labelKeys)

			labelIndex := 0
			for _, labelKey := range labelKeys {
				labelItem := skuBuy.Labels[labelKey]
				if labelIndex > 2 {
					break
				}

				colIndex := "D"
				switch labelIndex {
				case 0:
					colIndex = "D"
				case 1:
					colIndex = "E"
				case 2:
					colIndex = "F"
				}

				xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", colIndex, rowIndex), labelItem.Name)

				labelIndex++
			}

			xlsx.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), skuBuy.Sales)
			xlsx.SetCellValue(sheet, fmt.Sprintf("H%d", rowIndex), skuBuy.TotalSettlement)

			m.activeLineRowsCount++
		}
	}

	return
}

/**
配送社区单
*/
type SendCommunitySheets struct {
	SendExcel *SendExcel

	// excel文件
	ExcelFile *excelize.File

	// 当前有效的sheet
	activeSheet string

	// 路线项行编号
	activeLineRowsCount uint32

	// sheet数目
	sheetCount uint32

	// 配送社区
	SendCommunity *cidl.GroupBuyingSendCommunity
}

func (m *SendCommunitySheets) newSheet() (sheet string, err error) {
	xlsx := m.ExcelFile
	m.sheetCount++
	sheet = fmt.Sprintf("社群-%d-%s-%d", m.SendCommunity.GroupId, m.SendCommunity.GroupName, m.sheetCount)
	index := xlsx.NewSheet(sheet)
	templateSheetIndex := xlsx.GetSheetIndex(m.SendExcel.sendCommunitySheetTemplateName())
	err = xlsx.CopySheet(templateSheetIndex, index)
	if err != nil {
		log.Warnf("copy send community sheet failed. %s", err)
		return
	}

	m.activeSheet = sheet

	// 画基本文字
	sendCommunity := m.SendCommunity
	xlsx.SetCellValue(sheet, "A1", sendCommunity.OrganizationName)
	xlsx.SetCellValue(sheet, "F2", fmt.Sprintf("NO:%s", m.SendExcel.TicketNumber))
	xlsx.SetCellValue(sheet, "A2", fmt.Sprintf("地址:%s", sendCommunity.OrganizationAddress))
	xlsx.SetCellValue(sheet, "A4", fmt.Sprintf("社群名称:%s", sendCommunity.GroupName))
	xlsx.SetCellValue(sheet, "D4", fmt.Sprintf("客户电话:%s", sendCommunity.GroupManagerMobile))
	xlsx.SetCellValue(sheet, "F4", fmt.Sprintf("制单人员:%s", sendCommunity.AuthorName))
	xlsx.SetCellValue(sheet, "A5", fmt.Sprintf("客户地址:%s", sendCommunity.GroupAddress))
	xlsx.SetCellValue(sheet, "D5", fmt.Sprintf("联系人:%s", sendCommunity.GroupManagerName))
	xlsx.SetCellValue(sheet, "F5", fmt.Sprintf("送货日期:%s", m.SendExcel.Date.Format("2006-01-02")))
	xlsx.SetCellValue(sheet, "E24", fmt.Sprintf("小写金额:￥%.2f", m.SendCommunity.SettlementAmount))

	return
}

func (m *SendCommunitySheets) InitSheets() (err error) {
	xlsx := m.ExcelFile
	sendCommunity := m.SendCommunity
	for _, taskBuy := range *sendCommunity.Statistics { // 逐个任务

		var skuIds []string
		for skuId, _ := range taskBuy.Sku {
			skuIds = append(skuIds, skuId)
		}
		sort.Strings(skuIds)

		for _, skuId := range skuIds {
			skuBuy := taskBuy.Sku[skuId]
			if m.activeLineRowsCount >= 17 || m.activeSheet == "" {
				_, err = m.newSheet()
				if err != nil {
					log.Warnf("new sheet failed. %s", err)
					return
				}

				m.activeLineRowsCount = 0
			}

			sheet := m.activeSheet
			rowIndex := m.activeLineRowsCount + 7

			xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), taskBuy.TaskTitle)

			var labelKeys []string
			for labelKey, _ := range skuBuy.Labels {
				labelKeys = append(labelKeys, labelKey)
			}
			sort.Strings(labelKeys)

			var labelNames []string
			for _, labelKey := range labelKeys {
				labelItem := skuBuy.Labels[labelKey]
				labelNames = append(labelNames, labelItem.Name)
			}

			xlsx.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), strings.Join(labelNames, "-"))
			xlsx.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), skuBuy.Sales)
			xlsx.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), skuBuy.SettlementPrice)
			xlsx.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), skuBuy.TotalSettlement)

			m.activeLineRowsCount++
		}
	}

	return
}
