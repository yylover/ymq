package ymqd

import (
	"crypto/md5"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"
)

type Options struct {
	ID                       int64         `flag:"node-id" cfg:"id"`
	Verbose                  bool          `flag:"verbose"`
	LogPrefix                string        `flag:"log-prefix"`
	TCPAddress               string        `flag:"tcp-address"`
	HTTPAddress              string        `flag:"http-address" cfg:"http-address"`
	HTTPSAddress             string        `flag:"https-address"`
	LookupTCPAddress         []string      `flag:"lookupd-tcp-address" cfg:"nsqlookupd_tcp_addresses"`
	HTTPClientConnectTimeout time.Duration `flag:"http-client-connect-timeout" cfg:"http_client_connect_timeout"`
	HTTPClientRequestTimeout time.Duration `flag:"http-client-request-timeout" cfg:"http_client_request_timeout"`

	//diskqueue Options
	MemQueueSize    int64         `flag:"mem-queue-size"`
	MaxBytesPerFile int64         `flag:"max-bytes-per-file"`
	SyncEvery       int64         `flag:"sync-every"`
	SyncTimeout     time.Duration `flag:"sync-timeout"`

	//msg and command option
	MsgTimeout    time.Duration
	MaxMsgTimeout time.Duration
	MaxMsgSize    int64
	MaxBodySize   int64
	MaxReqTimeout time.Duration
	ClientTimeout time.Duration

	//compression
	DeflateEnabled  bool
	MaxDeflateLevel int
	SnappyEnabled   bool

	Logger Logger
}

func NewOptions() (*Options, error) {

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(hostname)
	h := md5.New()
	io.WriteString(h, hostname)

	// fmt.Printf("md5 :%x\n", md5.Sum([]byte(hostname)))
	// fmt.Printf("md5 :%x\n", h.Sum(nil))
	defaultId := int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024)
	options := &Options{
		ID:        defaultId,
		LogPrefix: "[ymqd]",

		TCPAddress:               "0.0.0.0:8480",
		HTTPAddress:              "0.0.0.0:8481",
		HTTPSAddress:             "0.0.0.0:8482",
		HTTPClientConnectTimeout: 2 * time.Second,
		HTTPClientRequestTimeout: 2 * time.Second,

		LookupTCPAddress: make([]string, 0),

		MemQueueSize:    10000,
		MaxBytesPerFile: 100 * 1024 * 1024,
		SyncEvery:       2500,
		SyncTimeout:     2 * time.Second,

		MsgTimeout:    60 * time.Second,
		MaxMsgTimeout: 15 * time.Minute,
		MaxMsgSize:    1024 * 1024,
		MaxBodySize:   5 * 1024 * 1024,
		MaxReqTimeout: 1 * time.Hour,
		ClientTimeout: 60 * time.Second,

		DeflateEnabled:  true,
		MaxDeflateLevel: 6,
		SnappyEnabled:   true,
	}

	return options, nil
}
