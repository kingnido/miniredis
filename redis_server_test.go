package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestRedisServer(t *testing.T) {
	// test server receiving ZADD, with n concurrent requests

	server := httptest.NewServer(commandHandler(NewRedisCmd(NewRedis())))
	defer server.Close()

	newRequest := func(cmd string) *http.Request {
		r, err := http.NewRequest("POST", server.URL+"/", strings.NewReader(cmd))
		if err != nil {
			t.Fatal(err)
		}

		return r
	}

	responseBodyFromRequest := func(cmd string) (string, error) {
		resp, err := http.DefaultClient.Do(newRequest(cmd))
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		if err != nil {
			return "", err
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(buf), nil
	}

	n := 50
	m := 200
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
				cmd := fmt.Sprintf("zadd %s %d %d", key, x, x)
				_, err := responseBodyFromRequest(cmd)
				if err != nil {
					t.Errorf("error on request: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()

	cmd := fmt.Sprintf("zcard %s", key)

	str, err := responseBodyFromRequest(cmd)
	if err != nil {
		t.Errorf("could not get cardinality for set: '%v'", err)
		return
	}

	got, err := strconv.Atoi(str)
	if err != nil {
		t.Errorf("int expected, got %s.", str)
		return
	}

	if got != n*m {
		t.Errorf("missing elements. expected: %d, got %d.", n*m, got)
		return
	}

	cmd = fmt.Sprintf("zrange key %d %d", start, stop)

	str, err = responseBodyFromRequest(cmd)
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
