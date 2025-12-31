/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package mq

import (
	"bailu/pkg/store"
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func Test_simple_mq(t *testing.T) {
	client := store.NewRedisClient(0)
	publisher := NewPublisher(client, "annmq")
	_, err := publisher.SendMsg(map[string]string{"msg": "hello"})
	if err != nil {
		t.Error(err)
	}

	consumer := NewRedisConsumer(&RedisStreamConsumerConfig{
		Client:              client,
		StreamName:          "annmq",
		ConsumerGroupName:   "trmqcg",
		StartCursor:         "0-0",
		BufferSize:          50,
		ReadSize:            50,
		ReadBlockTime:       30 * 1000,
		DataRecoverInterval: 10,
		PendingReadSize:     50,
		DeadMsgTime:         30 * 1000,
	})
	ctx, cancelFunc := context.WithCancel(context.Background())
	msgChan, err := consumer.GetMsgChan(ctx)
	if err != nil {
		t.Error(err)
	}

	go func() {

		for i := 0; i < 30; i++ {
			_, err := publisher.SendMsg(map[string]string{"msg": "hello" + strconv.Itoa(i)})
			if err != nil {
				t.Error(err)
			}
			time.Sleep(1 * time.Second)
		}

		cancelFunc()
		consumer.Wait()
	}()

	for msgs := range msgChan {
		msgStr := ""
		msgStr += fmt.Sprint("收到消息：[\n")
		for _, msg := range msgs {
			msgStr += fmt.Sprintf("%v,\n", msg)
		}
		msgStr += fmt.Sprint("]\n")
		fmt.Println(msgStr)

		err := consumer.AckMsg(ctx, msgs)
		if err != nil {
			t.Error(err)
		}
	}

	fmt.Println("end ...")

}
