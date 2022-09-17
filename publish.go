package gost

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
)

var policyFormat = `{
	"Version": "2012-10-17",
	"Statement": [
		{
		"Effect": "%s",
		"Principal": "*",
		"Action": "s3:GetObject",
		"Resource": "arn:aws:s3:::%s/%s/*",
		"Sid": ""
		}
	]
}`

// Publish data and make it publicly available
func (s *Store) Publish(ctx context.Context, filename string, contentType string, data []byte) (location string, err error) {
	_, err = s.client.PutObject(ctx, s.bucket, "public/"+filename, bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType})
	location = s.client.EndpointURL().String() + "/" + s.bucket + "/public/" + filename
	if err != nil {
		log.Println("Cannot publish object:", err)
	}
	return
}

// Delete published data
func (s *Store) Unpublish(ctx context.Context, filename string) (err error) {
	err = s.client.RemoveObject(ctx, s.bucket, "public/"+filename, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println("Cannot unpublish object:", err)
	}
	return
}

// Programmatically set up bucket folder /public to be publicly readable
func (s *Store) AllowPublic(ctx context.Context) (err error) {
	policy := fmt.Sprintf(policyFormat, "Allow", s.bucket, "public")
	err = s.client.SetBucketPolicy(ctx, s.bucket, policy)
	if err != nil {
		log.Println("Cannot set bucket policy:", err)
	}
	return
}

// Programmatically set up bucket folder /public to be private
func (s *Store) DenyPublic(ctx context.Context) (err error) {
	policy := fmt.Sprintf(policyFormat, "Deny", s.bucket, "public")
	err = s.client.SetBucketPolicy(ctx, s.bucket, policy)
	if err != nil {
		log.Println("Cannot set bucket policy:", err)
	}
	return
}

func (s *Store) IsPublic(ctx context.Context) (isPublic bool, err error) {
	policy, err := s.client.GetBucketPolicy(ctx, s.bucket)
	if err != nil {
		log.Println("Cannot get bucket policy:", err)
		return
	}
	isPublic = policy == fmt.Sprintf(policyFormat, "Allow", s.bucket, "public")
	return
}
