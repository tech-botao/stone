# Stone

> 存储相关类

## 实现功能

### gocache 

> 内存缓存， 很快， 关闭服务会消失
> 无需初始化

```golang

// 存储数据
gocache.Save(key, obj)
gocache.Result(key, castFunc)


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

obj := new(Object).Get("ID123")

// Counter
gocache.Increment(key)
gocache.Decrement(key)
gocache.GetCnt(key)

// Do
err := gocache.Do(func() err {
    // do something
    return nil
})


```

### redis

> KV store, 数据持久化，高IO
> 支持Json编码
> 通过OnError(err)自定义错误句柄
> 通过Encoder, Decoder可以自定义编解码方法

```golang
	o := obj{
		A: "123",
		B: 123,
	}

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
	Decr(counterKey)
	GetCounter(counterKey)
```

