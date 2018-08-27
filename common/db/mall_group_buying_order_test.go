package db

import (
	"fmt"
	"os"
	"testing"
	//"time"

	"github.com/mz-eco/mz/settings"
)

func TestMain(m *testing.M) {
	settings.LoadFrom("../../", "")
	os.Exit(m.Run())
}



func TestMallGroupBuyingOrder_CommunityAllowCancelOrderList(t *testing.T) {
	dbGroupBuyingOrder := NewMallGroupBuyingOrder()
	groupId := uint32(219)
	orderlist, _:= dbGroupBuyingOrder.CommunityAllowCancelOrderList(groupId)
	for orderId,_ := range orderlist {
		fmt.Println("%v",orderId)
	}
}


/*
func TestMallGroupBuyingOrder_TaskTodayCount(t *testing.T) {
	dbGroupBuyingOrder := NewMallGroupBuyingOrder()
	org_id := 15
	team_ids := []uint32{6}
	count, _:= dbGroupBuyingOrder.TaskTodayCount(uint32(org_id),[]uint32(team_ids))
	fmt.Println(count)
}
func TestMallGroupBuyingOrder_TaskList(t *testing.T) {
	dbGroupBuyingOrder := NewMallGroupBuyingOrder()
	org_id := 15
	team_ids := []uint32{6}
	tasklist, _:= dbGroupBuyingOrder.TaskTodayList(uint32(org_id),[]uint32(team_ids),1,10,false)
	fmt.Println(len(tasklist))
	for _,task := range tasklist {
		fmt.Println("%v",task)
	}
}

func TestMallGroupBuyingOrder_GetTask(t *testing.T) {
	dbGroupBuyingOrder := NewMallGroupBuyingOrder()
	taskId := 220;
	task, _:= dbGroupBuyingOrder.GetTask(uint32(taskId))
	fmt.Println(task)
}


func TestMallGroupBuyingOrder_CommunityBuyListByTaskIdLineIdsOrderSkuId(t *testing.T) {
	dbGroupBuyingOrder := NewMallGroupBuyingOrder()
	taskId := 220
	lineIds := []uint32{136,137}
	buyList, _:= dbGroupBuyingOrder.CommunityBuyListByTaskIdLineIdsOrderSkuId(uint32(taskId),lineIds)
	for _,buy := range buyList{
		fmt.Println(buy.TaskId,buy.GroupId)
	}
}

func TestMallGroupBuyingOrder_UpdateTaskSelectedLines(t *testing.T) {
	dbGroupBuyingOrder := NewMallGroupBuyingOrder()
	taskId := 220
	lineIds := []uint32{136,137}
	 dbGroupBuyingOrder.UpdateTaskSelectedLines(uint32(taskId),lineIds)
}

func TestMallGroupBuyingOrder_GetTaskLineList(t *testing.T) {
	dbGroupBuyingOrder := NewMallGroupBuyingOrder()
	taskId := 220;
	lines, _:= dbGroupBuyingOrder.GetTaskLineList(uint32(taskId))
	for _, line := range lines {
		fmt.Println(line)
	}
	
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
*/
