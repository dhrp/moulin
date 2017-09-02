package main

import (
	"fmt"
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

func (c *RedClient) del(key string) int {
	count, err := c.clientpool.Cmd("DEL", key).Int()
	if err != nil {
		log.Panic(err)
	}
	return count
}

func (c *RedClient) lpush(key string, value string) (int, error) {
	newLength, err := c.clientpool.Cmd("LPUSH", key, value).Int()
	if err != nil {
		log.Panic(err)
		return 0, err
	}

	log.Println("pushed " + value + " to list: " + key)
	return newLength, err
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

func (c *RedClient) zrem(set string, value string) (int, error) {

	log.Println("Doing: ZREM " + set + " " + value)
	count, err := c.clientpool.Cmd("ZREM", set, value).Int()
	if err != nil {
		log.Panic(err)
		return 0, fmt.Errorf("ZREM " + set + " " + value + "failed!")
	}
	return count, nil
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
