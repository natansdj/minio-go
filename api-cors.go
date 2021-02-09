package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/pkg/s3utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Request server for current bucket policy.
func (c Client) GetCors(bucketName string) (string, error) {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return "", err
	}

	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	urlValues.Set("cors", "")

	// Execute GET on bucket to list objects.
	resp, err := c.executeMethod(context.Background(), "GET", requestMetadata{
		bucketName:       bucketName,
		queryValues:      urlValues,
		contentSHA256Hex: emptySHA256Hex,
	})

	defer closeResponse(resp)
	if err != nil {
		return "", err
	}

	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return "", httpRespToErrorResponse(resp, bucketName, "")
		}
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ret := string(buf)
	return ret, err
}

// Saves a new bucket policy.
func (c Client) PutCors(bucketName string, body []byte) error {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}

	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	urlValues.Set("cors", "")

	// Content-length is mandatory for put body request
	bodyReader := strings.NewReader(string(body))
	b, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return err
	}

	header := http.Header{}
	header.Add("Content-Type", "application/xml")

	reqMetadata := requestMetadata{
		bucketName:       bucketName,
		queryValues:      urlValues,
		contentBody:      bodyReader,
		contentLength:    int64(len(b)),
		contentMD5Base64: sumMD5Base64(b),
		customHeader:     header,
	}
	fmt.Println(reqMetadata)
	fmt.Println(string(body))

	// Execute PUT to upload a new bucket body.
	resp, err := c.executeMethod(context.Background(), "PUT", reqMetadata)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusNoContent {
			return httpRespToErrorResponse(resp, bucketName, "")
		}
	}
	return nil
}

// Remove lifecycle from a bucket.
func (c Client) RemoveCors(bucketName string) error {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}
	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	urlValues.Set("cors", "")

	// Execute DELETE on objectName.
	resp, err := c.executeMethod(context.Background(), "DELETE", requestMetadata{
		bucketName:       bucketName,
		queryValues:      urlValues,
		contentSHA256Hex: emptySHA256Hex,
	})
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
