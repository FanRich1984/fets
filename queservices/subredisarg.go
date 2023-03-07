package queservices

import "github.com/go-redis/redis/v8"

type SubRedisArg struct {
	StreamKey    string                // 消息队列键值
	GroupName    string                // 消费组名称
	ConsumerName string                // 消费组里消者名称
	ReadArg      *redis.XReadGroupArgs // 读取参数
}
