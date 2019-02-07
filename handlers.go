package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func runServer(redisCmd *RedisCmd, port string, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		fmt.Println("starting server on port", port)
		err := http.ListenAndServe(port, commandHandler(redisCmd))
		if err != nil {
			fmt.Println("error on listen:", err)
		}
	}()
}

func commandHandler(redisCmd *RedisCmd) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "could not read request's body", http.StatusBadGateway)
			return
		}

		cmd := string(body[:])

		if str, err := redisCmd.Send(cmd); err == nil {
			fmt.Fprint(w, str)
		} else {
			http.Error(w, "ERROR: "+err.Error(), http.StatusBadRequest)
		}
	})
}
