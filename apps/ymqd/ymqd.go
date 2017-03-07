package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"ymq/internal/app"
	"ymq/internal/version"
	"ymq/ymqd"

	"github.com/BurntSushi/toml"
	options "github.com/mreiferson/go-options"
)

func ymqdFlagSet(options *ymqd.Options) *flag.FlagSet {
	flagset := flag.NewFlagSet("ymqd", flag.ExitOnError)

	flagset.Bool("version", false, "print the version string")
	flagset.Bool("verbose", false, "enable verbose logging")
	flagset.String("config", "", "path to config file")
	flagset.String("log-prefix", "[ymqd]", "log message prefix")
	flagset.Int64("node-id", options.ID, "unique part for messages IDs")

	flagset.String("https-address", options.HTTPSAddress, "<addr>:<port> to listen on the https client")
	flagset.String("http-address", options.HTTPAddress, "<addr>:<port> to listen on the http client")
	flagset.String("tcp-address", options.TCPAddress, "<addr>:<port> to listen on the tcp client")
	var lookupTcpAddress app.StringArray
	flagset.Var(&lookupTcpAddress, "lookupd-tcp-address", "lookupd TCP address (may be given multiple times)")
	flagset.Duration("http-client-connect-timeout", options.HTTPClientConnectTimeout, "timeout for HTTP connect")
	flagset.Duration("http-client-request-timeout", options.HTTPClientRequestTimeout, "timeout for HTTP request")

	flagset.String("msg-timeout", options.MsgTimeout.String(), "duration to wait before auto-requeing a message")
	flagset.Duration("max-msg-timeout", options.MaxMsgTimeout, "maximum duration before a message will timeout")
	flagset.Int64("max-msg-size", options.MaxMsgSize, "maximum size of a single message in bytes")
	flagset.Duration("max-req-timeout", options.MaxReqTimeout, "maximum requeuing timeout for a message")
	flagset.Int64("max-body-size", options.MaxBodySize, "maximum size of a single command body")

	flagset.Int64("mem-queue-size", options.MemQueueSize, "number of messages to keep in memory (per topic/channel)")
	flagset.Int64("max-bytes-per-file", options.MaxBytesPerFile, "number of bytes per diskqueue file before rolling")
	flagset.Int64("sync-every", options.SyncEvery, "number of messages per diskqueue fsync")
	flagset.Duration("sync-timeout", options.SyncTimeout, "duration of time per diskqueue fsync")

	return flagset
}

type config map[string]interface{}

func main() {
	// s := ymqd.NewHTTPServer()
	opts, err := ymqd.NewOptions()
	if err != nil {
		log.Fatal(err)
	}

	flagset := ymqdFlagSet(opts)
	flagset.Parse(os.Args[1:])

	//打印版本号选项直接退出
	if flagset.Lookup("version").Value.(flag.Getter).Get().(bool) {
		fmt.Printf("%s\n", version.String("ymqd"))
		os.Exit(0)
	}

	var cfg config
	configFile := flagset.Lookup("config").Value.String()
	if configFile != "" {
		_, err := toml.DecodeFile(configFile, &cfg)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("options : %v\n", opts)
	options.Resolve(opts, flagset, cfg)

	fmt.Printf("options : %v\n", opts)
	fmt.Printf("config:%v\n", cfg)
}
