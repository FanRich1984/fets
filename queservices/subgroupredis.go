package queservices

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewSubServices(pRedis *redis.Client, arg SubGroupRedisArg) SubServices {
	return NewSubGroupRedis(pRedis, arg)
}

// 构造对象
func NewSubGroupRedis(pRedis *redis.Client, arg SubGroupRedisArg) *SubGroupRedis {
	return &SubGroupRedis{
		pRedis:     pRedis,
		streamKey:  arg.StreamKey,
		groupName:  arg.GroupName,
		subs:       make([]*SubRedis, 0, 2),
		fnCallBack: arg.FnCallBack,
	}
}

type SubGroupRedis struct {
	pRedis     *redis.Client                      // Go-Redis客户端
	streamKey  string                             // 消息通道的KEY
	groupName  string                             // 读取组的名称
	ctx        context.Context                    // 消费组里的上下文
	stop       context.CancelFunc                 // 执行停止方法
	subs       []*SubRedis                        // 这个消费组里的消费者
	fnCallBack func(map[string]interface{}) error // 处理消息队列的信息
}

// 创建消费组
func (sg *SubGroupRedis) CreateGroup() error {
	_, err := sg.pRedis.XGroupCreate(context.Background(), sg.streamKey, sg.groupName, "0").Result()
	return err
}

// 判断当前是否有这个Key的消息队列
func (sg *SubGroupRedis) HasStreamKey() (int64, error) {
	result, err := sg.pRedis.Exists(context.Background(), sg.streamKey).Result()
	return result, err
}

// 创建消费组里的消费者
func (sg *SubGroupRedis) CreateConsum(ctx context.Context) *SubRedis {
	// 生成消费者名称
	consumerName := fmt.Sprintf("%s_%d", sg.groupName, len(sg.subs))
	ptSub := &SubRedis{
		pRedis:    sg.pRedis,
		streamKey: sg.streamKey,
		groupName: sg.groupName,
		name:      consumerName,
		ctx:       ctx,
		readArg: &redis.XReadGroupArgs{
			Group:    sg.groupName,
			Consumer: consumerName,
			Streams:  []string{sg.streamKey, ">"},
			Count:    1,
			Block:    0,
			NoAck:    true,
		},
		fnCallBack: sg.fnCallBack,
	}
	sg.subs = append(sg.subs, ptSub)
	return ptSub
}

func (sg *SubGroupRedis) RunRead() {
	sg.ctx.Done()
}

// 查看当前消息组里有多少个消费者
func (sg *SubGroupRedis) ListConsumer() {
	info, err := sg.pRedis.XInfoConsumers(context.Background(), sg.streamKey, sg.groupName).Result()
	if err != nil {
		fmt.Println("XInfoConsumers error:", err.Error())
		return
	}
	fmt.Println("XInfoConsumers Info:", info)
}

// 查看组里的所有信息
func (sg *SubGroupRedis) ListGroups() {
	info, err := sg.pRedis.XInfoGroups(context.Background(), sg.streamKey).Result()
	if err != nil {
		fmt.Println("XInfoGroups error:", err.Error())
		return
	}
	fmt.Println("XInfoGroups error:", info)
}

// 查看当前消息对列的信息
func (sg *SubGroupRedis) InfoStream() {
	info, err := sg.pRedis.XInfoStream(context.Background(), sg.streamKey).Result()
	if err != nil {
		fmt.Println("XInfoStream error:", err.Error())
	} else {
		fmt.Println("XInfoStream info:", info)
	}
}

// 删除当前的消费组对象
func (sg *SubGroupRedis) Destroy() error {
	val, err := sg.pRedis.XGroupDestroy(context.Background(), sg.streamKey, sg.groupName).Result()
	fmt.Println("执行结果:", val)
	return err
}

// 停止所有消费者
func (sg *SubGroupRedis) Stop() {
	// 清除当前所有的消费者
	sg.subs = make([]*SubRedis, 0, 2)
	sg.stop()
}

func (sg *SubGroupRedis) Run() {
	for i := 0; i < len(sg.subs); i++ {
		go sg.subs[i].SyncRead()
	}
}

// 安全启动消费组里的消费者
func (sg *SubGroupRedis) SafeRun(consumers int) {
	// 判断当前是否有这个Key
	for {
		if v, err := sg.HasStreamKey(); err == nil {
			if v > 0 {
				break
			}
			fmt.Println("键值不存在等待2秒")
			time.Sleep(time.Second * 2)
		} else {
			fmt.Println("查看键值发生错误:", err.Error())
			return
		}
	}

	// 创建消费组
	sg.CreateGroup()
	sg.ctx, sg.stop = context.WithCancel(context.Background())
	// 创建消费组
	i := len(sg.subs)
	for ; i < consumers; i++ {
		sg.CreateConsum(sg.ctx)
	}
	// 执行读取
	sg.Run()
}
