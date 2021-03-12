package cache

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"testing"
)

/*
1.在当前包或父包创建config.toml
2.[redis]
    default = "192.168.174.129:6379,0,godz1995?idleTimeout=600"
3.直接g.Redis()获取Redis cli，配置会引用default项的内容
*/
func TestRedis(t *testing.T) {
	redis := g.Redis()
	_, _ = redis.Do("set", "nigger", "black nigger slave")
	value, _ := redis.DoVar("get", "k")
	fmt.Println(value)
}
