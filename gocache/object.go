package gocache

// test 用
import (
	"errors"
	"github.com/tech-botao/logger"
)

type Object struct {
	ID string
	B  int
	C  []int
}

func (o *Object) GocacheKey() string {
	return "package:Object:" + o.ID
}

func (o *Object) GocacheCast(v interface{}) error {
	if vv, ok := v.(Object); ok {
		*o = vv
		return nil
	} else {
		return errors.New("can not cast interface{} to Object")
	}
}

func (*Object) Get(ID string) *Object {
	o := Object{ID: ID}
	err := Result(o.GocacheKey(), o.GocacheCast)
	if err != nil {
		logger.Error("[Object] Get() error", err)
		return nil
	}
	return &o
}
