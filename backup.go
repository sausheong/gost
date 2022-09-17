package gost

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"log"

	"github.com/minio/minio-go/v7"
)

func backup(uid string) string {
	return "backup/" + base64.StdEncoding.EncodeToString([]byte(uid)) + ".gob"
}

// Backup all the data for a given unique ID
// There is only 1 backup per unique ID, each time Backup is called, the previous backup is overwritten
func (s *Store) Backup(ctx context.Context, uid string) (err error) {
	all, err := s.GetAll(ctx, uid)
	if err != nil {
		log.Println("Cannot get data during backup:", err)
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(all)
	if err != nil {
		log.Println("Cannot encode gob:", err)
	}
	s.client.PutObject(ctx, s.bucket, backup(uid), &buf, int64(buf.Len()),
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println("Cannot put object:", err)
	}
	return
}

// Load all the data from the backup for a given unique ID
// You can use this to restore data from a backup
// You can also use this to view the data in the backup without restoring it
func (s *Store) Load(ctx context.Context, uid string) (data map[string]any, err error) {
	obj, err := s.client.GetObject(ctx, s.bucket, backup(uid), minio.GetObjectOptions{})
	if err != nil {
		data = make(map[string]any)
		log.Println("Backup doesn't exist:", err)
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

// Restore data from the backup for a given unique ID
// This will overwrite the current data for the unique ID
func (s *Store) Restore(ctx context.Context, uid string) (err error) {
	all, err := s.Load(ctx, uid)
	if err != nil {
		log.Println("Cannot get data during restore:", err)
	}
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
