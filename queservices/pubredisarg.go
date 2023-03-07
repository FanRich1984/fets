package queservices

type PubRedisArg struct {
	StreamKey string // 消息队列长度
	MaxLen    int64  // 消息队列长度
}
