// This is a tool for load testing of a URL shortening microservice.
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	// set an address of a running service
	address := "localhost:37485"
	connNum := 1000
	var wg sync.WaitGroup
	startTime := time.Now()
	rand.Seed(time.Now().Unix())
	for i := 0; i < connNum; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup) {
			defer w.Done()
			val := time.Now().Unix()
			val += rand.Int63()
			str := strconv.FormatInt(val, 10)

			conn, err := net.Dial("tcp", address)
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			_, err = conn.Write([]byte(str + "\n"))
			if err != nil {
				panic(err)
			}
			_, err = bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				panic(err)
			}
		}(&wg)
	}
	wg.Wait()
	elapsed := time.Since(startTime)
	fmt.Println(elapsed.String())
}
