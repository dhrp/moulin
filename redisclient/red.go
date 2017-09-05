package main

import (
	"log"
	"strconv"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// RedClient Basic client class to group Redis functions
type RedClient struct {
	clientpool *pool.Pool
	host       string
	client     *redis.Client
}

func (c *RedClient) init() *RedClient {

	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		// set password with CONFIG SET requirepass "nevermind"
		if err = client.Cmd("AUTH", "nevermind").Err; err != nil {
			client.Close()
			return nil, err
		}
		return client, nil
	}

	client, err := pool.NewCustom("tcp", c.host, 10, df)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("redis client connected successfully with radix driver")
	c.clientpool = client

	return c
}

func (c *RedClient) popQueueAndSaveKeyToSet(queueID string, destinationSet string, expirationSec int) (string, error) {

	expiresAt := int64(time.Now().Unix()) + int64(expirationSec)
	score := strconv.FormatInt(expiresAt, 10)

	luaScript := `
		local value = redis.call('RPOP', KEYS[1]);
		local jsonObj = cjson.decode(value);
		local taskID = string.format("%.0f", jsonObj['id']);
		local destinationKey = KEYS[1] .. "." .. taskID;
		redis.call('SET', destinationKey, value);
		local destinationSet = ARGV[1];
		local count = redis.call('ZADD', destinationSet, ARGV[2], taskID);
		return {taskID, value};
		`

	result, err := c.clientpool.Cmd("EVAL", luaScript, 1, queueID, destinationSet, score).Array()
	if err != nil {
		log.Panic(err)
	}
	taskID, _ := result[0].Int()
	msg, _ := result[1].Str()

	log.Printf("COMPLETED %d, %s", taskID, msg)
	return msg, nil
}

func (c *RedClient) del(key string) int {
	count, err := c.clientpool.Cmd("DEL", key).Int()
	if err != nil {
		log.Panic(err)
	}
	return count
}

func (c *RedClient) lpush(key string, value string) (int, error) {

	log.Printf("Doing: LPUSH %s %s", key, value)
	newLength, err := c.clientpool.Cmd("LPUSH", key, value).Int()
	if err != nil {
		log.Panic(err)
	}
	return newLength, err
}

func (c *RedClient) rpop(key string) (string, error) {

	log.Println("Doing: RPOP " + key)
	resp := c.clientpool.Cmd("RPOP", key)
	if err := resp.Err; err != nil {
		return "", err
	}

	return resp.Str()
}

func (c *RedClient) brpop(key string) string {

	log.Println("Doing: BRPOP " + key)
	val, err := c.clientpool.Cmd("BRPOP", key, 0).List()
	if err != nil {
		log.Panic(err)
	}
	msg := val[1] // [0] is the name of the queue / list
	return msg
}

func (c *RedClient) brpoplpush(from string, to string) string {

	log.Printf("Doing: BRPOP %s %s 0"+from, to)
	msg, err := c.clientpool.Cmd("BRPOPLPUSH", from, to, 0).Str()
	if err != nil {
		log.Panic(err)
	}
	return msg
}

func (c *RedClient) set(key string, value string) (bool, error) {

	log.Println("Doing: SET " + key + " " + value)
	_, err := c.clientpool.Cmd("SET", key, value).Str()
	if err != nil {
		log.Panic(err)
		return false, err
	}
	return true, nil
}

// // ZADD Q_working_set <now>+300 queue-id.task-id
func (c *RedClient) zadd(set string, score string, member string) (int, error) {

	log.Printf("Doing: ZADD %s %s %s", set, score, member)
	count, err := c.clientpool.Cmd("ZADD", set, score, member).Int()
	if err != nil {
		log.Panic(err)
	}
	return count, nil
}

// // ZADD Q_working_set <now>+300 queue-id.task-id
func (c *RedClient) zaddUpdate(set string, score string, member string) (int, error) {

	log.Printf("Doing: ZADD %s XX CH %s %s", set, score, member)
	// ZADD the element. XX says only update existing items, CH means return us the amount
	// of *changed* elements. So 1 is good (found an item, AND changed a score. 0 is bad)
	count, err := c.clientpool.Cmd("ZADD", set, "XX", "CH", score, member).Int()
	if err != nil {
		log.Panic(err)
		return 0, err
	}
	return count, nil
}

func (c *RedClient) zrevrange(set string, from int, to int) ([]string, error) {

	log.Printf("Doing: ZREVRANGE %s %d %d", set, from, to)
	members, err := c.clientpool.Cmd("ZREVRANGE", set, from, to).List()
	if err != nil {
		log.Panic(err)
	}
	return members, nil
}

// Do an atomic from sorted list; to sorted list operation
func (c *RedClient) moveMemberFromSetToSet(from string, to string, member string) (int, error) {

	luaScript := `
	redis.call('ZREM', KEYS[1], ARGV[1]);
	local count = redis.call('ZADD', KEYS[2], ARGV[1], ARGV[2]);
	return count
	`

	timestamp := int64(time.Now().Unix())
	score := strconv.FormatInt(timestamp, 10)

	count, err := c.clientpool.Cmd("EVAL", luaScript, 2, from, to, score, member).Int()
	if err != nil {
		log.Panic(err)
	}
	return count, nil
}

// func load -- Lua for Load
