package excel

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mz-eco/mz/settings"
)

func TestMain(m *testing.M) {
	settings.LoadFrom("../../", "")
	os.Exit(m.Run())
}

func TestIndentExcel_AddRow(t *testing.T) {
	excel, err := NewIndentExcel()
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i <= 29; i++ {
		err = excel.AddRow(
			time.Now(),
			"美味鸡爪",
			[3]string{"18.9一袋"},
			18.9,
			15.9,
			12.9,
			8574,
			32,
			100,
			100,
		)
		if err != nil {
			t.Error(err)
			return
		}

	}

	excel.BeforeSave()
	err = excel.SaveAs("F:\\17buy\\backend\\src\\business\\group_buying\\common\\excel\\test_add_row_1.xlsx")
	if err != nil {
		t.Error(err)
		return
	}

}

func TestIndentExcel_WriteToBuffer(t *testing.T) {
	excel, err := NewIndentExcel()
	if err != nil {
		t.Error(err)
		return
	}

	err = excel.AddRow(
		time.Now(),
		"美味鸡爪",
		[3]string{"18.9一袋"},
		18.9,
		15.9,
		12.9,
		8574,
		32,
		100,
		100,
	)
	if err != nil {
		t.Error(err)
		return
	}

	url, err := excel.SaveToQiniu("bilimall/indent/", "666768.xlsx")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(url)
}
