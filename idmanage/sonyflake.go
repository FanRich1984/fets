// 雪花ID生成器
// 根据cfg.ini里 idmanage 的配置信息实例化多个雪花生成ID，如果配置文件里没有这个对象
// 则生成默认的1个雪花生成对象，机器ID为本地主机的IP地址的后两位

package idmanage

import (
	"fmt"
	"time"

	"github.com/FanRich1984/fets/vatools"

	"github.com/sony/sonyflake"
)

// 定议一个全局可以访问的雪花ID生成对象
var ObFlakId *FlakeId

const (
	SPEED_YEAR  int = 2023
	SPEED_MONTH int = 1
	SPEED_DAY   int = 1
)

func init() {
	// 获取配置信息
	c := NewCfg()
	// 先判断是否有Machid
	if c.MachineId > 0 {
		if c.MachineId > 100 {
			c.MachineId = 100
		}
		// 判断数量
		if c.Quantity < 1 {
			c.Quantity = 1
		} else if c.Quantity > 10 {
			c.Quantity = 10
		}
		startStep := c.MachineId * 10
		var i int16
		ObFlakId = &FlakeId{
			sonfyFlaks: make([]*sonyflake.Sonyflake, 0, int(startStep)),
			machineId:  c.MachineId,
		}
		for ; i < c.Quantity; i++ {
			tMachineId := uint16(startStep + i)
			st := sonyflake.Settings{
				StartTime: time.Date(SPEED_YEAR, time.Month(SPEED_MONTH), SPEED_DAY, 0, 0, 0, 0, time.UTC),
				MachineID: func() (uint16, error) {
					return tMachineId, nil
				},
			}
			ObFlakId.sonfyFlaks = append(ObFlakId.sonfyFlaks, sonyflake.NewSonyflake(st))
		}
	} else {
		// 小于1则只开启单例
		ObFlakId = &FlakeId{
			sonfyFlaks: []*sonyflake.Sonyflake{
				sonyflake.NewSonyflake(sonyflake.Settings{
					StartTime: time.Date(SPEED_YEAR, time.Month(SPEED_MONTH), SPEED_DAY, 0, 0, 0, 0, time.UTC),
				}),
			},
			machineId: 0,
		}
	}
}

// SONY雪发ID管理器
type FlakeId struct {
	sonfyFlaks []*sonyflake.Sonyflake
	machineId  int16
}

// 获取雪花ID
func (flakid *FlakeId) GetId() (uint64, error) {
	c := len(flakid.sonfyFlaks)
	if c < 1 {
		return 0, fmt.Errorf("no flakeid object")
	} else if c == 1 {
		return flakid.sonfyFlaks[0].NextID()
	}
	rndVal := vatools.CRnd(0, c-1)
	return flakid.sonfyFlaks[rndVal].NextID()
}

// 获取当前生成对象的机器ID号
func (flakid *FlakeId) MachineID() int16 {
	return flakid.machineId
}

// 获取当前拥有的雪花ID生成对象总数
func (flakid *FlakeId) CountObject() int {
	return len(flakid.sonfyFlaks)
}
