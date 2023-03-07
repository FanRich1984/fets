package queservices

type SubGroupRedisArg struct {
	StreamKey  string // 消息队列键值
	GroupName  string // 消费组名称
	FnCallBack func(map[string]interface{}) error
}
