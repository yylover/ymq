package ymqd

import "sync"

type Consumer interface {
	Pause()
	UnPause()
	Close() error
	Empty()
	Stats() error
	TimedOutMessage()
}

type Channel struct {
	messageCount uint64
	timeoutCount uint64
	requeueCount uint64

	sync.RWMutex

	topicName string
	name      string
	ctx       *context

	memoryMsgChan chan *Message
	exitFlag      int32
	exitMutext    sync.RWMutex

	clients map[int64]*Consumer
}

// NewChannel 创建新的Channel
func NewChannel(topicName, channelName string, ctx *context) (*Channel, error) {

	c := &Channel{
		topicName:     topicName,
		name:          channelName,
		ctx:           ctx,
		clients:       make(map[int64]*Consumer, 100),
		memoryMsgChan: make(chan *Message), //大小
	}

	return c, nil
}
