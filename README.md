# go-aws

**go-aws** is a basic, high-level library to access the Amazon AWS API. Right now, it only works with S3.

## Getting started

First, instanciate the client with your access key ID and secret (you can get those from [this page][security-credentials]):

```go
aws_client := aws.New("<AWSAccessKeyID>", "<AWSAccessKeySecret>")
client := aws_client.NewS3()
```

Then you can upload a file:

```go
file_contents, err := ioutil.ReadFile(file_path)
if err != nil {
	panic("Can't read file")
}

content, response, err := client.Put(bucket, file_name, content_type, file_contents)
```

Download it:

```go
content, response, err = client.Get(bucket, file_name)
```

Or get the URL of that file:

```go
url, err := client.GetURL(bucket, file_name, expires_in)
```

## Installation

Install this package by typing `go get github.com/jimmysawczuk/go-aws` in your terminal. You can then use it in your import statement like so:

```go
import (
	"github.com/jimmysawczuk/go-aws"
)
```

## License

	The MIT License (MIT)
	Copyright (C) 2013 by Jimmy Sawczuk

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in
	all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE


  [security-credentials]: https://portal.aws.amazon.com/gp/aws/securityCredentials