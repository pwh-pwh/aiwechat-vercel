package db

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
)

var (
	ChatDbInstance ChatDb        = nil
	RedisClient    *redis.Client = nil
	Cache          sync.Map
)

const (
	PROMPT_KEY = "prompt"
	MSG_KEY    = "msg"
	MODEL_KEY  = "model"
	TODO_KEY   = "todo"
)

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
	RedisClient = client
	return &RedisChatDb{
		client: client,
	}, nil
}

func (r *RedisChatDb) GetMsgList(botType string, userId string) ([]Msg, error) {
	result, err := r.client.Get(context.Background(), fmt.Sprintf("%v:%v:%v", MSG_KEY, botType, userId)).Result()
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
	msgTime := os.Getenv("MSG_TIME")
	//转换为数字
	msgT, err := strconv.Atoi(msgTime)
	if err != nil || msgT <= 0 {
		msgT = 30
	}
	r.client.Set(context.Background(), fmt.Sprintf("%v:%v:%v", MSG_KEY, botType, userId), res, time.Minute*time.Duration(msgT))
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

func GetValueWithMemory(key string) (string, bool) {
	value, ok := Cache.Load(key)
	if ok {
		return value.(string), ok
	}
	return "", false
}

func SetValueWithMemory(key string, val any) {
	Cache.Store(key, val)
}

func DeleteKeyWithMemory(key string) {
	Cache.Delete(key)
}

func GetValue(key string) (val string, err error) {
	val, flag := GetValueWithMemory(key)
	if !flag {
		if RedisClient == nil {
			return
		}
		val, err = RedisClient.Get(context.Background(), key).Result()
		SetValueWithMemory(key, val)
		return
	}
	return
}

func SetValue(key string, val any, expires time.Duration) (err error) {
	SetValueWithMemory(key, val)

	if RedisClient == nil {
		return
	}
	if expires == 0 {
		expires = time.Minute * 30
	}

	err = RedisClient.Set(context.Background(), key, val, expires).Err()

	return
}

func DeleteKey(key string) {
	DeleteKeyWithMemory(key)
	if RedisClient == nil {
		return
	}
	RedisClient.Del(context.Background(), key)
}

func DeleteMsgList(botType string, userId string) {
	RedisClient.Del(context.Background(), fmt.Sprintf("%v:%v:%v", MSG_KEY, botType, userId))
}

func SetPrompt(userId, botType, prompt string) {
	SetValue(fmt.Sprintf("%s:%s:%s", PROMPT_KEY, userId, botType), prompt, 0)
}

func GetPrompt(userId, botType string) (string, error) {
	return GetValue(fmt.Sprintf("%s:%s:%s", PROMPT_KEY, userId, botType))
}

func RemovePrompt(userId, botType string) {
	DeleteKey(fmt.Sprintf("%s:%s:%s", PROMPT_KEY, userId, botType))
}

// todolist format: "todo1|todo2|todo3"
func GetTodoList(userId string) (string, error) {
	tListStr, err := GetValue(fmt.Sprintf("%s:%s", TODO_KEY, userId))
	if err != nil && RedisClient == nil {
		return "", err
	}
	if tListStr == "" {
		return "todolist为空", nil
	}
	split := strings.Split(tListStr, "|")
	var sb strings.Builder
	for index, todo := range split {
		sb.WriteString(fmt.Sprintf("%d. %s\n", index+1, todo))
	}
	return sb.String(), nil
}

func AddTodoList(userId string, todo string) error {
	todoList, err := GetValue(fmt.Sprintf("%s:%s", TODO_KEY, userId))
	if err != nil && RedisClient == nil {
		return err
	}
	if todoList == "" {
		todoList = todo
	} else {
		todoList = fmt.Sprintf("%s|%s", todoList, todo)
	}
	return SetValue(fmt.Sprintf("%s:%s", TODO_KEY, userId), todoList, 0)
}

func DelTodoList(userId string, todoIndex int) error {
	todoList, err := GetValue(fmt.Sprintf("%s:%s", TODO_KEY, userId))
	if err != nil && RedisClient == nil {
		return err
	}
	todoList = strings.Split(todoList, "|")[todoIndex-1]
	return SetValue(fmt.Sprintf("%s:%s", TODO_KEY, userId), todoList, 0)
}

func SetModel(userId, botType, model string) error {
	if model == "" {
		DeleteKey(fmt.Sprintf("%s:%s:%s", MODEL_KEY, userId, botType))
		return nil
	} else {
		return SetValue(fmt.Sprintf("%s:%s:%s", MODEL_KEY, userId, botType), model, 0)
	}
}

func GetModel(userId, botType string) (string, error) {
	return GetValue(fmt.Sprintf("%s:%s:%s", MODEL_KEY, userId, botType))
}
