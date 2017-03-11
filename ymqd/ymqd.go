package ymqd

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
	"ymq/internal/util"
)

type errStore struct {
	err error
}

// Ymqd 应用程序
type Ymqd struct {
	clientIDSequence int64
	sync.RWMutex

	topicMap  map[string]*Topic
	startTime time.Time

	tcpListener   net.Listener
	httpListener  net.Listener
	httpsListener net.Listener

	poolSize int

	notifyChan           chan interface{}
	optsNotificationChan chan struct{}
	exitChan             chan int
	waitGroup            util.WaitGroupWrapper

	ops *Options
}

// NewYmqd 创建新的服务器
func NewYmqd(ops *Options) (*Ymqd, error) {
	//TODO 认为Ymqd的初始化是不需要加锁的，因为全局只有一个app
	ymq := new(Ymqd)
	ymq.ops = ops

	if ops.Logger == nil {
		ops.Logger = log.New(os.Stderr, ops.LogPrefix, log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	ymq.startTime = time.Now()
	ymq.topicMap = make(map[string]*Topic)
	ymq.exitChan = make(chan int)
	ymq.notifyChan = make(chan interface{})
	ymq.optsNotificationChan = make(chan struct{})

	return ymq, nil
}

//Main 主线程loop
func (y *Ymqd) Main() error {

	httplistener, err := NewHTTPServer(nil)
	if err != nil {
		return err
	}
	// y.httpListener = httplistener

	y.waitGroup.Wrap(func() {
		httplistener.Serve()
	})
	//等待http server //TODO 放到后台进行

	return nil
}

func (y *Ymqd) logf(f string, args ...interface{}) {
	y.ops.Logger.Output(2, fmt.Sprintf(f, args...))
}

func readFlag() {
	var ip = flag.Int("flagname", 1234, "help message for flagname")
	flag.Parse()

	fmt.Printf("ip: %d\n", ip)
}

// Exit 程序退出
func (y *Ymqd) Exit() {

}

// GetTopic 线程安全操作，返回一个*Topic指针对象(可能创建的)
func (y *Ymqd) GetTopic(topicName string) *Topic {
	y.RLock()
	topic, ok := y.topicMap[topicName]
	y.RUnlock()

	if ok {
		return topic
	}
	// 要修改y.topic 要再重新判断一次
	y.Lock()
	topic, ok = y.topicMap[topicName]
	if ok {
		y.Unlock()
		return topic
	}

	topic, err := NewTopic(&context{y})
	if err != nil {
		return nil
	}
	y.topicMap[topicName] = topic
	y.Unlock()

	return topic
}

// GetExistingTopic 获取已经存在的Topic
func (y *Ymqd) GetExistingTopic(topicName string) (*Topic, error) {
	y.RLock()
	defer y.RUnlock()
	if topic, ok := y.topicMap[topicName]; ok {
		return topic, nil
	}

	return nil, errors.New("there is not topic named:" + topicName)
}

// DeleteExistingTopic 删除已经存在的Topic
func (y *Ymqd) DeleteExistingTopic(topicName string) error {
	y.RLock()
	defer y.Unlock()

	if _, ok := y.topicMap[topicName]; ok {
		y.RUnlock()
		delete(y.topicMap, topicName)
		return nil
	}

	//不能直接调用delete(map, key) ,要先对topic进行资源试岗

	return fmt.Errorf("There is not a topic named:%s", topicName)
}

//Notify 通知?
func (y *Ymqd) Notify(v interface{}) {

}

// channels 获取所有的channels
func (y *Ymqd) channels() []*Channel {
	var channels []*Channel
	y.RLock()
	defer y.RUnlock()

	for _, topic := range y.topicMap {
		topic.RLock()
		for _, c := range topic.channelMap {
			channels = append(channels, c)
		}
		topic.RUnlock()
	}

	return channels
}
