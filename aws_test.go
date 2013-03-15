package aws

import (
	"testing"
	"time"
)

var sampleAWSID string = "AKIAIOSFODNN7EXAMPLE"
var sampleAWSSecret string = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

var aws_client Client = New(sampleAWSID, sampleAWSSecret)
var s3_client S3Client = aws_client.NewS3()

func TestHTTPSignature(t *testing.T) {
	request_time, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", "Tue, 27 Mar 2007 19:36:42 +0000")

	req := S3Request{
		verb:         "GET",
		request_time: request_time,
		request_uri:  "/photos/puppy.jpg",
		request_host: "https://johnsmith.s3.amazonaws.com",
		signing_uri:  "/johnsmith/photos/puppy.jpg",
	}

	signature := req.sign(&s3_client, false)

	correct_signature := "bWq2s1WEIj+Ydj0vQ697zp+IXMU="

	if signature != correct_signature {
		t.Fail()
	}
}

func TestQueryStringSignature(t *testing.T) {
	request_time := time.Unix(1175139620, 0)

	req := S3Request{
		verb:         "GET",
		request_time: request_time,
		request_uri:  "/photos/puppy.jpg",
		request_host: "https://johnsmith.s3.amazonaws.com",
		signing_uri:  "/johnsmith/photos/puppy.jpg",
	}

	signature := req.sign(&s3_client, true)

	correct_signature := "NpgCjnDzrM+WFzoENXmpNDUsSn8="

	if signature != correct_signature {
		t.Fail()
	}
}
