package ymqd

// Client client对象
type Client struct {
}

// NewClient 创建新的Client
func NewClient() (*Client, error) {
	client := new(Client)
	return client, nil
}

// Pause 暂停
func (c *Client) Pause() {

}

// UnPause 取消暂停
func (c *Client) UnPause() {

}

// Close 关闭
func (c *Client) Close() error {
	return nil
}

// Empty 清空
func (c *Client) Empty() {

}

// Stats 状态
func (c *Client) Stats() error {
	return nil
}

// TimedOutMessage 超时消息
func (c *Client) TimedOutMessage() {

}
