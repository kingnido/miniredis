package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestRedisCmd(t *testing.T) {
	cmd := NewRedisCmd(NewRedis())
	n := 100
	m := 100
	key := "key"

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			for j := 0; j < m; j++ {
				x := i*m + j
				cmd.Send(fmt.Sprintf("zadd %s %d %d", key, x, x))
			}
		}(i)
	}

	wg.Wait()

	t.Log(cmd.Send("zcard key"))
	t.Log(cmd.Send("zrange key 100 -100"))
}
