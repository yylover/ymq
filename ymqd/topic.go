package ymqd

import (
	"sync"
	"ymq/internal/util"
)

type Topic struct {
	messageCount uint64
	sync.RWMutex

	name          string
	channelMap    map[string]*Channel
	memoryMsgChan chan *Message
	exitChan      chan int
	waitGroup     util.WaitGroupWrapper
	exitFlag      int32

	ctx *context
}

func NewTopic(ctx *context) (*Topic, error) {
	topic := &Topic{
		messageCount:  0,
		name:          "",
		channelMap:    make(map[string]*Channel),
		memoryMsgChan: make(chan *Message),
		exitChan:      make(chan int),
		exitFlag:      0,
		ctx:           ctx,
	}

	topic.waitGroup.Warp(topic.messagePump)

	return topic, nil
}

func (t *Topic) getExistingChannelOrCreate(name string) (*Channel, error) {
	return nil, nil
}

func (t *Topic) messagePump() {

}

func (t *Topic) exit(deleted bool) error {

	return nil
}
