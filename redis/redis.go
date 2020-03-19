package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"github.com/tech-botao/atom"
	"github.com/tech-botao/logger"
	"os"
	"time"
)

type (
	Config struct {
		URL         string
		OPT         *redis.Options
		ExpiredTime time.Duration
	}
	Client struct {
		ctx     context.Context
		conn    *redis.Client
		config  Config
		OnError func(error)
		Encoder func(interface{}) ([]byte, error)
		Decoder func([]byte, interface{}) error
	}
	Builder struct {
		config Config
	}
)

func NewBuilder() *Builder {
	return &Builder{
		config: Config{
			ExpiredTime: -1,
		},
	}
}

// For example redis://password:@localhost:6379/1
func (b *Builder) URL(url string) *Builder {
	b.config.URL = url
	opt, err := redis.ParseURL(url)
	if err != nil {
		logger.Panic("[stone.redis] config error, url="+url, err)
	}
	b.config.OPT = opt
	return b
}

func (b *Builder) ExpiredTime(d time.Duration) *Builder {
	b.config.ExpiredTime = d
	return b
}

func (b *Builder) Build(ctx context.Context) *Client {
	conn := redis.NewClient(b.config.OPT)
	c := Client{
		ctx:     ctx,
		conn:    conn,
		config:  b.config,
		OnError: DefaultError,
		Encoder: DefaultEncoder,
		Decoder: DefaultDecoder,
	}
	if _, err := c.conn.Ping().Result(); err != nil {
		logger.Panic("[stone.redis] connect error", err)
		return nil
	}
	return &c
}

var _url = ""
var _expr = -1

func SetUrl(url string) {
	_url = url
}

func SetExpired(expr int) {
	_expr = expr
}

// 初始化这个包
func Setup(ctx context.Context) {
	if len(_url) == 0 {
		_url = os.Getenv("REDIS_URL")
		if len(_url) == 0 {
			logger.Panic("[stone.redis] setting env like, export REDIS_URL=redis://{password}:@{host}:{port}/{db}\n", nil)
		}
	}
	C = NewBuilder().
		URL(_url).
		ExpiredTime(time.Duration(_expr) * time.Second).
		Build(ctx)
}

var C *Client

func DefaultError(err error) {
	logger.Error("[stone.redis] error", err)
}

func DefaultEncoder(input interface{}) ([]byte, error) {
	return atom.EncodeToBytes(input)
}

func DefaultDecoder(b []byte, output interface{}) error {
	return atom.DecodeFromByte(b, &output)
}

// ======= 外部接口 =========
func Save(key string, input interface{}) error {
	return SaveWith(key, input, C.config.ExpiredTime)
}

func SaveWith(key string, input interface{}, ttl time.Duration) error {
	var err error
	var b []byte
	defer func() {
		if err != nil {
			C.OnError(err)
		}
	}()
	if input == nil {
		err = errors.New("[stone.redis] data is nil, key =" + key)
		return err
	}
	b, err = C.Encoder(input)
	cmd := C.conn.Set(key, b, ttl)
	err = cmd.Err()
	return err
}

func Result(key string, out interface{}) error {
	cmd := C.conn.Get(key)
	var err error
	var b []byte
	defer func() {
		if err != nil {
			C.OnError(err)
		}
	}()

	if cmd.Err() != nil {
		err = cmd.Err()
		return err
	}
	b, err = cmd.Bytes()
	if err != nil {
		return err
	}
	err = C.Decoder(b, &out)
	return err
}

func Keys(key string) []string {
	cmd := C.conn.Keys(key + "*")
	if cmd.Err() != nil {
		C.OnError(cmd.Err())
		return []string{}
	}
	return cmd.Val()
}

func ListAs(prefix string, next func([]byte) error) {
	pipe := C.conn.Pipeline()
	cmds := make([]*redis.StringCmd, 0)
	keys := Keys(prefix+"*")
	for _, key := range keys {
		cmds = append(cmds, pipe.Get(key))
	}
	_, err := pipe.Exec()
	if err != nil {
		C.OnError(err)
	}
	for _, cmd := range cmds {
		b, err := cmd.Bytes()
		if err != nil {
			C.OnError(err)
			return
		}
		err = next(b)
		if err != nil {
			C.OnError(err)
			return
		}
	}
}

// ============= Counter ==============

func Incr(key string) (int64, error) {
	cmd := C.conn.Incr(key)
	return cmd.Result()
}

func GetCounter(key string) (int64, error) {
	cmd := C.conn.IncrBy(key, 0)
	return cmd.Val(), cmd.Err()
}

func Decr(key string) (int64, error) {
	cmd := C.conn.Decr(key)
	return cmd.Result()
}

func Del(key string) error {
	cmd := C.conn.Del(key)
	return cmd.Err()
}

// prefix
func DelAs(prefix string) error {
	return C.conn.Del(Keys(prefix+"*")...).Err()
}

func Dump(key string) error {
	cmd := C.conn.Dump(key)
	s, err := cmd.Result()
	if err != nil {
		C.OnError(err)
		return err
	} else {
		logger.Info("[stone.redis] dump for key:"+key, s)
		return nil
	}
}
