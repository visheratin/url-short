package main

import (
	"bufio"
	"net"
	"strconv"
	"strings"

	"github.com/visheratin/url-short/config"
	"github.com/visheratin/url-short/convert"
	"github.com/visheratin/url-short/log"
	"github.com/visheratin/url-short/storage"
)

// startListener creates a Converter instance with an appropriate
// permanent storage and starts listening for incoming connections.
func startListen() {
	config, err := config.Config()
	if err != nil {
		log.Log().Error.Println(err)
		panic(err)
	}
	storage, err := storage.Instance()
	if err != nil {
		log.Log().Error.Println(err)
	}
	prefix := config.Prefix
	converter := convert.NewConverter(config.CodeLength, storage)
	internalPort := strconv.Itoa(config.Port)
	internalAddress := ":" + internalPort
	ln, err := net.Listen("tcp", internalAddress)
	if err != nil {
		log.Log().Error.Println(err)
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Log().Error.Println(err)
			return
		}
		go handleConnection(conn, converter, prefix)
	}
}

// handleConnection reads the incoming string, figures whether it is a link that
// need to be encoded or an encoded link. Depending on the type of the incoming
// string, the method performs forward or backward encoding and sends the
// result to the client.
func handleConnection(c net.Conn, converter *convert.Converter, prefix string) {
	rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	defer c.Close()
	input, err := rw.ReadString('\n')
	if err != nil {
		log.Log().Error.Println(err)
		return
	}
	input = strings.Trim(input, "\n")
	prefixIdx := len(prefix)
	var result string
	if len(input) >= prefixIdx && input[0:prefixIdx] == prefix {
		code := strings.Replace(input, prefix, "", -1)
		code = strings.Trim(code, "/")
		result = converter.Extract(code)
	} else {
		code, err := converter.Load(input)
		if err != nil {
			return
		}
		result = prefix + code
	}
	result = result + "\n"
	_, err = rw.WriteString(result)
	if err != nil {
		log.Log().Error.Println(err)
		return
	}
	err = rw.Flush()
	if err != nil {
		log.Log().Error.Println(err)
		return
	}
}
