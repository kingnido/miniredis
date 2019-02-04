package main

import (
	"errors"
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
	}

	return "", errors.New("unknown command")
}

func (r *RedisCmd) get(params []string) (string, error) {
	if len(params) != 1 {
		return "", errors.New("invalid command")
	}
	return r.redis.Get(params[0])
}

func (r *RedisCmd) set(params []string) (string, error) {
	switch len(params) {
	case 2:
		if err := r.redis.Set(params[0], params[1]); err != nil {
			return "", err
		}
		return "OK", nil
	case 4:
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
	return "", errors.New("invalid command")
}

func (r *RedisCmd) del(params []string) (string, error) {
	if len(params) != 1 {
		return "", errors.New("invalid command")
	}

	if err := r.redis.Del(params[0]); err != nil {
		return "", err
	}

	return "OK", nil
}

func (r *RedisCmd) dbsize(params []string) (string, error) {
	if len(params) != 0 {
		return "", errors.New("invalid command")
	}

	return strconv.Itoa(r.redis.DbSize()), nil
}

func (r *RedisCmd) incr(params []string) (string, error) {
	if len(params) != 1 {
		return "", errors.New("invalid command")
	}

	i, err := r.redis.Incr(params[0])
	if err != nil {
		return "", err
	}

	return strconv.Itoa(i), nil
}
