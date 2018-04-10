package redis

import (
	"github.com/go-redis/redis"

	//"encoding/json"
	"fmt"
	"sync"
	"time"
)

var Session *redis.Client

func Connect() {
	Session = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//	value := map[string]interface{}{"xtv": user{Info: "info", Age: 4}}
	//	data, _ := json.Marshal(value)

	//	err := redis.Session.Set("user", data, 0).Err()
	//	if err != nil {
	//		fmt.Println("error 1 ")
	//	}

	//	val, _ := redis.Session.Get("user").Result()
	//	fmt.Println(val)

	//	var value map[string]user
	//	json.Unmarshal([]byte(val), &value)

	//	value["big"] = user{Info: "boo", Age: 3}
	//	data, _ := json.Marshal(value)
	//	redis.Session.Set("user", data, 0)

	//	fmt.Println(value)

	//	mm := map[string]int{"x": 100}
	//	c := class{M: mm}
	//	data, _ := json.Marshal(c)
	//	fmt.Println(string(data))

	//	fields := make(map[string]string)
	//	fields["b"] = string(data)
	//	redis.Session.HMSetMap("HM", fields)

	//	data, exist := redis.Session.HGet("HM", "a").Bytes()
	//	if exist == nil {
	//		c := class{}
	//		json.Unmarshal(data, &c)

	//		fmt.Println(c)
	//	}
}

func Run() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if time.Now().Minute() == 30 {
				//go Update_DB()
			}

			fmt.Println(time.Now().Second())
			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}
