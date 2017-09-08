package main

import (
	"errors"
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

func (c *RedClient) checkExpired(set string) (string, error) {

	now := int64(time.Now().Unix())

	log.Printf("Doing: ZRANGEBYSCORE %s -inf %d LIMIT 0 1", set, now)
	resp := c.clientpool.Cmd("ZRANGEBYSCORE", set, "-inf", now, "LIMIT", 0, 1)
	if err := resp.Err; err != nil {
		log.Panic(err)
	}

	list, err := resp.List()
	if err != nil {
		log.Panic(err)
	}
	if len(list) != 1 {
		return "", errors.New("Not exactly one expired item returned")
	}
	return list[0], nil
}

func (c *RedClient) fetchAndUpdateExpired(set string, expirationSec int) (string, error) {

	now := int64(time.Now().Unix())
	expiresAt := now + int64(expirationSec)
	score := strconv.FormatInt(expiresAt, 10)

	// ZRANGEBYSCORE test.sorted_sets.future -inf 1504764565 LIMIT 0 1
	// redis.call('ZADD', destinationSet, ARGV[2], taskID);

	luaScript := `
		local members = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, 1)
		if members[1] == nil then return nil end;
		local member = members[1]
		redis.call('ZADD', KEYS[1], ARGV[2], member);
		return member
	`

	log.Printf("Doing: ZRANGEBYSCORE %s -inf %d LIMIT 0 1", set, now)
	log.Printf("Doing: ZADD %s <returnedmember> %s", set, score)
	// resp := c.clientpool.Cmd("ZRANGEBYSCORE", set, "-inf", now, "LIMIT", 0, 1)
	resp := c.clientpool.Cmd("EVAL", luaScript, 1, set, now, score)
	if err := resp.Err; err != nil {
		log.Panic(err)
	}

	if resp.IsType(redis.Nil) {
		return "", errors.New("No expired members retrieved")
	}

	member, err := resp.Str()
	if err != nil {
		log.Panic(err)
	}

	return member, nil
}

func (c *RedClient) popQueueAndSaveKeyToSet(queueID string, expirationSec int) (string, error) {

	expiresAt := int64(time.Now().Unix()) + int64(expirationSec)
	score := strconv.FormatInt(expiresAt, 10)

	// TODO: Check if this works on a redis cluster, specifically because we
	// SET an item on a key that was not passed (we found it)

	luaScript := `
		local queueID = ARGV[1];
		local score = ARGV[2];
		local receivedList = KEYS[1];
		local destinationSet = KEYS[2];

		local value = redis.call('RPOP', receivedList);
		if value == false then
			local errorTable = {};
			errorTable["err"] = "No item in queue";
			return errorTable;
		end

		local jsonObj = cjson.decode(value);
		local taskID = jsonObj['id'];
		local destinationKey = queueID .. "." .. taskID;
		redis.call('SET', destinationKey, value);
		local count = redis.call('ZADD', destinationSet, score, taskID);
		return {taskID, destinationKey, value};
		`

	receivedList := queueID + ".received"
	destinationSet := queueID + ".running"

	resp := c.clientpool.Cmd("EVAL", luaScript, 2, receivedList, destinationSet, queueID, score)
	if err := resp.Err; err != nil {
		if err.Error() == "No item in queue" {
			return "", err
		} else {
			log.Panic(err)
		}
	}

	result, _ := resp.Array()
	taskID, _ := result[0].Str()
	destinationKey, _ := result[1].Str()
	msg, _ := result[2].Str()

	log.Printf("COMPLETED: Set task (%s) on key %s, body: %s", taskID, destinationKey, msg)
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

	log.Printf("Doing: BRPOPLPUSH %s %s 0", from, to)
	msg, err := c.clientpool.Cmd("BRPOPLPUSH", from, to, 0).Str()
	if err != nil {
		log.Panic(err)
	}
	return msg
}

func (c *RedClient) get(key string) (string, error) {

	log.Printf("Doing: GET %s", key)
	resp := c.clientpool.Cmd("GET", key)
	if err := resp.Err; err != nil {
		log.Panic(err)
	}
	if resp.IsType(redis.Nil) {
		return "", errors.New("Nothing found at key")
	}
	return resp.Str()
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

	log.Printf("Doing: EVAL <luascript> %s %s %s %s", from, to, score, member)
	count, err := c.clientpool.Cmd("EVAL", luaScript, 2, from, to, score, member).Int()
	if err != nil {
		log.Panic(err)
	}
	return count, nil
}

// func load -- Lua for Load
