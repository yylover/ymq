package ymqd

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"ymq/internal/util"
)

// Topic 类
type Topic struct {
	messageCount uint64 //可以用以字节对齐
	sync.RWMutex        //锁

	name           string                //topic name
	channelMap     map[string]*Channel   //map
	memoryMsgChan  chan *Message         //内存队列
	exitChan       chan int              //退出信号
	waitGroup      util.WaitGroupWrapper //线程控制
	exitFlag       int32                 //是否已经退出
	chanUpdateChan chan int

	ephemeral      bool         //临时队列
	deleteCallback func(*Topic) //回调
	deleter        sync.Once    //保证只执行一次

	paused    int32 //暂停
	pauseChan chan bool

	guidGenerator *guidFactory
	logger        Logger   //logger
	ctx           *context //context
}

// NewTopic 创建新的Topic
func NewTopic(ctx *context) (*Topic, error) {
	topic := &Topic{
		messageCount:  0,
		name:          "",
		channelMap:    make(map[string]*Channel),
		memoryMsgChan: make(chan *Message),
		exitChan:      make(chan int),
		exitFlag:      0,
		ctx:           ctx,
		logger:        &LoggerLocal{},
		guidGenerator: newGUIDFactory(0),
	}

	topic.waitGroup.Wrap(topic.messagePump)

	return topic, nil
}

func (t *Topic) messagePump() {
	var msg *Message

	for t.exitFlag == 0 {
		select {
		case msg = <-t.memoryMsgChan:

		}

		t.logger.Output(1, fmt.Sprintf("New Message:%v", msg))
		if msg != nil {
			//复制n份，发送给各个channel
		}
	}

}

func (t *Topic) exit(deleted bool) error {
	atomic.StoreInt32(&t.exitFlag, 1)
	return nil
}

//判断是否正在退出
func (t *Topic) exiting() bool {
	return atomic.LoadInt32(&t.exitFlag) == 1
}

// GetChannel 获取name对应的channel， 如果不存在，自动创建
func (t *Topic) GetChannel(name string) *Channel {
	return nil
}

// getOrCreateChannel 获取name对应的channel,不存在就创建
func (t *Topic) getOrCreateChannel(name string) (*Channel, error) {
	var err error

	t.Lock()
	channel, ok := t.channelMap[name]
	if !ok {
		//Createa
		channel, err = NewChannel(t.name, name, t.ctx)
		if err != nil {

		}
		t.channelMap[name] = channel
	}
	t.Unlock()
	return channel, nil
}

//GetExistingChannel 获取已经存在的channel
func (t *Topic) GetExistingChannel(name string) (*Channel, bool) {
	return nil, true
}

// DeleteExistingChannel 删除已经存在的channel
func (t *Topic) DeleteExistingChannel(name string) error {
	return nil
}

// PutMessage 发送消息
func (t *Topic) PutMessage(msg *Message) error {

	err := t.put(msg)

	return err
}

// PutMessages 批量发送消息
func (t *Topic) PutMessages(msgs []*Message) error {
	var err error
	for _, value := range msgs {
		err = t.put(value)
		if err != nil {
			return err
		}
	}

	return nil
}

// 发送消息
func (t *Topic) put(msg *Message) error {
	//阻塞
	t.memoryMsgChan <- msg
	return nil
}

//Depth 返回消息总数
func (t *Topic) Depth() int64 {
	return int64(len(t.memoryMsgChan))
}

//Delete 清空所有的topic,channel数据，并关闭
func (t *Topic) Delete() error {
	return nil
}

// Close 持久化所有的数据并关闭
func (t *Topic) Close() error {
	return nil
}

// Empty 清空所有的Topic 数据
func (t *Topic) Empty() error {
	return nil
}

// Flush 将所有的数据清空
func (t *Topic) Flush() error {
	return nil
}

//Pause 暂停
func (t *Topic) Pause() error {
	return t.dopause(true)
}

// UnPause 取消暂停
func (t *Topic) UnPause() error {
	return t.dopause(false)
}

func (t *Topic) dopause(pause bool) error {
	if pause {
		atomic.StoreInt32(&t.paused, 1)
	} else {
		atomic.StoreInt32(&t.paused, 0)
	}

	//通知各个协程
	return nil
}

//IsPaused 判断是否暂停
func (t *Topic) IsPaused() bool {
	return atomic.LoadInt32(&t.paused) == 1
}

// GenerateID 生成id,失败时尝试重新生成id
func (t *Topic) GenerateID() MessageID {

retry:
	id, err := t.guidGenerator.NewGUID()
	if err != nil {
		t.logger.Output(2, err.Error())
		time.Sleep(time.Microsecond)
		goto retry
	}

	return id.Hex()
}
