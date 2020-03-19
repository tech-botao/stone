package gocache

import (
	"context"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/tech-botao/logger"
	"time"
)

type (
	// 这个是一个容器，以便以后进行扩展
	Client struct {
		ctx   context.Context
		cache *cache.Cache
	}

	Receiver interface {
		GocacheKey() string
		GocacheCast(v interface{}) error
	}
)

var C *Client

func New(ctx context.Context) *Client {
	return &Client{
		ctx:   ctx,
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func init() {
	C = New(context.Background())
	Test()
}

// SetContext 可以修改相关Context
func SetContext(ctx context.Context) *Client {
	C.ctx = ctx
	return C
}

// SetCache 可以替换GoCache
func SetCache(cache2 *cache.Cache) *Client {
	C.cache = cache2
	return C
}

// Save 保存的是指针
func Save(key string, v interface{}) {
	C.cache.Set(key, v, cache.NoExpiration)
}

// Update 保存数据
func Update(key string, v interface{}) {
	C.cache.Set(key, v, cache.NoExpiration)
}

// Increment, key:cnt
func Increment(key string) (int64, error) {
	key = cntIndex(key)
	if _, ok := C.cache.Get(key); !ok{
		C.cache.Set(key, int64(1), cache.NoExpiration)
		return 1, nil
	}
	return C.cache.IncrementInt64(key, int64(1))
}

// Decrement, key:cnt
func Decrement(key string) (int64, error) {
	key = cntIndex(key)
	if _, ok := C.cache.Get(key); !ok{
		C.cache.Set(key, int64(0), cache.NoExpiration)
		return 0, nil
	}
	return C.cache.DecrementInt64(key, int64(1))
}

func cntIndex(key string) string {
	return key + ":cnt"
}

// 得到累计值
func GetCnt(key string) int64 {
	key = cntIndex(key)
	v, ok := Get(key)
	if !ok {
		return 0
	}
	return v.(int64)
}

func Get(key string) (interface{}, bool) {
	 x, ok := C.cache.Get(key)
	 if !ok {
	 	logger.Error("[gocache] get error", errors.New("not data, key =" + key))
	 	return nil, false
	 }
	 if x == nil {
		 logger.Error("[gocache] get error", errors.New("data is nil, key =" + key))
		 return nil, false
	 }
	 return x, true
}

func Result(key string, cast func(v interface{}) error) error {
	v, ok := Get(key)
	if !ok {
		return errors.New("[gocache] not data, key =" + key)
	}
	return cast(v)
}

func Delete(key string) {
	C.cache.Delete(key)
}

// 一个外部接口， 可以直接调用gocache的方法
func Do(fn func() error) error {
	return fn()
}

// Test
func Test() {
	if C == nil {
		logger.Panic("[gocache] init fail", fmt.Errorf("C is nil"))
	}
	key := "gocache.test"
	value := time.Now().Unix()

	Save(key, value)
	if got, ok := Get(key); !ok {
		logger.Panic("[gocache] init fail", fmt.Errorf("key [%s] not found", key))
		if got != value {
			logger.Panic("[gocache] init fail",fmt.Errorf("value not right, want: %d, got:%d", value, got))
		}
	}
}

