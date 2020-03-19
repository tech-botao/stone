package redis

import (
	"context"
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/tech-botao/atom"
	"github.com/tech-botao/logger"
	"reflect"
	"testing"
)

func init() {
	// 1. 初始方式， 手动
	SetUrl("redis://:@localhost:6379/0")
	SetExpired(-1)

	// 2. 初始方式， 设置环境变量
	// os.Setenv("REDIS_URL", "redis://:@localhost:6379/0")
	// Setup(context.Background())
	Setup(context.Background())
}

type obj struct {
	A string
	B int
}

func (o obj) Index() string {
	return "obj:" + o.A
}

func ExampleRedisFlow() {

	o := obj{A: "123", B: 123,}

	err := Save(o.Index(), o)
	if err != nil {
		logger.Error("[stone.redis] Save()", err)
		return
	}

	var o2 obj
	err = Result(o.Index(), &o2)
	if err != nil {
		logger.Error("[stone.redis] Result()", err)
		return
	}
	logger.Info("o2", o2)

	prefix := "list:obj"
	for i := 0; i < 100; i++ {
		o2.B++
		err = Save(fmt.Sprintf("%s:%d", prefix, i), o2)
		if err != nil {
			pp.Println(err)
			return
		}
	}

	var o3 = make([]obj, 0)
	var temp obj
	ListAs(prefix, func(b []byte) error {
		err := atom.DecodeFromByte(b, &temp)
		if err != nil {
			return err
		}
		o3 = append(o3, temp)
		return nil
	})
	logger.Info("[test redis] ListAs()", o3)

	// Incre
	counterKey := "counter"
	Incr(counterKey)
	Incr(counterKey)
	Incr(counterKey)
	fmt.Println(GetCounter(counterKey))
	Decr(counterKey)
	Decr(counterKey)
	Decr(counterKey)
	fmt.Println(GetCounter(counterKey))

	pp.Println(Dump(counterKey))

	// output:
	// 3 <nil>
	// 0 <nil>
}

func TestBuilder_Build(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestBuilder_DB(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestBuilder_ExpiredTime(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestBuilder_Host(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestBuilder_Port(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestConfig_GetAddr(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestDefaultDecoder(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestDefaultEncoder(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestDefaultError(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestNewBuilder(t *testing.T) {
	t.Skip("too simple dont test")
}

func TestResult(t *testing.T) {
	type args struct {
		index string
		out   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"int_ok", args{"test:1", 1}, false},
		{"obj_ok", args{"test:object:1", map[string]interface{}{"a": "123"}}, false},
		{"array_ok", args{"test:array:1", []string{"a", "123"}}, false},
		{"nil_ng", args{"test:array:nil", nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Save(tt.args.index, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var out interface{}
			if err := Result(tt.args.index, &out); (err != nil) != tt.wantErr {
				t.Errorf("Result() error = %v, wantErr %v", err, tt.wantErr)
			}
			if reflect.DeepEqual(tt.args.out.(interface{}), out) == false {
				//pp.Println(tt.args.out, out)
				t.Errorf("Result() expect = %v, out = %v", tt.args.out, out)
			}
		})
	}
}

func TestSave(t *testing.T) {
	t.Skip("include test Result() dont test")
}

func TestSaveWith(t *testing.T) {
	t.Skip("include test Result() dont test")
}
