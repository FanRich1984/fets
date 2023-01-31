package idmanage

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 测试并发生成10000个雪花ID对象
func TestSonyFlake(t *testing.T) {
	wg := new(sync.WaitGroup)
	startTime := time.Now().UnixMilli()
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			_, _ = ObFlakId.GetId()
			// fmt.Println(ObFlakId.GetId())
			wg.Done()
		}()
	}
	wg.Wait()
	endTime := time.Now().UnixMilli()
	fmt.Println(endTime - startTime)

	fmt.Println(time.UTC.String(), time.Local.String())
}

// 测试当前雪花ID生成对象里的机器ID
func TestSonyFlakeMachineId(t *testing.T) {
	c := NewCfg()
	fmt.Println(c.MachineId, c.Quantity)
	fmt.Println(ObFlakId.machineId, " Count:", ObFlakId.CountObject())
}
