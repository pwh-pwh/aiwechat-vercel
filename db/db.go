package db

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/pwh-pwh/aiwechat-vercel/chat"
	"github.com/sashabaranov/go-openai"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
)

var ChatDbInstance ChatDb = nil

func init() {
	db, err := GetChatDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	ChatDbInstance = db
}

type Msg struct {
	Role string
	Msg  string
}

type ChatMsg interface {
	openai.ChatCompletionMessage | chat.QwenMessage | chat.SparkMessage
}

func GetMsgListWithDb[T ChatMsg](botType, userId string, msg T, f func(msg T) Msg, f2 func(msg Msg) T) []T {
	if ChatDbInstance != nil {
		list, err := ChatDbInstance.GetMsgList(botType, userId)
		if err == nil {
			list = append(list, f(msg))
		}
		r := make([]T, 0)
		for _, msg := range list {
			r = append(r, f2(msg))
		}
		return r
	}
	return []T{msg}
}

func SaveMsgListWithDb[T ChatMsg](botType, userId string, msgList []T, f func(msg T) Msg) {
	if ChatDbInstance != nil {
		go func() {
			list := make([]Msg, 0)
			for _, msg := range msgList {
				list = append(list, f(msg))
			}
			ChatDbInstance.SetMsgList(botType, userId, list)
		}()
	}
}

type ChatDb interface {
	GetMsgList(botType string, userId string) ([]Msg, error)
	SetMsgList(botType string, userId string, msgList []Msg)
}

type RedisChatDb struct {
	client *redis.Client
}

func NewRedisChatDb(url string) (*RedisChatDb, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("redis url error")
	}
	options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	client := redis.NewClient(options)
	return &RedisChatDb{
		client: client,
	}, nil
}

func (r *RedisChatDb) GetMsgList(botType string, userId string) ([]Msg, error) {
	result, err := r.client.Get(context.Background(), fmt.Sprintf("%v:%v", botType, userId)).Result()
	if err != nil {
		return nil, err
	}
	var msgList []Msg
	err = sonic.Unmarshal([]byte(result), &msgList)
	if err != nil {
		return nil, err
	}
	return msgList, nil
}

func (r *RedisChatDb) SetMsgList(botType string, userId string, msgList []Msg) {
	res, err := sonic.Marshal(msgList)
	if err != nil {
		fmt.Println(err)
		return
	}
	r.client.Set(context.Background(), fmt.Sprintf("%v:%v", botType, userId), res, time.Minute*30)
}

func GetChatDb() (ChatDb, error) {
	kvUrl := os.Getenv("KV_URL")
	if kvUrl == "" {
		return nil, errors.New("请配置KV_URL")
	} else {
		db, err := NewRedisChatDb(kvUrl)
		if err != nil {
			return nil, err
		}
		return db, nil
	}
}
