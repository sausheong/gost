package gost

import (
	"context"
	"encoding/gob"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Store is the main struct for the database
type Store struct {
	client *minio.Client
	bucket string
}

// Create a new store
func NewStore(key string, secret string, endpoint string, useSSL bool, bucket string) (s *Store, err error) {
	s = &Store{
		bucket: bucket,
	}
	ctx := context.Background()
	s.client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(key, secret, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("Cannot initiate client:", err)
	}

	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatalln("Cannot check if bucket exists:", err)
	}
	if !exists {
		// Make gost bucket
		err = s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln("Gost bucket doesn't exist but we can't make the bucket either:", err)
		}
	}
	return
}

// Register a struct to be stored in the database
func Register(data any) {
	gob.Register(data)
}
