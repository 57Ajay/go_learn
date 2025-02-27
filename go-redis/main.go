package main

import (
	"context"
	"fmt"
	// "strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func redisClient() *redis.Client {
	// Start the redis server, I ma using the docker to run redis server
	// steps to start a redis server using docker

	// 1. docker pull redis:latest
	// 2. docker run --name redis -p 6379:6379 -d redis
	// for redis-cli use this command: -> dockr exec -it redis(this should be the name of the redisClient) redis-cli

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func redisOps(rdb *redis.Client) {
	ctx := context.Background()

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

func redisStrings(rdb *redis.Client) {
	// SET key value [EX seconds] [PX milliseconds] [NX|XX] [KEEPTTL]
	// Options (Important to note):
	// EX seconds: Set expire time in seconds.
	// PX milliseconds: Set expire time in milliseconds.
	// NX: Set only if the key does not already exist (SET if Not eXists).
	// XX: Set only if the key does already exist (SET if eXists).
	// KEEPTTL: Retain the existing time to live (TTL).

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

func advRedisString(rdb *redis.Client) {
	/* Advanced redisStrings operations */
	// First GETSET (key doesn't exist)
	oldValue, err := rdb.GetSet(context.Background(), "myKey", "Initial Value").Result()
	if err != nil && err != redis.Nil { // Check for errors, but redis.Nil is expected initially
		panic(err)
	}
	if err == redis.Nil {
		fmt.Println("GETSET (first time), old value: nil") // Key didn't exist before
	} else {
		fmt.Println("GETSET (first time), old value:", oldValue) // Should not happen on first run
	}
	val1, _ := rdb.Get(context.Background(), "myKey").Result() // Verify the key is set
	fmt.Println("Current value of myKey:", val1)

	// Second GETSET (key now exists)
	oldValue2, err := rdb.GetSet(context.Background(), "myKey", "New Value").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("GETSET (second time), old value:", oldValue2) // Should be "Initial Value"
	val2, _ := rdb.Get(context.Background(), "myKey").Result()
	fmt.Println("Current value of myKey:", val2)

	// ----------------------- PSETEX key milliseconds value ------------
	err = rdb.Set(context.Background(), "shortLivedKey", "This expires in 5 seconds (milliseconds)", time.Duration(5000)*time.Millisecond).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("Set 'shortLivedKey' with 5000 millisecond expiry using rdb.Set with EX option")

	ttl, err := rdb.TTL(context.Background(), "shortLivedKey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Initial TTL for 'shortLivedKey': %v\n", ttl)

	time.Sleep(6 * time.Second)

	val, err := rdb.Get(context.Background(), "shortLivedKey").Result()
	if err == redis.Nil {
		fmt.Println("'shortLivedKey' has expired and is gone (GET returns nil)")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("Value of 'shortLivedKey' after waiting (should not happen in this example):", val)
	}
}

func redisList(rdb *redis.Client) {
	listName := "myList"

	length1, err := rdb.LPush(context.Background(), listName, "item1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("LPUSH '%s' with 'item1', new length: %d\n", listName, length1)

	length2, err := rdb.LPush(context.Background(), listName, "item2", "item3").Result() // Push multiple
	if err != nil {
		panic(err)
	}
	fmt.Printf("LPUSH '%s' with 'item2', 'item3', new length: %d\n", listName, length2)

	// Get the list range to verify order
	items, err := rdb.LRange(context.Background(), listName, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("List after LPUSH operations:", items)

	length3, err := rdb.RPush(context.Background(), listName, "itemA").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("RPUSH '%s' with 'itemA', new length: %d\n", listName, length3)

	length4, err := rdb.RPush(context.Background(), listName, "itemB", "itemC").Result() // Push multiple
	if err != nil {
		panic(err)
	}
	fmt.Printf("RPUSH '%s' with 'itemB', 'itemC', new length: %d\n", listName, length4)

	itemsAfterRPush, err := rdb.LRange(context.Background(), listName, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("List after RPUSH operations:", itemsAfterRPush) // Will show items from LPUSH and RPUSH in order

	// ---------- LPOP --------------

	poppedItem1, err := rdb.LPop(context.Background(), listName).Result()
	if err == redis.Nil {
		fmt.Println("LPOP: List is empty or key doesn't exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("LPOP from '%s', popped item: %s\n", listName, poppedItem1)
	}

	poppedItem2, err := rdb.LPop(context.Background(), listName).Result()
	if err == redis.Nil {
		fmt.Println("LPOP (second time): List is empty or key doesn't exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("LPOP from '%s' (second time), popped item: %s\n", listName, poppedItem2)
	}

	remainingItems, _ := rdb.LRange(context.Background(), listName, 0, -1).Result() // Get remaining list
	fmt.Println("List after LPOP operations:", remainingItems)

	// ---------- RPOP --------------

	poppedItemR1, err := rdb.RPop(context.Background(), listName).Result()
	if err == redis.Nil {
		fmt.Println("RPOP: List is empty or key doesn't exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("RPOP from '%s', popped item: %s\n", listName, poppedItemR1)
	}

	poppedItemR2, err := rdb.RPop(context.Background(), listName).Result()
	if err == redis.Nil {
		fmt.Println("RPOP (second time): List is empty or key doesn't exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("RPOP from '%s' (second time), popped item: %s\n", listName, poppedItemR2)
	}

	remainingItemsAfterRPop, _ := rdb.LRange(context.Background(), listName, 0, -1).Result()
	fmt.Println("List after RPOP operations:", remainingItemsAfterRPop)

	// ----------LLEN------------

	length, err := rdb.LLen(context.Background(), listName).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("LLEN of '%s': %d\n", listName, length)

	emptyListLength, err := rdb.LLen(context.Background(), "nonExistentList").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("LLEN of 'nonExistentList': %d\n", emptyListLength) // Should be 0

	// ----------LINDEX------------

	elementAtIndex0, err := rdb.LIndex(context.Background(), listName, 0).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if err == redis.Nil {
		fmt.Println("LINDEX 0: Index out of range")
	} else {
		fmt.Printf("LINDEX 0: %s\n", elementAtIndex0)
	}

	elementAtIndexNeg1, err := rdb.LIndex(context.Background(), listName, -1).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if err == redis.Nil {
		fmt.Println("LINDEX -1: Index out of range")
	} else {
		fmt.Printf("LINDEX -1: %s\n", elementAtIndexNeg1)
	}

	elementAtIndex10, err := rdb.LIndex(context.Background(), listName, 10).Result() // Out of range index
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if err == redis.Nil {
		fmt.Println("LINDEX 10: Index out of range (returned nil as redis.Nil)")
	} else {
		fmt.Printf("LINDEX 10: %s\n", elementAtIndex10) // Will not be reached in this case
	}

	// ----------LINSERT key BEFORE|AFTER pivot value--------------

	insertResult1, err := rdb.LInsertBefore(context.Background(), listName, "item2", "newItemBefore2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("LINSERT BEFORE 'item2', new length: %d\n", insertResult1)

	insertResult2, err := rdb.LInsertAfter(context.Background(), listName, "item3", "newItemAfter3").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("LINSERT AFTER 'item3', new length: %d\n", insertResult2)

	insertResultNotFound, err := rdb.LInsertBefore(context.Background(), listName, "nonExistentItem", "anotherNewItem").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("LINSERT BEFORE 'nonExistentItem' (pivot not found), result: %d (should be -1)\n", insertResultNotFound)

	currentList, _ := rdb.LRange(context.Background(), listName, 0, -1).Result()
	fmt.Println("List after LINSERT operations:", currentList)

	// ----------LREM key count value--------------
	rdb.LPush(context.Background(), listName, "apple", "banana", "apple", "orange", "apple", "grape", "apple")

	removeCount1, err := rdb.LRem(context.Background(), listName, 2, "apple").Result() // Remove up to 2 from head
	if err != nil {
		panic(err)
	}
	fmt.Printf("LREM count 2 'apple', removed: %d\n", removeCount1)

	removeCount2, err := rdb.LRem(context.Background(), listName, -1, "apple").Result() // Remove up to 1 from tail
	if err != nil {
		panic(err)
	}
	fmt.Printf("LREM count -1 'apple', removed: %d\n", removeCount2)

	removeAllCount, err := rdb.LRem(context.Background(), listName, 0, "apple").Result() // Remove all
	if err != nil {
		panic(err)
	}
	fmt.Printf("LREM count 0 'apple', removed: %d\n", removeAllCount)

	currentListAfterRem, _ := rdb.LRange(context.Background(), listName, 0, -1).Result()
	fmt.Println("List after LREM operations:", currentListAfterRem)

	// ------------LSET------------------
	err = rdb.LSet(context.Background(), listName, 1, "mango").Err()
	if err != nil {
		panic(err)
	} // Handle error if index is out of range
	fmt.Println("LSET index 1 to 'mango' - successful (no error)")

	// Example of error handling for out-of-range index (commented out to avoid panic in example)
	// err = rdb.LSet(context.Background(), listName, 5, "outOfRange").Err()
	// if err != nil {
	//     fmt.Println("LSET index 5 (out of range) error:", err) // Handle the error
	// }

	currentListAfterLSet, _ := rdb.LRange(context.Background(), listName, 0, -1).Result()
	fmt.Println("List after LSET operation:", currentListAfterLSet)

	// ------------LTRIM------------------

	err = rdb.LTrim(context.Background(), listName, 0, 4).Err() // Trim to first 5 elements
	if err != nil {
		panic(err)
	}
	fmt.Println("LTRIM 0 4 - successful (no error)")
	trimmedList, _ := rdb.LRange(context.Background(), listName, 0, -1).Result()
	fmt.Println("List after LTRIM 0 4:", trimmedList)

	err = rdb.LTrim(context.Background(), listName, 1, -3).Err() // Trim with more complex range
	if err != nil {
		panic(err)
	}
	fmt.Println("LTRIM 1 -3 - successful (no error)")
	trimmedList2, _ := rdb.LRange(context.Background(), listName, 0, -1).Result()
	fmt.Println("List after LTRIM 1 -3:", trimmedList2)

}

func advRedisList1(rdb *redis.Client) {

	listName := "task_queue"

	// Simulate a worker waiting for tasks (blocking LPop)
	go func() {
		fmt.Println("Worker: Waiting for tasks...")
		res, err := rdb.BLPop(context.Background(), 10*time.Second, listName).Result() // Block for max 10 seconds
		if err == redis.Nil {
			fmt.Println("Worker: Timeout reached, no task received")
		} else if err != nil {
			panic(err)
		} else {
			fmt.Printf("Worker: Received task from list '%s': %s\n", res[0], res[1]) // res is a []string: [key, value]
		}
	}()

	time.Sleep(2 * time.Second) // Wait for a bit

	// Simulate a producer adding a task
	taskValue := "Process order #123"
	pushResult, err := rdb.RPush(context.Background(), listName, taskValue).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Producer: Pushed task '%s' to '%s', list length: %d\n", taskValue, listName, pushResult)

	time.Sleep(5 * time.Second) // Wait a bit longer to allow worker to potentially process
	fmt.Println("advRedisList: Done.")

}

func advRedisList2(rdb *redis.Client) {
	sourceQueue := "task_queuee"
	processingQueue := "processing_queue"

	// Simulate a reliable worker using BRPOPLPUSH
	go func() {
		fmt.Println("Reliable Worker: Starting to process tasks...")
		for { // Keep processing tasks indefinitely
			res, err := rdb.BRPopLPush(context.Background(), sourceQueue, processingQueue, 10*time.Second).Result()
			if err == redis.Nil {
				fmt.Println("Reliable Worker: Timeout reached, no new tasks. Checking again...")
				continue // Check for tasks again after timeout
			} else if err != nil {
				// fmt.Println("This ran")
				panic(err)
			}

			task := res // Task value is the result (string)
			fmt.Printf("Reliable Worker: Started processing task: %s\n", task)

			time.Sleep(3 * time.Second) // Simulate task processing time

			// Simulate successful task completion - remove from processing queue
			removeCount, err := rdb.LRem(context.Background(), processingQueue, 1, task).Result()
			if err != nil {
				fmt.Printf("Error removing processed task from processing queue: %v\n", err)
				// In a real system, you'd need robust error handling and potentially retry/dead-letter queue logic
			} else {
				fmt.Printf("Reliable Worker: Successfully processed and removed task: %s, removed count: %d\n", task, removeCount)
			}
		}
	}()

	time.Sleep(2 * time.Second) // Wait a bit

	// Simulate a producer adding tasks to the source queue
	tasks := []string{"Task A", "Task B", "Task C"}
	for _, task := range tasks {
		pushResult, err := rdb.RPush(context.Background(), sourceQueue, task).Result()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Producer: Pushed task '%s' to '%s', list length: %d\n", task, sourceQueue, pushResult)
		time.Sleep(1 * time.Second) // Add tasks with a small delay
	}

	time.Sleep(15 * time.Second) // Let workers process for a while
	fmt.Println("advRedisList2: Producer finished adding tasks. Let workers continue...")
	time.Sleep(10 * time.Second) // Wait a bit longer before exiting
	fmt.Println("BLPOPRPUSH: Done.")
}

func redisSets(rdb *redis.Client) {

	setName := "mySet"

	// SADD key member [member ...]
	addedCount1, err := rdb.SAdd(context.Background(), setName, "apple").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SADD '%s' with 'apple', added count: %d\n", setName, addedCount1) // Should be 1 (new member)

	addedCount2, err := rdb.SAdd(context.Background(), setName, "banana", "orange", "grape").Result() // Add multiple
	if err != nil {
		panic(err)
	}
	fmt.Printf("SADD '%s' with 'banana', 'orange', added count: %d\n", setName, addedCount2) // Should be 2 (new members)

	addedCount3, err := rdb.SAdd(context.Background(), setName, "apple").Result() // Add duplicate
	if err != nil {
		panic(err)
	}
	fmt.Printf("SADD '%s' with duplicate 'apple', added count: %d\n", setName, addedCount3) // Should be 0 (duplicate ignored)

	members, err := rdb.SMembers(context.Background(), setName).Result() // Get all members
	if err != nil {
		panic(err)
	}
	fmt.Println("Set members after SADD operations:", members) // Order is not guaranteed

	// SREM key member [member ...]
	removedCount1, err := rdb.SRem(context.Background(), setName, "banana").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SREM '%s' remove 'banana', removed count: %d\n", setName, removedCount1) // Should be 1

	removedCount2, err := rdb.SRem(context.Background(), setName, "grape").Result() // Remove non-existent
	if err != nil {
		panic(err)
	}
	fmt.Printf("SREM '%s' remove 'grape' (non-existent), removed count: %d\n", setName, removedCount2) // Should be 0

	membersAfterRem, err := rdb.SMembers(context.Background(), setName).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Set members after SREM operations:", membersAfterRem)

	// SISMEMBER key member
	isMember1, err := rdb.SIsMember(context.Background(), setName, "apple").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SISMEMBER '%s' contains 'apple': %t\n", setName, isMember1) // Will be true or false

	isMember2, err := rdb.SIsMember(context.Background(), setName, "grape").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SISMEMBER '%s' contains 'grape': %t\n", setName, isMember2) // Will be false

	isMemberEmptySet, err := rdb.SIsMember(context.Background(), "nonExistentSet", "apple").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SISMEMBER 'nonExistentSet' contains 'apple': %t\n", isMemberEmptySet) // Will be false (treated as empty)

	// SCARD key
	cardinality, err := rdb.SCard(context.Background(), setName).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SCARD of '%s': %d\n", setName, cardinality)

	emptySetCardinality, err := rdb.SCard(context.Background(), "emptySet").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SCARD of 'emptySet': %d\n", emptySetCardinality)
}

func advRedisSetOps(rdb *redis.Client) {
	// SINTER key [key ...]
	rdb.SAdd(context.Background(), "set1", "a", "b", "c", "d")
	rdb.SAdd(context.Background(), "set2", "c", "d", "e", "f")
	rdb.SAdd(context.Background(), "set3", "c", "d", "g")

	intersectionMembers, err := rdb.SInter(context.Background(), "set1", "set2", "set3").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("SINTER set1, set2, set3:", intersectionMembers)

	intersectionEmptySet, err := rdb.SInter(context.Background(), "set1", "nonExistentSet").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("SINTER set1, nonExistentSet:", intersectionEmptySet)

	// SUNION key [key ...]
	rdb.SAdd(context.Background(), "setA", "apple", "banana", "orange")
	rdb.SAdd(context.Background(), "setB", "banana", "grape", "kiwi")

	unionMembers, err := rdb.SUnion(context.Background(), "setA", "setB").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("SUNION setA, setB:", unionMembers)

	unionWithEmptySet, err := rdb.SUnion(context.Background(), "setA", "nonExistentSet").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("SUNION setA, nonExistentSet:", unionWithEmptySet)

	// SDIFF key [key ...]

	rdb.SAdd(context.Background(), "setX", "p", "q", "r", "s", "t")
	rdb.SAdd(context.Background(), "setY", "r", "s", "u", "v")
	rdb.SAdd(context.Background(), "setZ", "s", "t", "w")

	diffMembers, err := rdb.SDiff(context.Background(), "setX", "setY", "setZ").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("SDIFF setX, setY, setZ:", diffMembers)

	diffWithEmptySet, err := rdb.SDiff(context.Background(), "setX", "nonExistentSet").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("SDIFF setX, nonExistentSet:", diffWithEmptySet) // Should be members of setX

	diffEmptySetWithSet, err := rdb.SDiff(context.Background(), "nonExistentSet", "setX").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("SDIFF nonExistentSet, setX:", diffEmptySetWithSet) // Should be empty list

	// --------------ADVANCED SET OPERATIONS----------------
	fmt.Println("\n----------------ADVANCED SET OPERATIONS----------------")
	// SINTERSTORE
	storeResult1, err := rdb.SInterStore(context.Background(), "intersection_result", "set1", "set2", "set3").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SINTERSTORE into 'intersection_result', result set size: %d\n", storeResult1)
	storedIntersection, _ := rdb.SMembers(context.Background(), "intersection_result").Result()
	fmt.Println("Stored intersection members:", storedIntersection)

	// SUNIONSTORE
	storeResult2, err := rdb.SUnionStore(context.Background(), "union_result", "setA", "setB").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SUNIONSTORE into 'union_result', result set size: %d\n", storeResult2)
	storedUnion, _ := rdb.SMembers(context.Background(), "union_result").Result()
	fmt.Println("Stored union members:", storedUnion)

	// SDIFFSTORE
	storeResult3, err := rdb.SDiffStore(context.Background(), "diff_result", "setX", "setY", "setZ").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SDIFFSTORE into 'diff_result', result set size: %d\n", storeResult3)
	storedDiff, _ := rdb.SMembers(context.Background(), "diff_result").Result()
	fmt.Println("Stored difference members:", storedDiff)

	// SINTERCARD
	cardinality1, err := rdb.SInterCard(context.Background(), 0, "set1", "set2", "set3").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SINTERCARD set1, set2, set3: cardinality = %d\n", cardinality1)

	cardinalityLimited, err := rdb.SInterCard(context.Background(), 0, "set1", "set2", "set3").Result() // Limit to 1
	if err != nil {
		panic(err)
	}
	fmt.Printf("SINTERCARD ... LIMIT 1: cardinality = %d (limited)\n", cardinalityLimited)

	cardinalityEmptySet, err := rdb.SInterCard(context.Background(), 0, "set1", "nonExistentSet").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SINTERCARD set1, nonExistentSet: cardinality = %d\n", cardinalityEmptySet) // Should be 0
}

func redisHashes(rdb *redis.Client) {

	hashKey := "user:1001"

	// HSET key field value [field value ...]
	addedCount1, err := rdb.HSet(context.Background(), hashKey, "name", "John Doe").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("HSET '%s' field 'name', added count: %d\n", hashKey, addedCount1) // Should be 1 (new field)

	addedCount2, err := rdb.HSet(context.Background(), hashKey, "email", "john.doe@example.com", "age", 30).Result() // Multiple fields
	if err != nil {
		panic(err)
	}
	fmt.Printf("HSET '%s' multiple fields, added count: %d\n", hashKey, addedCount2) // Should be 2 (new fields)

	addedCount3, err := rdb.HSet(context.Background(), hashKey, "name", "Johnny Doe").Result() // Overwrite existing field
	if err != nil {
		panic(err)
	}
	fmt.Printf("HSET '%s' overwrite 'name', added count: %d\n", hashKey, addedCount3) // Should be 0 (field overwritten, not added)

	hashData, err := rdb.HGetAll(context.Background(), hashKey).Result() // Get all fields and values
	if err != nil {
		panic(err)
	}
	fmt.Println("Hash data after HSET operations:", hashData) // Will be a map[string]string

	// HGET key field

	nameValue, err := rdb.HGet(context.Background(), hashKey, "name").Result()
	if err == redis.Nil {
		fmt.Println("HGET 'name': field or hash key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("HGET 'name': %s\n", nameValue)
	}

	ageValue, err := rdb.HGet(context.Background(), hashKey, "age").Result()
	if err == redis.Nil {
		fmt.Println("HGET 'age': field or hash key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("HGET 'age': %s\n", ageValue)
	}

	cityValue, err := rdb.HGet(context.Background(), hashKey, "city").Result() // Non-existent field
	if err == redis.Nil {
		fmt.Println("HGET 'city' (non-existent): field or hash key does not exist (returned nil)")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("HGET 'city': %s\n", cityValue) // Will not be reached in this case
	}

	// HMGET key field [field ...]

	fieldValues, err := rdb.HMGet(context.Background(), hashKey, "name", "age", "city", "email").Result() // Request multiple fields
	if err != nil {
		panic(err)
	}
	fmt.Println("HMGET 'name', 'age', 'city', 'email':", fieldValues) // fieldValues is a []interface{}

	// Process the returned values - need to type assert to string if expect string values
	for i, val := range fieldValues {
		fieldName := []string{"name", "age", "city", "email"}[i] // Corresponding field name for index i
		if val == nil {
			fmt.Printf("HMGET: Field '%s' not found or hash key doesn't exist (value is nil)\n", fieldName)
		} else {
			stringValue, ok := val.(string)
			if ok {
				fmt.Printf("HMGET: Field '%s': %s\n", fieldName, stringValue)
			} else {
				fmt.Printf("HMGET: Field '%s': Value is not a string type (unexpected)\n", fieldName)
			}
		}
	}

	// HGETALL key

	allHashData, err := rdb.HGetAll(context.Background(), hashKey).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("HGETALL of '%s':'%v'", hashKey, allHashData) // allHashData is map[string]string

	emptyHashData, err := rdb.HGetAll(context.Background(), "emptyHash").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("HGETALL of 'emptyHash':", emptyHashData)

	// HKEYS key

	fieldNames, err := rdb.HKeys(context.Background(), hashKey).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("HKEYS of '%s':'%v'", hashKey, fieldNames) // fieldNames is []string

	fmt.Println()
	emptyHashKeys, err := rdb.HKeys(context.Background(), "emptyHash").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("HKEYS of 'emptyHash':", emptyHashKeys)

	// HVALS key

	fieldValues_, err := rdb.HVals(context.Background(), hashKey).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("HVALS of '%s':'%v'", hashKey, fieldValues_) // fieldValues is []string

	fmt.Println()
	emptyHashValues, err := rdb.HVals(context.Background(), "emptyHash").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("HVALS of 'emptyHash':", emptyHashValues)

	// HDEL key field [field ...]

	deletedCount1, err := rdb.HDel(context.Background(), hashKey, "name").Result()
	if err != nil {
		fmt.Println("Error deleting field 'name':", err)
	}
	fmt.Printf("HDEL '%s' field 'name', deleted count: %d\n", hashKey, deletedCount1)

	hashDataAfterDel, err := rdb.HGetAll(context.Background(), hashKey).Result()
	if err != nil {
		fmt.Println("Error getting hash data after HDEL:", err)
	}
	fmt.Println("Hash data after HDEL 'name':", hashDataAfterDel)

	// HEXISTS key feild

	exists, err := rdb.HExists(context.Background(), hashKey, "name").Result()
	if err != nil {
		fmt.Println("Error checking if field 'name' exists:", err)
	}
	fmt.Printf("HEXISTS '%s' field 'name': %t\n", hashKey, exists)

	notExists, err := rdb.HExists(context.Background(), hashKey, "city").Result()
	if err != nil {
		fmt.Println("Error checking if field 'city' exists:", err)
	}
	fmt.Printf("HEXISTS '%s' field 'city': %t\n", hashKey, notExists)

	// HLEN key

	hashLength, err := rdb.HLen(context.Background(), hashKey).Result()
	if err != nil {
		fmt.Println("Error getting hash length:", err)
	}
	fmt.Printf("HLEN of '%s': %d\n", hashKey, hashLength)

	// HINCRBY key field increment

	hashKey = "product:stats"

	rdb.HSet(context.Background(), hashKey, "views", "100") // Initialize 'views'

	newValue1, err := rdb.HIncrBy(context.Background(), hashKey, "views", 1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("HINCRBY '%s' field 'views' by 1, new value: %d\n", hashKey, newValue1)

	newValue2, err := rdb.HIncrBy(context.Background(), "user:profile", "login_count", 5).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("HINCRBY 'user:profile' field 'login_count' by 5, new value: %d\n", newValue2)

	loginCountValue, _ := rdb.HGet(context.Background(), "user:profile", "login_count").Result()
	fmt.Println("Verified 'login_count' value:", loginCountValue)

	// Example of error handling (invalid integer value in hash field) - not executed to avoid panic in example
	// rdb.HSet(context.Background(), "item:data", "price", "invalid")
	// _, err = rdb.HIncrBy(context.Background(), "item:data", "price", 1).Result()
	//
	//	if err != nil {
	//	    fmt.Println("HINCRBY error (invalid integer):", err) // Handle the error
	//	}
	//
	//
	// ----------> HINCRBYFLOAT key field increment

	hashKey = "product:price"

	rdb.HSet(context.Background(), hashKey, "current_price", "99.99")

	newValueStr1, err := rdb.HIncrByFloat(context.Background(), hashKey, "current_price", 10.5).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("HINCRBYFLOAT '%s' field 'current_price' by 10.5, new value (string): %f\n", hashKey, newValueStr1)

	fmt.Printf("HINCRBYFLOAT 'current_price', new value (float64): %f\n", newValueStr1)

	newValueStr2, err := rdb.HIncrByFloat(context.Background(), "sensor:readings", "temperature", 0.25).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("HINCRBYFLOAT 'sensor:readings' field 'temperature' by 0.25, new value (string): %s:->'%f'\n", hashKey, newValueStr2)

	// Example of error handling (invalid float value) - not executed to avoid panic in example
	// rdb.HSet(context.Background(), "item:data", "weight", "not_a_number")
	// _, err = rdb.HIncrByFloat(context.Background(), "item:data", "weight", 0.1).Result()
	// if err != nil {
	//     fmt.Println("HINCRBYFLOAT error (invalid float):", err) // Handle the error
	// }

}

func redisSortedSets(rdb *redis.Client) {

	// ZADD key [NX|XX] [CH] [INCR] score member [score member ...]
	sortedSetName := "leaderboard"

	addedCount1, err := rdb.ZAdd(context.Background(), sortedSetName, redis.Z{Score: 100, Member: "PlayerA"}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZADD '%s' with PlayerA, added count: %d\n", sortedSetName, addedCount1) // Should be 1

	addedCount2, err := rdb.ZAdd(context.Background(), sortedSetName,
		redis.Z{Score: 150, Member: "PlayerB"},
		redis.Z{Score: 120, Member: "PlayerC"},
	).Result() // Add multiple members
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZADD '%s' with PlayerB, PlayerC, added count: %d\n", sortedSetName, addedCount2) // Should be 2

	addedCount3, err := rdb.ZAdd(context.Background(), sortedSetName, redis.Z{Score: 130, Member: "PlayerB"}).Result() // Update score
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZADD '%s' update PlayerB score, added count: %d\n", sortedSetName, addedCount3) // Should be 0 (score updated, not a new member)

	rangeWithScores, err := rdb.ZRangeWithScores(context.Background(), sortedSetName, 0, -1).Result() // Get range with scores
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sorted Set after ZADD operations (ZRANGE WITHSCORES): %v", rangeWithScores)

	// ZRANGE key start stop [WITHSCORES]

	top3Players, err := rdb.ZRange(context.Background(), sortedSetName, 0, 2).Result() // Get top 3 (lowest scores)
	if err != nil {
		panic(err)
	}
	fmt.Println("ZRANGE 0 2 (top 3 players by score):", top3Players)

	allPlayersWithScores, err := rdb.ZRangeWithScores(context.Background(), sortedSetName, 0, -1).Result() // Get all with scores
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZRANGE 0 -1 WITHSCORES (all players with scores): %v", allPlayersWithScores)

	outOfRangePlayers, err := rdb.ZRange(context.Background(), sortedSetName, 10, 20).Result() // Out of range
	if err != nil {
		panic(err)
	}
	fmt.Println("ZRANGE 10 20 (out of range):", outOfRangePlayers)

	// ZREVRANGE key start stop [WITHSCORES]

	top3PlayersDesc, err := rdb.ZRevRange(context.Background(), sortedSetName, 0, 2).Result() // Top 3 (highest scores)
	if err != nil {
		panic(err)
	}
	fmt.Println("ZREVRANGE 0 2 (top 3 players by score - descending):", top3PlayersDesc)

	allPlayersWithScoresDesc, err := rdb.ZRevRangeWithScores(context.Background(), sortedSetName, 0, -1).Result() // All with scores descending
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZREVRANGE 0 -1 WITHSCORES (all players with scores - descending): %v", allPlayersWithScoresDesc)

	// ZSCORE key member

	playerAScore, err := rdb.ZScore(context.Background(), sortedSetName, "PlayerA").Result()
	if err == redis.Nil {
		fmt.Println("ZSCORE PlayerA: Member or sorted set does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("ZSCORE PlayerA: %f\n", playerAScore) // Score as string
	}

	nonExistentPlayerScore, err := rdb.ZScore(context.Background(), sortedSetName, "NonExistentPlayer").Result() // Member not found
	if err == redis.Nil {
		fmt.Println("ZSCORE NonExistentPlayer: Member or sorted set does not exist (returned nil)")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("ZSCORE NonExistentPlayer: %f\n", nonExistentPlayerScore) // Will not be reached in this case
	}

	nonExistentSetScore, err := rdb.ZScore(context.Background(), "nonExistentSortedSet", "PlayerA").Result() // Sorted set doesn't exist
	if err == redis.Nil {
		fmt.Println("ZSCORE from 'nonExistentSortedSet': Member or sorted set does not exist (returned nil)")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Printf("ZSCORE from 'nonExistentSortedSet': %f\n", nonExistentSetScore) // Will not be reached in this case
	}

	// ZREM key member [member ...]

	removedCount1, err := rdb.ZRem(context.Background(), sortedSetName, "PlayerC", "PlayerD").Result() // Remove existing and non-existent
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZREM '%s' members 'PlayerC', 'PlayerD', removed count: %d\n", sortedSetName, removedCount1) // Should be 1

	rangeAfterRem, err := rdb.ZRangeWithScores(context.Background(), sortedSetName, 0, -1).Result() // Get range after removal
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sorted Set after ZREM: %v", rangeAfterRem)

	// ZCARD key

	card, err := rdb.ZCard(context.Background(), sortedSetName).Result()
	if err != nil {
		fmt.Println("Error getting sorted set cardinality:", err)
	}
	fmt.Printf("ZCARD of '%s': %d\n", sortedSetName, card)

	// ZCOUNT key min max

	sortedSetName = "scores"

	rdb.ZAdd(context.Background(), sortedSetName,
		redis.Z{Score: 80, Member: "UserA"},
		redis.Z{Score: 90, Member: "UserB"},
		redis.Z{Score: 100, Member: "UserC"},
		redis.Z{Score: 110, Member: "UserD"},
		redis.Z{Score: 120, Member: "UserE"},
	)

	count1, err := rdb.ZCount(context.Background(), sortedSetName, "90", "110").Result() // Inclusive range
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZCOUNT '%s' range [90, 110] (inclusive): %d\n", sortedSetName, count1) // Should be 3

	count2, err := rdb.ZCount(context.Background(), sortedSetName, "90", "(110").Result() // Exclusive max
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZCOUNT '%s' range [90, (110) (exclusive max): %d\n", sortedSetName, count2) // Should be 2

	count3, err := rdb.ZCount(context.Background(), sortedSetName, "-inf", "100").Result() // -inf to 100
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZCOUNT '%s' range [-inf, 100]: %d\n", sortedSetName, count3) // Should be 3

	count4, err := rdb.ZCount(context.Background(), sortedSetName, "150", "+inf").Result() // 150 to +inf
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZCOUNT '%s' range [150, +inf]: %d\n", sortedSetName, count4) // Should be 0

	emptySetCount, err := rdb.ZCount(context.Background(), "emptySortedSet", "0", "100").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZCOUNT 'emptySortedSet' range [0, 100]: %d\n", emptySetCount) // Should be 0

	// ZRANK key member

	rank1, err := rdb.ZRank(context.Background(), sortedSetName, "UserC").Result()
	if err != nil {
		fmt.Println("Error getting rank for 'UserC':", err)
	}
	fmt.Printf("ZRANK '%s' member 'UserC': %d\n", sortedSetName, rank1)

	// ZREVANK key member

	rank2, err := rdb.ZRevRank(context.Background(), sortedSetName, "UserC").Result()
	if err != nil {
		fmt.Println("Error getting rank for 'UserC':", err)
	}
	fmt.Printf("ZREVRANK '%s' member 'UserC': %d\n", sortedSetName, rank2)

}

func advRedisSortedSets(rdb *redis.Client) {
	// ZREMRANGEBYRANK key start stop
	sortedSetName := "leaderboard"

	// Re-populate leaderboard for example
	rdb.ZAdd(context.Background(), sortedSetName,
		redis.Z{Score: 100, Member: "PlayerA"},
		redis.Z{Score: 120, Member: "PlayerC"},
		redis.Z{Score: 130, Member: "PlayerB"},
	)

	removedCount1, err := rdb.ZRemRangeByRank(context.Background(), sortedSetName, 0, 1).Result() // Remove ranks 0 and 1
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZREMRANGEBYRANK '%s' ranks [0, 1], removed count: %d\n", sortedSetName, removedCount1) // Should be 2

	rangeAfterRankRem, err := rdb.ZRangeWithScores(context.Background(), sortedSetName, 0, -1).Result() // Get range after rank removal
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sorted Set after ZREMRANGEBYRANK (ranks 0-1): %v", rangeAfterRankRem)

	removedCountAll, err := rdb.ZRemRangeByRank(context.Background(), sortedSetName, 0, -1).Result() // Remove all
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZREMRANGEBYRANK '%s' ranks [0, -1] (all), removed count: %d\n", sortedSetName, removedCountAll) // Should be 1 (PlayerB remaining)

	emptySetRemCount, err := rdb.ZRemRangeByRank(context.Background(), "emptySortedSet", 0, 5).Result() // Remove from empty set
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZREMRANGEBYRANK 'emptySortedSet' ranks [0, 5], removed count: %d\n", emptySetRemCount) // Should be 0

	// ZREMRANGEBYSCORE key min max

	rdb.ZAdd(context.Background(), sortedSetName,
		redis.Z{Score: 50, Member: "ProductA"},
		redis.Z{Score: 75, Member: "ProductB"},
		redis.Z{Score: 100, Member: "ProductC"},
		redis.Z{Score: 125, Member: "ProductD"},
	)

	removedCount1, err = rdb.ZRemRangeByScore(context.Background(), sortedSetName, "75", "100").Result() // Remove scores [75, 100] inclusive
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZREMRANGEBYSCORE '%s' scores [75, 100], removed count: %d\n", sortedSetName, removedCount1) // Should be 2

	rangeAfterScoreRem, err := rdb.ZRangeWithScores(context.Background(), sortedSetName, 0, -1).Result() // Get range after score removal
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sorted Set after ZREMRANGEBYSCORE (scores 75-100): %v", rangeAfterScoreRem)

	removedCountInf, err := rdb.ZRemRangeByScore(context.Background(), sortedSetName, "-inf", "60").Result() // Remove scores up to 60
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZREMRANGEBYSCORE '%s' scores [-inf, 60], removed count: %d\n", sortedSetName, removedCountInf) // Should be 1

	removedCountExclusive, err := rdb.ZRemRangeByScore(context.Background(), sortedSetName, "(120", "+inf").Result() // Remove scores > 120 (exclusive)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZREMRANGEBYSCORE '%s' scores (120, +inf], removed count: %d\n", sortedSetName, removedCountExclusive) // Should be 1

	emptySetRemScoreCount, err := rdb.ZRemRangeByScore(context.Background(), "emptySortedSet", "0", "100").Result() // Remove from empty set
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZREMRANGEBYSCORE 'emptySortedSet' scores [0, 100], removed count: %d\n", emptySetRemScoreCount) // Should be 0

	//  ZRANGEBYLEX key min max [LIMIT offset count] and ZREVRANGEBYLEX key max min [LIMIT offset count]
	//  i have some doubts to clear regarding that so some code might not be correct.

	sortedSetName = "dictionary"

	// Populate dictionary set (all with score 0)
	rdb.ZAddArgs(context.Background(), sortedSetName, redis.ZAddArgs{
		Members: []redis.Z{
			{Score: 0, Member: "apple"},
			{Score: 0, Member: "banana"},
			{Score: 0, Member: "apricot"},
			{Score: 0, Member: "cherry"},
			{Score: 0, Member: "date"},
			{Score: 0, Member: "elderberry"},
		},
	}).Result()

	// Get all items lexicographically
	lexRangeAll, err := rdb.ZRangeByLex(context.Background(), sortedSetName, &redis.ZRangeBy{Min: "-", Max: "+"}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("ZRANGEBYLEX - + (all lexicographically):", lexRangeAll)

	// Get range "[ap" to "(ba"
	lexRangePartial, err := rdb.ZRangeByLex(context.Background(), sortedSetName, &redis.ZRangeBy{Min: "[ap", Max: "(ba"}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("ZRANGEBYLEX [ap (ba :", lexRangePartial)

	// Get range "[banana" to "+" with LIMIT 0 2
	lexRangeLimited, err := rdb.ZRangeByLex(context.Background(), sortedSetName, &redis.ZRangeBy{
		Min:    "[banana",
		Max:    "+",
		Offset: 0,
		Count:  2,
	}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("ZRANGEBYLEX [banana + LIMIT 0 2:", lexRangeLimited)

	// Reverse lexicographical order
	revLexRangeAll, err := rdb.ZRevRangeByLex(context.Background(), sortedSetName, &redis.ZRangeBy{Min: "+", Max: "-"}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("ZREVRANGEBYLEX + - (all reverse lexicographically):", revLexRangeAll)

	// ZUNIONSTORE destination numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM|MIN|MAX]

	rdb.ZAdd(context.Background(), "setZ1", redis.Z{Score: 1, Member: "member1"}, redis.Z{Score: 2, Member: "member2"}, redis.Z{Score: 3, Member: "member3"})
	rdb.ZAdd(context.Background(), "setZ2", redis.Z{Score: 2, Member: "member2"}, redis.Z{Score: 3, Member: "member3"}, redis.Z{Score: 4, Member: "member4"})

	storeResult1, err := rdb.ZUnionStore(context.Background(), "union_result", &redis.ZStore{
		Keys: []string{"setZ1", "setZ2"},
	}).Result() // Default AGGREGATE SUM
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZUNIONSTORE into 'union_result', result set size: %d\n", storeResult1)
	unionResult, _ := rdb.ZRangeWithScores(context.Background(), "union_result", 0, -1).Result()
	fmt.Printf("Stored union members (SUM aggregation): %v", unionResult)

	storeResult2, err := rdb.ZUnionStore(context.Background(), "weighted_union_min", &redis.ZStore{
		Keys:      []string{"setZ1", "setZ2"},
		Weights:   []float64{0.5, 2},
		Aggregate: "MIN",
	}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nZUNIONSTORE into 'weighted_union_min', result set size: %d\n", storeResult2)
	weightedUnionMinResult, _ := rdb.ZRangeWithScores(context.Background(), "weighted_union_min", 0, -1).Result()
	fmt.Printf("Stored union members (WEIGHTS 0.5, 2, AGGREGATE MIN): %v", weightedUnionMinResult)

	// ZINTERSTORE destination numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM|MIN|MAX]

	rdb.ZAdd(context.Background(), "setInt1", redis.Z{Score: 1, Member: "itemA"}, redis.Z{Score: 2, Member: "itemB"}, redis.Z{Score: 3, Member: "itemC"})
	rdb.ZAdd(context.Background(), "setInt2", redis.Z{Score: 2, Member: "itemB"}, redis.Z{Score: 3, Member: "itemC"}, redis.Z{Score: 4, Member: "itemD"})

	storeResult3, err := rdb.ZInterStore(context.Background(), "intersect_result", &redis.ZStore{
		Keys: []string{"setInt1", "setInt2"},
	}).Result() // Default AGGREGATE SUM
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZINTERSTORE into 'intersect_result', result set size: %d\n", storeResult3)
	intersectResult, _ := rdb.ZRangeWithScores(context.Background(), "intersect_result", 0, -1).Result()
	fmt.Printf("Stored intersection members (SUM aggregation): %v", intersectResult)

	storeResult4, err := rdb.ZInterStore(context.Background(), "weighted_intersect_max", &redis.ZStore{
		Keys:      []string{"setInt1", "setInt2"},
		Weights:   []float64{1, 0.5},
		Aggregate: "MAX",
	}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ZINTERSTORE into 'weighted_intersect_max', result set size: %d\n", storeResult4)
	weightedIntersectMaxResult, _ := rdb.ZRangeWithScores(context.Background(), "weighted_intersect_max", 0, -1).Result()
	fmt.Println("Stored intersection members (WEIGHTS 1, 0.5, AGGREGATE MAX):", weightedIntersectMaxResult)

}

func main() {
	rdb := redisClient()
	// fmt.Println("-------Common Redis Operations-------")
	// redisOps(rdb)
	// fmt.Println("\n-------Redis String Operations-------")
	// redisStrings(rdb)
	// fmt.Println("\n-------Advanced Redis String Operations-------")
	// advRedisString(rdb)
	// fmt.Println("\n-------Redis List Operations-------")
	// redisList(rdb)
	// fmt.Println("\n-------Advanced Redis List Operations-------")
	// advRedisList1(rdb)
	// fmt.Println("")
	// advRedisList2(rdb)
	// fmt.Println("\n-------Redis Sets Operations-------")
	// redisSets(rdb)
	// fmt.Println("\n-------Advanced Redis Sets Operations-------")
	// advRedisSetOps(rdb)
	// fmt.Println("\n-------Redis Hashes Operations-------")
	// redisHashes(rdb)
	// fmt.Println("\n-------Redis Sorted Sets Operations-------")
	// redisSortedSets(rdb)
	fmt.Println("\n-------Advanced Redis Sorted Sets Operations-------")
	advRedisSortedSets(rdb)
}
