package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func redisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func redisOps() {
	ctx := context.Background()
	rdb := redisClient()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Redis Ping: ", pong)

	// --- SET command ---
	err = rdb.Set(ctx, "Name", "Ajay", 0).Err()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// --- GET command ---
	val, err := rdb.Get(ctx, "Name").Result()
	if err == redis.Nil {
		fmt.Println("Key does not exist")
		return
	}
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Name: ", val)

	// --- DEL command ---
	effected, err := rdb.Del(ctx, "Name").Result()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Deleted: ", effected) // 1 if Key exists,
	// 0 if Key does not exist

	// --- EXISTS command ---
	exists, err := rdb.Exists(ctx, "Name").Result()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Exists: ", exists) // Similer response to Del

}

func redisStrings() {
	// SET key value [EX seconds] [PX milliseconds] [NX|XX] [KEEPTTL]
	// Options (Important to note):
	// EX seconds: Set expire time in seconds.
	// PX milliseconds: Set expire time in milliseconds.
	// NX: Set only if the key does not already exist (SET if Not eXists).
	// XX: Set only if the key does already exist (SET if eXists).
	// KEEPTTL: Retain the existing time to live (TTL).

	rdb := redisClient()

	// -----SET KEY-----

	err := rdb.Set(context.Background(), "myString", "Hello String", 0).Err() // No expiry
	if err != nil {
		panic(err)
	}

	err = rdb.Set(context.Background(), "anotherString", "Value with expiry", time.Duration(60)*time.Second).Err() // Expire in 60 seconds
	if err != nil {
		panic(err)
	}

	set, err := rdb.SetNX(context.Background(), "newString", "Only set if new", 0).Result() // SET NX
	if err != nil {
		panic(err)
	}
	fmt.Println("SETNX success:", set) // Will be true if set, false if key existed

	// -----MSET and MGET-----

	err = rdb.MSet(context.Background(), "key1", "value1", "key2", "value2").Err()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("MSET success")

	mget, err := rdb.MGet(context.Background(), "jeff", "key1", "key2").Result()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("MGET success: ", mget)
	for i, v := range mget {
		if v == nil {
			fmt.Println("Key not found: ", i)
			continue
		}
		fmt.Println("Value: ", v)
	}

	// -----Append-----
	length, err := rdb.Append(context.Background(), "myString", " Redis!").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("APPEND new length:", length)
	val, _ := rdb.Get(context.Background(), "myString").Result()
	fmt.Println("Updated myString:", val)

	// -----StrLen-----
	length, err = rdb.StrLen(context.Background(), "myString").Result()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("StrLen: ", length)

	// -----GETRANGE key start end and SETRANGE key offset value-----
	substring, err := rdb.GetRange(context.Background(), "myString", 0, 3).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("GETRANGE 0 3:", substring)
	substring2, err := rdb.GetRange(context.Background(), "myString", -6, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("GETRANGE -6 -1:", substring2)

	_, err = rdb.SetRange(context.Background(), "myString", 6, "Redis").Result()
	if err != nil {
		panic(err)
	}
	val, _ = rdb.Get(context.Background(), "myString").Result()
	fmt.Println("SETRANGE updated myString:", val)

	_, err = rdb.SetRange(context.Background(), "emptyKey", 10, "Hello").Result()
	if err != nil {
		panic(err)
	}
	val, _ = rdb.Get(context.Background(), "emptyKey").Result()
	fmt.Println("SETRANGE on emptyKey:", val)

	// -----INCR, INCRBY, DECR, DECRBY-----

	rdb.Set(context.Background(), "counter", "10", 0) // Initialize counter

	newValue, err := rdb.Incr(context.Background(), "counter").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("INCR counter, new value:", newValue)

	newValue2, err := rdb.Incr(context.Background(), "nonExistentCounter").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("INCR nonExistentCounter, new value:", newValue2)

	// Example of error handling (if value is not an integer) - not executed here to avoid panic in example
	// rdb.Set(context.Background(), "invalidCounter", "hello", 0)
	// _, err = rdb.Incr(context.Background(), "invalidCounter").Result()
	// if err != nil {
	//     fmt.Println("INCR invalidCounter error:", err) // Handle the error
	// }

	newValue, err = rdb.Decr(context.Background(), "counter").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("DECR counter, new value:", newValue)

	newValue2, err = rdb.Decr(context.Background(), "anotherCounter").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("DECR anotherCounter, new value:", newValue2)

	newValue, err = rdb.IncrBy(context.Background(), "score", 50).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("INCRBY score by 50, new value:", newValue)

	newValue2, err = rdb.IncrBy(context.Background(), "score", -20).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("INCRBY score by -20, new value:", newValue2)

	rdb.Set(context.Background(), "stock", "100", 0)

	newValue, err = rdb.DecrBy(context.Background(), "stock", 10).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("DECRBY stock by 10, new value:", newValue)

	rdb.Set(context.Background(), "price", "10.50", 0) // Initialize price

	newValueStr, err := rdb.IncrByFloat(context.Background(), "price", 2.3).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("INCRBYFLOAT price by 2.3, new value (string):", newValueStr)

	fmt.Println("INCRBYFLOAT price, new value (float64):", newValueStr)

	newValueStr2, err := rdb.IncrByFloat(context.Background(), "floatCounter", 1.5).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("INCRBYFLOAT floatCounter by 1.5, new value (string):", newValueStr2)

}

func main() {
	redisOps()
	redisStrings()
}
