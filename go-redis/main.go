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

func main() {
	rdb := redisClient()
	fmt.Println("-------Common Redis Operations-------")
	redisOps(rdb)
	fmt.Println("\n-------Redis String Operations-------")
	redisStrings(rdb)
	fmt.Println("\n-------Advanced Redis String Operations-------")
	advRedisString(rdb)
	fmt.Println("\n-------Redis List Operations-------")
	redisList(rdb)
}
