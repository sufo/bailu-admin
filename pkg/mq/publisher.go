/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package mq

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var MaxLen = int64(100 * 1) // 最多存储消息数，依据FIFO原则，自动删除超过最长长度的消息

var Publisher *RedisStreamPublisher

type RedisStreamPublisher struct {
	Client *redis.Client
	//stream名称
	StreamName string
}

func NewPublisher(client *redis.Client, streamName string) *RedisStreamPublisher {
	return &RedisStreamPublisher{client, streamName}
}

// SendMsg 发送消息
func (p *RedisStreamPublisher) SendMsg(msgMap any) (any, error) {

	//defer publisher.client.Close()
	//args := []string{p.streamName, "*"}
	//for key, val := range msgMap {
	//	args = append(args, key, val)
	//}

	cmd := p.Client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: p.StreamName,
		MaxLen: MaxLen,
		Approx: true,
		ID:     "",
		Values: msgMap,
	})
	return cmd.Result()
}

//func (p *RedisStreamPublisher) Close() error {
//	p.closeMutex.Lock()
//	defer p.closeMutex.Unlock()
//
//	if p.closed {
//		return nil
//	}
//	p.closed = true
//
//	if err := p.client.Close(); err != nil {
//		return err
//	}
//	return nil
//}
