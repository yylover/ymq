package ymqd

import (
	"flag"
	"fmt"
)

// Ymqd
type Ymqd struct {
}

// NewYmqd 创建新的服务器
func NewYmqd() (*Ymqd, error) {
	ymq := new(Ymqd)
	return ymq, nil
}

func readFlag() {
	var ip = flag.Int("flagname", 1234, "help message for flagname")
	flag.Parse()

	fmt.Printf("ip: %d\n", ip)
}
