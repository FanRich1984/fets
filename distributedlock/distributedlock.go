package distributedlock

import (
	"context"
	"log"
	"time"

	"github.com/FanRich1984/fets/vatools"
	"github.com/go-redis/redis/v8"
)

// 实列化分布式锁
// 通过传入锁的名称和Redis缓存来构造分布式锁对象
func NewDistributeLock(lkName string, pRedis *redis.Client) *DistributedLock {
	return &DistributedLock{
		pRedis: pRedis,
		lkName: lkName,
		ctx:    context.Background(),
	}
}

// 分布式锁对象
type DistributedLock struct {
	pRedis *redis.Client   // Redis对象
	lkName string          // 锁名称
	ctx    context.Context // 上锁的上下文
}

// 传入要加锁的关键字上锁
func (dl *DistributedLock) Lock(keyName string) bool {
	blnOk, err := dl.pRedis.SetNX(dl.ctx, dl.lkName+keyName, 1, 1*time.Second).Result()
	if err != nil {
		// LOG发生错误
		log.Println("DistributedLock err:", err.Error())
		return false
	}
	return blnOk
}

// 同步锁着并等待获取分布锁
func (dl *DistributedLock) SyncLock(keyName string) bool {
	for i := 0; i < 1000; i++ {
		if dl.Lock(keyName) {
			return true
		}
		time.Sleep(time.Duration(vatools.CRnd(2000, 20000)) * time.Microsecond)
	}
	// 在指定等待时间内没有获取需要的分布锁退出
	return false
}

// 对指定的关键字的锁进行解锁
func (dl *DistributedLock) UnLock(keyName string) bool {
	set, err := dl.pRedis.Del(dl.ctx, dl.lkName+keyName).Result()
	if err != nil {
		log.Println("DistributedUnlock err:", err.Error(), set)
	}
	return true
}
