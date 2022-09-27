package gost

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	"github.com/minio/minio-go/v7"
)

// Put an object in the database, with an associated a unique ID
func (s *Store) PutObject(ctx context.Context, uid string, obj any) (err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(&obj)
	if err != nil {
		log.Println("Cannot encode gob:", err)
		return
	}
	s.client.PutObject(ctx, s.bucket, uid, &buf, int64(buf.Len()),
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println("Cannot put object:", err)
	}
	return
}

// Get a specific piece of data for a given unique ID
func (s *Store) GetObject(ctx context.Context, uid string) (obj any, err error) {
	mObj, err := s.client.GetObject(ctx, s.bucket, uid, minio.GetObjectOptions{})
	if err != nil {
		log.Println("Cannot get object:", err)
		return
	}
	defer mObj.Close()
	decoder := gob.NewDecoder(mObj)
	err = decoder.Decode(&obj)
	if err != nil {
		log.Println("Cannot decode object:", err)
	}
	return
}

// Delete a specific piece of data for a given unique ID
func (s *Store) DeleteObject(ctx context.Context, uid string) (err error) {
	err = s.client.RemoveObject(ctx, s.bucket, uid, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println("Cannot delete object:", err)
	}
	return
}
