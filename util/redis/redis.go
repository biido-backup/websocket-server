package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var client *redis.Client


func ConnectRedis() {

	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	password := viper.GetString("redis.password")

	client = redis.NewClient(&redis.Options{
		Addr:host+":"+port,
		Password:password,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

}

func GetValueByKey(key string) []byte{

	value, err := client.Get(key).Bytes()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		panic(err)
	} else {
		return value
	}
}

func SetValueByKey(key string, value []byte){
	err := client.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
