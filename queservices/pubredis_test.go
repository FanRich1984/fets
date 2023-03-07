package queservices

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/FanRich1984/fets/idmanage"
	"github.com/FanRich1984/fets/vatools"
)

// 测试玩家添加
func TestPubredis(t *testing.T) {
	// 构造对象
	ptPub := &PubRedis{
		pRedis:    NewTmpRedis(),
		streamKey: "STREAM_GAME_ORDER",
		ctx:       context.Background(),
		maxLen:    1000000,
	}

	wg := new(sync.WaitGroup)
	startTime := time.Now().UnixMilli()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 500; i++ {
				id, _ := idmanage.ObFlakId.GetId()

				ob := TmpTestDto{
					UserId:          id,
					PlantId:         int64(id),
					PlantOrderId:    idmanage.CreateBigInvitationCode(id),
					PlantGameId:     int64(id),
					PlantGameTypeId: int64(id),
				}
				strJson, _ := vatools.Json(&ob)
				ptPub.Push(map[string]interface{}{
					"data": strJson,
				})
			}
			wg.Done()
		}()
	}

	wg.Wait()
	endTime := time.Now().UnixMilli()
	log.Println("Success ", endTime-startTime)
}

type TmpTestDto struct {
	UserId          uint64 `json:"user_id"`            // 用户ID
	PlantId         int64  `json:"plant_id"`           // 第三方游戏平台ID
	PlantOrderId    string `json:"plant_order_id"`     // 第三方游戏平台的订单ID
	PlantGameId     int64  `json:"plant_game_id"`      // 第三方游戏ID
	PlantGameTypeId int64  `json:"plant_game_type_id"` // 第三方游戏类型ID
}
