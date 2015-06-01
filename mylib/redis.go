// This module configures a redis connect for the application
package mylib

import (
	"gopkg.in/redis.v3"
	"github.com/revel/revel"
)

var (
	Redis *redis.Client
)

func Init() {
	// Read configuration.
	var found bool
	var redisAddr string
	var password string

	if redisAddr, found = revel.Config.String("redis.addr"); 
		!found {
			revel.ERROR.Fatal("No redis.addr found.")
		}

    password, _ = revel.Config.String("redis.password");

    Redis = redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: password, 
        DB:       0,  
    })

	err := Redis.Ping().Err()
	if err != nil {
		revel.ERROR.Fatal(err)
	}
}		

type RedisController struct {
	*revel.Controller
	Redis *redis.Client
}

func (c *RedisController) Begin() revel.Result {
	c.Redis = Redis
	return nil
}


func init() {
	revel.OnAppStart(Init)
	revel.InterceptMethod((*RedisController).Begin, revel.BEFORE)
}