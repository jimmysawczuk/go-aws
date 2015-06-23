package aws

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

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

var RootURL string = "s3.amazonaws.com"

type S3Client struct {
	key    string
	secret string

	http_client *http.Client
}

type S3Request struct {
	Verb string

	URI  string
	Host string

	Headers     http.Header
	ContentType string

	Content []byte

	content_md5  string
	request_time time.Time
	signing_uri  string
	signature    string
	client       *S3Client
}

// Creates a new S3Client using the AWS information in the Client
func (c Client) NewS3() S3Client {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := http.DefaultClient
	client.Transport = tr

	s := S3Client{
		key:         c.key,
		secret:      c.secret,
		http_client: client,
	}
	return s
}

func (s *S3Client) NewS3Request(verb, bucket, file_name string) (S3Request, error) {
	req := S3Request{
		Verb: verb,
		URI:  "/" + file_name,
		Host: "https://" + bucket + "." + RootURL,

		request_time: time.Now(),
		signing_uri:  "/" + bucket + "/" + file_name,
		client:       s,

		Headers: make(http.Header),
	}

	return req, nil
}

// Performs an authenticated GET request on file_name in bucket.
func (s *S3Client) Get(bucket, file_name string) (*bytes.Buffer, *http.Response, error) {

	req, err := s.NewS3Request("GET", bucket, file_name)
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating GET request: %s", err)
	}

	return req.Exec()
}

// Gets an expiring URL for `file_name` in `bucket`, that expires in `expires`. `expires` should
// be parseable by `time.ParseDuration`.
func (s *S3Client) GetURL(bucket, file_name string, expires string) (string, error) {

	expires_duration, err := time.ParseDuration(expires)
	if err != nil {
		return "", err
	}

	expires_time := time.Now().Add(expires_duration)

	req := S3Request{
		Verb: "GET",
		URI:  "/" + file_name,
		Host: "https://" + bucket + ".s3.amazonaws.com",

		request_time: expires_time,
		signing_uri:  "/" + bucket + "/" + file_name,
		client:       s,
	}

	signature := req.sign(true)

	vals := url.Values{
		"AWSAccessKeyId": []string{s.key},
		"Expires":        []string{strconv.FormatInt(expires_time.Unix(), 10)},
		"Signature":      []string{signature},
	}

	url := req.Host + req.URI + "?" + vals.Encode()

	return url, nil
}

// Uploads `content` of type `content_type` to `file_name` in `bucket`.
func (s *S3Client) Put(bucket, file_name, content_type string, content []byte) (*bytes.Buffer, *http.Response, error) {

	req, err := s.NewS3Request("PUT", bucket, file_name)
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating PUT request: %s", err)
	}

	req.ContentType = content_type
	req.Content = content

	req.sign(false)

	return req.Exec()
}

func (req *S3Request) sign(use_unix bool) string {

	formatted_time := req.request_time.Format(time.RFC1123Z)
	if use_unix {
		formatted_time = strconv.FormatInt(req.request_time.Unix(), 10)
	}

	canonical_headers := ""
	if len(req.Headers) > 0 {
		for header, values := range req.Headers {
			if strings.HasPrefix(strings.ToLower(header), "x-amz-") {
				canonical_headers += strings.ToLower(header) + ":" + strings.ToLower(strings.Join(values, ",")) + "\n"
			}
		}
	}

	string_to_sign := req.Verb + "\n" +
		req.content_md5 + "\n" +
		req.ContentType + "\n" +
		formatted_time + "\n" +
		canonical_headers +
		req.signing_uri

	hash := hmac.New(sha1.New, bytes.NewBufferString(req.client.secret).Bytes())

	buf := bytes.NewBufferString(string_to_sign)

	hash.Write(buf.Bytes())

	encoder := base64.StdEncoding
	result := encoder.EncodeToString(hash.Sum([]byte{}))

	req.signature = req.client.key + ":" + result

	return result
}

func (this *S3Request) Exec() (buf *bytes.Buffer, resp *http.Response, err error) {

	if this.signature == "" {
		this.sign(false)
	}

	buf = new(bytes.Buffer)

	content_buffer := bytes.NewBuffer(this.Content)

	req, err := http.NewRequest(this.Verb, this.Host+this.URI, content_buffer)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "AWS "+this.signature)
	req.Header.Add("Date", this.request_time.Format(time.RFC1123Z))
	if this.ContentType != "" {
		req.Header.Add("Content-Type", this.ContentType)
	}

	for header, vals := range this.Headers {
		for _, val := range vals {
			req.Header.Add(header, val)
		}
	}

	resp, err = this.client.http_client.Do(req)
	if err != nil {
		return
	}

	buf.ReadFrom(resp.Body)

	return buf, resp, nil
}
