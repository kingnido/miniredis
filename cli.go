package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func runCli(redisCmd *RedisCmd, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("redis > ")
		for scanner.Scan() {
			input := scanner.Text()
			if len(input) > 0 {
				if str, err := redisCmd.Send(input); err != nil {
					fmt.Printf("error: '%v'\n", err)
				} else {
					fmt.Printf("%s\n", str)
				}
			}
			fmt.Print("redis > ")
		}

		fmt.Println("exiting cli...")
	}()
}
