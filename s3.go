package aws

import (
	"bytes"
	"strconv"

	// related to signing request
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"time"

	// making request
	"crypto/tls"
	"net/http"
	"net/url"
)

type S3Client struct {
	key    string
	secret string
}

type S3Request struct {
	verb string

	request_uri  string
	request_host string

	content_md5  string
	content_type string
	request_time time.Time

	content []byte

	signing_uri string

	signature string
}

// Creates a new S3Client using the AWS information in the Client
func (c Client) NewS3() S3Client {
	s := S3Client{key: c.key, secret: c.secret}
	return s
}

// Performs an authenticated GET request on file_name in bucket.
func (s *S3Client) Get(bucket, file_name string) (*bytes.Buffer, *http.Response, error) {

	s3_req := S3Request{
		verb:         "GET",
		request_uri:  "/" + file_name,
		request_host: "https://" + bucket + ".s3.amazonaws.com",
		request_time: time.Now(),
		signing_uri:  "/" + bucket + "/" + file_name,
	}

	s3_req.sign(s, false)
	return s3_req.do()
}

// Gets an expiring URL for `file_name` in `bucket`, that expires in `expires`. `expires` should
// be parseable by `time.ParseDuration`.
func (s *S3Client) GetURL(bucket, file_name string, expires string) (string, error) {

	expires_duration, err := time.ParseDuration(expires)
	if err != nil {
		return "", err
	}

	expires_time := time.Now().Add(expires_duration)

	s3_req := S3Request{
		verb:         "GET",
		request_uri:  "/" + file_name,
		request_host: "https://" + bucket + ".s3.amazonaws.com",
		request_time: expires_time,
		signing_uri:  "/" + bucket + "/" + file_name,
	}

	signature := s3_req.sign(s, true)

	vals := url.Values{
		"AWSAccessKeyId": []string{s.key},
		"Expires":        []string{strconv.FormatInt(expires_time.Unix(), 10)},
		"Signature":      []string{signature},
	}

	url := s3_req.request_host + s3_req.request_uri + "?" + vals.Encode()

	return url, nil
}

// Uploads `content` of type `content_type` to `file_name` in `bucket`.
func (s *S3Client) Put(bucket, file_name, content_type string, content []byte) (*bytes.Buffer, *http.Response, error) {

	s3_req := S3Request{
		verb:         "PUT",
		request_uri:  "/" + file_name,
		request_host: "https://" + bucket + ".s3.amazonaws.com",
		request_time: time.Now(),
		signing_uri:  "/" + bucket + "/" + file_name,

		content_type: content_type,
		content:      content,
	}

	s3_req.sign(s, false)
	return s3_req.do()
}

func (req *S3Request) sign(c *S3Client, use_unix bool) string {

	formatted_time := req.request_time.Format(time.RFC1123Z)
	if use_unix {
		formatted_time = strconv.FormatInt(req.request_time.Unix(), 10)
	}

	string_to_sign := req.verb + "\n" +
		req.content_md5 + "\n" +
		req.content_type + "\n" +
		formatted_time + "\n" +
		req.signing_uri

	hash := hmac.New(sha1.New, bytes.NewBufferString(c.secret).Bytes())

	buf := bytes.NewBufferString(string_to_sign)

	hash.Write(buf.Bytes())

	encoder := base64.StdEncoding
	result := encoder.EncodeToString(hash.Sum([]byte{}))

	req.signature = c.key + ":" + result

	return result
}

func (s3_req *S3Request) do() (buf *bytes.Buffer, resp *http.Response, err error) {

	buf = new(bytes.Buffer)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := http.DefaultClient
	client.Transport = tr

	content_buffer := bytes.NewBuffer(s3_req.content)

	req, err := http.NewRequest(s3_req.verb, s3_req.request_host+s3_req.request_uri, content_buffer)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "AWS "+s3_req.signature)
	req.Header.Add("Date", s3_req.request_time.Format(time.RFC1123Z))
	if s3_req.content_type != "" {
		req.Header.Add("Content-Type", s3_req.content_type)
	}

	resp, err = client.Do(req)
	if err != nil {
		return
	}

	buf.ReadFrom(resp.Body)

	return buf, resp, nil
}
