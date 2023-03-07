package queservices

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func NewPubRedis(pRedis *redis.Client, pArg PubRedisArg) *PubRedis {
	ctx, stop := context.WithCancel(context.Background())
	return &PubRedis{
		pRedis:    pRedis,
		streamKey: pArg.StreamKey,
		ctx:       ctx,
		stop:      stop,
	}
}

// 基于Redis的消息队列发布服务对象
type PubRedis struct {
	pRedis    *redis.Client      // Redis客户端
	streamKey string             // 发布到队列的Key
	ctx       context.Context    // 发送消息的上下文
	maxLen    int64              // 消息队列长度
	stop      context.CancelFunc // 停止函数
}

// 将消息发布到Redis的消息队列里
func (pr *PubRedis) Push(values map[string]interface{}) error {
	res, err := pr.pRedis.XAdd(pr.ctx, &redis.XAddArgs{
		Stream: pr.streamKey,
		MaxLen: pr.maxLen,
		Approx: true,
		Values: values,
		ID:     "*",
	}).Result()

	fmt.Println("XADD Result :", res, err)

	if err != nil {
		log.Println("PubRedis ", pr.streamKey, " Error:", err.Error())
	}

	return err
}
