package main

import (
	"flag"
	"fmt"

	"ymq/ymqd"
)

func readFlag() {
	var ip = flag.Int("flagname", 1234, "help message for flagname")
	flag.Parse()

	fmt.Printf("ip: %d\n", ip)
}

func main() {

	readFlag()

	server := ymqd.NewHttpServer()

	fmt.Println("hello to my Topic")

}
