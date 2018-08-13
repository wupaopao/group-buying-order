package db

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

func TestMallGroupBuyingOrder_TxLockInventorySurplus(t *testing.T) {
	dbGroupBuyingOrder := NewMallGroupBuyingOrder()
	tx, err := dbGroupBuyingOrder.DB.Begin()
	if err != nil {
		t.Error(err)
		return
	}

	taskId := uint32(130)
	skuIds := []string{
		"1:1-2:1",
		"1:1-2:2",
	}
	_, err = dbGroupBuyingOrder.TxLockInventorySurplus(tx, taskId, skuIds)
	if err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}

	go func() {
		_, err := dbGroupBuyingOrder.SubtractInventorySurplus(taskId, skuIds[0], 1)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println("go sub with no tx ", time.Now())
	}()

	time.Sleep(5 * time.Second)
	_, err = dbGroupBuyingOrder.TxSubtractInventorySurplus(tx, taskId, skuIds[0], 0)
	if err != nil {
		tx.Rollback()
		t.Error(err)
		return
	}

	tx.Commit()

	fmt.Println("sub with tx ", time.Now())

	_, err = dbGroupBuyingOrder.SubtractInventorySurplus(taskId, skuIds[0], 1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("sub with no tx ", time.Now())
}
