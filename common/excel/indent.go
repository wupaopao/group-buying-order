package excel

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"common/file"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/mz-eco/mz/log"
)

type GroupBuyingExcelFile struct {
	ExcelFile *excelize.File
}

func (m *GroupBuyingExcelFile) SaveAs(fileName string) (err error) {
	return m.ExcelFile.SaveAs(fileName)
}

func (m *GroupBuyingExcelFile) WriteToBuffer() (buffer *bytes.Buffer, err error) {
	buffer = bytes.NewBuffer([]byte(""))
	err = m.ExcelFile.Write(buffer)
	return
}

func (m *GroupBuyingExcelFile) SaveToQiniu(prefix string, fileName string) (url string, err error) {
	buffer, err := m.WriteToBuffer()
	if err != nil {
		log.Warnf("write to buffer failed. %s")
		return
	}

	qiniu, err := file.GetQiniuPrivateBucket()
	if err != nil {
		log.Warnf("get qiniu pubilc bucket failed. %s", err)
		return
	}

	url, err = qiniu.UploadFileStream(buffer, int64(buffer.Len()), prefix, fileName)
	if err != nil {
		log.Warnf("save file to qiniu failed. %s", err)
		return
	}

	urlPrefix := qiniu.GetUrlPrefix()
	url = urlPrefix + url

	return
}

// 模板
var (
	indentTemplateBytes []byte
)

func init() {
	template, err := excelize.OpenFile("./data/excel/indent_template.xlsx")
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

	indentTemplateBytes = buffer.Bytes()
}

/**
生成订货单xlsx文件
*/
type IndentExcel struct {
	GroupBuyingExcelFile

	// 单据号
	TickerNumber string

	// 时间
	Date time.Time

	// 当前有效的sheet
	activeSheet string

	// 路线项行编号
	activeLineRowsCount uint32

	// sheet数目
	sheetCount uint32
}

func NewIndentExcel() (indentExcel *IndentExcel, err error) {
	return NewIndentExcelFromReader(bytes.NewReader(indentTemplateBytes))
}

func NewIndentExcelFromReader(reader io.Reader) (indentExcel *IndentExcel, err error) {
	excel, err := excelize.OpenReader(reader)
	if err != nil {
		log.Warnf("open reader failed. %s", err)
		return
	}

	indentExcel = &IndentExcel{
		GroupBuyingExcelFile: GroupBuyingExcelFile{
			ExcelFile: excel,
		},
	}

	return
}

func (m *IndentExcel) templateSheetName() string {
	return "模板-订货单"
}

func (m *IndentExcel) newSheet() (sheet string, err error) {
	xlsx := m.ExcelFile
	m.sheetCount++
	sheet = fmt.Sprintf("订货单-%d", m.sheetCount)

	index := xlsx.NewSheet(sheet)
	templateSheetIndex := xlsx.GetSheetIndex(m.templateSheetName())
	err = xlsx.CopySheet(templateSheetIndex, index)
	if err != nil {
		log.Warnf("copy template sheet failed. %s", err)
		return
	}

	m.activeSheet = sheet

	return
}

func (m *IndentExcel) AddRow(
	startTime time.Time,
	title string,
	specification [3]string,
	buyPrice float64,
	settlePrice float64,
	costPrice float64,
	salesCount uint32,
	communityCount uint32,
	totalCost float64,
	totalSettle float64,

) (err error) {

	if m.activeLineRowsCount >= 28 || m.activeSheet == "" {
		_, err = m.newSheet()
		if err != nil {
			log.Warnf("new sheet failed. %s", err)
			return
		}

		m.activeLineRowsCount = 0
	}

	rowIndex := m.activeLineRowsCount + 2
	xlsx := m.ExcelFile
	sheet := m.activeSheet

	xlsx.SetCellValue(sheet, fmt.Sprintf("A%d", rowIndex), startTime)
	xlsx.SetCellValue(sheet, fmt.Sprintf("B%d", rowIndex), title)
	xlsx.SetCellValue(sheet, fmt.Sprintf("C%d", rowIndex), specification[0])
	xlsx.SetCellValue(sheet, fmt.Sprintf("D%d", rowIndex), specification[1])
	xlsx.SetCellValue(sheet, fmt.Sprintf("E%d", rowIndex), specification[2])
	xlsx.SetCellValue(sheet, fmt.Sprintf("F%d", rowIndex), buyPrice)
	xlsx.SetCellValue(sheet, fmt.Sprintf("G%d", rowIndex), settlePrice)
	xlsx.SetCellValue(sheet, fmt.Sprintf("H%d", rowIndex), costPrice)
	xlsx.SetCellValue(sheet, fmt.Sprintf("I%d", rowIndex), salesCount)
	xlsx.SetCellValue(sheet, fmt.Sprintf("J%d", rowIndex), communityCount)
	xlsx.SetCellValue(sheet, fmt.Sprintf("K%d", rowIndex), totalCost)
	xlsx.SetCellValue(sheet, fmt.Sprintf("L%d", rowIndex), totalSettle)

	m.activeLineRowsCount++

	return
}

func (m *IndentExcel) BeforeSave() {
	xlsx := m.ExcelFile
	xlsx.DeleteSheet(m.templateSheetName())
}
