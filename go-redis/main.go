package main

import (
	"context"
	"fmt"
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

func main() {
	redisOps()
}
