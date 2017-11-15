package redis

import (
	"time"

	"gopkg.in/redis.v5"
)

var (
	redisConn *redis.Client
)

func init() {
	redisConn = getRedis()
}
func GetDigInfo(key string) (val string, err error) {

	val, err = redisConn.Get(key).Result()
	if err != nil {
		return
	}
	return
	//	log.Println("the return is :", val)

	//	val1, err := redisConn.Get("key5").Result()
	//	if err == redis.Nil {
	//		log.Println("this key is not exists")
	//	} else {
	//		log.Println("this value of key is :", val1)
	//	}

}
func SetDigInfo(key string, value string) (err error) {
	err = redisConn.Set(key, value, 1*time.Minute).Err()
	if err != nil {
		return
	}
	return
}
func factory() *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
		PoolSize: 50,
	})
}
func getRedis() *redis.Client {
	return factory()
}
