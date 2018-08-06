// This is a microservice for URL shortening.
// It works over TCP connections for both generating
// codes for input links and extracting initial links
// from the encoded links.
package main

import "github.com/visheratin/url-short/config"

func main() {
	err := config.Init("config.json")
	if err != nil {
		panic(err)
	}
	startListen()
}
