package aws

type Client struct {
	key    string
	secret string
}

func New(key, secret string) Client {
	c := Client{key: key, secret: secret}

	return c
}
