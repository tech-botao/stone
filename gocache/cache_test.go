package gocache

import (
	"github.com/k0kubun/pp"
	"github.com/patrickmn/go-cache"
	"github.com/tech-botao/logger"
	"testing"
)


func ExampleFLow() {

	c := cache.New(cache.DefaultExpiration, cache.NoExpiration)

	u := Object{
		ID: "A",
		B:  100,
		C:  []int{1, 2, 3, 4, 5, 6},
	}

	err := c.Add("foo:bar", u, cache.NoExpiration)
	if err != nil {
		logger.Error("[cache] add error", err)
	}
	// 对一个已有的Key进行操作会报错
	//err = C.Add("foo:bar", u, cache.NoExpiration)
	//if err != nil {
	//	logger.Error("[cache] add2 error", err)
	//}

	c.Set("foo:bar:set", u, cache.NoExpiration)
	u.B = 200
	err = c.Replace("foo:bar:set", u, cache.NoExpiration)
	if err != nil {
		logger.Error("[cache] replace error", err)
	}
	_, _ = pp.Println(c.GetWithExpiration("foo:bar:set"))

	logger.Info("items", c.Items())
	c.Flush() // 清空一起的大招
	logger.Info("items", c.Items())

	// output:

}

func ExampleCache() {

	u := Object{
		ID: "A",
		B:  100,
		C:  []int{1, 2, 3, 4, 5, 6},
	}

	Save(u.GocacheKey(), u)
	var u2 Object

	err := Result(u.GocacheKey(), u2.GocacheCast)
	pp.Println(u2, err)

	pp.Println(Increment(u.GocacheKey()))
	pp.Println(Increment(u.GocacheKey()))
	pp.Println(Increment(u.GocacheKey()))

	pp.Println(Decrement(u.GocacheKey()))
	pp.Println(Decrement(u.GocacheKey()))
	pp.Println(Decrement(u.GocacheKey()))

	pp.Println(GetCnt(u.GocacheKey()))
	// penetrate

	var u3 Object
	err = Result(u.GocacheKey(), u3.GocacheCast)
	pp.Println(u3, err)

	u4 := new(Object).Get("A")
	pp.Println("U4", u4)

	_ = Do(func() error {
		pp.Println(C.cache.ItemCount())
		return nil
	})
	// output:

}

func TestDecrement(t *testing.T) {
	t.Skip("test for example Cache")
}

func TestDelete(t *testing.T) {
	t.Skip("pure wrap cache.Delete()")
}

func TestGet(t *testing.T) {
	t.Skip("pure wrap cache.Delete()")
}

func TestIncrement(t *testing.T) {
	t.Skip("pure wrap cache.Increment()")
}

func TestNew(t *testing.T) {
	t.Skip("pure wrap cache.New()")
}

func TestSave(t *testing.T) {
	type args struct {
		u Object
	}
	tests := []struct {
		name string
		args args
	}{
		{"ok", args{Object{ID: "1", B: 0, C: nil,}}},
		{"ok", args{Object{ID: "1", B: 20, C: []int{1, 2, 3, 4, 5},}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Save(tt.args.u.GocacheKey(), tt.args.u)
		})
	}
}

func TestSetCache(t *testing.T) {
	t.Skip("pure accessor")
}

func TestSetContext(t *testing.T) {
	t.Skip("pure accessor")
}

func TestUpdate(t *testing.T) {
	t.Skip("pure wrap save")
}
