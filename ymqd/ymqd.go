package ymqd

import (
	"flag"
	"fmt"
)

type Ymqd struct {
}

func NewYmqd() (*Ymqd, error) {
	ymq := new(Ymqd)
	return ymq, nil
}

func readFlag() {
	var ip = flag.Int("flagname", 1234, "help message for flagname")
	flag.Parse()

	fmt.Printf("ip: %d\n", ip)
}
