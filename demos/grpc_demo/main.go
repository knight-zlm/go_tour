package main

import "flag"

var port string

func init() {
	flag.StringVar(&port, "p", "8008", "启动端口号")
	flag.Parse()
}

// protoc --go_out=plugins=grpc:. ./proto/*.proto
func main() {

}
