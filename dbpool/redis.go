package dbpool

import (
	"strings"
	"time"
	"vote_cj/config"

	"github.com/garyburd/redigo/redis"
)

type RedisOption struct {
	Db int
}

var Redispool *redis.Pool

func init() {
	Redispool = &redis.Pool{ //实例化一个连接池
		MaxIdle:     config.Redis.MaxIdle,                    //最初的连接数量
		MaxActive:   config.Redis.MaxActive,                  //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: time.Duration(config.Redis.IdleTimeout), //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			var sbd strings.Builder
			sbd.WriteString(config.Redis.Host)
			sbd.WriteString(":")
			sbd.WriteString(config.Redis.Port)
			return redis.Dial("tcp", sbd.String(), redis.DialDatabase(config.Redis.Index), redis.DialPassword(config.Redis.Password))
		},
	}
}
