package lwqRedis

import (
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

/**
 * @author miku
 * @date 2019/11/30
 */
var RedisPool *redis.Pool

func InitRedis() (err error) {
	RedisPool = &redis.Pool{
		MaxIdle:     0,
		MaxActive:   100,
		IdleTimeout: 10,
		Dial: func() (redis.Conn, error) {
			//这里将 conn 类型赋值给接口
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}

	conn := RedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("PING")
	if err != nil {
		logs.Error("ping redis failed , err:%v", err)
		return
	}
	return
}
