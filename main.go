package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	redis, _ := NewRedis()
	redisCmd := NewRedisCmd(redis)

	fmt.Print("redis > ")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "quit" {
			return
		}

		s, e := redisCmd.Send(input)
		fmt.Printf("string: '%s', error: '%v'\n", s, e)

		fmt.Print("redis > ")
	}
	fmt.Println("exiting...")
}
