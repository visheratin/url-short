package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func BenchmarkListener(b *testing.B) {
	// need to start server before testing and specify its address.
	address := "127.0.0.1:37485"
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		str := time.Now().String()
		time.Sleep(time.Millisecond)
		b.StartTimer()
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
	}
}
