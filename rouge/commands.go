package rouge

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
	"github.com/pkg/errors"
)

func (c *Client) checkExpired(set string) (string, error) {

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

func (c *Client) fetchAndUpdateExpired(set string, expirationSec int) (string, error) {

	now := int64(time.Now().Unix())
	expiresAt := now + int64(expirationSec)
	score := strconv.FormatInt(expiresAt, 10)

	luaScript := `
		local members = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, 1)
		if members[1] == nil then return nil end;
		local member = members[1]
		redis.call('ZADD', KEYS[1], ARGV[2], member);
		return member
	`

	log.Printf("Doing: EVAL <luascript for get expired> with set: %v, now: %v score: %v", set, now, score)
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

func (c *Client) popQueueAndSaveKeyToSet(queueID, receivedList, targetSet string, expirationSec int) (string, error) {

	expiresAt := int64(time.Now().Unix()) + int64(expirationSec)
	score := strconv.FormatInt(expiresAt, 10)

	// TODO: Check if this works on a redis cluster, specifically because we
	// SET an item on a key that was not passed (we found it)

	// what we can do, is get the expected taskID from the previous brpoplpush
	// and compare that to what we get back. That guarantees that we can pass the task

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

	log.Printf("Doing: EVAL <luascript> 2 %s %s %s %s", receivedList, targetSet, queueID, score)
	resp := c.clientpool.Cmd("EVAL", luaScript, 2, receivedList, targetSet, queueID, score)
	if err := resp.Err; err != nil {
		if err.Error() == "No item in queue" {
			fmt.Println(err.Error())
			return "", err
		}
		log.Panic(err)
	}

	var taskID, destinationKey, msg string
	var err error

	result, _ := resp.Array()
	taskID, _ = result[0].Str()
	destinationKey, _ = result[1].Str()
	msg, err = result[2].Str()

	if err != nil {
		log.Panic("error in getting string from popqueueandsavetoset result")
	}

	if msg == "" {
		log.Panic("message was empty, which cannot be ?!?")
	}

	log.Printf("Completed: Set task (%s) with value %s on key %s", taskID, msg, destinationKey)
	return msg, nil
}

func (c *Client) del(key string) int {
	count, err := c.clientpool.Cmd("DEL", key).Int()
	if err != nil {
		log.Panic(err)
	}
	return count
}

// flush flushes the entire database. Use with care!
func (c *Client) flushdb() bool {
	r := c.clientpool.Cmd("FLUSHDB")
	if r.Err != nil {
		log.Panic(r.Err)
	}
	return true
}

func (c *Client) lpush(key string, value string) (int, error) {

	log.Printf("Doing: LPUSH %s %s", key, value)
	newLength, err := c.clientpool.Cmd("LPUSH", key, value).Int()
	if err != nil {
		log.Panic(err)
	}
	return newLength, err
}

func (c *Client) rpop(key string) (string, error) {

	log.Println("Doing: RPOP " + key)
	resp := c.clientpool.Cmd("RPOP", key)
	if err := resp.Err; err != nil {
		return "", err
	}

	return resp.Str()
}

func (c *Client) brpop(key string) string {

	log.Println("Doing: BRPOP " + key)
	val, err := c.clientpool.Cmd("BRPOP", key, 0).List()
	if err != nil {
		log.Panic(err)
	}
	msg := val[1] // [0] is the name of the queue / list
	return msg
}

func (c *Client) brpoplpush(from string, to string) (string, error) {

	log.Printf("Doing: BRPOPLPUSH %s %s 0", from, to)
	res := c.clientpool.Cmd("BRPOPLPUSH", from, to, 0)
	if err := res.Err; err != nil {
		return "", err
	}
	return res.Str()
}

func (c *Client) lpoprpush(from string, to string) (int, error) {

	log.Printf("Puting member back onto the incoming queue")
	luaScript := `
		local popped = redis.call('LPOP', KEYS[1]);
		local result = redis.call('RPUSH', KEYS[2], popped);
		return result
	`
	log.Printf("Doing: LPOPRPUSH (luascript) %s %s", from, to)

	res := c.clientpool.Cmd("EVAL", luaScript, 2, from, to)
	if err := res.Err; err != nil {
		return -1, err
	}
	if res.IsType(redis.Nil) {
		return -1, errors.New("No received members put back")
	}
	return res.Int()
}

func (c *Client) get(key string) (string, error) {

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

func (c *Client) getListLength(key string) (int, error) {

	log.Printf("Doing: LLEN %s", key)
	resp := c.clientpool.Cmd("LLEN", key)
	if err := resp.Err; err != nil {
		log.Panic(err)
	}
	if resp.IsType(redis.Nil) {
		return 0, errors.New("Nothing found at key")
	}
	return resp.Int()
}

func (c *Client) set(key string, value string) (bool, error) {

	log.Println("Doing: SET " + key + " " + value)
	_, err := c.clientpool.Cmd("SET", key, value).Str()
	if err != nil {
		log.Panic(err)
		return false, err
	}
	return true, nil
}

// // ZADD Q_working_set <now>+300 queue-id.task-id
func (c *Client) zadd(set string, score string, member string) (int, error) {

	log.Printf("Doing: ZADD %s %s %s", set, score, member)
	count, err := c.clientpool.Cmd("ZADD", set, score, member).Int()
	if err != nil {
		log.Panic(err)
	}
	return count, nil
}

// // ZADD Q_working_set <now>+300 queue-id.task-id
func (c *Client) zaddUpdate(set string, score string, member string) (int, error) {

	// ToDo: instead of returning an error if the item had been updated in the same
	// second, return success. -- will havo to be EVAL GET then ZADD
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

func (c *Client) zrange(set string, from int, to int) ([]string, error) {

	log.Printf("Doing: ZRANGE %s %d %d", set, from, to)
	members, err := c.clientpool.Cmd("ZRANGE", set, from, to).List()
	if err != nil {
		log.Panic(err)
	}
	return members, nil
}

func (c *Client) zrevrange(set string, from int, to int) ([]string, error) {

	log.Printf("Doing: ZREVRANGE %s %d %d", set, from, to)
	members, err := c.clientpool.Cmd("ZREVRANGE", set, from, to).List()
	if err != nil {
		log.Panic(err)
	}
	return members, nil
}

func (c *Client) lrange(list string, from int, to int) ([]string, error) {
	var members []string
	log.Printf("Doing: LRANGE %s %d %d", list, from, to)
	r := c.clientpool.Cmd("LRANGE", list, from, to)
	if err := r.Err; err == nil {
		members, err = r.List()
		return members, err
	}
	return members, nil
}

func (c *Client) zcount(set string, from string, to string) (int, error) {
	log.Printf("Doing: ZCOUNT %s %s %s", set, from, to)
	return c.clientpool.Cmd("ZCOUNT", set, from, to).Int()
}

func (c *Client) zrangebyscore(set string, from string, to string, limit int) (lst []string, err error) {
	resp := c.clientpool.Cmd("ZRANGEBYSCORE", set, from, to, "LIMIT", 0, limit)
	if resp.Err != nil {
		return lst, resp.Err
	}
	return resp.List()
}

// Do an atomic from sorted list; to sorted list operation
func (c *Client) moveMemberFromSetToSet(from string, to string, member string) (bool, error) {

	var removed, added int
	var err error

	luaScript := `
		local removed = redis.call('ZREM', KEYS[1], ARGV[2]);
		local count = redis.call('ZADD', KEYS[2], ARGV[1], ARGV[2]);
		return {removed, count}
	`

	timestamp := int64(time.Now().Unix())
	score := strconv.FormatInt(timestamp, 10)

	log.Printf("Doing: EVAL <luascript> %s %s %s %s", from, to, score, member)
	r := c.clientpool.Cmd("EVAL", luaScript, 2, from, to, score, member)
	if r.Err != nil {
		return false, errors.Wrap(r.Err, "lua script error")
	}
	lst, err := r.Array()
	if err != nil {
		return false, errors.Wrap(err, "lua script error: didn't get an array return")
	}

	if removed, _ = lst[0].Int(); removed == 0 {
		log.Printf("Item %v was not removed from source set!", member)
		return false, errors.New(fmt.Sprintf("item %v was not removed from source set", member))
	}
	if added, _ = lst[1].Int(); added == 0 {
		log.Printf("Item %v already existed in target set!", member)
		return false, errors.New(fmt.Sprintf("item %v already existed in target set", member))
	}

	return true, nil
}

// func load -- Lua for Load
