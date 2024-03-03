package db

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
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

type ChatDb interface {
	GetMsgList(userId string) ([]Msg, error)
	SetMsgList(userId string, msgList []Msg)
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

func (r *RedisChatDb) GetMsgList(userId string) ([]Msg, error) {
	result, err := r.client.Get(context.Background(), userId).Result()
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

func (r *RedisChatDb) SetMsgList(userId string, msgList []Msg) {
	res, err := sonic.Marshal(msgList)
	if err != nil {
		fmt.Println(err)
		return
	}
	r.client.Set(context.Background(), userId, res, time.Minute*30)
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
