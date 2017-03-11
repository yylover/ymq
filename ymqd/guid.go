package ymqd

import (
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type guid int64

/**
 * [伪毫秒时间序列][10位nodeId序列][12位sequence 序列]
 * [*************][**********][************]
 * @type {[type]}
 */
const (
	nodeIDBits     = uint64(10)
	sequenceBits   = uint64(12)
	nodeIDShift    = sequenceBits
	timestampShift = sequenceBits + nodeIDBits
	sequenceMask   = int64(-1) ^ (int64(-1) << sequenceBits)

	// ( 2012-10-28 16:23:42 UTC ).UnixNano() >> 20
	twepoch = int64(1288834974288)
)

//ErrTimeBackwards 时间后退了
var ErrTimeBackwards = errors.New("timestamp backwards")

//ErrSequenceExpired Sequence 后退了
var ErrSequenceExpired = errors.New("Sequence has expired")

//ErrIDBackwards id 后退了
var ErrIDBackwards = errors.New("ID went backwards")

type guidFactory struct {
	sync.Mutex

	nodeID        int64
	sequence      int64
	lastTimeStamp int64
	lastID        guid
}

//NewGUIDFactory 创建对象
func newGUIDFactory(nodeID int64) *guidFactory {
	return &guidFactory{
		nodeID: nodeID,
	}
}

//NewGUID 生成GUID
func (f *guidFactory) NewGUID() (guid, error) {
	//假的毫秒时间
	f.Lock()

	ts := time.Now().UnixNano() >> 20

	if ts < f.lastTimeStamp {
		f.Unlock()
		return 0, ErrTimeBackwards
	}

	if f.lastTimeStamp == ts { //时间戳相同，通过sequence来区分
		f.sequence = (f.sequence + 1) & sequenceMask //通过位运算实现模运算
		if f.sequence == 0 {                         //超出个数
			f.Unlock()
			return 0, ErrSequenceExpired
		}
	} else {
		f.sequence = 0
	}

	f.lastTimeStamp = int64(ts)

	//timestamp + nodeId + sequence 生成唯一的guid
	id := guid(((ts - twepoch) << timestampShift) |
		(f.nodeID << nodeIDShift) |
		f.sequence)

	if id <= f.lastID {
		f.Unlock()
		return 0, ErrIDBackwards
	}

	f.Unlock()
	return id, nil
}

// 将int64转换为[8]byte
func (g guid) Hex() MessageID {
	var h MessageID
	var b [8]byte //数组

	b[7] = byte(g >> 56)
	b[6] = byte(g >> 48)
	b[5] = byte(g >> 40)
	b[4] = byte(g >> 32)
	b[3] = byte(g >> 24)
	b[2] = byte(g >> 16)
	b[1] = byte(g >> 8)
	b[0] = byte(g)

	hex.Encode(h[:], b[:])

	return h
}
