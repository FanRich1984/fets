package queservices

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type SubRedis struct {
	pRedis     *redis.Client
	streamKey  string
	groupName  string
	name       string
	ctx        context.Context
	readArg    *redis.XReadGroupArgs
	blnRun     bool                               // 是否在运行中
	fnCallBack func(map[string]interface{}) error // 收到消息的回调处理
}

// 开始侦听消息通道数据
func (sr *SubRedis) SyncRead() {
	if sr.blnRun {
		log.Println(sr.name, "已在运行中")
		return
	}

	sr.blnRun = true
	log.Println(sr.name, "开始侦听", sr.streamKey, "消息...")
	for {
		select {
		case <-sr.ctx.Done():
			fmt.Println(sr.name, "收到终止消息")
			sr.blnRun = false
			return
		default:
			// 读取消息队里有消息，如果读取发生错误则等待2秒后再次尝试
			if err := sr.Read(); err != nil {
				time.Sleep(time.Second * 2)
				// 处理错误发生时
				// log.Println()
				fmt.Println("读取Redis里的Stream时发生错误:", err.Error())
				// 判断是不是不存在这个GROUP信息
				strErrInfo := err.Error()
				if strings.Contains(strErrInfo, "NOGROUP") {
					// 执行创建组
					fmt.Println("执行了重新创建组:", sr.streamKey, sr.groupName)
					sr.pRedis.XGroupCreate(context.Background(), sr.streamKey, sr.groupName, "0").Result()
				}
			}
		}
	}
}

// 读取消息队列里的消息
func (sr *SubRedis) Read() error {
	vals, err := sr.pRedis.XReadGroup(sr.ctx, sr.readArg).Result()
	if err != nil {
		// fmt.Println("读取消息队列里的数据发生错误:", err.Error())
		return err
	}

	for i := 0; i < len(vals); i++ {
		msgs := vals[i].Messages
		for j := 0; j < len(msgs); j++ {
			msg := msgs[j]
			if sr.fnCallBack != nil {
				err := sr.fnCallBack(msg.Values)
				if err != nil {
					log.Println(sr.streamKey, sr.groupName, sr.name, "-执行回调发生错误:", err.Error())
				}
			}
		}
	}

	return nil
}

// 列出当前所有未确认的消息
func (sr *SubRedis) ListPending() error {
	info, err := sr.pRedis.XPending(context.Background(), sr.streamKey, sr.groupName).Result()
	fmt.Println(info, err)
	return err
}
