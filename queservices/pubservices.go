package queservices

// 发布接口
type PubServices interface {
	// 发布接口对象
	Push(values map[string]interface{}) error
}
