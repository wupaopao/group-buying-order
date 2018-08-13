package excel

import (
	"testing"
	"time"

	"business/group_buying/common/db"
)

func TestSendLineSheets_AddLineRow(t *testing.T) {
	var err error
	sendExcel, err := NewSendExcel()
	if err != nil {
		t.Error(err)
		return
	}

	sendExcel.TicketNumber = "20180131"
	sendExcel.Date = time.Now()
	sendExcel.OrganizationName = "味罗天下"

	lineSheets := sendExcel.GetSendLineSummarySheets()

	for i := 0; i < 20; i++ {
		lineSheets.AddLineRow("配送路线一", 68, 189473)
	}

	err = sendExcel.SaveAs("F:\\17buy\\backend\\src\\business\\group_buying\\common\\excel\\test_add_line_row.xlsx")
	if err != nil {
		t.Error(err)
		return
	}

}

func TestSendExcel_AddSendLineSheets(t *testing.T) {
	var err error
	sendExcel, err := NewSendExcel()
	if err != nil {
		t.Error(err)
		return
	}

	sendExcel.TicketNumber = "20180131"
	sendExcel.Date = time.Now()
	sendExcel.OrganizationName = "味罗天下"

	dbGroupBuying := db.NewMallGroupBuying()
	sendLine, err := dbGroupBuying.GetSendLine("94ce79636922e9ebc1b287e9455819bc", 112)
	if err != nil {
		t.Error(err)
		return
	}

	err = sendExcel.AddSendLineSheets(sendLine)
	if err != nil {
		t.Error(err)
		return
	}

	err = sendExcel.SaveAs("F:\\17buy\\backend\\src\\business\\group_buying\\common\\excel\\test_send_line.xlsx")
	if err != nil {
		t.Error(err)
		return
	}

}

func TestSendExcel_AddSendCommunitySheets(t *testing.T) {
	var err error
	sendExcel, err := NewSendExcel()
	if err != nil {
		t.Error(err)
		return
	}

	sendExcel.TicketNumber = "20180131"
	sendExcel.Date = time.Now()
	sendExcel.OrganizationName = "味罗天下"

	dbGroupBuying := db.NewMallGroupBuying()
	sendCommunity, err := dbGroupBuying.GetSendCommunity("ad1d189dd9da5d61446bab8b3176b530", 219)
	if err != nil {
		t.Error(err)
		return
	}

	err = sendExcel.AddSendCommunitySheets(sendCommunity)
	if err != nil {
		t.Error(err)
		return
	}

	err = sendExcel.SaveAs("F:\\17buy\\backend\\src\\business\\group_buying\\common\\excel\\test_send_community.xlsx")
	if err != nil {
		t.Error(err)
		return
	}

}
