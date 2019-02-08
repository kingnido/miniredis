package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestRedisCmd(t *testing.T) {
	cmd := NewRedisCmd(NewRedis())
	n := 100
	m := 100
	start := 100
	stop := -100
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

	if str, err := cmd.Send("zcard key"); err != nil {
		t.Errorf("could not get cardinality for set: '%v'", err)
	} else {
		got, err := strconv.Atoi(str)

		if err != nil {
			t.Errorf("int expected, got %s.", str)
			return
		}

		if got != n*m {
			t.Errorf("missing elements. expected: %d, got %d.", n*m, got)
			return
		}
	}

	str, err := cmd.Send(fmt.Sprintf("zrange key %d %d", start, stop))
	if err != nil {
		t.Errorf("error not expected: %v", err)
	}

	var result []string

	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		t.Errorf("error not expected: %v", err)
	}

	for i := 0; i < len(result); i++ {
		if got, _ := strconv.Atoi(result[i]); got != start+i {
			t.Errorf("missing element. expected: %d, got %s.", i+start, result[i])
		}
	}
}
