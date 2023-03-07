package queservices

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/FanRich1984/fets/vatools"
)

func TestSubRedis(t *testing.T) {
	ptSubGroup := NewSubGroupRedis(NewTmpRedis(), SubGroupRedisArg{
		StreamKey: "STREAM_GAME_ORDER",
		GroupName: "SAVE_ORDER_GROUP",
		FnCallBack: func(mpInfo map[string]interface{}) error {
			fmt.Println("回调处理:", mpInfo)
			return nil
		},
	})

	go ptSubGroup.SafeRun(2)

	fmt.Println("waiting...")
	time.Sleep(time.Second * 2)

	fmt.Println("停止当前运行的消费者")
	ptSubGroup.Stop()
	time.Sleep(time.Second * 3)

	fmt.Println("再次运行3个消费者")
	go ptSubGroup.SafeRun(2)
	time.Sleep(time.Second * 3)

	// 写入数据
	ptPub := NewPubRedis(NewTmpRedis(), PubRedisArg{
		StreamKey: "STREAM_GAME_ORDER",
		MaxLen:    4,
	})
	ptPub.Push(map[string]interface{}{"data": fmt.Sprint("测试的数据", vatools.CRnd(1, 200))})

	time.Sleep(time.Second * 5)
	fmt.Println("执行终止")
	ptSubGroup.Stop()
	time.Sleep(time.Second * 2)
	fmt.Println("Success")
}

func TestSubRedisIsKey(t *testing.T) {
	ptSubGroup := NewSubGroupRedis(NewTmpRedis(), SubGroupRedisArg{
		StreamKey: "STREAM_GAME_ORDER",
		GroupName: "SAVE_ORDER_GROUP",
		FnCallBack: func(mpInfo map[string]interface{}) error {
			fmt.Println("回调处理:", mpInfo)
			return nil
		},
	})
	fmt.Println(ptSubGroup.HasStreamKey())
}

func TestSubRedisPending(t *testing.T) {
	ptSubGroup := NewSubGroupRedis(NewTmpRedis(), SubGroupRedisArg{
		StreamKey: "STREAM_GAME_ORDER",
		GroupName: "SAVE_ORDER_GROUP",
		FnCallBack: func(mpInfo map[string]interface{}) error {
			fmt.Println("回调处理:", mpInfo)
			return nil
		},
	})

	pConsumer := ptSubGroup.CreateConsum(context.Background())
	pConsumer.ListPending()
}
