package main

import (
	"flag"
	"fmt"
)

func main() {

	var version string
	var d bool
	// version = flag.String("version", "1.9", "print verion")
	// version = flag.String("v", "1.9", "print version v")
	//
	fmt.Printf("args :%v\n", flag.Args())

	flag.StringVar(&version, "version", "1.9", "print version")
	flag.StringVar(&version, "v", "1.9", "print version")
	flag.BoolVar(&d, "daemon", false, "print daemoned ")
	flag.BoolVar(&d, "d", false, "print daemoned")
	flag.Parse()

	fmt.Printf("version :%s", version)
	fmt.Printf("version :%v\n", d)

}
