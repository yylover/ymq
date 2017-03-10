package ymqd

const (
	//MsgIDLength 消息长度
	MsgIDLength       = 16
	minValidMsgLength = MsgIDLength + 8 + 2 //Timestamp + Attempts
)

// MessageID 消息ID结构
type MessageID [MsgIDLength]byte

//Message 结构
type Message struct {
	ID MessageID

	body         []byte
	retryedCount int32 //失败重发次数
}
