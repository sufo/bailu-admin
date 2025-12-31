/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package mq

import (
	"bailu/pkg/log"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"sync"
	"time"
)

const (
	DefaultStartCursor = "0-0"

	DefaultBufferSize = 50

	DefaultReadSzie = 50

	DefaultBlockTime = 30 * 1000

	DefaultDeadMsgTime = 30 * 1000

	DefaultDataRecoverInterval = 5

	//DefaultClaimInterval = 5

	DefaultClaimBatchSize = int64(100)

	DefaultMaxIdleTime = 60

	DefaultPendingReadSize = 50

	DefaultCheckConsumersInterval = 300
	DefaultConsumerTimeout        = 600
)

type StreamMsg struct {
	Stream string
	ID     string
	Values map[string]interface{}
}

var Consumer *RedisStreamConsumer

// RedisStreamConsumer 消费者
type RedisStreamConsumer struct {
	client *redis.Client
	// stream名称
	streamName string
	// 消费者组名称
	consumerGroupName string
	// 消费开始时游标的位置，0-0指定开头，传入ID指定特定ID，传入$指定结尾，当消费者组已经存在时不会进行新建，但是仍然会移动游标
	startCursor string
	// 消费通道缓冲区大小
	bufferSize int
	// 单次读取消息数量
	ReadSize int
	// 读取stream时阻塞时间，单位毫秒
	ReadBlockTime int
	// 数据修复协程运行间隔，单位秒
	dataRecoverInterval int
	// How many pending messages are claimed at most each claim interval.
	ClaimBatchSize int64
	// How long should we treat a pending message as claimable.
	MaxIdleTime int //单位秒
	// 单次读取pending消息数量
	PendingReadSize int
	// 未ack消息时间超出多久时视为死信进行重新消费，单位毫秒
	DeadMsgTime int
	// 用于等待所有内部协程终止
	wg sync.WaitGroup
	// 保存msgChan
	msgChan                   chan []*StreamMsg
	ShouldClaimPendingMessage func(redis.XPendingExt) bool
}

// RedisStreamConsumerConfig 消费者的配置信息
type RedisStreamConsumerConfig struct {
	Client *redis.Client
	// stream名称
	StreamName string
	// 消费者组名称
	ConsumerGroupName string
	// 消费开始时游标的位置，0-0指定开头，传入ID指定特定ID，传入$指定结尾，当消费者组已经存在时不会进行新建，但是仍然会移动游标
	StartCursor string
	// 消费通道缓冲区大小
	BufferSize int
	// 单次读取消息数量
	ReadSize int
	// 读取stream时阻塞时间，单位毫秒
	ReadBlockTime int
	// 数据修复协程运行间隔，单位秒
	DataRecoverInterval int

	// How many pending messages are claimed at most each claim interval.
	ClaimBatchSize int64
	// How long should we treat a pending message as claimable.
	MaxIdleTime int //单位秒

	// 单次读取pending消息数量
	PendingReadSize int
	// 未ack消息时间超出多久时视为死信进行重新消费，单位毫秒
	DeadMsgTime int

	// If this is set, it will be called to decide whether a pending message that
	// has been idle for more than MaxIdleTime should actually be claimed.
	// If this is not set, then all pending messages that have been idle for more than MaxIdleTime will be claimed.
	// This can be useful e.g. for tasks where the processing time can be very variable -
	// so we can't just use a short MaxIdleTime; but at the same time dead
	// consumers should be spotted quickly - so we can't just use a long MaxIdleTime either.
	// In such cases, if we have another way for checking consumers' health, then we can
	// leverage that in this callback.
	ShouldClaimPendingMessage func(redis.XPendingExt) bool
}

func NewRedisConsumer(config *RedisStreamConsumerConfig) *RedisStreamConsumer {
	config.setDefault()
	//return (*RedisStreamConsumer)(unsafe.Pointer(config))
	return &RedisStreamConsumer{
		client:                    config.Client,
		streamName:                config.StreamName,
		consumerGroupName:         config.ConsumerGroupName,
		startCursor:               config.StartCursor,
		bufferSize:                config.BufferSize,
		ReadSize:                  config.ReadSize,
		ReadBlockTime:             config.ReadBlockTime,
		dataRecoverInterval:       config.DataRecoverInterval,
		ClaimBatchSize:            config.ClaimBatchSize,
		MaxIdleTime:               config.MaxIdleTime,
		PendingReadSize:           config.PendingReadSize,
		DeadMsgTime:               config.DeadMsgTime,
		ShouldClaimPendingMessage: config.ShouldClaimPendingMessage,
	}
}

func (cc *RedisStreamConsumerConfig) setDefault() {
	if cc.BufferSize == 0 {
		cc.BufferSize = DefaultBufferSize
	}
	if cc.ReadSize == 0 {
		cc.ReadSize = DefaultReadSzie
	}
	if cc.ReadBlockTime == 0 {
		cc.ReadBlockTime = DefaultBlockTime
	}
	if cc.DataRecoverInterval == 0 {
		cc.DataRecoverInterval = DefaultDataRecoverInterval
	}
	if cc.ClaimBatchSize == 0 {
		cc.ClaimBatchSize = DefaultClaimBatchSize
	}
	if cc.MaxIdleTime == 0 {
		cc.MaxIdleTime = DefaultMaxIdleTime
	}
	if cc.StartCursor == "" {
		cc.StartCursor = DefaultStartCursor
	}
	if cc.DeadMsgTime == 0 {
		cc.DeadMsgTime = DefaultDeadMsgTime
	}
	if cc.PendingReadSize == 0 {
		cc.PendingReadSize = DefaultPendingReadSize
	}
}

// Wait 阻塞等待内部协程完成
func (c *RedisStreamConsumer) Wait() {
	c.wg.Wait()
	close(c.msgChan)
}

// GetMsgChan 读取消息
func (c *RedisStreamConsumer) GetMsgChan(ctx context.Context) (<-chan []*StreamMsg, error) {

	if c.msgChan != nil {
		return c.msgChan, nil
	}

	c.msgChan = make(chan []*StreamMsg, c.bufferSize)

	client := c.client

	//defer utils.CloseRedis(redis, consumer.logger)
	// 创建消费者组
	var err error
	_, err = client.XGroupCreate(ctx, c.streamName, c.consumerGroupName, c.startCursor).Result()
	// 隐藏已经存在消费者的错误，继续运行
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return nil, err
	}

	// 移动游标
	_, err = client.XGroupSetID(ctx, c.streamName, c.consumerGroupName, c.startCursor).Result()
	if err != nil {
		return nil, err
	}

	c.wg.Add(1)
	// mq 监听协程
	go func() {
		err := c.watchStream(context.WithValue(ctx, "name", "streamWatch"))
		if err != nil {
			fmt.Println("stream监听协程报错退出")
		} else {
			fmt.Println("stream监听协程退出")
		}
		c.wg.Done()
	}()

	c.wg.Add(1)
	go func() {
		err := c.dataRecover(context.WithValue(ctx, "name", "dataRecover"))
		if err != nil {
			fmt.Println("数据修复协程报错退出")
		} else {
			fmt.Println("数据修复协程退出")
		}
		c.wg.Done()
	}()

	return c.msgChan, nil
}

// 监听stream，读取消息
func (c *RedisStreamConsumer) watchStream(ctx context.Context) error {
	//defer utils.CloseRedis(redis, consumer.logger)
	exist := false
	for !exist {
		r, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    c.consumerGroupName,         // 消费者组的名称
			Consumer: "readConsumer",              // 消费者的名称
			Streams:  []string{c.streamName, ">"}, // Stream的名称和ID 最后的>符号代表这个消费者从最新的没被确认消费过的ID开始消费, 你也可以填任何具体id
			Count:    int64(c.ReadSize),           // 一次要读取的消息数量

			//a. 当stream中有数据时,拉取指定条数的数据,
			//b. 当stream中没有数据时,阻塞住,一旦有数据就会立刻拉取
			//c. BLOCK 0: 其中0表示超时时间无限,即一直阻塞,如果填1000就是阻塞1000ms,超时后没有获取到数据就返回nil
			Block: time.Duration(c.ReadBlockTime) * time.Millisecond, // 阻塞时间，0表示不阻塞
		}).Result()
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				// 此时发生异常，如果是超时超过创建时配置的超时，那么连接会关闭，需要重新申请连接
				fmt.Println(err)
				//utils.CloseRedis(redis, consumer.logger)
				//redis, err = consumer.pool.GetResource()
				//if err != nil {
				//	consumer.logger.Println(err)
				//	return err
				//}
			} else {
				fmt.Println(err)
				//return err
			}
			continue
		} else {

			result := c.GetXGroupReadResult(r)
			if len(result) == 0 {
				fmt.Println("消息为空")
			} else {
				// 推送给消费协程进行消费
				for _, xgroupreadResult := range result {
					// 写入channel
					c.msgChan <- xgroupreadResult
				}
			}
		}
		if MaxLen != 0 {
			c.xTrimMaxLen(ctx)
		}
		// 监听关闭
		select {
		case <-ctx.Done():
			exist = true
			break
		default:
			break
		}
	}
	return nil
}

// 数据修复，死信消费
func (c *RedisStreamConsumer) dataRecover(ctx context.Context) error {
	//defer utils.CloseRedis(redis, consumer.logger)

	// 每隔一段读取未进行ack的数据，将等待时间超时的数据重新推送到消费者队列中
	tick := time.NewTicker(time.Duration(c.dataRecoverInterval) * time.Second)
	exits := false
OUTER_LOOP:
	for !exits {

		if now, ok := <-tick.C; ok {
			fmt.Println(now)

			xps, err := c.client.XPendingExt(ctx,
				&redis.XPendingExtArgs{
					c.streamName, c.consumerGroupName,
					time.Duration(c.MaxIdleTime) * time.Second,
					"-", "+", int64(c.ClaimBatchSize), "",
				}).Result()
			if err != nil {
				log.L.Error("xpendingext fail", err)
				continue
			}

			for _, xp := range xps {
				shouldClaim := xp.Idle >= time.Duration(c.MaxIdleTime)
				if shouldClaim && c.ShouldClaimPendingMessage != nil {
					shouldClaim = c.ShouldClaimPendingMessage(xp)
				}
				if shouldClaim {
					// assign the ownership of a pending message to the current consumer
					xm, err := c.client.XClaim(ctx, &redis.XClaimArgs{
						Stream:   c.streamName,
						Group:    c.consumerGroupName,
						Consumer: "recoverConsumer",
						// this is important: it ensures that 2 concurrent subscribers
						// won't claim the same pending message at the same time
						MinIdle:  time.Duration(c.MaxIdleTime),
						Messages: []string{xp.ID},
					}).Result()
					if err != nil {
						log.L.Error("xpendingext fail", err)
						continue OUTER_LOOP
					}
					if len(xm) > 0 {
						var ms = make([]*StreamMsg, 0)
						for _, msg := range xm {
							ms = append(ms, &StreamMsg{c.streamName, msg.ID, msg.Values})
						}
						c.msgChan <- ms
					}
				}
			}
		}
		// 监听关闭
		select {
		case <-ctx.Done():
			tick.Stop()
			exits = true
			break
		default:
			break
		}
	}
	return nil
}

// AckMsg ack消息
func (c *RedisStreamConsumer) AckMsg(ctx context.Context, msgs []*StreamMsg) error {
	//redis, _ := consumer.pool.GetResource()
	//defer utils.CloseRedis(redis, consumer.logger)

	// 发送xack
	args := make([]string, 0)
	for _, msg := range msgs {
		args = append(args, msg.ID)
	}
	_, err := c.client.XAck(ctx, c.streamName, c.consumerGroupName, args...).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisStreamConsumer) GetXGroupReadResult(streams []redis.XStream) [][]*StreamMsg {
	var xStream = make([][]*StreamMsg, 0)
	for _, stream := range streams {
		streamName := stream.Stream
		fmt.Println("获取消息 长度", len(stream.Messages))
		var msgs = make([]*StreamMsg, 0)
		for _, message := range stream.Messages {
			fmt.Println("读取到Stream: streamName=", streamName, " messageID=", message.ID, " messageValues=", message.Values)
			msg := &StreamMsg{streamName, message.ID, message.Values}
			msgs = append(msgs, msg)
		}
		xStream = append(xStream, msgs)
	}
	return xStream
}

func (c *RedisStreamConsumer) xTrimMaxLen(ctx context.Context) {
	// 修剪Stream
	trimCmd := c.client.XTrimMaxLen(ctx, c.streamName, MaxLen)
	_, err := trimCmd.Result()
	if err != nil {
		fmt.Println("修剪Stream失败:", err)
		//return
	}
	//fmt.Println("修剪Stream结果:", trimmed)
}

func (c *RedisStreamConsumer) HandlerMsg(ctx context.Context) {
	msgChan, err := c.GetMsgChan(ctx)
	if err != nil {
		log.L.Error(err)
	}
	for msgs := range msgChan {
		msgStr := ""
		msgStr += fmt.Sprint("收到消息：[\n")
		for _, msg := range msgs {
			msgStr += fmt.Sprintf("%v,\n", msg)
		}
		msgStr += fmt.Sprint("]\n")
		fmt.Println(msgStr)

		err := c.AckMsg(ctx, msgs)
		if err != nil {
			log.L.Error(err)
		}
	}
}
