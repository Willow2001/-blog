package session

import (
"bytes"
"encoding/gob"
"fmt"
"log"
"strconv"
"sync"
"time"

"github.com/go-redis/redis"
uuid "github.com/satori/go.uuid"
)

// redisSession redis session对象
type redisSession struct {
	// redis session id 对象
	id string
	// session 数据对象
	data map[string]interface{}
	// session 数据是否有更新
	modifyFlag bool
	// 过期时间
	expired int
	rwLock  sync.RWMutex
	client  *redis.Client
}

func NewRedisSession(id string, client *redis.Client) (session Session) {
	session = &redisSession{
		id:     id,
		data:   make(map[string]interface{}, 8),
		client: client,
	}
	return
}

func (r *redisSession) ID() string {
	return r.id
}

func (r *redisSession) Load() (err error) {
	data, err := r.client.Get(r.id).Bytes()
	if err != nil {
		log.Printf("get session data from redis by %s failed, err: %v\n", r.id, err)
		return
	}

	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err = dec.Decode(&r.data)
	if err != nil {
		log.Printf("gob decode session data failed, err: %v\n", err)
		return
	}
	return
}

func (r *redisSession) Get(key string) (value interface{}, err error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	value, ok := r.data[key]
	if !ok {
		err = fmt.Errorf("invalid key")
		return
	}
	return
}

func (r *redisSession) Set(key string, value interface{}) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	r.data[key] = value
	r.modifyFlag = true
}

func (r *redisSession) Del(key string) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	delete(r.data, key)
	r.modifyFlag = true
}

func (r *redisSession) SetExpired(expired int) {
	r.expired = expired
}

func (r *redisSession) Save() {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	if !r.modifyFlag {
		return
	}
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(r.data)
	if err != nil {
		log.Fatalf("gob encode r.data failed, err: %v\n", err)
		return
	}

	r.client.Set(r.id, buf.Bytes(), time.Second*time.Duration(r.expired))
	log.Printf("set data %v to redis.\n", buf.Bytes())
	r.modifyFlag = false
}

// redisSessionMgr redis Session管理器对象
type redisSessionMgr struct {
	session map[string]Session
	rwLock  sync.RWMutex
	client  *redis.Client
}

// NewRedisSessionMgr Redis SessionMgr类构造函数
func NewRedisSessionMgr() *redisSessionMgr {
	return &redisSessionMgr{
		session: make(map[string]Session, 1024),
	}
}

func (r *redisSessionMgr) Init(addr string, options ...string) (err error) {
	var (
		password string
		db       int
	)
	if len(options) == 1 {
		password = options[0]
	}

	if len(options) == 2 {
		password = options[0]
		db, err = strconv.Atoi(options[1])
		if err != nil {
			log.Fatalln("invalid redis DB param")
		}
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err = r.client.Ping().Result()
	if err != nil {
		return
	}
	return nil
}

func (r *redisSessionMgr) GetSession(sessionID string) (sd Session, err error) {
	sd = NewRedisSession(sessionID, r.client)
	err = sd.Load()

	if err != nil {
		return
	}

	r.rwLock.RLock()
	r.session[sessionID] = sd
	r.rwLock.RUnlock()
	return
}

func (r *redisSessionMgr) CreateSession() (sd Session) {
	sessionID := uuid.NewV4().String()
	sd = NewRedisSession(sessionID, r.client)
	r.session[sd.ID()] = sd
	return
}

func (r *redisSessionMgr) Clear(sessionID string) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	delete(r.session, sessionID)
}
