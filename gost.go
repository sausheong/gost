package gost

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"log"

	"github.com/minio/minio-go/v7"
)

// get the name of the object to store
func name(uid string) string {
	return "data/" + base64.StdEncoding.EncodeToString([]byte(uid)) + ".gob"
}

// Put a piece of data in the database, with a unique ID
// Each piece of data is associated with a key
func (s *Store) Put(ctx context.Context, uid string, key string, data any) (err error) {
	all, err := s.GetAll(ctx, uid)
	if err != nil {
		log.Println("Cannot get data during put:", err)
		return
	}
	all[key] = data
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(all)
	if err != nil {
		log.Println("Cannot encode gob:", err)
		return
	}
	s.client.PutObject(ctx, s.bucket, name(uid), &buf, int64(buf.Len()),
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println("Cannot put object:", err)
	}

	return
}

// Get all the data for a given unique ID
func (s *Store) GetAll(ctx context.Context, uid string) (data map[string]any, err error) {
	obj, err := s.client.GetObject(ctx, s.bucket, name(uid), minio.GetObjectOptions{})
	if err != nil {
		data = make(map[string]any)
		log.Println("Data doesn't exist:", err)
		return
	}
	defer obj.Close()
	_, err = obj.Stat()
	if err != nil {
		data = make(map[string]any)
		log.Println("Object doesn't exist:", err, obj)
		return
	}

	decoder := gob.NewDecoder(obj)
	err = decoder.Decode(&data)
	if err != nil {
		log.Println("Cannot decode data:", err)
	}
	return
}

// Get a specific piece of data for a given unique ID
func (s *Store) Get(ctx context.Context, uid string, key string) (data any, err error) {
	all, err := s.GetAll(ctx, uid)
	if err != nil {
		return
	}
	data = all[key]
	return
}

// Delete a specific piece of data for a given unique ID
func (s *Store) Delete(ctx context.Context, uid string, key string) (err error) {
	all, err := s.GetAll(ctx, uid)
	if err != nil {
		return
	}
	delete(all, key)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(all)
	if err != nil {
		log.Println("Cannot encode gob:", err)
	}
	s.client.PutObject(ctx, s.bucket, name(uid), &buf, int64(buf.Len()),
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println("Cannot put object:", err)
	}
	return
}

// Delete all data for a given unique ID
func (s *Store) DeleteAll(ctx context.Context, uid string) (err error) {
	empty := make(map[string]any)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(empty)
	if err != nil {
		log.Println("Cannot encode gob:", err)
	}
	s.client.PutObject(ctx, s.bucket, name(uid), &buf, int64(buf.Len()),
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println("Cannot put object:", err)
	}
	return
}
