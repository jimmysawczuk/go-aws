package aws

type Client struct {
	key    string
	secret string
}

// Create a new AWS client with the AWS Access Key ID key and secret secret.
func New(key, secret string) Client {
	c := Client{key: key, secret: secret}

	return c
}
