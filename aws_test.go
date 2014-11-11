package aws

import (
	"net/http"
	"testing"
	"time"
)

var sampleAWSID string = "AKIAIOSFODNN7EXAMPLE"
var sampleAWSSecret string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

var sample_aws_client Client = New(sampleAWSID, sampleAWSSecret)
var sample_s3_client S3Client = sample_aws_client.NewS3()

func TestHTTPSignature(t *testing.T) {
	request_time, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", "Tue, 27 Mar 2007 19:36:42 +0000")

	req := S3Request{
		Verb: "GET",
		URI:  "/photos/puppy.jpg",
		Host: "https://johnsmith.s3.amazonaws.com",

		request_time: request_time,
		signing_uri:  "/johnsmith/photos/puppy.jpg",
		client:       &sample_s3_client,
	}

	signature := req.sign(false)

	correct_signature := "bWq2s1WEIj+Ydj0vQ697zp+IXMU="

	if signature != correct_signature {
		t.Fail()
	}
}

func TestQueryStringSignature(t *testing.T) {
	request_time := time.Unix(1175139620, 0)

	req := S3Request{
		Verb: "GET",
		URI:  "/photos/puppy.jpg",
		Host: "https://johnsmith.s3.amazonaws.com",

		request_time: request_time,
		signing_uri:  "/johnsmith/photos/puppy.jpg",
		client:       &sample_s3_client,
	}

	signature := req.sign(true)

	correct_signature := "NpgCjnDzrM+WFzoENXmpNDUsSn8="

	if signature != correct_signature {
		t.Fail()
	}
}

func TestPutWithHeaders(t *testing.T) {
	request_time, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", "Tue, 27 Mar 2007 19:36:42 +0000")

	req := S3Request{
		Verb: "PUT",
		URI:  "/photos/puppy.jpg",
		Host: "https://johnsmith.s3.amazonaws.com",

		request_time: request_time,
		signing_uri:  "/johnsmith/photos/puppy.jpg.gz",
		client:       &sample_s3_client,

		ContentType: "image/jpeg",
		Content:     []byte("hello, world"),

		Headers: http.Header{
			"Content-Encoding": []string{"gzip"},
		},
	}

	signature := req.sign(false)

	correct_signature := "yaq0CosnkNUgQ0aciRRgEyKtz3c="

	if signature != correct_signature {
		t.Errorf(signature)
		t.Fail()
	}
}
