## Redis

> Wrap for redis

## Key Rule

> key最好尽可能的短，并且需要能看懂
> 保持下面的格式的话，可以通过*读取List

```golang
keyFormat = redis:model:index
```

## TODO

### 支持其它数据结构

```golang
// ======= TODO =========
// PIPELINE > 提高多个数据读取和存储的速度
// XADD > 时间序列的数据结构 (https://redis.io/topics/streams-intro)
// LIST > 保存数组(https://redis.io/commands/lrange)
// HSET > Map
```

### 支持Pub/Sub

