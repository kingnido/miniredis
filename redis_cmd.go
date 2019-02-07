package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RedisCmd struct {
	redis *Redis
}

func NewRedisCmd(redis *Redis) *RedisCmd {
	return &RedisCmd{redis}
}

func (r *RedisCmd) Send(cmd string) (string, error) {
	s := strings.Fields(cmd)
	if len(s) == 0 {
		return "", errors.New("empty command")
	}

	switch strings.ToLower(s[0]) {
	case "get":
		return r.get(s[1:])
	case "set":
		return r.set(s[1:])
	case "del":
		return r.del(s[1:])
	case "dbsize":
		return r.dbsize(s[1:])
	case "incr":
		return r.incr(s[1:])
	case "zadd":
		return r.zadd(s[1:])
	case "zcard":
		return r.zcard(s[1:])
	case "zrank":
		return r.zrank(s[1:])
	case "zrange":
		return r.zrange(s[1:])
	case "help":
		return r.help()
	}

	return "", errors.New("unknown command")
}

// string or error
func (r *RedisCmd) get(params []string) (string, error) {
	if len(params) != 1 {
		return "", errors.New(fmt.Sprintf("expecting %d params, got %d", 1, len(params)))
	}

	s, err := r.redis.Get(params[0]) // key
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(s)
	return string(j), nil
}

// OK or error
func (r *RedisCmd) set(params []string) (string, error) {
	switch len(params) {
	case 2:
		// key val
		r.redis.Set(params[0], params[1])
		return "OK", nil

	case 4:
		// key val ex delta
		if strings.ToLower(params[2]) != "ex" {
			return "", errors.New("invalid command")
		}

		i, err := strconv.Atoi(params[3])
		if err != nil {
			return "", errors.New("expiration is not a number")
		}

		if err := r.redis.SetExpire(params[0], params[1], i); err != nil {
			return "", err
		}

		return "OK", nil
	}

	return "", errors.New(fmt.Sprintf("expecting 2 or 4d params, got %d", len(params)))
}

// OK or error
func (r *RedisCmd) del(params []string) (string, error) {
	if len(params) != 1 {
		return "", errors.New(fmt.Sprintf("expecting %d params, got %d", 1, len(params)))
	}

	if err := r.redis.Del(params[0]); err != nil {
		return "", err
	}

	return "OK", nil
}

// int or error
func (r *RedisCmd) dbsize(params []string) (string, error) {
	if len(params) != 0 {
		return "", errors.New(fmt.Sprintf("expecting %d params, got %d", 0, len(params)))
	}

	i := r.redis.DbSize()

	j, _ := json.Marshal(i)
	return string(j), nil
}

// int or error
func (r *RedisCmd) incr(params []string) (string, error) {
	if len(params) != 1 {
		return "", errors.New(fmt.Sprintf("expecting %d params, got %d", 1, len(params)))
	}

	i, err := r.redis.Incr(params[0])
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(i)
	return string(j), nil
}

// OK or error
func (r *RedisCmd) zadd(params []string) (string, error) {
	if len(params) != 3 {
		return "", errors.New(fmt.Sprintf("expecting %d params, got %d", 3, len(params)))
	}

	i, err := strconv.Atoi(params[1])
	if err != nil {
		return "", errors.New("score is not a number")
	}

	_, err = r.redis.ZAdd(params[0], i, params[2])

	if err != nil {
		return "", err
	}

	return "OK", nil
}

// int or error
func (r *RedisCmd) zcard(params []string) (string, error) {
	if len(params) != 1 {
		return "", errors.New(fmt.Sprintf("expecting %d params, got %d", 1, len(params)))
	}

	i, err := r.redis.ZCard(params[0])
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(i)
	return string(j), nil
}

// int or error
func (r *RedisCmd) zrank(params []string) (string, error) {
	if len(params) != 2 {
		return "", errors.New(fmt.Sprintf("expecting %d params, got %d", 2, len(params)))
	}

	i, err := r.redis.ZRank(params[0], params[1])
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(i)
	return string(j), nil
}

// list of string or error
func (r *RedisCmd) zrange(params []string) (string, error) {
	if len(params) != 3 {
		return "", errors.New(fmt.Sprintf("expecting %d params, got %d", 3, len(params)))
	}

	start, err := strconv.Atoi(params[1])
	if err != nil {
		return "", errors.New("start is not a number")
	}

	stop, err := strconv.Atoi(params[2])
	if err != nil {
		return "", errors.New("stop is not a number")
	}

	s, err := r.redis.ZRange(params[0], start, stop)
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(s)
	return string(j), nil
}

func (r *RedisCmd) help() (string, error) {
	return `Commands:
	get <key:string>
		returns value for key.
		returns error if key does not exist or value is not string

	set <key:string> <value:string>
		set value for key

	set <key:string> <value:string> ex <delta:int>
		set value for key with expiration time in seconds
		returns error if delta is not int or is less than 1

	del <key:string>
		delete key
		returns error if key does not exist

	dbsize
		returns the number of keys

	incr <key:string>
		increments value for key
		returns the new value as int
		returns error if value is not convertable to int

	zadd <key:string> <score:int> <member:string>
		adds a member with score to set for key
		create a set for key, if key does not exist
		if member is updated, it will be placed in the proper position according to the new score
		returns error if value for key is not a set

	zcard <key:string>
		returns the number of members for key
		returns error if key does not exist or value is not a set

	zrank <key:string> <member:string>
		returns position of the member in the set for key
		returns error if key does not exist or member does not exist or value is not a set

	zrange <key:string> <int:start> <int:stop>
		returns list of members in positions start to stop in the set for key
		invalid indexes in the start - stop range are ignored
		returns error if key does not exist or value is not a set`, nil
}
