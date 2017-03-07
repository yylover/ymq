package ymqd

type Client struct {
}

func NewClient() (*Client, error) {
	client := new(Client)
	return client, nil
}
