package main

import (
	"flag"
	"sync"
)

func main() {
	// initialize redis structures
	redis := NewRedis()
	redisCmd := NewRedisCmd(redis)

	// set terminal options
	portPtr := flag.String("port", "", "run server on port")
	cliPrt := flag.Bool("cli", false, "run cli")
	flag.Parse()

	wg := &sync.WaitGroup{}

	// run server?
	port := ":" + *portPtr
	if len(port) > 1 {
		runServer(redisCmd, port, wg)
	}

	// run cli?
	if *cliPrt {
		runCli(redisCmd, wg)
	}

	wg.Wait()
}
