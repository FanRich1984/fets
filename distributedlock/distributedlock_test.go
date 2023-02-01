package distributedlock

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/FanRich1984/fets/vatools"
	"github.com/go-redis/redis/v8"
)

// 分布式锁的测试
func TestDistributedLock(t *testing.T) {
	var iCount int
	keyName := "9911001011"

	// 构造Redis
	pRedis := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
		DB:      1,
	})
	// Ping
	pRedis.Ping(context.Background())

	// 构造分布式锁
	pUserLk := NewDistributeLock("LK_USER", pRedis)
	startTime := time.Now().UnixMilli()
	wg := new(sync.WaitGroup)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			// 不使用分布锁增加
			// writeKeyVal(pRedis)
			// iCount++

			// 使用分布锁增加值
			if pUserLk.SyncLock(keyName) {
				iCount++
				writeKeyVal(pRedis)
				fmt.Println(idx, ":", iCount)
				pUserLk.UnLock(keyName)
			} else {
				fmt.Println(idx, "被锁中执行失败")
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	endTime := time.Now().UnixMilli()
	fmt.Println("Success : ", endTime-startTime)
}

func TestWriteRedisExpTime(t *testing.T) {
	// 构造Redis
	pRedis := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
		DB:      1,
	})
	// Ping
	pRedis.Ping(context.Background())
	pRedis.Set(context.Background(), "TMP_TEST_1", "val", 60*time.Second)
}

func writeKeyVal(pRedis *redis.Client) {
	v := pRedis.Get(context.Background(), "TMP_TEST").Val()
	iCount := vatools.SInt(v)
	fmt.Println("当前Key值:", vatools.SInt(v))
	iCount++
	pRedis.Set(context.Background(), "TMP_TEST", iCount, 0)
	fmt.Println("在Redis里累加后的值:", pRedis.Get(context.Background(), "TMP_TEST").Val())
}
