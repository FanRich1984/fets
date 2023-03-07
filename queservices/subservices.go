package queservices

// 处理消息对象
type SubServices interface {
	// 安全运行读取消息通道
	SafeRun(consumers int)
	// 停止
	Stop()
}
